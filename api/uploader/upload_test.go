package uploader_test

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/internal/signature"
)

// TestUploader_VerifyApiResponseSignature tests API response signature verification.
func TestUploader_VerifyApiResponseSignature(t *testing.T) {
	const publicID1 = "b8sjhoslj8cq8ovoa0ma"
	const publicID2 = "z5sjhoskl2cq8ovoa0mv"
	const version1 = "1555337587"
	const version2 = "1555337588"

	// Test valid signature
	urlParams := make(url.Values)
	urlParams.Set("public_id", publicID1)
	urlParams.Set("version", version1)
	correctSignature, err := api.SignParametersUsingAlgoAndVersion(urlParams, uploadAPI.Config.Cloud.APISecret,
		uploadAPI.Config.Cloud.GetSignatureAlgorithm(), 1)
	if err != nil {
		t.Error(err)
	}

	isValid := uploadAPI.VerifyApiResponseSignature(publicID1, version1, correctSignature)
	assert.True(t, isValid, "The response signature is valid for the same parameters")

	// Test invalid signature with wrong version
	urlParams.Set("version", version2)
	newVersionSignature, err := api.SignParametersUsingAlgoAndVersion(urlParams, uploadAPI.Config.Cloud.APISecret,
		uploadAPI.Config.Cloud.GetSignatureAlgorithm(), 1)
	if err != nil {
		t.Error(err)
	}

	isValid = uploadAPI.VerifyApiResponseSignature(publicID1, version1, newVersionSignature)
	assert.False(t, isValid, "The response signature is invalid for the wrong version")

	// Test invalid signature with wrong resource
	urlParams.Set("version", version1)
	urlParams.Set("public_id", publicID2)
	anotherResourceSignature, err := api.SignParametersUsingAlgoAndVersion(urlParams, uploadAPI.Config.Cloud.APISecret,
		uploadAPI.Config.Cloud.GetSignatureAlgorithm(), 1)
	if err != nil {
		t.Error(err)
	}

	isValid = uploadAPI.VerifyApiResponseSignature(publicID1, version1, anotherResourceSignature)
	assert.False(t, isValid, "The response signature is invalid for the wrong resource")
}

// TestUploader_VerifyApiResponseSignatureWithAmpersand tests signature verification with & characters.
func TestUploader_VerifyApiResponseSignatureWithAmpersand(t *testing.T) {
	const testSecret = "hdcixPpR2iKERPwqvH6sHdK9cyac"

	tempConfig := uploadAPI.Config
	tempConfig.Cloud.APISecret = testSecret
	tempAPI := &uploader.API{Config: tempConfig}
	publicIDWithAmpersand := "callback?a=1&tags=hello,world"
	version := "1568810420"

	urlParams := make(url.Values)
	urlParams.Set("public_id", publicIDWithAmpersand)
	urlParams.Set("version", version)
	v1Signature, err := api.SignParametersUsingAlgoAndVersion(urlParams, testSecret, signature.SHA1, 1)
	if err != nil {
		t.Error(err)
	}
	isValid := tempAPI.VerifyApiResponseSignature(publicIDWithAmpersand, version, v1Signature)
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
