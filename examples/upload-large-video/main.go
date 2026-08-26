// Upload a large video. Chunking is automatic in this SDK — there is no separate
// UploadLarge method — so this is the same Upload call used for images, with
// ResourceType set to video.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./upload-large-video
//
// Downloads a ~9 MB sample video to the working directory on first run so the example
// needs no arguments. In a real project you would upload the user's file and set
// UploadTimeout to something appropriate for your largest expected asset.
//
// Docs: docs/upload-large-video.md
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sampleVideoURL = "https://res.cloudinary.com/demo/video/upload/dog.mp4"
	localPath      = "dog.mp4"
	publicID       = "examples/uploaded-large-video"

	// Cloudinary requires chunks of at least 5 MB; stay at or above that.
	chunkSize = 6 * 1024 * 1024

	uploadTimeoutSeconds = 600
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

	// UploadTimeout defaults to 0, which falls back to API.Timeout (60s) — too short
	// for a large video on a slow link.
	cld.Config.API.ChunkSize = chunkSize
	cld.Config.API.UploadTimeout = uploadTimeoutSeconds

	ctx := context.Background()

	if err := ensureLocalVideo(); err != nil {
		return err
	}

	result, err := cld.Upload.Upload(ctx, localPath, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: api.Video, // required: the default is image
		Overwrite:    api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if result.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", result.Error.Message)
	}

	// A chunked upload that stalled reports success with an empty result, so assert
	// the response actually describes a finished asset.
	if result.PublicID == "" || result.SecureURL == "" {
		return fmt.Errorf("upload did not complete: %+v", result.Response)
	}

	fmt.Println("asset_id: ", result.AssetID)
	fmt.Println("public_id:", result.PublicID)
	fmt.Println("bytes:    ", result.Bytes)
	fmt.Println("url:      ", result.SecureURL)
	if result.PlaybackURL != "" {
		fmt.Println("hls:      ", result.PlaybackURL)
	}

	return nil
}

// ensureLocalVideo downloads the demo sample once, so the example runs with no arguments.
func ensureLocalVideo() error {
	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	fmt.Println("downloading sample video...")

	resp, err := http.Get(sampleVideoURL)
	if err != nil {
		return fmt.Errorf("download sample: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download sample: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write sample: %w", err)
	}

	return nil
}
