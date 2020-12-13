package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strings"
	"testing"
)

const folder = "go-folder"

func TestUploader_CreateZip(t *testing.T) {
	UploadTestAsset(t, publicID)
	UploadTestAsset(t, publicID2)

	params := CreateArchiveParams{
		Tags: tags,
	}

	resp, err := uploadApi.CreateZip(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.AssetID == "" {
		t.Error(resp)
	}
}

func TestUploader_DownloadZipUrl(t *testing.T) {
	UploadTestAsset(t, publicID)
	UploadTestAsset(t, publicID2)

	params := CreateArchiveParams{
		Tags:           tags,
		TargetPublicId: "goArchive",
	}

	arURL, err := uploadApi.DownloadZipUrl(params)

	if err != nil {
		t.Error(err)
	}

	if arURL == "" {
		t.Error(arURL)
	}
}

func TestUploader_DownloadFolder(t *testing.T) {

	UploadTestAsset(t, folder+"/"+publicID)
	UploadTestAsset(t, folder+"/"+publicID2)

	params := CreateArchiveParams{
		Tags:           tags,
		TargetPublicId: "goArchive",
	}

	folderURL, _ := uploadApi.DownloadFolder(folder, params)
	assert.Contains(t, folderURL, GenerateArchive)
	assert.Contains(t, folderURL, folder)
	assert.Contains(t, folderURL, params.TargetPublicId)
	assert.Contains(t, folderURL, api.All)
	assert.Contains(t, folderURL, url.QueryEscape(strings.Join(tags, ",")))

	folderURL, _ = uploadApi.DownloadFolder(folder, CreateArchiveParams{ResourceType: api.Image})
	assert.Contains(t, folderURL, api.Image)
}

func TestUploader_DownloadBackedUpAsset(t *testing.T) {
	params := DownloadABackedUpAssetParams{
		AssetID:   "b71b23d9c89a81a254b88a91a9dad8cd",
		VersionID: "0e493356d8a40b856c4863c026891a4e",
	}

	downloadURL, _ := uploadApi.DownloadBackedUpAsset(params)
	assert.Contains(t, downloadURL, "asset_id")
	assert.Contains(t, downloadURL, "version_id")
	assert.Contains(t, downloadURL, DownloadBackup)
}
