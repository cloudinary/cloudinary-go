# Import and call the SDK

## Install

```bash
go get github.com/cloudinary/cloudinary-go/v2
```

The `/v2` suffix is part of the module path and of every import path. Omitting it resolves
the abandoned v1 module.

## Import and construct

```go
package main

import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func main() {
	cld, err := cloudinary.New() // reads CLOUDINARY_URL
	if err != nil {
		log.Fatalf("cloudinary: %v", err)
	}

	ctx := context.Background()

	result, err := cld.Upload.Upload(ctx, "https://res.cloudinary.com/demo/image/upload/sample.jpg",
		uploader.UploadParams{PublicID: "examples/uploaded-sample"})
	if err != nil {
		log.Fatalf("upload transport: %v", err)
	}
	if result.Error.Message != "" {
		log.Fatalf("upload rejected: %s", result.Error.Message)
	}

	log.Println(result.SecureURL)
}
```

Note the package name is `cloudinary`, not `cloudinary_go` — the directory in the module
path (`cloudinary-go/v2`) does not match it, so some editors will not auto-import it
correctly. Write the import explicitly.

## The three API surfaces

One `*cloudinary.Cloudinary` value carries everything:

| Field | Package | Use for |
|---|---|---|
| `cld.Upload` | `api/uploader` | uploads, destroy, rename, tags, context, archives |
| `cld.Admin` | `api/admin` | asset listing/details, search, folders, presets, metadata fields |
| `cld.Config` | `config` | the resolved configuration, readable and writable |

URL building hangs directly off the client: `cld.Image`, `cld.Video`, `cld.File`,
`cld.Media`.

You can also construct the sub-APIs standalone — `uploader.New()` and `admin.New()` — when
a component only needs one of them.

## Conventions this SDK follows

- **Every network call takes a `context.Context` first.** Pass the request's context so
  cancellation and deadlines propagate; see
  [Serve uploads over HTTP](serve-uploads-over-http.md).
- **Parameters are structs, not option maps.** `uploader.UploadParams`,
  `admin.AssetParams`, and so on. Unset fields are omitted from the request, so the zero
  value is "not specified" rather than "send empty".
- **Optional booleans are `*bool`.** Use `api.Bool(true)` to distinguish "false" from
  "unset". Passing a plain `false` is indistinguishable from omitting the field.
- **Results are typed structs with a raw escape hatch.** Every result also has a
  `Response interface{}` holding the decoded JSON, so a field the struct does not model
  yet is still reachable.

## Conventions this SDK does *not* follow

- **API errors are not `error` values.** `err` is transport-level only; check
  `result.Error.Message` as well. There are no sentinel errors and no typed error
  hierarchy, so `errors.Is` and `errors.As` have nothing to match against. See
  [Handle errors](handle-errors.md).
- **Field naming differs between packages.** The uploader calls it `ResourceType`; the
  admin API calls the same concept `AssetType`. The compiler will catch it.

## Related

- [Configure Cloudinary](configure-cloudinary.md)
- [Handle errors](handle-errors.md)
- [Go SDK guide](https://cloudinary.com/documentation/go_integration.md)
