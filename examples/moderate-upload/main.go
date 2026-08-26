// Upload an asset into a manual moderation queue, list the queue, and record a
// decision.
//
// IMPORTANT: a "pending" asset is delivered normally — its URL returns HTTP 200. The
// moderation status is metadata your application gates on; Cloudinary does not
// withhold the asset by default. Blocking non-approved assets is a product-environment
// setting configured by Cloudinary support, not an upload parameter.
//
// This example uses "manual" moderation, which needs no add-on. The automatic kinds
// (aws_rek, webpurify, google_video_moderation, perception_point, duplicate) each
// require a human to register the add-on in the console first, and some also require
// accepting the provider's terms of service. Neither step has an API.
// See https://cloudinary.com/documentation/cloudinary_add_ons.md
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./moderate-upload
//
// In a real project steps 2 and 3 are separate requests, days apart: a reviewer UI
// lists the queue, and a human action triggers the update.
//
// Docs: docs/moderate-upload.md
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sourceURL = "https://res.cloudinary.com/demo/image/upload/sample.jpg"
	publicID  = "examples/moderated-upload"

	moderationKind = "manual"
	maxResults     = 50
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

	// 1. Upload into the moderation queue. The asset starts as "pending".
	uploaded, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:   publicID,
		Moderation: moderationKind,
		Overwrite:  api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}

	// Moderation is a slice and is empty for non-moderated assets — guard the index.
	if len(uploaded.Moderation) > 0 {
		fmt.Printf("uploaded: kind=%s status=%s\n",
			uploaded.Moderation[0].Kind, uploaded.Moderation[0].Status)
	}
	fmt.Println("deliverable right now:", uploaded.SecureURL)

	// 2. A review UI lists what is waiting. Status is a plain string here, while
	//    UpdateAssetParams.ModerationStatus is typed — hence the conversion.
	queue, err := cld.Admin.AssetsByModeration(ctx, admin.AssetsByModerationParams{
		Kind:       moderationKind,
		Status:     string(api.Pending),
		MaxResults: maxResults,
	})
	if err != nil {
		return fmt.Errorf("list queue transport: %w", err)
	}
	if queue.Error.Message != "" {
		return fmt.Errorf("list queue rejected: %s", queue.Error.Message)
	}
	fmt.Printf("pending review: %d asset(s)\n", len(queue.Assets))

	// 3. Record the reviewer's decision. UpdateAsset needs the public ID.
	updated, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID:         publicID,
		ModerationStatus: api.Approved,
	})
	if err != nil {
		return fmt.Errorf("update transport: %w", err)
	}
	if updated.Error.Message != "" {
		return fmt.Errorf("update rejected: %s", updated.Error.Message)
	}

	fmt.Println("approved: ", updated.SecureURL)
	fmt.Println("Gate delivery on your own copy of this status, not on the URL failing.")

	return nil
}
