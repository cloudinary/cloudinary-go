# AGENTS.md — cloudinary-go

## What this package is (one line)
Official Cloudinary Go **server-side** SDK: upload assets, build transformation/delivery URLs, and call the Admin API from your backend — the SDK that holds your `API_SECRET`.

## When to use this / when NOT to use this
- **Use this when:** you are in a Go server runtime (HTTP handlers, workers, CLIs, build scripts) and need to upload assets, sign delivery URLs, or administer assets — anything that requires the `API_SECRET` to stay off the client.
- **Do NOT use this when:** you only need the **AI Content Analysis API** as a focused typed client → use [`analysis-go`](https://github.com/cloudinary/analysis-go); or you're doing **account provisioning** (sub-accounts, users, access keys) → use [`account-provisioning-go`](https://github.com/cloudinary/account-provisioning-go); or you want a no-code/autonomous-agent path → use the Cloudinary MCP server.
- **Sibling packages:** `analysis-go` and `account-provisioning-go` are narrow clients for a single API surface. This package is the **full server SDK** (upload + transform URLs + Admin API together) and is the default for general Cloudinary work in Go.

## Setup
```bash
go get github.com/cloudinary/cloudinary-go/v2
```
Required configuration / credentials:
```bash
export CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME
```
`cloudinary.New()` reads `CLOUDINARY_URL` from the environment automatically.

## Minimal runnable example
```go
package main

import (
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func main() {
	cld, _ := cloudinary.New() // reads CLOUDINARY_URL
	ctx := context.Background()

	resp, _ := cld.Upload.Upload(ctx, "my_picture.jpg", uploader.UploadParams{})
	fmt.Println(resp.SecureURL)

	image, _ := cld.Image("sample.jpg")
	image.Transformation = "c_fill,h_150,w_100"
	url, _ := image.String()
	fmt.Println(url)
}
```

## Build / test commands (run these after editing)
```bash
go mod download          # fetch dependencies
go build ./...           # compile everything
go test ./...            # run the suite (CI uses gotestsum ./..., equivalent)
go vet ./...             # static checks
```
Notes:
- **Tests need `CLOUDINARY_URL`** — much of the suite hits the live Upload/Admin APIs, so set it to a real (test) cloud before running. CI allocates one via `scripts/get_test_cloud.sh`; locally, export your own. See `TEST.md` for test-writing conventions.
- If you edit param/option structs, regenerate fluent setters: `make generate`. Don't hand-edit generated `*_setters.go` files.
- CI (`.github/workflows/test.yaml`) matrix runs Go 1.20–1.24, tests only (via `gotestsum`). There is no golangci-lint step and no `.golangci.yml` in the repo — `go vet` and `gofmt` are the de-facto gate.

## Conventions & gotchas
- **The `/v2` import path is mandatory.** Module is `github.com/cloudinary/cloudinary-go/v2`; importing without `/v2` resolves to the unmaintained 1.x and will not compile against current code.
- Format with `gofmt` / `go fmt ./...` before committing.
- Public entry point is the `cloudinary` package (`cloudinary.New()`); upload lives under `api/uploader`, admin under `api/admin`. Import these subpackages explicitly rather than reaching into deep internal paths.
- Generated setters: edit the generator under `gen/generate_setters/`, not the output. Run `make generate` after changing option structs.
- Signed uploads and the Admin API require the `API_SECRET` — keep it server-side; never embed it in a client binary or browser bundle.
- **Version support:** SDK 2.8+ requires **Go 1.20+**; SDK 2.7 / 1.x support Go 1.13–1.19. See the table in `README.md`.

## Canonical docs (leave the repo for depth)
- Go SDK guide: https://cloudinary.com/documentation/go_integration
- Transformations: https://cloudinary.com/documentation/go_media_transformations
- Upload: https://cloudinary.com/documentation/go_image_and_video_upload
- Admin API: https://cloudinary.com/documentation/go_asset_administration
- API reference (pkg.go.dev): https://pkg.go.dev/github.com/cloudinary/cloudinary-go/v2
- MCP server (agent/no-code path): https://github.com/cloudinary/mcp-servers

## Agent / MCP note
If the capability you need is also exposed via the Cloudinary MCP servers, prefer the MCP tool for autonomous task execution and use this SDK for generated Go code. See cloudinary/mcp-servers.

## Commit / PR conventions
- Branch off `main`; rebase on `upstream/main` before opening a PR (see `CONTRIBUTING.md`).
- PR description must clearly state the bug/feature and reference the related issue number.
- Provide tests covering new code (`TEST.md`), and ensure CI tests pass before requesting review.
