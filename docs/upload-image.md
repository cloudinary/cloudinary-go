# Upload an image

## When to use

Server-side upload of a local file, a remote URL, a base64 string, or a stream into your
Cloudinary product environment. For uploads started in a browser, see
[Sign a browser upload](sign-browser-upload.md). For accepting a file from an incoming
request, see [Serve uploads over HTTP](serve-uploads-over-http.md).

## Complete flow

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const sourceURL = "https://res.cloudinary.com/demo/image/upload/sample.jpg"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	cld, err := cloudinary.New() // reads CLOUDINARY_URL
	if err != nil {
		return err
	}

	ctx := context.Background()

	result, err := cld.Upload.Upload(ctx, sourceURL, uploader.UploadParams{
		PublicID:  "examples/uploaded-sample", // stable, addressable ID; omit for a random one
		Overwrite: api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if result.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", result.Error.Message)
	}

	fmt.Println(result.PublicID)  // examples/uploaded-sample
	fmt.Println(result.AssetID)   // 32-char hex, immutable
	fmt.Println(result.SecureURL) // canonical delivery URL of the original
	fmt.Println(result.Width, result.Height, result.Format, result.Bytes)
	return nil
}
```

Both checks are required: a rejected upload returns `err == nil` with the reason in
`result.Error.Message`. See [Handle errors](handle-errors.md).

## What the first argument accepts

`Upload` takes `interface{}` and dispatches on the dynamic type. Accepted:

| Value | Treated as |
|---|---|
| `string` that parses as a URL | remote URL, fetched by Cloudinary |
| `string` matching `data:...;base64,...` | inline base64 payload |
| any other `string` | local file path |
| `*os.File` | open file; chunked if larger than `ChunkSize` |
| `*multipart.FileHeader` | an uploaded file from an HTTP request |
| `*io.SectionReader` | sized reader; chunked if larger than `ChunkSize` |
| `io.Reader` | streamed in a single request, never chunked |

Anything else returns `invalid file parameter of unsupported type %T` as a real `error`.

The string cases are disambiguated by parsing, not by a flag — so a local path that happens
to look like a URL is treated as a URL. Pass an `*os.File` when the input is untrusted.

> A plain `io.Reader` is never chunked, because its size is unknown. For a large stream,
> wrap it in an `*io.SectionReader` so the SDK can size it. See
> [Upload a large video](upload-large-video.md).

## Overwrite behaviour

Verified against a live cloud: re-uploading to an existing public ID **overwrites by
default** and reports `Overwritten: true` with a new `Version`. There is no error and no
"existing asset" response.

```go
fmt.Println(result.Overwritten) // true on the second upload to the same public ID
fmt.Println(result.Version)     // changes; use it for cache busting
```

Pass `Overwrite: api.Bool(false)` to refuse instead. Note `api.Bool` is required to express
`false`: the field is a `*bool`, so a plain zero value is indistinguishable from "not
specified" and the default applies.

## Result fields to keep

Store `AssetID`. It never changes; `PublicID` changes when an asset is renamed or moved.

```go
fmt.Println(result.AssetID) // e.g. "abcdef0123456789abcdef0123456789"
```

Look assets up later with `cld.Admin.AssetByAssetID`. The response carries the current
`PublicID`, which is what delivery URLs, `UpdateAsset`, and `Destroy` take — so the lookup
is the step that connects the stored handle to the acting one:

```go
details, err := cld.Admin.AssetByAssetID(ctx, admin.AssetByAssetIDParams{AssetID: storedAssetID})
// ... check err and details.Error.Message ...
image, _ := cld.Image(details.PublicID) // delivery needs the public ID
```

For batch reads, `cld.Admin.AssetsByIDs` takes **public** IDs.

## Size limits

Two separate limits apply, and they are unrelated:

- **The per-request ceiling.** Files above `cld.Config.API.ChunkSize` (default 20 MB) are
  split into chunks automatically — this SDK has no separate `upload_large` method. See
  [Upload a large video](upload-large-video.md).
- **Your product environment's maximum asset size**, which varies by plan. Chunking does
  **not** raise it; an asset over the environment maximum is rejected however it is sent.

Read the real values rather than assuming:

```go
usage, err := cld.Admin.Usage(ctx, admin.UsageParams{})
// ... check err and usage.Error.Message ...
fmt.Println(usage.MediaLimits.ImageMaxSizeBytes)
fmt.Println(usage.MediaLimits.VideoMaxSizeBytes)
fmt.Println(usage.MediaLimits.ImageMaxPx, usage.MediaLimits.AssetMaxTotalPx)
```

If an asset exceeds the environment maximum, compress or resize it before uploading, or
upgrade the plan.

## Troubleshooting

- `must provide CLOUDINARY_URL` — configuration missing; see
  [Configure Cloudinary](configure-cloudinary.md).
- `result.Error.Message` is `Invalid Signature ...` — wrong API secret for this cloud name.
- `result.Error.Message` is `Error in loading <url> - ERR_DNS_FAIL 0` — the remote URL must
  be publicly reachable from Cloudinary. A URL that resolves on your machine but not from
  the internet fails here.
- `File size too large` — the asset exceeds your product environment's maximum. Chunking
  will not help; see [Size limits](#size-limits).
- Upload appears to succeed but `result.PublicID` is empty — you checked only `err`. Check
  `result.Error.Message`.

## Related

- Runnable example: `examples/upload-image/main.go` (in the repository, not the module cache)
- [Transform and deliver an image](transform-and-deliver-image.md)
- [Handle errors](handle-errors.md)
- [Upload guide](https://cloudinary.com/documentation/go_image_and_video_upload.md)
- [Upload API reference](https://cloudinary.com/documentation/image_upload_api_reference.md)
