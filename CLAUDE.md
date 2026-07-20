@AGENTS.md

# CLAUDE.md — cloudinary-go

## What this repo is

Official Cloudinary Go **server-side** SDK: upload assets, build transformation/delivery URLs, and call the Admin API from backend code. Module path: `github.com/cloudinary/cloudinary-go/v2` (the `/v2` suffix is mandatory — dropping it resolves to the unmaintained 1.x line).

## Key constraints / gotchas

- **`/v2` import path is mandatory.** Import `github.com/cloudinary/cloudinary-go/v2` everywhere; subpackages are `api/uploader` and `api/admin`.
- **`API_SECRET` stays server-side.** Never embed it in a client binary or browser bundle; use the signed-upload pattern for browser-direct flows.
- **Tests need `CLOUDINARY_URL`.** Most of the suite hits the live Upload/Admin APIs; set `export CLOUDINARY_URL=cloudinary://API_KEY:API_SECRET@CLOUD_NAME` before running. CI uses `scripts/get_test_cloud.sh`.
- **Do not hand-edit `*_setters.go` files.** They are generated. Edit the generator under `gen/generate_setters/`, then run `make generate`.
- **Go floor:** SDK 2.8+ requires Go 1.20+. SDK 2.7 / 1.x support Go 1.13–1.19.
- **CI gate:** `.github/workflows/test.yaml` runs `gotestsum ./...` on Go 1.20–1.24. There is no golangci-lint step; `go vet` and `gofmt` are the de-facto gate.
- **Format before committing:** `go fmt ./...`

## Verified build / test commands

```bash
go mod download          # fetch dependencies
go build ./...           # compile everything
go vet ./...             # static checks
go test ./...            # full suite (requires CLOUDINARY_URL pointing at a real test cloud)
make generate            # regenerate fluent setters after editing param/option structs
```
