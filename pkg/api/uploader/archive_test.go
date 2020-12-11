package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strings"
	"testing"
)

const folder = "go-folder"

func TestUploader_CreateArchive(t *testing.T) {
	UploadTestAsset(t, publicID)
	UploadTestAsset(t, publicID2)

	params := CreateArchiveParams{
		Tags: tags,
	}

	resp, err := uploadApi.CreateArchive(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.AssetID == "" {
		t.Error(resp)
	}
}

func TestUploader_DownloadArchiveUrl(t *testing.T) {
	UploadTestAsset(t, publicID)
	UploadTestAsset(t, publicID2)

	params := CreateArchiveParams{
		Tags:           tags,
		TargetPublicId: "goArchive",
	}

	arURL, err := uploadApi.DownloadArchiveUrl(params)

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
