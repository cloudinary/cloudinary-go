package uploader_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"github.com/cloudinary/cloudinary-go/v2/internal/signature"
)

var oAuthTokenConfig, _ = config.NewFromOAuthToken(cldtest.CloudName, "MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4")

const onSuccessStr = "current_asset.update({tags: [\"autocaption\"]});"

// Acceptance test cases for user agent and user platform
func getUserAgentTestCases() []UploadAPIAcceptanceTestCase {
	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test User Agent",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method:  "POST",
				URI:     "/auto/upload",
				Headers: &map[string]string{"User-Agent": api.UserAgent},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
		{
			Name: "Upload Test User Agent With User Platform",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				api.UserPlatform = "Test/1.2.3"
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method:  "POST",
				URI:     "/auto/upload",
				Headers: &map[string]string{"User-Agent": fmt.Sprintf("Test/1.2.3 %s", api.UserAgent)},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for OAuth Authorization
func getAuthorizationTestCases() []UploadAPIAcceptanceTestCase {
	return []UploadAPIAcceptanceTestCase{
		{
			Name:   "Upload Test OAuth Authorization",
			Config: oAuthTokenConfig,
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method:  "POST",
				URI:     "/auto/upload",
				Headers: &map[string]string{"Authorization": "Bearer MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4"},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for folder decoupling
func getFolderDecouplingTestCases() []UploadAPIAcceptanceTestCase {
	body := "asset_folder=asset_folder&display_name=test&file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7&folder=folder%2Ftest&public_id_prefix=fd_public_id_prefix&timestamp=123456789&unique_display_name=true&unsigned=true&use_asset_folder_as_public_id_prefix=true&use_filename_as_display_name=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Folder Decoupling",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					PublicIDPrefix:                 "fd_public_id_prefix",
					DisplayName:                    "test",
					UniqueDisplayName:              api.Bool(true),
					Folder:                         "folder/test",
					AssetFolder:                    "asset_folder",
					UseAssetFolderAsPublicIDPrefix: api.Bool(true),
					UseFilenameAsDisplayName:       api.Bool(true),
					Unsigned:                       api.Bool(true),
					Timestamp:                      123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &body,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for auto transcription
func getAutoTranscriptionTestCases() []UploadAPIAcceptanceTestCase {
	bodyEmpty := "auto_transcription=%7B%7D" +
		"&file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&timestamp=123456789" +
		"&unsigned=true"
	bodyTranslate := "auto_transcription=%7B%22translate%22%3A%5B%22en%22%5D%7D" +
		"&file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&timestamp=123456789" +
		"&unsigned=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Auto Transcription Empty",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					AutoTranscription: &api.AutoTranscription{},
					Unsigned:          api.Bool(true),
					Timestamp:         123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &bodyEmpty,
			},
			ExpectedCallCount: 1,
		},
		{
			Name: "Upload Test Auto Transcription Translate",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					AutoTranscription: &api.AutoTranscription{Translate: []string{"en"}},
					Unsigned:          api.Bool(true),
					Timestamp:         123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &bodyTranslate,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for auto video details
func getAutoVideoDetailsTestCases() []UploadAPIAcceptanceTestCase {
	bodyEmpty := "auto_video_details=%7B%7D" +
		"&file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&timestamp=123456789" +
		"&unsigned=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Auto Video Details Empty",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					AutoVideoDetails: &api.AutoVideoDetails{},
					Unsigned:         api.Bool(true),
					Timestamp:        123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &bodyEmpty,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for handling of boolean values
func getBooleanValuesTestCases() []UploadAPIAcceptanceTestCase {
	body := "file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&timestamp=123456789&unique_filename=false&unsigned=true&use_filename=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Boolean Values",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					UniqueFilename: api.Bool(false),
					UseFilename:    api.Bool(true),
					Unsigned:       api.Bool(true),
					Timestamp:      123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &body,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for handling of various values.
func getVariousValuesTestCases() []UploadAPIAcceptanceTestCase {
	body := "file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&on_success=current_asset.update%28%7Btags%3A+%5B%22autocaption%22%5D%7D%29%3B" +
		"&timestamp=123456789" +
		"&unsigned=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Various Values",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					Timestamp: 123456789,
					Unsigned:  api.Bool(true),
					OnSuccess: onSuccessStr,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &body,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for handling upload configuration.
func getUploadConfigTestCases() []UploadAPIAcceptanceTestCase {
	body := "api_key=key&" +
		"file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7" +
		"&signature=06dbb2b05a8fdce468026e104596e6b8009dee3d08a6a62956ee7fef9aef8c74" +
		"&timestamp=123456789"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Configuration",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				uploadAPI.Config.Cloud.SignatureAlgorithm = signature.SHA256
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					Timestamp: 123456789,
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    "/auto/upload",
				Body:   &body,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestUploadAPI_Acceptance(t *testing.T) {
	t.Parallel()
	testUploadAPIByTestCases(getUserAgentTestCases(), t)
	testUploadAPIByTestCases(getAuthorizationTestCases(), t)
	testUploadAPIByTestCases(getFolderDecouplingTestCases(), t)
	testUploadAPIByTestCases(getAutoTranscriptionTestCases(), t)
	testUploadAPIByTestCases(getBooleanValuesTestCases(), t)
	testUploadAPIByTestCases(getVariousValuesTestCases(), t)
	testUploadAPIByTestCases(getUploadConfigTestCases(), t)
}
