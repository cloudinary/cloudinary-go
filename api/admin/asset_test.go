package admin_test

import (
	"github.com/cloudinary/cloudinary-go/api"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestAsset_Asset(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	resp, err := adminAPI.Asset(ctx, admin.AssetParams{
		PublicID:              cldtest.PublicID,
		Exif:                  api.Bool(true),
		Colors:                api.Bool(true),
		Faces:                 api.Bool(true),
		QualityAnalysis:       api.Bool(true),
		ImageMetadata:         api.Bool(true),
		Phash:                 api.Bool(true),
		Pages:                 api.Bool(true),
		AccessibilityAnalysis: api.Bool(true),
		CinemagraphAnalysis:   api.Bool(true),
		Coordinates:           api.Bool(true),
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

func TestAsset_AssetByAssetID(t *testing.T) {
	asset := cldtest.UploadTestAsset(t, cldtest.PublicID)

	t.Run("", func(t *testing.T) {
		resp, err := adminAPI.AssetByAssetID(ctx, admin.AssetByAssetIDParams{AssetID: asset.AssetID})
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
		resp, err := adminAPI.AssetByAssetID(ctx, admin.AssetByAssetIDParams{
			AssetID: asset.AssetID,
			Colors:  api.Bool(true),
			Exif:    api.Bool(true),
			Faces:   api.Bool(true),
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
