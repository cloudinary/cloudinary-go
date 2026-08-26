# Configure Cloudinary

## When to use

Once, at startup, before any upload, admin, or URL-generation call.

**Prerequisite:** a cloud name, API key, and API secret. If you do not have them, see
[Get Cloudinary credentials](get-credentials.md) — `npx @cloudinary/cloud` provisions a
working cloud with no signup.

## Recommended: environment variable

```bash
export CLOUDINARY_URL=cloudinary://<api_key>:<api_secret>@<cloud_name>
```

```go
cld, err := cloudinary.New() // reads CLOUDINARY_URL
if err != nil {
	log.Fatalf("cloudinary: %v", err)
}
log.Println(cld.Config.Cloud.CloudName)
```

`cloudinary.New()` is the **only** constructor that validates its input: it returns
`must provide CLOUDINARY_URL` when the variable is unset or empty. Prefer it.

`CLOUDINARY_URL` is the only environment variable read. There is no
`CLOUDINARY_CLOUD_NAME` / `CLOUDINARY_API_KEY` fallback and no merging of individual
variables — setting them has no effect (verified). This differs from some other Cloudinary
SDKs, where a separate cloud-name variable can override the URL.

## Alternatives, and their sharp edge

Pick one:

```go
cldFromURL, errURL := cloudinary.NewFromURL("cloudinary://key:secret@cloud")
cldFromParams, errParams := cloudinary.NewFromParams("cloud", "key", "secret")
cldFromOAuth, errOAuth := cloudinary.NewFromOAuthToken("cloud", token) // OAuth instead of key/secret
```

**These take the values as given.** `NewFromURL` and `NewFromParams` are for callers who
already hold their credentials, so an incomplete set surfaces at the first API call rather
than at construction. Validate the resolved config yourself when you use them — see
[Choosing a constructor](handle-errors.md#choosing-a-constructor).

## Instance-scoped, not process-global

Configuration lives on the `*cloudinary.Cloudinary` value rather than in a process-wide
global, so two clients with different clouds coexist safely:

```go
prod, _ := cloudinary.NewFromURL(prodURL)
staging, _ := cloudinary.NewFromURL(stagingURL)
```

Mutate a live client's config directly when you need to change behaviour:

```go
cld.Config.URL.Secure = true            // https (already the default)
cld.Config.API.Timeout = 120            // seconds, for API calls
cld.Config.API.UploadTimeout = 600      // seconds, uploads only; 0 means "use Timeout"
cld.Config.API.ChunkSize = 20_000_000   // bytes; keep at or above 5 MB
```

Set these at startup. They are plain struct fields with no synchronisation, so write them
before other goroutines begin issuing calls.

> **Keep `ChunkSize` at or above Cloudinary's 5 MB minimum.** The default already does. See
> [Upload a large video](upload-large-video.md#chunk-size).

## Defaults worth knowing

Read off the `config` package; all verified against a live cloud:

| Field | Default | Effect |
|---|---|---|
| `URL.Secure` | `true` | URLs are `https://` |
| `URL.Analytics` | `true` | appends a `?_a=` tracking parameter to generated URLs |
| `URL.ForceVersion` | `true` | injects a `v1` path segment for public IDs containing `/` |
| `API.Timeout` | `60` | seconds, per API call |
| `API.UploadTimeout` | `0` | unset — uploads fall back to `Timeout` |
| `API.ChunkSize` | `20000000` | bytes; files above this are uploaded in chunks |

To drop the analytics parameter: `cld.Config.URL.Analytics = false`.

## Validate configuration early

Fail at startup rather than on the first user request:

```go
func newCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.New()
	if err != nil {
		return nil, err
	}
	if _, err := cld.Admin.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("cloudinary unreachable: %w", err)
	}
	return cld, nil
}
```

`Ping` is a cheap Admin API call that confirms the credentials actually work — which
constructing a client does not. Remember to check `result.Error.Message` too if you want
to catch a rejected ping rather than only a transport failure.

## Troubleshooting

- `must provide CLOUDINARY_URL` — the variable is unset or empty. Note this is returned by
  `cloudinary.New()` only.
- `must provide API Secret` — the client was built without a secret (for example from a
  URL with no password component) and you attempted a signed call. URL generation still
  works; uploads and Admin calls do not.
- `Invalid Signature ...` on every call — the key and secret do not belong to this cloud
  name. Re-copy all three. This is how a wrong secret reports itself; nothing names the
  secret directly.
- Calls succeed against the wrong environment — check which `CLOUDINARY_URL` the process
  actually loaded, and log `cld.Config.Cloud.CloudName` at startup.

## Related

- [Get Cloudinary credentials](get-credentials.md) — if you do not have an account yet
- [Handle errors](handle-errors.md)
- [Go SDK guide](https://cloudinary.com/documentation/go_integration.md)
