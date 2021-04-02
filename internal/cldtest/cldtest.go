package cldtest

import (
	"context"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

// LogoURL is the URL of the publicly available logo.
const LogoURL = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"

// Base64Image us a base64 encoded test image.
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

// PublicID is the test public ID.
const PublicID = "go_test_image"

// PublicID2 is another test public ID.
const PublicID2 = "go_test_image_2"

// VideoPublicID is the public ID of the test video.
const VideoPublicID = "go_test_video"

// Tag1 is the test tag.
const Tag1 = "go_tag1"

// Tag2 is another test tag.
const Tag2 = "go_tag2"

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
