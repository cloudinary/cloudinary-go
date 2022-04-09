package uploader_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"testing"
)

var oAuthTokenConfig, _ = config.NewFromOAuthToken("TEST", "MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4")

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
				Uri:     "/auto/upload",
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
				Uri:     "/auto/upload",
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
				Uri:     "/auto/upload",
				Headers: &map[string]string{"Authorization": "Bearer MTQ0NjJkZmQ5OTM2NDE1ZTZjNGZmZjI4"},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for folder decoupling
func getFolderDecouplingTestCases() []UploadAPIAcceptanceTestCase {
	body := "asset_folder=asset_folder&display_name=test&file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7&folder=folder%2Ftest&public_id_prefix=fd_public_id_prefix&timestamp=0001-01-01T00%3A00%3A00Z&unsigned=true&use_filename_as_display_name=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Folder Decoupling",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					PublicIDPrefix:           "fd_public_id_prefix",
					DisplayName:              "test",
					Folder:                   "folder/test",
					AssetFolder:              "asset_folder",
					UseFilenameAsDisplayName: api.Bool(true),
					Unsigned:                 api.Bool(true),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				Uri:    "/auto/upload",
				Body:   &body,
			},
			ExpectedCallCount: 1,
		},
	}
}

// Acceptance test cases for handling of boolean values
func getBooleanValuesTestCases() []UploadAPIAcceptanceTestCase {
	body := "file=data%3Aimage%2Fgif%3Bbase64%2CR0lGODlhAQABAIAAAAAAAP%2F%2F%2FyH5BAEAAAAALAAAAAABAAEAAAIBRAA7&timestamp=0001-01-01T00%3A00%3A00Z&unique_filename=false&unsigned=true&use_filename=true"

	return []UploadAPIAcceptanceTestCase{
		{
			Name: "Upload Test Boolean Values",
			RequestTest: func(uploadAPI *uploader.API, ctx context.Context) (interface{}, error) {
				return uploadAPI.Upload(ctx, cldtest.Base64Image, uploader.UploadParams{
					UniqueFilename: api.Bool(false),
					UseFilename:    api.Bool(true),
					Unsigned:       api.Bool(true),
				})
			},
			ResponseTest: func(response interface{}, t *testing.T) {},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				Uri:    "/auto/upload",
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
	testUploadAPIByTestCases(getBooleanValuesTestCases(), t)
}
