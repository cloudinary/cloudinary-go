package asset_test

import (
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/asset"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

const authTokenKey = "00112233FF99"
const authTokenAltKey = "CCBB2233FF00"

const duration = 300
const startTime = 11111111

const authTokenTestImage = "sample.jpg"
const authTokenTestConfigACL = "/*/t_foobar"
const authTokenTestPath = "http://res.cloudinary.com/test123/image/upload/v1486020273/sample.jpg"

var authTokenConfig = config.AuthToken{
	Duration:  duration,
	StartTime: startTime,
	Key:       authTokenKey,
	ACL:       "/image/*",
}

func TestAsset_AuthToken_GenerateWithStartTimeAndDuration(t *testing.T) {
	newConfig := authTokenConfig
	a := asset.AuthToken{Config: &newConfig}
	a.Config.StartTime = 1111111111

	expectedToken := "__cld_token__=st=1111111111~exp=1111111411~acl=%2fimage%2f*" +
		"~hmac=1751370bcc6cfe9e03f30dd1a9722ba0f2cdca283fa3e6df3342a00a7528cc51"

	assert.Equal(t, a.Generate(""), expectedToken)
}

func TestAsset_AuthToken_MustProvideExpirationOrDuration(t *testing.T) {
	a := asset.AuthToken{Config: &config.AuthToken{Key: authTokenKey}}

	assert.Panics(t, func() { a.Generate("") })
}

func TestAsset_AuthToken_NoStartTimeRequired(t *testing.T) {
	a := asset.AuthToken{Config: &config.AuthToken{Key: authTokenKey, Expiration: startTime + duration}}

	expectedToken := "__cld_token__=exp=11111411~hmac=470d32e3ee9b872d64bd00d974c559d96892398c8542ff33ce2f2647ee1bf7a4"
	assert.Equal(t, expectedToken, a.Generate(authTokenTestImage))
}

func TestAsset_AuthToken_ShouldIgnoreUrlIfAclIsProvided(t *testing.T) {
	a := asset.AuthToken{Config: &authTokenConfig}
	aclToken := a.Generate("")
	aclTokenUrlIgnored := a.Generate(authTokenTestImage)

	a.Config.ACL = ""

	urlToken := a.Generate(authTokenTestImage)

	assert.NotEqual(t, aclToken, urlToken)
	assert.Equal(t, aclToken, aclTokenUrlIgnored)
}

func TestAsset_AuthToken_EscapeToLower(t *testing.T) {
	a := asset.AuthToken{Config: &authTokenConfig}
	a.Config.ACL = ""

	expected := "__cld_token__=st=11111111~exp=11111411~hmac=7ffc0fd1f3ee2622082689f64a65454da39d94c297bcf498b682aa65a0d2ce0a"

	assert.Equal(t, expected, a.Generate("Encode these :~@#%^&{}[]\\\"';/\", but not those $!()_.*"))
}

func TestAsset_AuthToken_URLWithoutACL(t *testing.T) {
	i, _ := asset.Image(authTokenTestImage, nil)

	i.AuthToken.Config = &authTokenConfig
	i.AuthToken.Config.ACL = ""

	i.DeliveryType = api.Authenticated
	i.Version = 1486020273
	i.Config.URL.SignURL = true

	expected := "__cld_token__=st=11111111~exp=11111411~hmac=8db0d753ee7bbb9e2eaf8698ca3797436ba4c20e31f44527e43b6a6e995cfdb3"

	iURL, _ := i.String()
	assert.Contains(t, iURL, expected)
}
