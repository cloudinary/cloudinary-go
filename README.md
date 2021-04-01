Cloudinary Go SDK
==========

Cloudinary is a cloud service that offers a solution to a web application's entire image management pipeline.

Easily upload images to the cloud. Automatically perform smart image resizing, cropping and conversion without
installing any complex software. Integrate Facebook or Twitter profile image extraction in a snap, in any dimension and
style to match your website's graphics requirements. Images are seamlessly delivered through a fast CDN, and much-much
more.

Cloudinary offers comprehensive APIs and administration capabilities and is easy to integrate with any web application,
existing or new.

Cloudinary provides URL and HTTP based APIs that can be easily integrated with any Web development framework.

For Go, Cloudinary provides a module for simplifying the integration even further.

## Setup ######################################################################

To install Cloudinary Go SDK, use `go get`:

```
go get github.com/cloudinary/cloudinary-go
```

## Usage

### Configuration

Each request for building a URL of a remote cloud resource must have the `CloudName` parameter set. Each request to our
secure APIs (e.g., image uploads, eager sprite generation) must have the `ApiKey` and `ApiSecret` parameters set.
See [API, URLs and access identifiers](http://cloudinary.com/documentation/api_and_access_identifiers) for more details.

Setting the `CloudName`, `ApiKey` and `ApiSecret` parameters can be done by initializing the Cloudinary object, or by
using the CLOUDINARY_URL environment variable / system property.

The entry point of the library is the Cloudinary struct.
```go
cld, _ := cloudinary.New()
```
Here's an example of setting the configuration parameters programatically:
```go
cld, _ := cloudinary.NewFromParams("n07t21i7", "123456789012345", "abcdeghijklmnopqrstuvwxyz12")
```
Another example of setting the configuration parameters by providing the CLOUDINARY_URL value:
```go
cld, _ := cloudinary.NewFromURL("cloudinary://123456789012345:abcdeghijklmnopqrstuvwxyz12@n07t21i7")
```
### Upload

Assuming you have your Cloudinary configuration parameters defined (`CloudName`, `ApiKey`, `ApiSecret`), uploading to
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
	// Alternatively you can use cloudinary.NewFromParams() or cloudinary.NewFromUrl().
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

	// Now we can use Admin Api to see the details about the asset.
	// The request can be customised by providing AssetParams.
	asset, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "logo"})
	if err != nil {
		log.Fatalf("Failed to get asset details, %v\n", err)
	}

	// Print some basic information about the asset.
	log.Printf("Public ID: %v, URL: %v\n", asset.PublicID, asset.SecureURL)

	// Cloudinary also provides a very flexible Search Api for filtering and retrieving
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

## Additional resources ##########################################################

Additional resources are available at:

* [Package documentation](https://pkg.go.dev/github.com/cloudinary/cloudinary-go)
* [Website](https://cloudinary.com)
* [Interactive demo](https://demo.cloudinary.com/default)
* [Knowledge Base](https://support.cloudinary.com/hc/en-us)
* [Documentation](https://cloudinary.com/documentation)
* [Upload API documentation](https://cloudinary.com/documentation/upload_images)

## Support

You can [open an issue through GitHub](https://github.com/cloudinary/cloudinary-go/issues).

Stay tuned for updates, tips and tutorials: [Blog](https://cloudinary.com/blog)
, [Twitter](https://twitter.com/cloudinary), [Facebook](https://www.facebook.com/Cloudinary).

## Join the Community ###########################################################

Impact the product, hear updates, test drive new features and more!
Join [here](https://www.facebook.com/groups/CloudinaryCommunity).

## Staying up to date  ##########################################################

To update Cloudinary Go SDK to the latest version, use `go get -u github.com/cloudinary/cloudinary-go`.

## Contributing  ################################################################

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue.

## License ######################################################################

Released under the MIT license. 
