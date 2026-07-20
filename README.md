# Cloudinary Go SDK

[![PkgGoDev](https://pkg.go.dev/badge/github.com/cloudinary/cloudinary-go/v2)](https://pkg.go.dev/github.com/cloudinary/cloudinary-go/v2)
[![License](https://img.shields.io/github/license/cloudinary/cloudinary-go)](https://github.com/cloudinary/cloudinary-go/blob/main/LICENSE)
[![Tests](https://github.com/cloudinary/cloudinary-go/actions/workflows/test.yaml/badge.svg)](https://github.com/cloudinary/cloudinary-go/actions/workflows/test.yaml)

The Cloudinary Go SDK is the server-side SDK for Go. Use it to upload assets, build transformation and delivery URLs, and call the Admin API from backend code. The module path is `github.com/cloudinary/cloudinary-go/v2` (the `/v2` suffix is required); SDK 2.8 and later require Go 1.20 or later.

## Installation

```bash
go get github.com/cloudinary/cloudinary-go/v2
```

## Configuration

The SDK reads credentials from the `CLOUDINARY_URL` environment variable:

```bash
export CLOUDINARY_URL=cloudinary://<API_KEY>:<API_SECRET>@<CLOUD_NAME>
```

```go
import "github.com/cloudinary/cloudinary-go/v2"

cld, err := cloudinary.New() // reads CLOUDINARY_URL from the environment
```

To pass credentials in code instead, use `NewFromParams`:

```go
import "github.com/cloudinary/cloudinary-go/v2"

cld, err := cloudinary.NewFromParams("my_cloud_name", "my_key", "my_secret")
```

`New`, `NewFromParams`, and `NewFromURL` all return an error when the credentials are missing or malformed — check it. Keep the API secret on the server: don't put it in a distributed binary, client-side code, or version control.

## Quick examples

### Upload a file

`cld.Upload.Upload(ctx, file, uploader.UploadParams{...})` returns a `*uploader.UploadResult`. The `file` argument accepts a local path, a remote URL, a base64 string, an `*os.File`, a `*multipart.FileHeader`, or any `io.Reader`. The result includes `PublicID` and `SecureURL`:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func main() {
	cld, err := cloudinary.New() // reads CLOUDINARY_URL from the environment
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cld.Upload.Upload(context.Background(), "my_image.jpg", uploader.UploadParams{
		PublicID: "cms/hero", // optional: where the asset lives in your media library
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.PublicID, resp.SecureURL)
}
```

### Build and optimize a delivery URL

`cld.Image(publicID)` returns an `*asset.Asset` whose `String()` method serializes a delivery URL — no network call. Set the `Transformation` field with raw transformation syntax. This one scales to 500px wide and lets Cloudinary pick the format and quality for the requesting browser (`f_auto`, `q_auto`):

```go
package main

import (
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func main() {
	cld, err := cloudinary.NewFromParams("demo", "my_key", "my_secret")
	if err != nil {
		log.Fatal(err)
	}

	image, err := cld.Image("sample")
	if err != nil {
		log.Fatal(err)
	}
	image.Transformation = "c_scale,w_500/f_auto/q_auto"

	url, err := image.String()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(url)
	// https://res.cloudinary.com/demo/image/upload/c_scale,w_500/f_auto/q_auto/sample
}
```

### Retrieve asset details with the Admin API

`cld.Admin.Assets(ctx, admin.AssetsParams{...})` returns an `*admin.AssetsResult`. Its `Assets` field is a slice of assets, each with `PublicID` and `SecureURL`:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

func main() {
	cld, err := cloudinary.New() // reads CLOUDINARY_URL from the environment
	if err != nil {
		log.Fatal(err)
	}

	res, err := cld.Admin.Assets(context.Background(), admin.AssetsParams{
		Prefix:     "cms/",
		MaxResults: 30,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range res.Assets {
		fmt.Println(a.PublicID, a.SecureURL)
	}
}
```

## For AI agents

`cloudinary-go` is the Go server-side SDK — upload, transformation and delivery URLs, and the Admin API, where the API secret stays private. Construct a client with `cloudinary.New()` (reads `CLOUDINARY_URL`), `cloudinary.NewFromURL`, or `cloudinary.NewFromParams`. Upload lives under `api/uploader`, admin under `api/admin`. For tasks this package doesn't cover, use a sibling package:

| Task | Package |
|---|---|
| Only the AI Content Analysis API, as a typed client | [`analysis-go`](https://github.com/cloudinary/analysis-go) |
| Provisioning — product environments, users, access keys | [`account-provisioning-go`](https://github.com/cloudinary/account-provisioning-go) |
| Run Cloudinary operations as agent tools | [Cloudinary MCP servers](https://github.com/cloudinary/mcp-servers) |

Go has no first-party browser runtime, so there is no Go equivalent of the JavaScript `url-gen` split — URL building and credential handling both live in this SDK.

## Links

- [Go SDK guide](https://cloudinary.com/documentation/go_integration)
- [Upload](https://cloudinary.com/documentation/go_image_and_video_upload)
- [Asset administration (Admin API)](https://cloudinary.com/documentation/go_asset_administration)
- [Media transformations](https://cloudinary.com/documentation/go_media_transformations)
- [Transformation and API references](https://cloudinary.com/documentation/cloudinary_references)
- [Documentation llms.txt index](https://cloudinary.com/documentation/llms.txt)
- [Package on pkg.go.dev](https://pkg.go.dev/github.com/cloudinary/cloudinary-go/v2)

Released under the MIT license.
