// Search for assets, read one by its immutable asset ID, and update its tags.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//   - Run ./upload-image first, so there is something to find.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./search-and-manage-assets
//
// In a real project you would store AssetID at upload time and look assets up by it,
// rather than searching for them. The Admin and Search APIs are rate-limited — keep
// them out of request paths.
//
// This example deliberately performs no deletion. See docs/search-and-manage-assets.md
// for the delete calls and why prefix deletion deserves care.
//
// Docs: docs/search-and-manage-assets.md
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
)

const (
	// A field query. Note that "folder:examples" matches nothing on dynamic-folder
	// environments even when assets are there — use asset_folder, or filter on
	// public_id as here.
	expression = `resource_type:image AND public_id:examples/*`

	maxResults = 30
	newTag     = "reviewed"
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

	query := search.Query{
		Expression: expression,
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: maxResults,
	}

	found, err := cld.Admin.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("search transport: %w", err)
	}
	if found.Error.Message != "" {
		return fmt.Errorf("search rejected: %s", found.Error.Message)
	}

	fmt.Printf("total matches: %d\n", found.TotalCount)
	for _, asset := range found.Assets {
		fmt.Printf("  %s  %s  %d bytes\n", asset.AssetID, asset.PublicID, asset.Bytes)
	}

	// A valid expression that matches nothing is not an error, so say so explicitly
	// rather than letting an empty list look like a failure.
	if len(found.Assets) == 0 {
		fmt.Println("No assets matched. Run ./upload-image first.")
		return nil
	}

	// Look the first result up by its immutable handle. AssetByAssetID is the only
	// asset-id entry point in this SDK; everything that mutates needs the public ID.
	assetID := found.Assets[0].AssetID
	details, err := cld.Admin.AssetByAssetID(ctx, admin.AssetByAssetIDParams{AssetID: assetID})
	if err != nil {
		return fmt.Errorf("lookup transport: %w", err)
	}
	if details.Error.Message != "" {
		return fmt.Errorf("lookup rejected: %s", details.Error.Message)
	}
	fmt.Printf("looked up %s -> public_id %s\n", assetID, details.PublicID)

	// Updates take the public ID read off the lookup above.
	updated, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID: details.PublicID,
		Tags:     api.CldAPIArray{newTag},
	})
	if err != nil {
		return fmt.Errorf("update transport: %w", err)
	}
	if updated.Error.Message != "" {
		return fmt.Errorf("update rejected: %s", updated.Error.Message)
	}
	fmt.Printf("tags now: %v\n", updated.Tags)

	// Pagination: feed NextCursor back until it comes back empty.
	if found.NextCursor != "" {
		query.NextCursor = found.NextCursor
		page2, err := cld.Admin.Search(ctx, query)
		if err != nil {
			return fmt.Errorf("search page 2 transport: %w", err)
		}
		if page2.Error.Message != "" {
			return fmt.Errorf("search page 2 rejected: %s", page2.Error.Message)
		}
		fmt.Printf("second page: %d asset(s)\n", len(page2.Assets))
	}

	return nil
}
