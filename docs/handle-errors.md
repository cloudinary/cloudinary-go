# Handle errors

## When to use

Read this before writing any call site. This SDK reports failures on two channels, and
checking both is what makes a call site correct.

## The rule

Every network method returns `(*Result, error)`. The two values report **different
classes of failure**, so check both:

```go
result, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: "id"})
if err != nil {
	// Transport, context cancellation/deadline, or JSON decoding failed.
	// The request may never have reached Cloudinary.
	return fmt.Errorf("upload transport: %w", err)
}
if result.Error.Message != "" {
	// The request reached Cloudinary and Cloudinary rejected it.
	return fmt.Errorf("upload rejected: %s", result.Error.Message)
}
```

**An API rejection arrives with `err == nil`**, because the request itself succeeded — the
transport worked and Cloudinary answered. The answer was a rejection, and it lands in
`result.Error.Message`. Check both and every failure is covered. Verified against a live
cloud:

| Condition | `err` | `result.Error.Message` |
|---|---|---|
| Unreachable remote source URL | `nil` | `Error in loading https://... - ERR_DNS_FAIL 0` |
| Wrong API secret | `nil` | `Invalid Signature <hash>. String to sign - '...'.` |
| Undefined metadata field | `nil` | `Metadata External IDs do not exist: ["no_such_field"]` |
| Asset not found | `nil` | `Resource not found - <public_id>` |
| Invalid moderation value | `nil` | `Moderation <value> moderation is not valid` |
| Duplicate metadata field | `nil` | `external id <id> already exists` |
| Context deadline exceeded | `context deadline exceeded` | `""` |

Nearly every result struct in both `api/uploader` and `api/admin` carries this
`Error api.ErrorResp` field, so the pattern holds across all ~70 admin methods and almost
all uploader methods.

**Two older uploader results carry the field differently.** `UpdateMetadataResult` and
`RenameResult` type `Error` as `interface{}`, so compare it against `nil` instead:

```go
updated, err := cld.Upload.UpdateMetadata(ctx, uploader.UpdateMetadataParams{
	Metadata:  api.CldAPIMap{"sku": "SKU-00042"},
	PublicIDs: []string{"examples/product-photo"},
})
if err != nil {
	return err
}
if updated.Error != nil { // interface{} — compare against nil
	return fmt.Errorf("rejected: %v", updated.Error)
}
```

The compiler tells you which form a given result needs, so you find out at build time.

**Check `err` first, before reading the result.** When the SDK catches a problem before
sending — a missing API secret, for example — it returns `(nil, err)`:

```go
// With a client built without an API secret:
result, err := cld.Upload.Upload(ctx, file, params)
// err     == "must provide API Secret"
// result  == nil        <- so check err before reading result
```

That is why the order in the snippet above matters: `if err != nil { return }` comes before
any access to `result`.

## A helper worth writing once

Result types share no common interface, so a helper cannot take "any result". Pass the two
failure channels in explicitly instead:

```go
// APIError reports whether a Cloudinary call failed, collapsing the SDK's two
// failure channels into one error value.
func APIError(errMessage string, err error) error {
	if err != nil {
		return fmt.Errorf("cloudinary transport: %w", err)
	}
	if errMessage != "" {
		return fmt.Errorf("cloudinary rejected the request: %s", errMessage)
	}
	return nil
}
```

Called as:

```go
result, err := cld.Upload.Upload(ctx, file, params)
if err := APIError(result.Error.Message, err); err != nil {
	return err
}
```

Passing `result.Error.Message` explicitly keeps it compile-checked per call site: a result
type that ever loses the field becomes a build failure rather than a silent skip.

## Matching on specific errors

- **`errors.Is` and `errors.As` apply to the `err` channel**, where they work as usual for
  `context.Canceled` and `context.DeadlineExceeded`. API rejections come through
  `api.ErrorResp`, a plain struct with a `Message string`. Since the wording is server-side
  and not contractual, branch on it only where handling genuinely differs (retry vs fail),
  and log the message rather than parsing it.
- **Delivery-URL failures carry their status on the HTTP response**, with the reason in the
  `x-cld-error` header — see [Troubleshoot errors](troubleshoot-errors.md). API results
  report the reason through `Error.Message`.

## Choosing a constructor

`cloudinary.New()` validates the configuration it reads and reports what is missing:

```go
cloudinary.New()                              // err: "must provide CLOUDINARY_URL" when unset
cloudinary.NewFromURL("not-a-cloudinary-url") // err == nil — accepts what it is given
cloudinary.NewFromParams("", "", "")          // err == nil — accepts what it is given
```

`NewFromURL` and `NewFromParams` are for callers who already hold their credentials and
assemble the client directly, so they take the values as given; an incomplete set surfaces at
the first API call as `Invalid Signature` or a request against an empty cloud name.
**Prefer `cloudinary.New()`.** When you need the others, validate the resolved config up
front:

```go
cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
if err != nil {
	return err
}
if cld.Config.Cloud.CloudName == "" || cld.Config.Cloud.APIKey == "" || cld.Config.Cloud.APISecret == "" {
	return errors.New("cloudinary: incomplete credentials")
}
```

## Reading fields the structs do not model

Every result carries `Response interface{}` holding the decoded JSON, which is how you reach
a field the typed struct omits (video `duration`, for example). Its dynamic type is a
**pointer** to the map:

```go
// The dynamic type is *map[string]interface{} — assert to the pointer form.
if raw, ok := result.Response.(*map[string]interface{}); ok {
	fmt.Println((*raw)["duration"])
}
```

Do not log `Response` wholesale: it includes your `api_key`.

## Logging

Errors carry the message only. To see the raw request and response, raise the log level:

```go
cld.Logger.SetLevel(logger.DEBUG)
```

Do not log the whole config or a whole request: the signed form parameters and the
configuration both contain your API secret. Log `result.Error.Message`.

## Related

- [Troubleshoot errors](troubleshoot-errors.md) — specific messages and what to do about them
- [Configure Cloudinary](configure-cloudinary.md)
- [Serve uploads over HTTP](serve-uploads-over-http.md) — mapping these to HTTP responses
