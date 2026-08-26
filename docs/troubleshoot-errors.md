# Troubleshoot errors

## First: check the right value

When a call appears to do nothing, check the second failure channel first. API rejections
arrive on the result, not on `err`:

```go
result, err := cld.Upload.Upload(ctx, file, params)
if err != nil { /* transport, context, decoding — result may be nil */ }
if result.Error.Message != "" { /* Cloudinary rejected it — CHECK THIS */ }
```

If a call appears to succeed but the result is empty, read
`result.Error.Message` before anything else. Full detail: [Handle errors](handle-errors.md).

To see the raw traffic:

```go
cld.Logger.SetLevel(logger.DEBUG)
```

## Errors by symptom

### `must provide CLOUDINARY_URL`
A real Go `error` from `cloudinary.New()`. The variable is unset or empty. Note the other
constructors do **not** validate — see
[Configure Cloudinary](configure-cloudinary.md#alternatives-and-their-sharp-edge).

### `must provide API Secret`
A real Go `error`, returned **with a `nil` result** — do not dereference the result before
checking `err`. The client was built without a secret and you attempted a signed call. URL
generation still works without one.

### `Invalid Signature <hash>. String to sign - '...'`
In `result.Error.Message`. Either the API secret does not match the cloud name, or (for
browser uploads) the client sent parameters that were not signed. The message includes the
string that was signed, which is the fastest way to spot the mismatched parameter. See
[Sign a browser upload](sign-browser-upload.md).

A wrong secret reports itself this way rather than naming the secret.

### `api_secret mismatch`
The Admin API's wording for the same condition an upload reports as `Invalid Signature`. Both
mean the key and secret do not belong to this cloud name — verified: the same bad credentials
produce `Invalid Signature` from `cld.Upload` and `api_secret mismatch` from `cld.Admin`. Do
not branch on one and miss the other.

### `Stale request`
A signature older than one hour. Sign at upload time, and check the server clock.

### `Resource not found - <public_id>`
In `result.Error.Message` for Admin calls, or as `x-cld-error` on a `404` delivery URL. Three
common causes, in order of likelihood:

1. **Wrong asset type.** Admin params default to `image`. A video needs
   `AssetType: api.Video`. The same public ID can exist as several types.
2. You passed an **asset ID** where a public ID is required. Only `AssetByAssetID` accepts
   asset IDs; see
   [Search and manage assets](search-and-manage-assets.md#read-and-update-a-single-asset).
3. The public ID genuinely differs — including its folder prefix.

### `Metadata External IDs do not exist: [...]`
The metadata field is not defined in this environment. **The whole upload failed** and no
asset was created. Define the field first — see
[Use structured metadata](use-structured-metadata.md).

### `external id <id> already exists`
Metadata field definitions are permanent and per-environment. Reuse the field.

### `Moderation <value> moderation is not valid`
Misspelled moderation value. See [Moderate an upload](moderate-upload.md#automatic-moderation).

### `Error in loading <url> - ERR_DNS_FAIL 0`
A remote source URL was not reachable **from Cloudinary's servers**. A URL that works in your
browser or on localhost still fails here. It must be publicly resolvable.

### `Rate limit exceeded`
Admin and Search APIs are rate-limited per hour. Retry with backoff, and keep Admin calls out
of request paths — caching what you read is what keeps you comfortably inside the limit.

Note that **an unsubscribed add-on can also surface as a rate-limit error** rather than a
permission error. If you hit this while using a moderation or analysis add-on, check the
add-on registration before assuming you are calling too often.

### `File size too large`
The asset exceeds your **product environment's** maximum, which is a plan limit, not a
per-request one. Chunking does not raise it. Read the real values from
`cld.Admin.Usage` — see [Upload an image](upload-image.md#size-limits) — then compress,
resize, or upgrade the plan.

### Upload "succeeds" but no asset exists, and the raw response shows `"done": false`
`cld.Config.API.ChunkSize` is below Cloudinary's 5 MB minimum. The SDK does not validate it,
so the failure is entirely silent: `err == nil`, empty `Error.Message`, no asset. Keep
`ChunkSize` at or above `5 * 1024 * 1024`. See
[Upload a large video](upload-large-video.md#chunk-size).

### A video will not transform or stream
Its URL contains `/image/upload/`. `UploadParams.ResourceType` defaults to `image`; re-upload
with `api.Video`.

### `423 Processing`
The asset is still being processed and is not yet available for the operation you requested —
common right after uploading a large video, or while an eager or add-on-driven transformation
is running. This is **transient**: retry with backoff rather than treating it as a failure.
For long jobs prefer `EagerAsync` with a `NotificationURL` over polling.

### `context deadline exceeded` / `context canceled`
These are genuine Go errors and the two cases where `errors.Is` works with this SDK:

```go
errors.Is(err, context.DeadlineExceeded) // your deadline elapsed
errors.Is(err, context.Canceled)         // the caller cancelled — often a disconnected client
```

For uploads, `Config.API.Timeout` (60 s default) applies unless
`Config.API.UploadTimeout` is set — it defaults to `0`, meaning "fall back to Timeout",
which is too short for a large video. See
[Serve uploads over HTTP](serve-uploads-over-http.md).

### A delivery URL returns `400` or `404`
The reason is in the `x-cld-error` response header (lowercase over HTTP/2):

```bash
curl -sI "https://res.cloudinary.com/<cloud>/image/upload/w_abc/sample.jpg" | grep -i x-cld-error
# x-cld-error: Invalid width in transformation: abc
```

Verified values:

| Status | `x-cld-error` | Meaning |
|---|---|---|
| 400 | `Invalid width in transformation: abc` | malformed transformation parameter |
| 400 | `Unknown transformation nope_xyz` | named transformation (`t_`) missing on this environment |
| 404 | `Resource not found - docs/nope-xyz` | wrong public ID, folder, or resource type in the path |

The header is CORS-exposed, so browser code can read it from a failed image load too.

### A delivery URL returns `401` with `x-cld-error: ACL deny`
Not a credentials problem. An **unclaimed Claimable Cloud restricts delivery to the IP it was
provisioned from** — uploads and Admin calls keep working while every delivery URL returns
`401`. A VPN reconnect or a changed egress IP triggers it on a cloud that worked yesterday.
Claim the cloud, or add the viewer IP. See
[Get Cloudinary credentials](get-credentials.md#two-limits-before-the-cloud-is-claimed).

### A search returns zero results for data you know exists
A valid expression that matches nothing is not an error. On dynamic-folder environments
`folder:` matches nothing — use `asset_folder:`. Leading wildcards and a bare `*` are
rejected outright. The search index also lags writes, so read-after-write should use
`AssetByAssetID`. See
[Search and manage assets](search-and-manage-assets.md#writing-expressions-that-match).

### `URL Suffix is not supported`
`Asset.Suffix` is set. It enables SEO short URLs rather than adding a file extension, and the
asset type/delivery type combination does not support it. Put the extension in the public ID
instead. See
[Setting the delivered format](transform-and-deliver-image.md#setting-the-delivered-format).

### `index out of range` reading moderation status
`UploadResult.Moderation` is a slice, empty for non-moderated assets. Check `len` first. See
[Moderate an upload](moderate-upload.md#where-the-status-lives).

### Stale content after re-uploading
CDN-cached URLs do not update instantly. Deliver with the new `Version` from the upload
result, which changes the URL immediately. `Invalidate: api.Bool(true)` also works but
propagates more slowly. See
[Transform and deliver an image](transform-and-deliver-image.md#cache-behaviour).

## Do not log these

- The whole result or `result.Response` — the raw response includes your `api_key`.
- `cld.Config` — it holds the API secret.

Log `result.Error.Message`.

## Still stuck

- **Platform status: https://status.cloudinary.com** — check this first. A widespread
  incident explains failures that look like a bug in your code.
- SDK bugs: https://github.com/cloudinary/cloudinary-go/issues
- Account issues: https://support.cloudinary.com
- [Troubleshooting index for agents](https://cloudinary.com/documentation/llms-troubleshooting.txt)

## Related

- [Handle errors](handle-errors.md) — the two failure channels, in depth
- [Configure Cloudinary](configure-cloudinary.md)
