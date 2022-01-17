package cldtest

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api/admin/metadata"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// LogoURL is the URL of the publicly available logo.
const LogoURL = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"

// VideoURL is the URL of the publicly available video.
const VideoURL = "https://res.cloudinary.com/demo/video/upload/dog.mp4"

// Base64Image us a base64 encoded test image.
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

// PublicID is the test public ID.
const PublicID = "go_test_image"

// PublicID2 is another test public ID.
const PublicID2 = "go_test_image_2"

// ImgExt is the extension of the image.
const ImgExt = ".png"

// VideoExt is the extension of the video.
const VideoExt = ".mp4"

// FileExt is the extension of the file.
const FileExt = ".bin"

// VideoPublicID is the public ID of the test video.
const VideoPublicID = "go_test_video"

// Folder is the test folder path.
const Folder = "test_folder"

// Tag1 is the test tag.
const Tag1 = "go_tag1"

// Tag2 is another test tag.
const Tag2 = "go_tag2"

// SEOName is a SEO friendly name.
const SEOName = "my_favorite_sample"

const Transformation = "c_scale,w_500"

const ApiVersion = "v1_1"

// ImageInFolder is the test public ID in folder.
var ImageInFolder = fmt.Sprintf("%s/%s", Folder, PublicID)

// ImageFilePath is a full path to the test image file.
var ImageFilePath = TestDataDir() + "cloudinary_logo.png"

// VideoFilePath is a full path to the test video file.
var VideoFilePath = TestDataDir() + "movie.mp4"

// TestSuffix is the unique test suffix.
var TestSuffix = GetTestSuffix()

// Tags are the test tags.
var Tags = []string{Tag1, Tag2}

// CldContext is the test context.
var CldContext = map[string]string{"go-context-key": "go-context-value"}

var ctx = context.Background()
var uploadAPI, _ = uploader.New()
var adminAPI, _ = admin.New()

var stringMetadataField = metadata.Field{
	Type:         metadata.StringFieldType,
	ExternalID:   UniqueID("string_md_field_id"),
	Label:        UniqueID("string_md_field_label"),
	DefaultValue: "Gopher",
	Validation:   metadata.StringLengthValidation(2, 6),
}

// UploadTestAsset uploads a test image asset for test purposes.
func UploadTestAsset(t *testing.T, publicID string) {
	params := uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      Tags,
	}

	resp, err := uploadAPI.Upload(ctx, LogoURL, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

// UploadTestVideoAsset uploads a test video asset for test purposes.
func UploadTestVideoAsset(t *testing.T, publicID string) {
	params := uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      Tags,
	}

	resp, err := uploadAPI.Upload(ctx, VideoFilePath, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func CreateStringMetadataField(t *testing.T, prefix string) string {
	stringMetadataField.ExternalID = UniqueID(prefix + "id")
	stringMetadataField.Label = UniqueID(prefix + "label")
	res, err := adminAPI.AddMetadataField(ctx, stringMetadataField)
	if err != nil {
		t.Error(err)
	}
	if res.Error.Message != "" {
		t.Error(res.Error.Message)
	}

	return res.ExternalID
}

// GetTestSuffix returns a unique test suffix.
func GetTestSuffix() string {
	testSuffix := os.Getenv("TRAVIS_JOB_ID")

	if testSuffix == "" {
		rand.Seed(time.Now().UnixNano())
		testSuffix = strconv.Itoa(rand.Intn(999999))
	}

	return testSuffix
}

// UniqueID returns a unique ID for the specified prefix.
func UniqueID(prefix string) string {
	return prefix + TestSuffix
}

// TestDataDir returns the full path to the directory with test files.
func TestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d) + "/cldtest/testdata/"
}

// GetServerMock Get HTTP server mock
func GetServerMock(fn TestFunction) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	})

	srv := httptest.NewServer(handler)

	return srv
}

// ApiResponseTest Test function for the response from the API.
type ApiResponseTest func(response interface{}, t *testing.T)

// TestFunction the test function.
type TestFunction func(w http.ResponseWriter, r *http.Request)

// ExpectedRequestParams are the expected request parameters
type ExpectedRequestParams struct {
	Method  string             // Expected HTTP method of the request
	Uri     string             // Expected URI
	Params  *url.Values        // Expected URI params
	Body    *string            // Expected HTTP body (for POST / PUT requests)
	Headers *map[string]string // Expected HTTP request headers
}

// GetTestHandler gets the test handler for HTTP server. Contains basic checks by expected request params.
func GetTestHandler(response string, t *testing.T, callCounter *int, ep ExpectedRequestParams) TestFunction {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != ep.Method {
			t.Errorf("HTTP method should be %s", ep.Method)
		}

		if ep.Params != nil && ep.Params.Encode() != r.URL.Query().Encode() {
			t.Errorf(
				"Expected query string: %s, got: %s\n",
				ep.Params.Encode(),
				r.URL.Query().Encode(),
			)
		}

		if ep.Headers != nil {
			for expectedName, expectedValue := range *ep.Headers {
				value, present := r.Header[expectedName]
				if !present {
					t.Errorf("Expected request header: '%s' not found\n", expectedName)
				}
				stringValue := strings.Join(value, ", ")
				if expectedValue != stringValue {
					t.Errorf("Expected request header %s value: %s, got: %s\n", expectedName, expectedValue, value)
				}

			}
		}

		expectedURI := "/" + ApiVersion + "/TEST" + ep.Uri
		if expectedURI != r.URL.Path {
			t.Errorf(
				"Expected request URI: %s, got: %s\n",
				expectedURI,
				r.URL.Path,
			)
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			if r.Body != nil && ep.Body != nil {
				bodyString, err := ioutil.ReadAll(r.Body)

				if err != nil {
					t.Error(err)
				}

				if string(bodyString) != *ep.Body {
					t.Errorf("Wrong request body. Expected: %s, given: %s", *ep.Body, string(bodyString))
				}
			}
		}

		*callCounter++
		_, err := io.WriteString(w, response)
		if err != nil {
			t.Error(err)
		}
	}
}
