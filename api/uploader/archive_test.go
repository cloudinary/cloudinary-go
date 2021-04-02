package uploader_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"github.com/stretchr/testify/assert"
)

const folder = "go-folder"

func TestUploader_CreateZip(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	cldtest.UploadTestAsset(t, cldtest.PublicID2)

	params := uploader.CreateArchiveParams{
		Tags: cldtest.Tags,
	}

	resp, err := uploadAPI.CreateZip(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.AssetID == "" {
		t.Error(resp)
	}
}

func TestUploader_DownloadZipURL(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	cldtest.UploadTestAsset(t, cldtest.PublicID2)

	params := uploader.CreateArchiveParams{
		Tags:           cldtest.Tags,
		TargetPublicID: "goArchive",
	}

	arURL, err := uploadAPI.DownloadZipURL(params)

	if err != nil {
		t.Error(err)
	}

	if arURL == "" {
		t.Error(arURL)
	}
}

func TestUploader_DownloadFolder(t *testing.T) {

	cldtest.UploadTestAsset(t, folder+"/"+cldtest.PublicID)
	cldtest.UploadTestAsset(t, folder+"/"+cldtest.PublicID2)

	params := uploader.CreateArchiveParams{
		Tags:           cldtest.Tags,
		TargetPublicID: "goArchive",
	}

	folderURL, _ := uploadAPI.DownloadFolder(folder, params)
	assert.Contains(t, folderURL, "generate_archive")
	assert.Contains(t, folderURL, folder)
	assert.Contains(t, folderURL, params.TargetPublicID)
	assert.Contains(t, folderURL, api.All)
	assert.Contains(t, folderURL, url.QueryEscape(strings.Join(cldtest.Tags, ",")))

	folderURL, _ = uploadAPI.DownloadFolder(folder, uploader.CreateArchiveParams{ResourceType: api.Image})
	assert.Contains(t, folderURL, api.Image)
}

func TestUploader_DownloadBackedUpAsset(t *testing.T) {
	params := uploader.DownloadBackedUpAssetParams{
		AssetID:   "b71b23d9c89a81a254b88a91a9dad8cd",
		VersionID: "0e493356d8a40b856c4863c026891a4e",
	}

	downloadURL, _ := uploadAPI.DownloadBackedUpAsset(params)
	assert.Contains(t, downloadURL, "asset_id")
	assert.Contains(t, downloadURL, "version_id")
	assert.Contains(t, downloadURL, "download_backup")
}
