// Build Cloudinary delivery URLs for a video: a scaled MP4, a poster frame, and an
// HLS manifest for adaptive bitrate streaming.
//
// cld.Video returns an asset whose String() gives you a delivery URL — everything a
// <video> tag needs. See docs/transform-and-deliver-video.md for composing the markup
// with html/template, or use the Cloudinary Video Player on the client.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./transform-and-deliver-video
//
// The first request to a new streaming URL is slow because the renditions are being
// generated. In production, request them eagerly at upload time with EagerAsync and a
// NotificationURL.
//
// Docs: docs/transform-and-deliver-video.md
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
	sourceURL = "https://res.cloudinary.com/demo/video/upload/dog.mp4"
	publicID  = "examples/delivered-video"
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

	// ResourceType must be video, or the asset lands under /image/upload/ and no
	// video transformation applies to it.
	uploaded, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: api.Video,
		Overwrite:    api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}

	// A scaled MP4 for direct playback.
	playback, err := videoURL(cld, publicID, "c_scale,w_640/q_auto")
	if err != nil {
		return err
	}
	fmt.Println("playback:", playback)

	// A poster frame two seconds in. The image extension goes in the public ID;
	// Asset.Suffix is a separate feature (SEO short URLs, private CDN).
	poster, err := videoURL(cld, publicID+".jpg", "so_2,c_fill,w_640")
	if err != nil {
		return err
	}
	fmt.Println("poster:  ", poster)

	// An HLS manifest. Use .mpd for DASH.
	hls, err := videoURL(cld, publicID+".m3u8", "sp_auto")
	if err != nil {
		return err
	}
	fmt.Println("hls:     ", hls)

	// Cloudinary often supplies a ready-made HLS URL on the upload result.
	if uploaded.PlaybackURL != "" {
		fmt.Println("playback_url:", uploaded.PlaybackURL)
	}

	// Video metadata is not on the typed struct — read it from the raw response.
	// Note Response is a *pointer* to the decoded map.
	if raw, ok := uploaded.Response.(*map[string]interface{}); ok {
		fmt.Printf("duration: %v seconds\n", (*raw)["duration"])
	}

	return nil
}

// videoURL builds a delivery URL for a video asset.
func videoURL(cld *cloudinary.Cloudinary, id, transformation string) (string, error) {
	video, err := cld.Video(id)
	if err != nil {
		return "", err
	}
	video.Transformation = transformation
	return video.String()
}
