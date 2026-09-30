// Accept a file from an HTTP request and forward it to Cloudinary.
//
// Shows the idiomatic net/http wiring: one shared client, the request context passed
// through so a disconnecting client cancels the outbound upload, the
// *multipart.FileHeader handed to the SDK unread, and the SDK's two failure channels
// mapped onto status codes.
//
// Prefer signing a direct browser upload when you can — it keeps large bodies off your
// service. See ./sign-browser-upload.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./serve-uploads-over-http
//	# then, in another terminal:
//	curl -F file=@some-image.jpg localhost:8080/upload
//
// In a real project this handler sits behind authentication, and the folder and
// overwrite policy come from the authenticated user rather than being hardcoded.
//
// Docs: docs/serve-uploads-over-http.md
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

const (
	listenAddr = ":8080"
	uploadPath = "/upload"

	maxUploadBytes    = 32 << 20 // reject larger bodies before allocating
	inMemoryFormBytes = 8 << 20  // above this, parts spill to temp files
	uploadTimeout     = 5 * time.Minute
	readHeaderTimeout = 10 * time.Second

	uploadFolder = "user-uploads"
)

type server struct {
	cld *cloudinary.Cloudinary
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	// Build the client once and share it: it is safe for concurrent use, and its
	// underlying http.Client pools connections.
	cld, err := cloudinary.New()
	if err != nil {
		return err
	}

	srv := &server{cld: cld}

	mux := http.NewServeMux()
	mux.Handle(uploadPath, http.MaxBytesHandler(http.HandlerFunc(srv.handleUpload), maxUploadBytes))

	httpServer := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("listening on %s — try: curl -F file=@some-image.jpg localhost%s%s",
		listenAddr, listenAddr, uploadPath)
	return httpServer.ListenAndServe()
}

func (s *server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(inMemoryFormBytes); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}
	// Parts above the in-memory threshold become temp files; without this the disk
	// fills up slowly.
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
	}()

	_, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing 'file' field", http.StatusBadRequest)
		return
	}

	// Derive the deadline from the request context so a disconnecting client cancels
	// the outbound upload too.
	ctx, cancel := context.WithTimeout(r.Context(), uploadTimeout)
	defer cancel()

	// Hand the SDK the *multipart.FileHeader itself. It opens and streams the file,
	// and chunks it automatically past Config.API.ChunkSize. Reading it into a []byte
	// first would defeat both.
	result, err := s.cld.Upload.Upload(ctx, header, uploader.UploadParams{
		Folder:       uploadFolder,
		ResourceType: api.Auto, // detect image / video / raw from the content
		Overwrite:    api.Bool(false),
	})

	switch {
	case errors.Is(err, context.Canceled):
		// The client gave up. Not our error, and nobody is listening.
		return
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "upload timed out", http.StatusGatewayTimeout)
		return
	case err != nil:
		log.Printf("cloudinary transport: %v", err)
		http.Error(w, "upload failed", http.StatusBadGateway)
		return
	}

	// The request completed. That does NOT mean it succeeded: an API rejection
	// arrives with err == nil and the reason in result.Error.Message.
	if result.Error.Message != "" {
		log.Printf("cloudinary rejected upload: %s", result.Error.Message)
		http.Error(w, "upload rejected", http.StatusBadRequest)
		return
	}

	// Return the delivery handles only — never the config or the raw response, which
	// carry the api_secret and api_key respectively.
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"asset_id\":%q,\"public_id\":%q,\"url\":%q}\n",
		result.AssetID, result.PublicID, result.SecureURL)
}
