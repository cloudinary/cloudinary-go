package admin_test

import (
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
)

func TestAssets_AssetTypes(t *testing.T) {
	resp, err := adminAPI.AssetTypes(ctx)

	if err != nil || len(resp.AssetTypes) < 1 {
		t.Error(err, resp)
	}
}

func TestAssets_Assets(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	resp, err := adminAPI.Assets(ctx, admin.AssetsParams{
		Tags:        api.Bool(true),
		Context:     api.Bool(true),
		Moderations: api.Bool(true),
		MaxResults:  1})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_AssetsByIDs(t *testing.T) {
	cldtest.UploadTestVideoAsset(t, cldtest.VideoPublicID)
	resp, err := adminAPI.AssetsByIDs(ctx, admin.AssetsByIDsParams{
		PublicIDs: []string{cldtest.PublicID},
		Tags:      api.Bool(true)})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}

	resp, err = adminAPI.AssetsByIDs(ctx, admin.AssetsByIDsParams{PublicIDs: []string{cldtest.VideoPublicID}, AssetType: api.Video})

	if err != nil || len(resp.Assets) != 1 {
		t.Error(err, resp)
	}
}

func TestAssets_AssetsByAssetFolder(t *testing.T) {
	cldtest.SkipFeature(t, cldtest.SkipDynamicFolders)

	asset1 := cldtest.UploadTestAsset(t, cldtest.PublicID)
	asset2 := cldtest.UploadTestAsset(t, cldtest.PublicID2)

	resp, err := adminAPI.AssetsByAssetFolder(ctx, admin.AssetsByAssetFolderParams{
		AssetFolder: cldtest.UniqueFolder,
	})

	if err != nil || len(resp.Assets) != 2 {
		t.Error(err, resp)
	}

	assert.Equal(t, asset1.AssetFolder, resp.Assets[0].AssetFolder)
	assert.Equal(t, asset2.AssetFolder, resp.Assets[1].AssetFolder)
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

func TestAssets_VisualSearchImageFile(t *testing.T) {
	uploadAPI, _ := uploader.New()

	asset1, err := uploadAPI.Upload(ctx, cldtest.LogoURL, uploader.UploadParams{VisualSearch: api.Bool(true)})
	if err != nil {
		t.Error(err, asset1)
	}

	// Get the data
	res, err := http.Get(cldtest.LogoURL)
	if err != nil {
		t.Error(err, res)
	}

	defer api.DeferredClose(res.Body)

	resp, err := adminAPI.VisualSearch(ctx, admin.VisualSearchParams{ImageFile: res.Body})
	if err != nil {
		t.Error(err, resp)
	}

	assert.GreaterOrEqual(t, len(resp.Assets), 1)
}
