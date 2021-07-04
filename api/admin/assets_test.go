package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestAssets_AssetTypes(t *testing.T) {
	resp, err := adminAPI.AssetTypes(ctx)

	if err != nil || len(resp.AssetTypes) < 1 {
		t.Error(err, resp)
	}
}

func TestAssets_Assets(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	resp, err := adminAPI.Assets(ctx, admin.AssetsParams{Tags: true, Context: true, Moderations: true, MaxResults: 1})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_AssetsByIDs(t *testing.T) {
	cldtest.UploadTestVideoAsset(t, cldtest.VideoPublicID)
	resp, err := adminAPI.AssetsByIDs(ctx, admin.AssetsByIDsParams{PublicIDs: []string{cldtest.PublicID}, Tags: true})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}

	resp, err = adminAPI.AssetsByIDs(ctx, admin.AssetsByIDsParams{PublicIDs: []string{cldtest.VideoPublicID}, AssetType: api.Video})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_RestoreAssets(t *testing.T) {
	resp, err := adminAPI.RestoreAssets(ctx, admin.RestoreAssetsParams{PublicIDs: []string{"api_test_restore_20891", "api_test_restore_94060"}})
	if err != nil {
		t.Error(err, resp)
	}
}

func TestAssets_DeleteAssets(t *testing.T) {
	resp, err :=
		adminAPI.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: []string{"api_test_restore_20891", "api_test_restore_94060"}})
	if err != nil {
		t.Error(err, resp)
	}
}
