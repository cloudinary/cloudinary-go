package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestAsset_Asset(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	resp, err := adminApi.Asset(ctx, admin.AssetParams{
		PublicID:              cldtest.PublicID,
		Exif:                  true,
		Colors:                true,
		Faces:                 true,
		QualityAnalysis:       true,
		ImageMetadata:         true,
		Phash:                 true,
		Pages:                 true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
		Coordinates:           true,
	})

	if err != nil || resp.PublicID == "" {
		t.Error(resp)
	}
}

func TestAsset_UpdateAsset(t *testing.T) {
	resp, err := adminApi.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID: cldtest.PublicID,
		Tags:     []string{"tagA", "tagB", "TagC"},
	})

	if err != nil || len(resp.Tags) != 3 {
		t.Error(resp, err)
	}
}
