# Sign a browser upload

## When to use

A browser or mobile app uploads directly to Cloudinary, but your server authorizes the
operation. The API secret stays on the server; the client gets only a signature. This keeps
large request bodies off your service entirely — prefer it over proxying uploads when you
can.

If you would rather receive the file and forward it, see
[Serve uploads over HTTP](serve-uploads-over-http.md).

For uploads with no server round-trip at all, use an
[unsigned upload preset](https://cloudinary.com/documentation/upload_presets.md) — and note
that it is deliberately restricted, because anyone who finds the preset name can use it.

## How long a signature lasts

A signature is valid for **1 hour** from the `timestamp` it was signed with. That window is
enforced server-side, so a signature minted at page load and used 90 minutes later is
rejected as stale. Generate one per upload, at upload time.

## Server: the signing endpoint

There is no `api_sign_request` equivalent method on the client — signing lives in the `api`
package as `api.SignParameters`:

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
)

type signResponse struct {
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	Folder    string `json:"folder"`
	APIKey    string `json:"api_key"`
	CloudName string `json:"cloud_name"`
}

func main() {
	cld, err := cloudinary.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	http.HandleFunc("/api/sign-upload", func(w http.ResponseWriter, r *http.Request) {
		const folder = "user-uploads"
		timestamp := time.Now().Unix()

		// Sign ONLY the parameters the client is allowed to use.
		params := url.Values{}
		params.Set("timestamp", strconv.FormatInt(timestamp, 10))
		params.Set("folder", folder)

		signature, err := api.SignParameters(params, cld.Config.Cloud.APISecret)
		if err != nil {
			log.Printf("signing failed: %v", err)
			http.Error(w, "could not sign upload", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(signResponse{
			Signature: signature,
			Timestamp: timestamp,
			Folder:    folder,
			APIKey:    cld.Config.Cloud.APIKey,
			CloudName: cld.Config.Cloud.CloudName,
		})
	})

	server := &http.Server{Addr: ":8080", ReadHeaderTimeout: 10 * time.Second}
	log.Println("listening on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
```

`api.SignParameters` sets `timestamp` itself if you leave it unset or zero — but set it
explicitly anyway, because you need the same value in the JSON response for the client to
send back.

Return `api_key` and `cloud_name` to the client; they are not secrets. **Never** return
`api_secret`.

## Client: use the signature

```js
const { signature, timestamp, folder, api_key, cloud_name } =
  await (await fetch('/api/sign-upload')).json();

const form = new FormData();
form.append('file', fileInput.files[0]);
form.append('api_key', api_key);
form.append('timestamp', timestamp);
form.append('signature', signature);
form.append('folder', folder); // exactly what the server signed

// 'auto' lets Cloudinary detect image / video / raw from the file itself
const response = await fetch(
  `https://api.cloudinary.com/v1_1/${cloud_name}/auto/upload`,
  { method: 'POST', body: form }
);
const asset = await response.json(); // public_id, secure_url, asset_id, ...
```

`auto` in the upload path means "detect the resource type from the content", so one endpoint
handles images, video, and raw files. It is deliberately excluded from the signature, along
with `file`, `api_key`, and `signature` itself.

## The signing rule

Every parameter the client sends **except** `file`, `api_key`, `signature`, and
`resource_type` must be part of the signed set, or Cloudinary rejects the upload with
`Invalid Signature`.

That cuts both ways, and it is the security boundary: to let the client choose a
`public_id`, a tag, or a transformation, you must add it to the signed parameters
server-side — which means your server decides whether that is allowed. Signing a
client-supplied value without validating it hands over control of where assets land.

The SDK excludes exactly those four keys internally when it signs its own requests, so the
rule is consistent between server-side uploads and this flow.

## Signature algorithm and version

Defaults, verified: algorithm **SHA-1**, signature version **2**. Version 2 percent-encodes
`&` in values to prevent parameter smuggling. Three helpers, in increasing specificity:

```go
signParams := url.Values{"timestamp": {"1700000000"}, "folder": {"user-uploads"}}

api.SignParameters(signParams, secret)                              // sha1, version 2
api.SignParametersUsingAlgo(signParams, secret, "sha256")            // sha256, version 2
api.SignParametersUsingAlgoAndVersion(signParams, secret, "sha1", 2) // explicit both
```

SHA-1 produces a 40-character hex signature, SHA-256 a 64-character one. Only override the
defaults if your product environment is configured for it — a mismatch fails as
`Invalid Signature`, indistinguishable from a wrong secret.

Read the environment's configured values from `cld.Config.Cloud.GetSignatureAlgorithm()` and
`GetSignatureVersion()` rather than hardcoding.

## Verifying what came back

Two verification helpers, both on the uploader, for trusting data that arrives from
Cloudinary rather than from your own call:

```go
// A webhook notification body.
ok := cld.Upload.VerifyNotificationSignature(body, timestamp, receivedSignature, 7200)

// An upload response relayed by an untrusted client.
ok = cld.Upload.VerifyApiResponseSignature(publicID, version, receivedSignature)
```

`VerifyNotificationSignature` takes a validity window in seconds and defaults to 7200 when
you pass `0` or less. Always verify webhooks before acting on them — the endpoint is public.

## Troubleshooting

- `Invalid Signature` — the client sent a parameter that was not signed, or sent a different
  value than the one signed. Compare the two sets exactly; a differing `folder` is the usual
  culprit. The same message also appears for a wrong API secret entirely.
- `Stale request` — the signature is over an hour old. Fetch it at upload time, not page
  load, and check your server clock: a skewed clock mints timestamps that are already stale
  on arrival.
- Uploads land in the wrong folder, or overwrite each other — the client is choosing
  parameters you signed blindly. Validate before signing.
- The upload succeeds but you never learn about it — the browser gets the response, your
  server does not. Set a `notification_url` on the signed parameters and verify the webhook.

## Related

- Runnable example: `examples/sign-browser-upload/main.go` (in the repository)
- [Serve uploads over HTTP](serve-uploads-over-http.md) — the proxying alternative
- [Generating authentication signatures](https://cloudinary.com/documentation/upload_images.md#generating_authentication_signatures)
- [Upload presets](https://cloudinary.com/documentation/upload_presets.md)
