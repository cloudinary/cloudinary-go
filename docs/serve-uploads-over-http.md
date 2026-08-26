# Serve uploads over HTTP

## When to use

Your Go service accepts a file from a browser or mobile client and forwards it to
Cloudinary. This is the idiomatic wiring for `net/http`: where the client goes, how the
request context flows, and how to map the SDK's two failure channels onto status codes.

If the client can upload straight to Cloudinary, prefer that — it keeps large bodies off
your service entirely. See [Sign a browser upload](sign-browser-upload.md).

## Build the client once

`cloudinary.New()` reads configuration and constructs an HTTP client. Do it at startup and
share the value: it is safe for concurrent use by multiple goroutines, and its underlying
`http.Client` pools connections. Constructing one per request wastes connections and
re-reads the environment on every call.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const maxUploadBytes = 32 << 20 // 32 MiB

type server struct {
	cld *cloudinary.Cloudinary
}

func main() {
	cld, err := cloudinary.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

	srv := &server{cld: cld}

	mux := http.NewServeMux()
	mux.Handle("/upload", http.MaxBytesHandler(http.HandlerFunc(srv.handleUpload), maxUploadBytes))

	httpServer := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Println("listening on :8080")
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func (s *server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form; keep small files in memory, spill larger ones to disk.
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll() // delete any spilled temp files
		}
	}()

	_, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing 'file' field", http.StatusBadRequest)
		return
	}

	// Give the upload its own deadline, derived from the request context so that a
	// disconnecting client cancels the outbound call too.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	// Pass the *multipart.FileHeader directly: the SDK opens it, streams it, and
	// chunks it if it exceeds Config.API.ChunkSize. Never read it fully into memory.
	result, err := s.cld.Upload.Upload(ctx, header, uploader.UploadParams{
		Folder:       "user-uploads",
		ResourceType: api.Auto, // detect image / video / raw from the file itself
		Overwrite:    api.Bool(false),
	})

	switch {
	case errors.Is(err, context.Canceled):
		// The client went away; nothing useful to write.
		return
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "upload timed out", http.StatusGatewayTimeout)
		return
	case err != nil:
		log.Printf("cloudinary transport: %v", err)
		http.Error(w, "upload failed", http.StatusBadGateway)
		return
	}

	// Reaching here means the request completed. It does NOT mean it succeeded:
	// an API rejection arrives with err == nil.
	if result.Error.Message != "" {
		log.Printf("cloudinary rejected upload: %s", result.Error.Message)
		http.Error(w, "upload rejected", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"asset_id":%q,"public_id":%q,"url":%q}`+"\n",
		result.AssetID, result.PublicID, result.SecureURL)
}
```

## The four rules

1. **Share one client.** Build it in `main`, store it on your handler struct or inject it.
   It is concurrency-safe for calls; only mutating `cld.Config` after startup is not.
2. **Pass `r.Context()`.** Cancellation propagates into the outbound HTTP request, so a
   client that disconnects mid-upload does not leave the transfer running. Derive a
   timeout from it rather than using `context.Background()`, which ignores disconnects.
3. **Pass the `*multipart.FileHeader`, not bytes.** The SDK opens and streams it, and
   chunks it automatically past `ChunkSize`. Reading it into a `[]byte` first defeats both
   and makes memory use scale with request size.
4. **Check both failure channels, in order.** `err` first (and never dereference `result`
   before it — a pre-flight failure returns a `nil` result), then
   `result.Error.Message`. See [Handle errors](handle-errors.md).

## Distinguishing cancellation from real failure

`context.Canceled` means the *client* gave up; it is not your error and should not be
logged as one or reported as 5xx. `context.DeadlineExceeded` means your own deadline
elapsed. These are the two cases where `errors.Is` genuinely works with this SDK, because
they come from the standard library rather than from Cloudinary — API rejections carry no
sentinel to match. The `switch` above separates all three.

## Bound the request size

`http.MaxBytesHandler` rejects oversized bodies before your handler allocates anything.
Set it below the smaller of your plan's asset limit and whatever your service can afford
to buffer; read the plan limit from
[`cld.Admin.Usage`](upload-image.md#size-limits) rather than guessing.

Note that `ParseMultipartForm` writes anything above its in-memory threshold to a temp
file, so `RemoveAll` in a `defer` is not optional — without it the disk fills up slowly.

## Do not put the API secret in the response

Return `SecureURL`, `PublicID`, and `AssetID`. Never echo the config, the full result
struct, or `result.Response` back to a client: the raw response includes your `api_key`,
and your configuration holds the secret. Store `AssetID` as the durable handle — see
[Upload an image](upload-image.md#result-fields-to-keep).

## Related

- [Sign a browser upload](sign-browser-upload.md) — keep large bodies off your service
- [Handle errors](handle-errors.md)
- [Upload an image](upload-image.md) — accepted source types
- [Configure Cloudinary](configure-cloudinary.md) — timeouts and chunk size
