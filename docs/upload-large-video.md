# Upload a large video

## When to use

Any upload large enough to benefit from chunked transfer — typically video.

## Chunking is automatic

`cld.Upload.Upload` handles files of any size. It inspects the source and switches to chunked
transfer on its own once the size exceeds `cld.Config.API.ChunkSize` (default 20 MB), so one
call covers both small and large uploads and there is no separate method to choose.

Chunking applies to sources whose size the SDK can determine:

| Source type | Chunked above `ChunkSize`? |
|---|---|
| local file path (`string`) | yes |
| `*os.File` | yes |
| `*multipart.FileHeader` | yes |
| `*io.SectionReader` | yes |
| plain `io.Reader` | **no** — size unknown, sent as one request |
| remote URL / base64 string | not applicable — Cloudinary fetches it server-side |

For a large stream of unknown length, buffer it to a temp file or wrap it in an
`*io.SectionReader` so the SDK can size it.

## Complete flow

```go
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	sampleVideoURL = "https://res.cloudinary.com/demo/video/upload/dog.mp4"
	localPath      = "dog.mp4"
	publicID       = "examples/uploaded-large-video"
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

	if err := ensureLocalVideo(); err != nil {
		return err
	}

	// Uploads over ChunkSize are chunked automatically. ResourceType must be set
	// for video: the default is image, and an mp4 stored as an image cannot be
	// streamed or transformed as video.
	result, err := cld.Upload.Upload(ctx, localPath, uploader.UploadParams{
		PublicID:     publicID,
		ResourceType: api.Video,
		Overwrite:    api.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if result.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", result.Error.Message)
	}

	fmt.Println(result.PublicID)
	fmt.Println(result.Bytes, "bytes")
	fmt.Println(result.SecureURL)   // playback URL
	fmt.Println(result.PlaybackURL) // HLS, when the environment provides one
	return nil
}

// ensureLocalVideo downloads the demo sample once, so the example runs with no arguments.
func ensureLocalVideo() error {
	if _, err := os.Stat(localPath); err == nil {
		return nil
	}

	resp, err := http.Get(sampleVideoURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sample download failed: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
```

## Set `ResourceType` for video

`UploadParams.ResourceType` defaults to `image`. For video, set `api.Video`:

```go
params := uploader.UploadParams{ResourceType: api.Video}
```

An mp4 uploaded as an image is stored and delivered under `/image/upload/`, and video
transformations and streaming profiles do not apply to it. Use `api.Auto` to let Cloudinary
detect the type from the file when the input is mixed or unknown.

> **Note the field name per API:** the uploader calls this `ResourceType`, the Admin API
> calls the same concept `AssetType` (`admin.AssetParams{AssetType: api.Video}`). Both
> default to image, so set the type on Admin lookups too — otherwise a video returns
> `Resource not found`.

## Chunk size

`cld.Config.API.ChunkSize` sets the chunk size in bytes; the default is 20,000,000.

> **Keep `ChunkSize` at or above Cloudinary's 5 MB minimum** (`5 * 1024 * 1024`). The
> default already satisfies this, so the common case needs nothing. Below the minimum the
> platform stops after the first chunk: verified with `ChunkSize` at 1 MB, a 9 MB upload
> returns `err == nil`, an empty `result.Error.Message`, a partial-progress body
> (`{"done": false, ...}`), and no asset. If you tune the value, confirm the result
> describes a finished asset:
>
> ```go
> if result.PublicID == "" || result.SecureURL == "" {
>     return fmt.Errorf("upload did not complete: %+v", result.Response)
> }
> ```

The last chunk is allowed to be smaller than the minimum; only the intermediate ones are
constrained.

## Timeouts

Chunks are uploaded sequentially over one call, so a large file needs a budget for the whole
transfer, not per chunk:

```go
cld.Config.API.UploadTimeout = 600 // seconds, uploads only
```

`UploadTimeout` defaults to `0`, which means uploads fall back to `API.Timeout` (60 s) — too
short for a large video on a slow link. Set it explicitly. A context deadline you pass in is
respected too, and is the better mechanism when the upload is serving a request.

## Asynchronous processing

Large videos may still be **processing** after the upload returns. For derived versions,
request eager transformations asynchronously and receive a webhook rather than polling:

```go
params := uploader.UploadParams{
	ResourceType:    api.Video,
	Eager:           "sp_hd/m3u8",
	EagerAsync:      api.Bool(true),
	NotificationURL: "https://example.com/cloudinary-webhook",
}
```

Verify the webhook with `cld.Upload.VerifyNotificationSignature` before trusting it.

## Troubleshooting

- Upload returns without error but no asset exists, and `result.Response` shows
  `"done": false` — `ChunkSize` is below the 5 MB minimum. See [Chunk size](#chunk-size).
- The video will not transform or stream, and its URL contains `/image/upload/` — it was
  uploaded with the default `ResourceType`. Re-upload with `api.Video`.
- `Resource not found` when looking the video up — the Admin API defaults to `image`; pass
  `AssetType: api.Video`.
- `context deadline exceeded` partway through a large upload — raise
  `Config.API.UploadTimeout`, or pass a context with a longer deadline.
- `File size too large` — the asset exceeds your product environment's maximum. Chunking
  does not raise that limit; see [Upload an image](upload-image.md#size-limits).

## Related

- Runnable example: `examples/upload-large-video/main.go` (in the repository)
- [Transform and deliver a video](transform-and-deliver-video.md) — what to do with it next
- [Upload an image](upload-image.md) — source types and size limits in full
- [Video upload guide](https://cloudinary.com/documentation/go_image_and_video_upload.md)
