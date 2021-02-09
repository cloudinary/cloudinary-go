package admin

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
)

func TestAssets_AssetTypes(t *testing.T) {
	resp, err := adminApi.AssetTypes(ctx)

	if err != nil || len(resp.AssetTypes) < 1 {
		t.Error(err, resp)
	}
}

func TestAssets_Assets(t *testing.T) {
	resp, err := adminApi.Assets(ctx, AssetsParams{Tags: true, Context: true, Moderations: true, MaxResults: 1})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_AssetsByIDs(t *testing.T) {
	resp, err := adminApi.AssetsByIDs(ctx, AssetsByIDsParams{PublicIDs: []string{"sample"}, Tags: true})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}

	resp, err = adminApi.AssetsByIDs(ctx, AssetsByIDsParams{PublicIDs: []string{"dog"}, AssetType: api.Video})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_RestoreAssets(t *testing.T) {
	resp, err := adminApi.RestoreAssets(ctx, RestoreAssetsParams{PublicIDs: []string{"api_test_restore_20891", "api_test_restore_94060"}})
	if err != nil {
		t.Error(err, resp)
	}
}

func TestAssets_DeleteAssets(t *testing.T) {
	resp, err := adminApi.DeleteAssets(ctx, DeleteAssetsParams{PublicIDs: []string{"api_test_restore_20891", "api_test_restore_94060"}})
	if err != nil {
		t.Error(err, resp)
	}
}
