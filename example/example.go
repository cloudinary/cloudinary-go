package main

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"log"
)

func main() {
	// Start by creating a new instance of Cloudinary using CLOUDINARY_URL environment variable.
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromURL().
	var cld, err = cloudinary.New()
	if err != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", err)
	}

	var ctx = context.Background()

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
	uploadResult, err := cld.Upload.Upload(
		ctx,
		"https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png",
		uploader.UploadParams{PublicID: "logo"})
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}

	log.Println(uploadResult.SecureURL)
	// Prints something like:
	// https://res.cloudinary.com/<your cloud name>/image/upload/v1615875158/logo.png

	// uploadResult contains useful information about the asset, like Width, Height, Format, etc.
	// See uploader.UploadResult struct for more details.

	// We can also build an image URL using the Public ID.
	image, err := cld.Image("logo")
	if err != nil {
		log.Fatalf("Failed to build image URL, %v\n", err)
	}

	// Image can be further transformed and optimized as follows:
	image.Transformation = "c_scale,w_500/f_auto/q_auto"
	// Here the image is scaled to the width of 500 pixes. Format and quality are set to "auto".

	imageURL, err := image.String()
	if err != nil {
		log.Fatalf("Failed to serialize image URL, %v\n", err)
	}

	log.Printf("Image URL: %s", imageURL)
	// Prints something like:
	// https://res.cloudinary.com/<your cloud name>/image/upload/c_scale,w_500/f_auto/q_auto/logo

	// Now we can use Admin API to see the details about the asset.
	// The request can be customised by providing AssetParams.
	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "logo"})
	if err != nil {
		log.Fatalf("Failed to get asset details, %v\n", err)
	}

	// Print some basic information about the asset.
	log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)

	// Cloudinary also provides a very flexible Search API for filtering and retrieving
	// information on all the assets in your account with the help of query expressions
	// in a Lucene-like query language.
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
