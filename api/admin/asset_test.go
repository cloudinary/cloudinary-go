package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestAsset_Asset(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	resp, err := adminAPI.Asset(ctx, admin.AssetParams{
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
	resp, err := adminAPI.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID: cldtest.PublicID,
		Tags:     []string{"tagA", "tagB", "TagC"},
	})

	if err != nil || len(resp.Tags) != 3 {
		t.Error(resp, err)
	}
}

func TestAssets_AssetsByAssetID(t *testing.T) {
	asset := cldtest.UploadTestAsset(t, cldtest.PublicID)

	t.Run("", func(t *testing.T) {
		resp, err := adminAPI.AssetByAssetIDs(ctx, admin.AssetByAssetIDParams{AssetID: asset.AssetID})
		if err != nil {
			t.Fatal(err)
		}

		if resp.Colors != nil {
			t.Error()
		}

		if resp.Exif != nil {
			t.Error()
		}

		if resp.Faces != nil {
			t.Error()
		}
	})

	t.Run("With Extra Info", func(t *testing.T) {
		resp, err := adminAPI.AssetByAssetIDs(ctx, admin.AssetByAssetIDParams{
			AssetID: asset.AssetID,
			Colors:  true,
			Exif:    true,
			Faces:   true,
		})

		if err != nil {
			t.Fatal(err)
		}

		if resp.Colors == nil {
			t.Error()
		}

		if resp.Exif == nil {
			t.Error()
		}

		if resp.Faces == nil {
			t.Error()
		}
	})
}
