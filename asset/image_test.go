package asset_test

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/asset"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsset_SimpleImage(t *testing.T) {
	i := getTestImage(t)

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s://%s/%s/image/upload/%s", i.Config.URL.Protocol(), i.Config.URL.SharedHost, i.Config.Cloud.CloudName, cldtest.PublicID))
}

func TestAsset_FetchImage(t *testing.T) {
	i, err := asset.Image(cldtest.LogoURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	i.DeliveryType = api.Fetch

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("image/fetch/%s", cldtest.LogoURL))
}

func TestAsset_ImageSEO(t *testing.T) {
	i, err := asset.Image(cldtest.PublicID+cldtest.ImgExt, nil)
	if err != nil {
		t.Fatal(err)
	}

	i.Suffix = cldtest.SEOName

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("images/%s/%s", cldtest.PublicID, cldtest.SEOName+cldtest.ImgExt))

	i.Config.URL.Shorten = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("iu/%s/%s", cldtest.PublicID, cldtest.SEOName+cldtest.ImgExt))

	i.Config.URL.UseRootPath = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s/%s/%s", i.Config.Cloud.CloudName, cldtest.PublicID, cldtest.SEOName+cldtest.ImgExt))

	i.PublicID = cldtest.PublicID

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s/%s/%s", i.Config.Cloud.CloudName, cldtest.PublicID, cldtest.SEOName))
}

func TestAsset_ImageSuffixForNotSupportedDeliveryType(t *testing.T) {
	i := getTestImage(t)

	i.Suffix = cldtest.SEOName
	i.DeliveryType = api.Public

	_, err := i.String()

	if err == nil || err.Error() != "failed to build URL: URL Suffix is not supported for image/public" {
		t.Fatal(err)
	}
}

func TestAsset_ImageTransformation(t *testing.T) {
	i := getTestImage(t)

	i.Transformation = cldtest.Transformation

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("image/upload/%s/%s", cldtest.Transformation, cldtest.PublicID))
}

func TestAsset_ImageAnalytics(t *testing.T) {
	i := getTestImage(t)

	assert.Contains(t, getAssetUrl(t, i), "?_a=")

	i.Config.URL.Analytics = false

	assert.NotContains(t, getAssetUrl(t, i), "?_a=")
}
