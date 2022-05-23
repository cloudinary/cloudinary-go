package asset_test

import (
	"github.com/cloudinary/cloudinary-go/v2/asset"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

const authTokenKey = "00112233FF99"
const authTokenAltKey = "CCBB2233FF00"

const duration = 300
const startTime = 1111111111

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
	a := asset.AuthToken{Config: &authTokenConfig}

	expectedToken := "__cld_token__=st=1111111111~exp=1111111411~acl=%2fimage%2f*" +
		"~hmac=1751370bcc6cfe9e03f30dd1a9722ba0f2cdca283fa3e6df3342a00a7528cc51"

	assert.Equal(t, a.Generate(""), expectedToken)
}

func TestAsset_AuthToken_MustProvideExpirationOrDuration(t *testing.T) {
	a := asset.AuthToken{Config: &config.AuthToken{Key: authTokenKey}}

	assert.Panics(t, func() { a.Generate("") })
}

func TestAsset_AuthToken_ShouldIgnoreUrlIfAclIsProvided(t *testing.T) {
	a := asset.AuthToken{Config: &authTokenConfig}
	aclToken := a.Generate("")
	aclTokenUrlIgnored := a.Generate(cldtest.PublicID)

	a.Config.ACL = ""

	urlToken := a.Generate(cldtest.PublicID)

	assert.NotEqual(t, aclToken, urlToken)
	assert.Equal(t, aclToken, aclTokenUrlIgnored)
}

func TestAsset_AuthToken_EscapeToLower(t *testing.T) {
	a := asset.AuthToken{Config: &authTokenConfig}
	a.Config.ACL = ""

	expected := "__cld_token__=st=1111111111~exp=1111111411~hmac=9ee78e220dd8099445b0640986d4255ff2ff4d04609c55c8812d2d2490a0d509"

	assert.Equal(t, expected, a.Generate("Encode these :~@#%^&{}[]\\\"';/\", but not those $!()_.*"))
}
