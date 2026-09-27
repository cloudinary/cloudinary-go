package uploader_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/internal/signature"
)

// TestUploader_VerifyApiResponseSignature tests API response signature verification.
// Expected signatures are fixed values, the same as other Cloudinary SDKs compute.
func TestUploader_VerifyApiResponseSignature(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"
	const publicID = "b8sjhoslj8cq8ovoa0ma"
	const version = "1555337587"

	tempConfig := uploadAPI.Config
	tempConfig.Cloud.APISecret = testSecret
	tempConfig.Cloud.SignatureAlgorithm = signature.SHA1
	tempAPI := &uploader.API{Config: tempConfig}

	const validSignature = "3e974ccf50d4d1b195d2046c95994fe50ff81267"

	assert.True(t, tempAPI.VerifyApiResponseSignature(publicID, version, validSignature),
		"The response signature is valid for the same parameters")
	assert.False(t, tempAPI.VerifyApiResponseSignature(publicID, "1555337588", validSignature),
		"The response signature is invalid for the wrong version")
	assert.False(t, tempAPI.VerifyApiResponseSignature("z5sjhoskl2cq8ovoa0mv", version, validSignature),
		"The response signature is invalid for the wrong resource")

	tempAPI.Config.Cloud.SignatureAlgorithm = signature.SHA256
	assert.True(t, tempAPI.VerifyApiResponseSignature(publicID, version,
		"80f837513994d089160b01dd3ee04e6b86f127f75fab820dca5c15057deddc48"),
		"The response signature is valid with SHA256")
}

// TestUploader_VerifyApiResponseSignatureWithAmpersand tests signature verification with & characters.
func TestUploader_VerifyApiResponseSignatureWithAmpersand(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"

	tempConfig := uploadAPI.Config
	tempConfig.Cloud.APISecret = testSecret
	tempConfig.Cloud.SignatureAlgorithm = signature.SHA1
	tempAPI := &uploader.API{Config: tempConfig}

	// Signature version 1 does not encode "&" in values.
	isValid := tempAPI.VerifyApiResponseSignature("callback?a=1&tags=hello,world", "1568810420",
		"8bda7afdf8ccf991581ddb4782598cf2fc915b8f")
	assert.True(t, isValid, "Should verify signature correctly with version 1")
}

// TestUploader_VerifyNotificationSignature tests webhook notification signature verification.
func TestUploader_VerifyNotificationSignature(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"
	body := `{"public_id":"b8sjhoslj8cq8ovoa0ma","version":"1555337587","width":"1000","height":"800"}`

	tempConfig := uploadAPI.Config
	tempConfig.Cloud.APISecret = testSecret
	tempConfig.Cloud.SignatureAlgorithm = signature.SHA1
	tempAPI := &uploader.API{Config: tempConfig}

	currentTimestamp := time.Now().Unix()
	validResponseTimestamp := currentTimestamp - 5000

	payload := fmt.Sprintf("%s%d", body, validResponseTimestamp)
	rawSignature, _ := signature.Sign(payload, testSecret, signature.SHA1)
	validSignature := hex.EncodeToString(rawSignature)

	// Test valid signature with sufficient time
	isValid := tempAPI.VerifyNotificationSignature(body, validResponseTimestamp, validSignature, 7200)
	assert.True(t, isValid, "The notification signature is valid for matching and not expired signature")

	// Test expired signature
	isValid = tempAPI.VerifyNotificationSignature(body, validResponseTimestamp, validSignature, 4000)
	assert.False(t, isValid, "The notification signature is invalid for matching but expired signature")

	// Test invalid signature
	isValid = tempAPI.VerifyNotificationSignature(body, validResponseTimestamp, validSignature+"chars", 7200)
	assert.False(t, isValid, "The notification signature is invalid for non matching and not expired signature")

	// Test invalid signature with expiration
	isValid = tempAPI.VerifyNotificationSignature(body, validResponseTimestamp, validSignature+"chars", 4000)
	assert.False(t, isValid, "The notification signature is invalid for non matching and expired signature")
}

// TestUploader_VerifyNotificationSignatureWithSHA256 tests notification signature verification with SHA256.
func TestUploader_VerifyNotificationSignatureWithSHA256(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"
	body := `{}`

	tempConfig := uploadAPI.Config
	tempConfig.Cloud.APISecret = testSecret
	tempConfig.Cloud.SignatureAlgorithm = signature.SHA256
	tempAPI := &uploader.API{Config: tempConfig}

	currentTimestamp := time.Now().Unix()
	validResponseTimestamp := currentTimestamp - 5000

	payload := fmt.Sprintf("%s%d", body, validResponseTimestamp)
	rawSignature, _ := signature.Sign(payload, testSecret, signature.SHA256)
	correctSignature := hex.EncodeToString(rawSignature)

	// Test valid signature with SHA256 algorithm
	isValid := tempAPI.VerifyNotificationSignature(body, validResponseTimestamp, correctSignature, 7200)
	assert.True(t, isValid, "The notification signature is valid with SHA256")
}
