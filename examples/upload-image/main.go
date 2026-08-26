// Upload an image to Cloudinary and print the fields worth keeping.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md; `npx
//     @cloudinary/cloud` provisions a working cloud with no signup.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./upload-image
//
// In a real project you would upload a file the user supplied — a local path, an
// *os.File, or a *multipart.FileHeader from an HTTP request — rather than the demo
// URL hardcoded below, and you would store result.AssetID in your own database.
//
// Docs: docs/upload-image.md
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
	publicID  = "examples/uploaded-sample"
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

	result, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: api.Bool(true),
	})
	// Check err first: a pre-flight failure returns a nil result.
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	// Then check the API's own verdict — a rejection arrives with err == nil.
	if result.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", result.Error.Message)
	}

	fmt.Println("asset_id: ", result.AssetID) // immutable; store this
	fmt.Println("public_id:", result.PublicID)
	fmt.Println("url:      ", result.SecureURL)
	fmt.Printf("size:      %dx%d %s, %d bytes\n",
		result.Width, result.Height, result.Format, result.Bytes)

	return nil
}
