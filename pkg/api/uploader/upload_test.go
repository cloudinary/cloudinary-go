package uploader

import (
	"testing"
)

func TestUploader_Upload(t *testing.T) {
	u, err := Create()
	if err != nil {
		t.Error(err)
	}

	params := UploadParams{
		PublicID:              "go_test_image",
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
	}

	resp, err := u.Upload("https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png", params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != "go_test_image" {
		t.Error(resp)
	}
}
