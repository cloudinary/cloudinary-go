// Build a responsive srcset for an image.
//
// This module generates delivery URLs, so a srcset is a loop over the widths you
// want plus the markup your template needs. URL generation is local: no network
// call, and only a cloud name is required.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./responsive-srcset
//
// In a real project the asset already exists, so you would skip the upload and build
// the srcset from a public ID you looked up. Pick widths that match your layout's
// breakpoints rather than these defaults, and set `sizes` to tell the browser how
// much space the image occupies.
//
// Docs: docs/transform-and-deliver-image.md
package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sourceURL = "https://res.cloudinary.com/demo/image/upload/sample.jpg"
	publicID  = "examples/responsive-image"
	sizes     = "(max-width: 800px) 100vw, 800px"
)

// widths to generate. Match these to your layout's breakpoints.
var widths = []int{400, 800, 1200, 1600}

const imgTag = `<img src="{{.Fallback}}" srcset="{{.SrcSet}}" sizes="{{.Sizes}}" alt="Sample">`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	cld, err := cloudinary.New()
	if err != nil {
		return err
	}
	ctx := context.Background()

	uploaded, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}

	srcset, err := buildSrcSet(cld, publicID, widths)
	if err != nil {
		return err
	}
	fmt.Println("srcset:", srcset)

	// A fallback for browsers that ignore srcset, at the middle width.
	fallback, err := scaledURL(cld, publicID, widths[len(widths)/2])
	if err != nil {
		return err
	}
	fmt.Println("src:   ", fallback)

	var buf bytes.Buffer
	err = template.Must(template.New("img").Parse(imgTag)).Execute(&buf, struct {
		Fallback, SrcSet, Sizes string
	}{fallback, srcset, sizes})
	if err != nil {
		return err
	}
	fmt.Println(buf.String())

	return nil
}

// buildSrcSet returns a srcset attribute value: one "<url> <width>w" entry per width.
func buildSrcSet(cld *cloudinary.Cloudinary, publicID string, widths []int) (string, error) {
	entries := make([]string, 0, len(widths))
	for _, w := range widths {
		url, err := scaledURL(cld, publicID, w)
		if err != nil {
			return "", err
		}
		entries = append(entries, fmt.Sprintf("%s %dw", url, w))
	}
	return strings.Join(entries, ", "), nil
}

// scaledURL builds one delivery URL scaled to the given width, with automatic
// format and quality.
func scaledURL(cld *cloudinary.Cloudinary, publicID string, width int) (string, error) {
	image, err := cld.Image(publicID)
	if err != nil {
		return "", err
	}
	image.Transformation = fmt.Sprintf("c_scale,w_%d/f_auto,q_auto", width)
	return image.String()
}
