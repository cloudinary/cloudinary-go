package admin

import (
	"testing"
)

func TestAsset_Asset(t *testing.T) {
	resp, err := adminApi.Asset(ctx, AssetParams{
		PublicID:              "sample",
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
	resp, err := adminApi.UpdateAsset(ctx, UpdateAssetParams{
		PublicID: "sample",
		Tags:     []string{"tagA", "tagB"},
	})

	if err != nil || len(resp.Tags) != 2 {
		t.Error(resp, err)
	}
}
