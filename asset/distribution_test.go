package asset_test

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/asset"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsset_SecureDistribution(t *testing.T) {
	i := getTestImage(t)

	assertURLProtocol(t, i, "https")

	i.Config.URL.Secure = false

	assertURLProtocol(t, i, "http")
}

func TestAsset_SecureDistributionFromConfig(t *testing.T) {
	i, err := asset.Image(cldtest.PublicID, &config.Configuration{URL: config.URL{Secure: false}})
	if err != nil {
		t.Fatal(err)
	}

	assertURLProtocol(t, i, "http")
}

func TestAsset_SecureCNameOverwrite(t *testing.T) {
	host := "something.else.com"

	i := getTestImage(t)

	i.Config.URL.SecureCName = host

	assert.Contains(t, getAssetUrl(t, i), host)
}

func TestAsset_SecureAkamai(t *testing.T) {
	i := getTestImage(t)
	i.Config.URL.PrivateCDN = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s-%s", i.Config.Cloud.CloudName, i.Config.URL.SharedHost))
	assert.NotContains(t, getAssetUrl(t, i), fmt.Sprintf("/%s/", i.Config.Cloud.CloudName))
}

func TestAsset_SecureNonAkamai(t *testing.T) {
	host := "something.cloudfront.net"

	i := getTestImage(t)
	i.Config.URL.PrivateCDN = true
	i.Config.URL.SecureCName = host

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s://%s/image/upload/", i.Config.URL.Protocol(), host))
}

func TestAsset_HTTPPrivateCDN(t *testing.T) {
	i := getTestImage(t)
	i.Config.URL.PrivateCDN = true
	i.Config.URL.Secure = false

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s-%s", i.Config.Cloud.CloudName, i.Config.URL.SharedHost))
}

func TestAsset_SecureSharedSubDomain(t *testing.T) {
	i := getTestImage(t)
	i.Config.URL.CDNSubDomain = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s-3.%s", i.Config.URL.SubDomain, i.Config.URL.Domain))
}

func TestAsset_SecureSubDomainTrue(t *testing.T) {
	i := getTestImage(t)
	i.Config.URL.CDNSubDomain = true
	i.Config.URL.SecureCDNSubDomain = true
	i.Config.URL.PrivateCDN = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s-%s-3.%s", i.Config.Cloud.CloudName, i.Config.URL.SubDomain, i.Config.URL.Domain))
}

func TestAsset_CName(t *testing.T) {
	host := "hello.com"

	i := getTestImage(t)
	i.Config.URL.Secure = false
	i.Config.URL.CName = host

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("%s://%s/%s/", i.Config.URL.Protocol(), host, i.Config.Cloud.CloudName))
}

func TestAsset_CNameSubDomain(t *testing.T) {
	host := "hello.com"

	i := getTestImage(t)
	i.Config.URL.Secure = false
	i.Config.URL.CName = host
	i.Config.URL.CDNSubDomain = true

	assert.Contains(t, getAssetUrl(t, i), fmt.Sprintf("a3.%s/%s/", host, i.Config.Cloud.CloudName))
}

func assertURLProtocol(t *testing.T, a *asset.Asset, protocol string) {
	url := getAssetUrl(t, a)

	assert.Regexp(t, fmt.Sprintf("^%s://", protocol), url)
}

func getAssetUrl(t *testing.T, a *asset.Asset) string {
	url, err := a.String()
	if err != nil {
		t.Fatal(err)
	}

	return url
}

func getTestImage(t *testing.T) *asset.Asset {
	i, err := asset.Image(cldtest.PublicID, nil)
	if err != nil {
		t.Fatal(err)
	}

	return i
}
