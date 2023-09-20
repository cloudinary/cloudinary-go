package asset_test

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/asset"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsset_Signature(t *testing.T) {
	i := getTestImage(t)

	i.DeliveryType = api.Authenticated
	i.Config.URL.SignURL = true

	assert.Regexp(t, "s--.{8}--", getAssetUrl(t, i))
}

func TestAsset_LongURLSignature(t *testing.T) {
	i := getTestImage(t)
	i.DeliveryType = api.Authenticated
	i.Config.URL.SignURL = true
	i.Config.URL.LongURLSignature = true

	assert.Regexp(t, "s--.{32}--", getAssetUrl(t, i))
}

func TestAsset_WithAuthToken(t *testing.T) {
	i, _ := asset.Image(authTokenTestImage, nil)
	i.DeliveryType = api.Authenticated
	i.Version = 1486020273

	i.Config.URL.SignURL = true

	localTokenConfig := authTokenConfig
	i.AuthToken.Config = &localTokenConfig

	i.AuthToken.Config.StartTime = 1111111111

	assert.Contains(t, getAssetUrl(t, i), "1751370bcc6cfe9e03f30dd1a9722ba0f2cdca283fa3e6df3342a00a7528cc51")

	assert.NotContains(t, getAssetUrl(t, i), "s--") // no simple signature
	assert.NotContains(t, getAssetUrl(t, i), "_a=") // no analytics

	i.AuthToken.Config.ACL = ""
	i.AuthToken.Config.StartTime = startTime

	i.Config.URL.PrivateCDN = true

	assert.Contains(t, getAssetUrl(t, i), "8db0d753ee7bbb9e2eaf8698ca3797436ba4c20e31f44527e43b6a6e995cfdb3")
}

func TestAsset_ForceVersion(t *testing.T) {
	i, err := asset.Image(cldtest.ImageInFolder, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, getAssetUrl(t, i), "v1")

	i.Config.URL.ForceVersion = false

	assert.NotContains(t, getAssetUrl(t, i), "v1")
}

func TestAsset_Video(t *testing.T) {
	v, err := asset.Video(cldtest.VideoPublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, getAssetUrl(t, v), fmt.Sprintf("video/upload/%s", cldtest.VideoPublicID))
}

func TestAsset_VideoSEO(t *testing.T) {
	f, err := asset.Video(cldtest.VideoPublicID+cldtest.VideoExt, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.Suffix = "my_favorite_video"

	assert.Contains(t, getAssetUrl(t, f), fmt.Sprintf("videos/%s/%s%s", cldtest.VideoPublicID, f.Suffix, cldtest.VideoExt))
}

func TestAsset_File(t *testing.T) {
	f, err := asset.File(cldtest.PublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, getAssetUrl(t, f), fmt.Sprintf("raw/upload/%s", cldtest.PublicID))
}

func TestAsset_FileSEO(t *testing.T) {
	f, err := asset.File(cldtest.PublicID+cldtest.FileExt, nil)
	if err != nil {
		t.Fatal(err)
	}
	f.Suffix = "my_favorite_sample"

	assert.Contains(t, getAssetUrl(t, f), fmt.Sprintf("files/%s/%s%s", cldtest.PublicID, f.Suffix, cldtest.FileExt))
}

func TestAsset_Media(t *testing.T) {
	m, err := asset.Media(cldtest.PublicID, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, getAssetUrl(t, m), fmt.Sprintf("image/upload/%s", cldtest.PublicID))
}
