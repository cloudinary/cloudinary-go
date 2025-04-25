package main

import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func main() {

	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}
	ctx := context.Background()

	uploadResult := uploadImage(cld, ctx)
	log.Println(uploadResult.SecureURL)

	imageURL := buildImageURL(cld)
	log.Printf("Image URL: %s", imageURL)

	getAssetDetails(cld, ctx)

	searchAssets(cld, ctx)
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
		"./1.jpg",
		uploader.UploadParams{PublicID: "logo"},
	)
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
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
		log.Fatalf("Failed to build image URL, %v\n", err)
	}

	image.Transformation = "c_scale,w_500/f_auto/q_auto"

	imageURL, err := image.String()
	if err != nil {
		log.Fatalf("Failed to serialize image URL, %v\n", err)
	}

	return imageURL
}

// Now we can use Admin API to see the details about the asset.
// The request can be customised by providing AssetParams.
func getAssetDetails(cld *cloudinary.Cloudinary, ctx context.Context) {
	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "logo"})
	if err != nil {
		log.Fatalf("Failed to get asset details, %v\n", err)
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
		log.Fatalf("Failed to search for assets, %v\n", err)
	}

	log.Printf("Assets found: %v\n", searchResult.TotalCount)

	for _, asset := range searchResult.Assets {
		log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)
	}
}
