package config_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/internal/signature"
	"github.com/stretchr/testify/assert"
)

var cldURL = "cloudinary://key:secret@test123"
var fakeOAuthToken = "MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4"

func TestConfiguration_CreateInstance(t *testing.T) {
	c, _ := config.New()

	if c.Cloud.CloudName == "" {
		t.Error("Please set up CLOUDINARY_URL environment variable to run the test.")
	}

	c, err := config.NewFromURL(cldURL + "?signature_algorithm=sha256")
	if err != nil {
		t.Error("Error: ", err)
	}

	assert.Equal(t, "test123", c.Cloud.CloudName)
	assert.Equal(t, signature.SHA256, c.Cloud.SignatureAlgorithm)

	c, err = config.NewFromURL("")
	if err == nil || err.Error() != "must provide CLOUDINARY_URL" {
		t.Error("Error expected, got: ", err)
	}

	c, _ = config.NewFromParams("test123", "key", "secret")

	assert.Equal(t, "test123", c.Cloud.CloudName)
	assert.Equal(t, "key", c.Cloud.APIKey)
	assert.Equal(t, "secret", c.Cloud.APISecret)

	c, _ = config.NewFromOAuthToken("test123", fakeOAuthToken)

	assert.Equal(t, "test123", c.Cloud.CloudName)
	assert.Equal(t, fakeOAuthToken, c.Cloud.OAuthToken)
	assert.Equal(t, "", c.Cloud.APIKey)
	assert.Equal(t, "", c.Cloud.APISecret)

	// check a few default values
	assert.EqualValues(t, signature.SHA1, c.Cloud.GetSignatureAlgorithm())
	assert.EqualValues(t, 60, c.API.Timeout)
	assert.EqualValues(t, true, c.URL.Secure)
}

func TestConfiguration_API(t *testing.T) {
	c, err := config.NewFromURL(cldURL +
		"?upload_prefix=https://test.prefix.com&timeout=59&upload_timeout=59&chunk_size=7357")
	if err != nil {
		t.Error("Error: ", err)
	}

	assert.Equal(t, "test123", c.Cloud.CloudName)

	assert.Equal(t, "https://test.prefix.com", c.API.UploadPrefix)
	assert.EqualValues(t, 59, c.API.Timeout)
	assert.EqualValues(t, 59, c.API.UploadTimeout)
	assert.EqualValues(t, 7357, c.API.ChunkSize)
}

func TestConfiguration_URL(t *testing.T) {
	c, err := config.NewFromURL(cldURL +
		"?cname=cname.com&secure_cname=secure.cname.com&secure=false&cdn_sub_domain=true" +
		"&secure_cdn_sub_domain=true&private_cdn=true&sign_url=true&long_url_signature=true" +
		"&shorten=true&use_root_path=true&force_version=false&analytics=false")
	if err != nil {
		t.Error("Error: ", err)
	}

	assert.Equal(t, "test123", c.Cloud.CloudName)

	assert.Equal(t, "cname.com", c.URL.CName)
	assert.Equal(t, "secure.cname.com", c.URL.SecureCName)
	assert.Equal(t, false, c.URL.Secure)
	assert.Equal(t, true, c.URL.CDNSubDomain)
	assert.Equal(t, true, c.URL.SecureCDNSubDomain)
	assert.Equal(t, true, c.URL.PrivateCDN)
	assert.Equal(t, true, c.URL.SignURL)
	assert.Equal(t, true, c.URL.LongURLSignature)
	assert.Equal(t, true, c.URL.Shorten)
	assert.Equal(t, true, c.URL.UseRootPath)
	assert.Equal(t, false, c.URL.ForceVersion)
	assert.Equal(t, false, c.URL.Analytics)
}

func TestConfiguration_AuthToken(t *testing.T) {
	c, err := config.NewFromURL(cldURL +
		"?key=key&ip=127.0.0.1&acl=*&start_time=1&expiration=3&duration=2")
	if err != nil {
		t.Error("Error: ", err)
	}

	assert.Equal(t, "test123", c.Cloud.CloudName)

	assert.Equal(t, "key", c.AuthToken.Key)
	assert.Equal(t, "127.0.0.1", c.AuthToken.IP)
	assert.Equal(t, "*", c.AuthToken.ACL)
	assert.EqualValues(t, 1, c.AuthToken.StartTime)
	assert.EqualValues(t, 3, c.AuthToken.Expiration)
	assert.EqualValues(t, 2, c.AuthToken.Duration)
}
