# Transform and deliver a video

## When to use

Generate CDN-backed delivery URLs for a video already in Cloudinary. URL generation is
**local** — no network call, only a cloud name required. Cloudinary derives the rendition on
first request, then serves it from CDN cache.

For images, see [Transform and deliver an image](transform-and-deliver-image.md).

## Building player markup

`cld.Video(publicID)` returns an `*asset.Asset`, and `String()` gives you a delivery URL.
That URL is everything a `<video>` tag needs, and Go's `html/template` composes the markup
in a few lines:

```go
package main

import (
	"bytes"
	"html/template"

	"github.com/cloudinary/cloudinary-go/v2"
)

const videoTag = `<video controls poster="{{.Poster}}" width="640">
  <source src="{{.MP4}}" type="video/mp4">
  <source src="{{.HLS}}" type="application/x-mpegURL">
</video>`

// videoMarkup renders a <video> tag with a poster frame and both MP4 and HLS sources.
func videoMarkup(cld *cloudinary.Cloudinary, publicID string) (string, error) {
	url := func(id, transformation string) (string, error) {
		a, err := cld.Video(id)
		if err != nil {
			return "", err
		}
		a.Transformation = transformation
		return a.String()
	}

	mp4, err := url(publicID, "c_scale,w_640/q_auto")
	if err != nil {
		return "", err
	}
	poster, err := url(publicID+".jpg", "so_2,c_fill,w_640")
	if err != nil {
		return "", err
	}
	hls, err := url(publicID+".m3u8", "sp_auto")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = template.Must(template.New("v").Parse(videoTag)).Execute(&buf, struct {
		MP4, Poster, HLS string
	}{mp4, poster, hls})
	return buf.String(), err
}

func main() {}
```

For a full-featured client-side player with analytics and adaptive switching built in, use
the [Cloudinary Video Player](https://cloudinary.com/documentation/cloudinary_video_player.md).
The rest of this page covers how to produce each URL that markup needs.

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

const publicID = "examples/delivered-video"

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

	// ResourceType must be video, or the asset lands under /image/upload/ and
	// no video transformation will apply to it.
	uploaded, err := cld.Upload.Upload(ctx, "https://res.cloudinary.com/demo/video/upload/dog.mp4",
		uploader.UploadParams{
			PublicID:     publicID,
			ResourceType: api.Video,
			Overwrite:    api.Bool(true),
		})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}

	// 1. A scaled MP4 for direct playback.
	playback, err := cld.Video(publicID)
	if err != nil {
		return err
	}
	playback.Transformation = "c_scale,w_640/q_auto"
	playbackURL, err := playback.String()
	if err != nil {
		return err
	}
	fmt.Println("playback:", playbackURL)

	// 2. A poster frame two seconds in — request an image format from the video.
	poster, err := cld.Video(publicID + ".jpg")
	if err != nil {
		return err
	}
	poster.Transformation = "so_2,c_fill,w_640"
	posterURL, err := poster.String()
	if err != nil {
		return err
	}
	fmt.Println("poster:", posterURL)

	// 3. An HLS manifest for adaptive bitrate streaming.
	hls, err := cld.Video(publicID + ".m3u8")
	if err != nil {
		return err
	}
	hls.Transformation = "sp_auto"
	hlsURL, err := hls.String()
	if err != nil {
		return err
	}
	fmt.Println("hls:", hlsURL)

	// Cloudinary may also hand you a ready-made HLS URL on the upload result.
	if uploaded.PlaybackURL != "" {
		fmt.Println("playback_url:", uploaded.PlaybackURL)
	}

	return nil
}
```

## Poster frame from a video

Ask for an image extension on a video asset and Cloudinary renders a still. `so_<seconds>`
(start offset) picks the frame:

```go
poster, _ := cld.Video("examples/delivered-video.jpg")
poster.Transformation = "so_2,c_fill,w_640"
// https://res.cloudinary.com/<cloud>/video/upload/so_2,c_fill,w_640/examples/delivered-video.jpg
```

Put the extension in the **public ID**, as above. `Asset.Suffix` is a different feature — it
produces the SEO short-URL form and requires a private CDN. See
[Setting the delivered format](transform-and-deliver-image.md#setting-the-delivered-format).

## Adaptive bitrate streaming

For anything longer than a short clip, deliver HLS or DASH so the player can switch
renditions rather than committing to one MP4:

```go
hls, _ := cld.Video("examples/delivered-video.m3u8")
hls.Transformation = "sp_auto" // or a named profile: sp_hd, sp_full_hd
// .../video/upload/sp_auto/examples/delivered-video.m3u8

dash, _ := cld.Video("examples/delivered-video.mpd")
dash.Transformation = "sp_auto"
```

Streaming profiles are per-environment. List what is available with
`cld.Admin.ListStreamingProfiles(ctx)`; create your own with
`cld.Admin.CreateStreamingProfile`.

`UploadResult.PlaybackURL` often already contains an `sp_auto` HLS URL for uploaded videos —
check it before constructing one.

## Result fields for video

`UploadResult` is shared across asset types, so `Width`, `Height`, `Bytes`, `Format`, and
`PlaybackURL` are on the struct. Video-specific fields like duration and frame rate come from
the raw response, which every result carries:

```go
// Response holds a *pointer* to the decoded map — assert to *map[string]interface{}.
if raw, ok := uploaded.Response.(*map[string]interface{}); ok {
	fmt.Println((*raw)["duration"])   // e.g. 13.4134
	fmt.Println((*raw)["frame_rate"]) // e.g. 29.97002997002997
	fmt.Println((*raw)["bit_rate"])   // e.g. 5424041
}
```

## What the SDK adds to your URL

Same as for images: `https://` by default, an `?_a=` analytics parameter, and a `v1`
placeholder segment when the public ID contains a `/` and no `Version` is set. See
[What the SDK adds to your URL](transform-and-deliver-image.md#what-the-sdk-adds-to-your-url).

## Cache behaviour

- A new transformation string means a new URL and a fresh derivation. The first request for
  a large rendition is slow because it is being generated.
- After re-uploading, deliver with the new `Version` from the upload result rather than
  waiting for CDN expiry.
- For long transcodes prefer eager transformations with `EagerAsync` and a
  `NotificationURL` over polling — see
  [Upload a large video](upload-large-video.md#asynchronous-processing).

## Troubleshooting

- The video will not transform or stream, and its URL contains `/image/upload/` — it was
  uploaded with the default `ResourceType` (image). Re-upload with `api.Video`. See
  [Upload a large video](upload-large-video.md).
- `404` with `x-cld-error: Resource not found` — check the resource type in the path as well
  as the public ID; the same ID can exist as both an image and a video.
- `Resource not found` from `cld.Admin.Asset` for a video that plays fine — the Admin API
  defaults to `image`; pass `AssetType: api.Video`.
- The first request to a streaming URL is very slow or times out — the renditions are being
  generated. Pre-generate them eagerly at upload time.
- `URL Suffix is not supported` — `Asset.Suffix` is set; put the extension in the public ID
  instead.

## Related

- Runnable example: `examples/transform-and-deliver-video/main.go` (in the repository)
- [Upload a large video](upload-large-video.md)
- [Transform and deliver an image](transform-and-deliver-image.md)
- Every transformation parameter and its accepted values:
  [Transformation reference](https://cloudinary.com/documentation/transformation_reference.md)
- [Video transformation guide](https://cloudinary.com/documentation/video_manipulation_and_delivery.md)
