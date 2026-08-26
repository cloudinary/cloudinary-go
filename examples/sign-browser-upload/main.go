// Sign an upload on your server so a browser can upload directly to Cloudinary.
//
// The API secret never leaves the server; the client receives only a signature, valid
// for one hour from the timestamp it was signed with. This keeps large request bodies
// off your service entirely.
//
// Prerequisites:
//   - CLOUDINARY_URL in the environment. See docs/get-credentials.md.
//     Copy .env.example to .env and fill it in, or export the variable directly.
//
// Run:
//
//	export $(grep CLOUDINARY_URL .env)
//	go run ./sign-browser-upload
//	# then, in another terminal:
//	curl -s localhost:8080/api/sign-upload
//
// In a real project this handler sits in your application, behind authentication, and
// signs only the parameters that user is allowed to set. Every parameter the client
// sends — except file, api_key, signature, and resource_type — must be signed, so
// signing a client-supplied value without validating it hands over control of where
// assets land.
//
// Docs: docs/sign-browser-upload.md
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

const (
	listenAddr = ":8080"
	signPath   = "/api/sign-upload"

	// The one folder this endpoint permits. Decided server-side, on purpose.
	uploadFolder = "user-uploads"

	readHeaderTimeout = 10 * time.Second
)

// signResponse is everything the browser needs. Note what is absent: api_secret.
type signResponse struct {
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	Folder    string `json:"folder"`
	APIKey    string `json:"api_key"`
	CloudName string `json:"cloud_name"`
}

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

	mux := http.NewServeMux()
	mux.HandleFunc(signPath, signUploadHandler(cld))

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("listening on %s — try: curl -s localhost%s%s", listenAddr, listenAddr, signPath)
	return server.ListenAndServe()
}

func signUploadHandler(cld *cloudinary.Cloudinary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Unix()

		// Sign ONLY what the client may use. api.SignParameters would set a timestamp
		// itself, but we need the same value in the response for the client to echo.
		params := url.Values{}
		params.Set("timestamp", strconv.FormatInt(timestamp, 10))
		params.Set("folder", uploadFolder)

		signature, err := api.SignParameters(params, cld.Config.Cloud.APISecret)
		if err != nil {
			log.Printf("signing failed: %v", err)
			http.Error(w, "could not sign upload", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(signResponse{
			Signature: signature,
			Timestamp: timestamp,
			Folder:    uploadFolder,
			APIKey:    cld.Config.Cloud.APIKey,
			CloudName: cld.Config.Cloud.CloudName,
		}); err != nil {
			log.Printf("writing response failed: %v", err)
		}
	}
}
