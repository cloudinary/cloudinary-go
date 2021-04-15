[![Build Status](https://travis-ci.com/cloudinary/cloudinary-go.svg)](https://travis-ci.com/cloudinary/cloudinary-go) 
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudinary/cloudinary-go)](https://goreportcard.com/report/github.com/cloudinary/cloudinary-go) 
[![PkgGoDev](https://pkg.go.dev/badge/github.com/cloudinary/cloudinary-go)](https://pkg.go.dev/github.com/cloudinary/cloudinary-go)

Cloudinary Go SDK
==========

Cloudinary is a cloud service for managing your web application's entire media management pipeline. Easily upload images and videos to the cloud using Cloudinary's comprehensive APIs and administration capabilities and easily integrate with your web application, existing or new.

For Go, Cloudinary provides a module for simplifying the integration even further. This Readme provides the basic information needed to get started with the Cloudinary Go SDK. For full documentation, take a look at the [SDK guide](https://cloudinary.com/documentation/go_integration). For a high-level introduction on Cloudinary and a step-by-step walk through of how to integrate Cloudinary in your application, see our [getting started guide](https://cloudinary.com/documentation/how_to_integrate_cloudinary).

## Setup

### Installation

To install the Cloudinary Go SDK, use the `go get` command:

```
go get github.com/cloudinary/cloudinary-go
```

### Configuration

For requests to our secure APIs (e.g., image uploads, asset administration) you must have the `APIKey` and `APISecret` parameters set.
You can find your account-specific configuration credentials in the **Dashboard** page of the [account console](https://cloudinary.com/console).

Setting your `CloudName`, `APIKey` and `APISecret` parameters can be done by initializing the Cloudinary object, or by
using the CLOUDINARY_URL environment variable / system property.

The entry point of the library is the Cloudinary struct.

```go
cld, _ := cloudinary.New()
```

Here's an example of setting the configuration parameters programatically:
```go
cld, _ := cloudinary.NewFromParams("n07t21i7", "123456789012345", "abcdeghijklmnopqrstuvwxyz12")
```

You can also set the configuration parameters by providing the CLOUDINARY_URL value:
```go
cld, _ := cloudinary.NewFromURL("cloudinary://123456789012345:abcdeghijklmnopqrstuvwxyz12@n07t21i7")
```

**Learn more**: [Go configuration](https://cloudinary.com/documentation/go_integration#configuration)


### Update 

To update the Cloudinary Go SDK to the latest version, use the `go get` command with the `-u` option:

```
go get -u github.com/cloudinary/cloudinary-go
```

### Logging

Cloudinary SDK logs errors using standard `go log` functions.

For details on redefining the logger or adjusting the logging level,  see [Logging](logger/README.md).

## Usage


### Upload 

Assuming you have your Cloudinary configuration parameters defined (`CloudName`, `APIKey`, `APISecret`), uploading to
Cloudinary is very simple.

The following example uploads a local JPG to the cloud:

```go
resp, err := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{})
```

The uploaded image is assigned a randomly generated public ID. The image is immediately available for a download through
a CDN:

```go
log.Println(resp.SecureURL)

// https://res.cloudinary.com/demo/image/upload/abcfrmo8zul1mafopawefg.jpg
```

You can also specify your own public ID:

```go
resp, err := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{PublicID: "sample_remote"});
if err != nil {...}
log.Println(resp.SecureURL)

// https://res.cloudinary.com/demo/image/upload/sample_remote.jpg
```

**Learn more**: [Go upload](https://cloudinary.com/documentation/go_image_and_video_upload)


### Complete SDK Example
```go
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

```

## Contributing 

Please feel free to submit issues, fork the repository and send pull requests!

For more information on how to contribute, take a look at the [contributing](CONTRIBUTING.md) page.

## Additional resources

Additional resources are available at:

* [Package reference documentation](https://pkg.go.dev/github.com/cloudinary/cloudinary-go)
* [SDK Documentation](https://cloudinary.com/documentation/go_integration)
* [Upload API documentation](https://cloudinary.com/documentation/upload_images)
* [Website](https://cloudinary.com)
* [Interactive demo](https://demo.cloudinary.com/default)
* [Knowledge Base](https://support.cloudinary.com/hc/en-us)

## Community and Support

Impact the product, hear updates, test drive new features and more!
Join [here](https://www.facebook.com/groups/CloudinaryCommunity).

You can [open an issue through GitHub](https://github.com/cloudinary/cloudinary-go/issues).

Stay tuned for updates, tips and tutorials: [Blog](https://cloudinary.com/blog)
, [Twitter](https://twitter.com/cloudinary), [Facebook](https://www.facebook.com/Cloudinary).

## License 

Released under the MIT license. 
