package uploader

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"os"
	"testing"
)

var ctx = context.Background()
var uploadApi, _ = Create()

const LogoUrl = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"
const LogoFilePath = "testdata/cloudinary_logo.png"
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
const publicID = "go_test_image"
const publicID2 = "go_test_image_2"
const tag1 = "go_tag1"
const tag2 = "go_tag2"

var tags = api.CldApiArray{tag1, tag2}
var cContext = api.CldApiMap{"go-context-key": "go-context-value"}

func TestUploader_UploadLocalPath(t *testing.T) {
	params := UploadParams{
		PublicID:              publicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
		Overwrite:             true,
	}

	resp, err := uploadApi.Upload(ctx, LogoFilePath, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadIOReader(t *testing.T) {
	file, err := os.Open(LogoFilePath)
	if err != nil {
		t.Error(fmt.Printf("unable to read a file: %v\n", err))
	}

	defer api.DeferredClose(file)

	params := UploadParams{
		PublicID:              publicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
	}

	resp, err := uploadApi.Upload(ctx, file, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadUrl(t *testing.T) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadBase64Image(t *testing.T) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
	}

	resp, err := uploadApi.Upload(ctx, Base64Image, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func UploadTestAsset(t *testing.T, publicID string) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      tags,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}
