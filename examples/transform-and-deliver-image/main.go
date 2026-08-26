// Build Cloudinary delivery URLs for an image.
//
// URL generation is local: no network call, and only a cloud name is required.
// This example uploads first so it has an asset of its own to address, then builds
// several URLs from it.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./transform-and-deliver-image
//
// In a real project the asset already exists, so you would skip the upload and build
// URLs from a public ID you looked up. Transformations are plain strings — use the
// transformation reference or the cloudinary-transformations skill to compose them.
//
// Docs: docs/transform-and-deliver-image.md
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sourceURL = "https://res.cloudinary.com/demo/image/upload/sample.jpg"
	publicID  = "examples/transformed-sample"

	thumbnail = "c_thumb,g_auto,h_200,w_200/f_auto,q_auto"
	banner    = "c_fill,g_auto,h_720,w_1280/co_white,l_text:Arial_64_bold:SALE,g_south_east,x_24,y_24/f_auto,q_auto"
)

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

	// A square, auto-focused, auto-format thumbnail.
	thumbURL, err := buildURL(cld, publicID, thumbnail, 0)
	if err != nil {
		return err
	}
	fmt.Println("thumbnail:", thumbURL)

	// Chained components: each runs on the output of the previous one, so order matters.
	bannerURL, err := buildURL(cld, publicID, banner, 0)
	if err != nil {
		return err
	}
	fmt.Println("banner:   ", bannerURL)

	// Pinning the version busts CDN caches after a re-upload to the same public ID.
	versionedURL, err := buildURL(cld, publicID, "f_auto,q_auto", uploaded.Version)
	if err != nil {
		return err
	}
	fmt.Println("versioned:", versionedURL)

	return nil
}

// buildURL returns a delivery URL for publicID. A version of 0 means "unset".
func buildURL(cld *cloudinary.Cloudinary, id, transformation string, version int) (string, error) {
	image, err := cld.Image(id)
	if err != nil {
		return "", err
	}
	image.Transformation = transformation
	if version != 0 {
		image.Version = version
	}
	return image.String()
}
