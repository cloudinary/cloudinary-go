package api_test

import (
	"net/url"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/internal/signature"
	"github.com/stretchr/testify/assert"
)

// TestAPISignRequestPreventsParameterSmuggling tests that parameter smuggling via & characters is prevented.
// Should prevent parameter smuggling via & characters in parameter values.
func TestAPISignRequestPreventsParameterSmuggling(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"

	// Test with notification_url containing & characters
	paramsWithAmpersand := url.Values{}
	paramsWithAmpersand.Set("cloud_name", "dn6ot3ged")
	paramsWithAmpersand.Set("timestamp", "1568810420")
	paramsWithAmpersand.Set("notification_url", "https://fake.com/callback?a=1&tags=hello,world")

	signatureWithAmpersand, err := api.SignParametersUsingAlgo(paramsWithAmpersand, testSecret, signature.SHA1)
	assert.NoError(t, err)

	// Test that attempting to smuggle parameters by splitting the notification_url fails
	paramsSmuggled := url.Values{}
	paramsSmuggled.Set("cloud_name", "dn6ot3ged")
	paramsSmuggled.Set("timestamp", "1568810420")
	paramsSmuggled.Set("notification_url", "https://fake.com/callback?a=1")
	paramsSmuggled.Set("tags", "hello,world") // This would be smuggled if & encoding didn't work

	signatureSmuggled, err := api.SignParametersUsingAlgo(paramsSmuggled, testSecret, signature.SHA1)
	assert.NoError(t, err)

	// The signatures should be different, proving that parameter smuggling is prevented
	assert.NotEqual(t, signatureWithAmpersand, signatureSmuggled,
		"Signatures should be different to prevent parameter smuggling")

	// Verify the expected signature for the properly encoded case
	const expectedSignature = "4fdf465dd89451cc1ed8ec5b3e314e8a51695704"
	assert.Equal(t, expectedSignature, signatureWithAmpersand)

	// Verify the expected signature for the smuggled parameters case
	const expectedSmuggledSignature = "7b4e3a539ff1fa6e6700c41b3a2ee77586a025f9"
	assert.Equal(t, expectedSmuggledSignature, signatureSmuggled)
} 