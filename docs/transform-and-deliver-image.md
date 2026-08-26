# Transform and deliver an image

## When to use

Generate CDN-backed delivery URLs that resize, crop, overlay, or optimize an image. URL
generation is **local** — no network call, and only a cloud name is required, not an API
secret. Cloudinary creates the derived asset on first request, then serves it from CDN
cache.

For video, see [Transform and deliver a video](transform-and-deliver-video.md).

## Transformations are strings

There is **no typed transformation builder in this SDK.** `Asset.Transformation` is a
`transformation.RawTransformation`, which is an alias for `string`. You assemble the
Cloudinary transformation syntax yourself:

```go
image, err := cld.Image("sample")
if err != nil {
	return err
}
image.Transformation = "c_thumb,g_auto,h_200,w_200/f_auto,q_auto"

url, err := image.String()
if err != nil {
	return err
}
fmt.Println(url)
// https://res.cloudinary.com/<cloud>/image/upload/c_thumb,g_auto,h_200,w_200/f_auto,q_auto/sample?_a=...
```

Because nothing validates the string at compile time, a typo becomes a runtime `400` from
the CDN rather than a build error. Two consequences worth planning around:

- Use the `cloudinary-transformations` skill or the
  [transformation reference](https://cloudinary.com/documentation/transformation_reference.md)
  to compose the string, rather than recalling parameter names.
- Test the URL, not the code. A bad parameter reports itself in the `x-cld-error` response
  header — see [Troubleshoot errors](troubleshoot-errors.md).

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

const publicID = "examples/transformed-sample"

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

	// Upload something to transform (URL building alone needs no upload).
	uploaded, err := cld.Upload.Upload(ctx, "https://res.cloudinary.com/demo/image/upload/sample.jpg",
		uploader.UploadParams{PublicID: publicID, Overwrite: api.Bool(true)})
	if err != nil {
		return fmt.Errorf("upload transport: %w", err)
	}
	if uploaded.Error.Message != "" {
		return fmt.Errorf("upload rejected: %s", uploaded.Error.Message)
	}

	// A square thumbnail, auto-focused, auto-format, auto-quality.
	thumb, err := cld.Image(publicID)
	if err != nil {
		return err
	}
	thumb.Transformation = "c_thumb,g_auto,h_200,w_200/f_auto,q_auto"
	thumbURL, err := thumb.String()
	if err != nil {
		return err
	}
	fmt.Println("thumbnail:", thumbURL)

	// Cache-busting: deliver the exact version the upload returned.
	versioned, err := cld.Image(publicID)
	if err != nil {
		return err
	}
	versioned.Version = uploaded.Version
	versioned.Transformation = "f_auto,q_auto"
	versionedURL, err := versioned.String()
	if err != nil {
		return err
	}
	fmt.Println("versioned:", versionedURL)

	return nil
}
```

## Chaining: order matters

A `/` separates transformation components, and each runs on the output of the previous one:

```go
image.Transformation = "c_fill,g_auto,h_720,w_1280/co_white,l_text:Arial_64_bold:SALE,g_south_east,x_24,y_24/f_auto,q_auto"
```

Reordering the components changes the result. When you want to hit an eagerly generated
derived asset, the serialized string must match **exactly** — same parameters, same order.

## What the SDK adds to your URL

Verified defaults, all controlled by `cld.Config.URL`:

| Behaviour | Default | Turn off with |
|---|---|---|
| `https://` scheme | on | `cld.Config.URL.Secure = false` |
| `?_a=` analytics parameter | on | `cld.Config.URL.Analytics = false` |
| `v1` path segment injected when the public ID contains `/` | on | `cld.Config.URL.ForceVersion = false` |

Worth knowing about the version placeholder. Verified:

```
cld.Image("sample")      -> .../image/upload/sample
cld.Image("folder/name") -> .../image/upload/v1/folder/name
```

A public ID with a slash gets `v1` unless you set a real `Version` or disable
`ForceVersion`. This is intentional — it keeps CDN paths stable for folder-like IDs — but it
means the URL is not a naive concatenation of your public ID.

### Setting the delivered format

The generated URL carries **no file extension** unless the public ID has one.
`cld.Image("sample")` delivers the original format. Three ways to change it:

```go
cld.Image("sample.jpg")                  // .../image/upload/sample.jpg — extension in the ID
image.Transformation = "f_webp"           // .../image/upload/f_webp/sample — forced format
image.Transformation = "f_auto"           // best format per requesting browser (preferred)
```

Use one of the three above rather than `Asset.Suffix`, which is a separate feature: it
enables Cloudinary's SEO "short URL" form, dropping the delivery type from the path, and
requires a private CDN distribution:

```go
image.Suffix = "jpg" // -> .../democloud/images/sample/jpg   (NOT .../image/upload/sample.jpg)
```

Combinations that do not support it report `URL Suffix is not supported for
<type>/<delivery>`. Leave it unset unless you have configured SEO suffixes on a private CDN.

## Responsive images

A `srcset` is one delivery URL per width, so build it with a loop and let the browser pick:

```go
// buildSrcSet returns a srcset attribute value: one "<url> <width>w" entry per width.
func buildSrcSet(cld *cloudinary.Cloudinary, publicID string, widths []int) (string, error) {
	entries := make([]string, 0, len(widths))
	for _, w := range widths {
		image, err := cld.Image(publicID)
		if err != nil {
			return "", err
		}
		image.Transformation = fmt.Sprintf("c_scale,w_%d/f_auto,q_auto", w)
		url, err := image.String()
		if err != nil {
			return "", err
		}
		entries = append(entries, fmt.Sprintf("%s %dw", url, w))
	}
	return strings.Join(entries, ", "), nil
}
```

Used in an `<img>` tag, with `sizes` describing how much space the image occupies:

```html
<img src="...w_800/f_auto,q_auto/sample"
     srcset="...w_400/... 400w, ...w_800/... 800w, ...w_1200/... 1200w"
     sizes="(max-width: 800px) 100vw, 800px"
     alt="Sample">
```

Choose widths that match your layout's breakpoints rather than a fixed ladder. `f_auto`
and `q_auto` handle format and compression per browser, so each entry differs only in width.

For art direction — a different crop per breakpoint rather than the same image rescaled —
generate a URL per crop and use `<picture>` with `<source media="...">`. `c_fill` with
`g_auto` picks the subject automatically at each aspect ratio.

## Generative and AI transformations

Background removal, generative fill, and similar edits are expressed in the same
transformation string as everything else:

```go
image.Transformation = "e_gen_remove:prompt_car/f_auto,q_auto"
```

Availability is account- and plan-dependent. Verify a given effect against
[generative AI transformations](https://cloudinary.com/documentation/generative_ai_transformations.md)
before relying on it; an unavailable one fails at delivery time, not at build time.

## Cache behaviour

- The same URL is served from CDN cache; a new transformation string means a new URL and a
  fresh derivation.
- After re-uploading to the same public ID, cached URLs do not update. Deliver with the new
  `Version` from the upload result, which changes the URL immediately — see the flow above.
- `Invalidate: api.Bool(true)` on upload purges the CDN copy, but propagation is not
  instant; a version bump is the deterministic fix.

## Signed and access-controlled URLs

For assets that should not be publicly guessable, sign the URL:

```go
cld.Config.URL.SignURL = true
```

The signature covers the transformation and public ID, so a client cannot alter either.
Token-based (time-limited) access uses `config.AuthToken` instead. See
[control access to media](https://cloudinary.com/documentation/control_access_to_media.md).

## Troubleshooting

- `400` with `x-cld-error: Invalid <param> in transformation: <value>` — a malformed
  transformation string. Verified example: `w_abc` returns
  `Invalid width in transformation: abc`.
- `400` with `x-cld-error: Unknown transformation <name>` — a named transformation
  (`t_<name>`) that does not exist on this environment.
- `404` with `x-cld-error: Resource not found - <public_id>` — wrong public ID, wrong
  folder, or wrong resource type in the URL path.
- `401` with `x-cld-error: ACL deny` on every URL from a working cloud — an unclaimed
  Claimable Cloud restricts delivery to its provisioning IP. Not a credentials problem; see
  [Get Cloudinary credentials](get-credentials.md).
- The URL contains an unexpected `v1` — see
  [What the SDK adds to your URL](#what-the-sdk-adds-to-your-url).

## Related

- Runnable examples: `examples/transform-and-deliver-image/main.go` and
  `examples/responsive-srcset/main.go` (in the repository)
- [Transform and deliver a video](transform-and-deliver-video.md)
- Every transformation parameter and its accepted values:
  [Transformation reference](https://cloudinary.com/documentation/transformation_reference.md)
- [Image transformation guide](https://cloudinary.com/documentation/go_media_transformations.md)
