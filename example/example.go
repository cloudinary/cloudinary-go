package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)


const (
	imageFilePath = "https://res.cloudinary.com/demo/image/upload/sample.jpg"
	videoFilePath = "https://res.cloudinary.com/demo/video/upload/dog.mp4"
)

func main() {

	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("failed to initialize Cloudinary: %v", err)
	}
	ctx := context.Background()

	uploadResult := uploadImage(cld, ctx)
	log.Println(uploadResult.SecureURL)

	imageURL := buildImageURL(cld)
	log.Printf("Image URL: %s", imageURL)

	getAssetDetails(cld, ctx)

	searchAssets(cld, ctx)

	// Generate responsive srcset for the "logo" image
	srcset, err := generateResponsiveSrcSet(cld, "logo")
	if err != nil {
		log.Fatalf("failed to build srcset: %v", err)
	}
	log.Printf("SrcSet: %s", srcset)

	// Upload a video with transformations applied on upload
	videoResult := uploadVideoWithTransformations(cld, ctx)
	log.Println(videoResult.SecureURL)

	// Delete a single asset by Public ID
	deleteAsset(cld, ctx, "logo")

	// Bulk delete multiple assets by Public IDs
	bulkDeleteAssets(cld, ctx, []string{"old_img1", "old_img2", "old_img3"})

	// List assets with pagination (first 5 per page)
	listAssetsWithPagination(cld, ctx, 5)
}

// Upload an image to your Cloudinary account from a specified URL.
//
// Alternatively you can provide a path to a local file on your filesystem,
// base64 encoded string, io.Reader and more.
//
// For additional information see:
// https://cloudinary.com/documentation/upload_images
//
// Upload can be greatly customized by specifying uploader.UploadParams,
// in this case we set the Public ID of the uploaded asset to "logo".
func uploadImage(cld *cloudinary.Cloudinary, ctx context.Context) *uploader.UploadResult {
	uploadResult, err := cld.Upload.Upload(
		ctx,
		imageFilePath,
		uploader.UploadParams{PublicID: "logo"},
	)
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}

	// uploadResult contains useful information about the asset, like Width, Height, Format, etc.
	// See uploader.UploadResult struct for more details.
	return uploadResult
}

// We can also build an image URL using the Public ID.
//
// Image can be further transformed and optimized as follows:
// Here the image is scaled to the width of 500 pixels. Format and quality are set to "auto".
func buildImageURL(cld *cloudinary.Cloudinary) string {
	image, err := cld.Image("logo")
	if err != nil {
		log.Fatalf("failed to build image URL: %v", err)
	}

	image.Transformation = "c_scale,w_500/f_auto/q_auto"

	imageURL, err := image.String()
	if err != nil {
		log.Fatalf("failed to serialize image URL: %v", err)
	}

	return imageURL
}

// Now we can use Admin API to see the details about the asset.
// The request can be customised by providing AssetParams.
func getAssetDetails(cld *cloudinary.Cloudinary, ctx context.Context) {
	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "logo"})
	if err != nil {
		log.Fatalf("failed to get asset details: %v", err)
	}

	// Print some basic information about the asset.
	log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)
}

// Cloudinary also provides a very flexible Search API for filtering and retrieving
// information on all the assets in your account with the help of query expressions
// in a Lucene-like query language.
func searchAssets(cld *cloudinary.Cloudinary, ctx context.Context) {
	searchQuery := search.Query{
		Expression: "resource_type:image AND uploaded_at>1d AND bytes<1m",
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 30,
	}

	searchResult, err := cld.Admin.Search(ctx, searchQuery)
	if err != nil {
		log.Fatalf("failed to search for assets: %v", err)
	}

	log.Printf("Assets found: %v\n", searchResult.TotalCount)

	for _, asset := range searchResult.Assets {
		log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)
	}
}

// Generate a responsive srcset string for a given public ID by building URLs at multiple widths.
func generateResponsiveSrcSet(cld *cloudinary.Cloudinary, publicID string) (string, error) {
	widths := []int{200, 400, 800, 1200}
	var parts []string

	for _, w := range widths {
		img, err := cld.Image(publicID)
		if err != nil {
			return "", fmt.Errorf("failed to initialize image %s: %w", publicID, err)
		}
		img.Transformation = fmt.Sprintf("c_scale,w_%d/f_auto/q_auto", w)

		url, err := img.String()
		if err != nil {
			return "", fmt.Errorf("failed to build URL for width %d: %w", w, err)
		}
		parts = append(parts, fmt.Sprintf("%s %dw", url, w))
	}

	return strings.Join(parts, ", "), nil
}

// Upload a video with transformations applied on upload to generate posters or clips.
func uploadVideoWithTransformations(cld *cloudinary.Cloudinary, ctx context.Context) *uploader.UploadResult {
	uploadResult, err := cld.Upload.Upload(
		ctx,
		videoFilePath,
		uploader.UploadParams{
			PublicID:     "promo_clip",
			Folder:       "videos/promos",
			ResourceType: "video",
			Eager:        "c_fill,h_360,w_640,b_black|c_crop,ar_16:9,e_volume:0.5,du_15",
			Tags:         []string{"video", "promo"},
		},
	)
	if err != nil {
		log.Fatalf("failed to upload video: %v", err)
	}

	// uploadResult.SecureURL points to the original video; uploadResult.Eager to derivatives.
	return uploadResult
}

// Delete a single asset by its Public ID.
func deleteAsset(cld *cloudinary.Cloudinary, ctx context.Context, publicID string) {
	_, err := cld.Upload.Destroy(
		ctx,
		uploader.DestroyParams{PublicID: publicID, ResourceType: "image"},
	)
	if err != nil {
		log.Fatalf("failed to delete asset %s: %v", publicID, err)
	}

	// Asset deleted successfully.
}

// Bulk delete multiple assets by their Public IDs.
func bulkDeleteAssets(cld *cloudinary.Cloudinary, ctx context.Context, publicIDs []string) {
	resp, err := cld.Admin.DeleteAssets(
		ctx,
		admin.DeleteAssetsParams{PublicIDs: publicIDs},
	)
	if err != nil {
		log.Fatalf("failed to bulk delete assets: %v", err)
	}

	// Print how many were deleted.
	log.Printf("Deleted assets count: %d", len(resp.Deleted))
}

// List all assets in pages of up to maxResults, using cursor-based pagination.
func listAssetsWithPagination(cld *cloudinary.Cloudinary, ctx context.Context, maxResults int) {
	nextCursor := ""

	for {
		page, err := cld.Admin.Assets(
			ctx,
			admin.AssetsParams{MaxResults: maxResults, NextCursor: nextCursor},
		)
		if err != nil {
			log.Fatalf("failed to list assets: %v", err)
		}

		for _, asset := range page.Assets {
			log.Printf("Public ID: %v, URL: %v", asset.PublicID, asset.SecureURL)
		}

		if page.NextCursor == "" {
			break // no more assets
		}
		nextCursor = page.NextCursor
	}
}
