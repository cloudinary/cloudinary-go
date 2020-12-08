package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"testing"
)

var ctx = context.Background()
var uploadApi, _ = Create()

const LogoUrl = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"
const publicID = "go_test_image"
const publicID2 = "go_test_image_2"
const tag1 = "go_tag1"
const tag2 = "go_tag2"

var tags = api.CldApiArray{tag1, tag2}
var cContext = api.CldApiMap{"go-context-key": "go-context-value"}

func TestUploader_Upload(t *testing.T) {

	params := UploadParams{
		PublicID:              publicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

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
