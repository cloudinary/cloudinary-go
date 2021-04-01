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

const LogoUrl = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
const PublicID = "go_test_image"
const PublicID2 = "go_test_image_2"
const VideoPublicID = "go_test_video"
const Tag1 = "go_tag1"
const Tag2 = "go_tag2"

var ImageFilePath = TestDataDir() + "cloudinary_logo.png"
var VideoFilePath = TestDataDir() + "movie.mp4"

var TestSuffix = GetTestSuffix()

var Tags = []string{Tag1, Tag2}
var CldContext = map[string]string{"go-context-key": "go-context-value"}

var ctx = context.Background()
var uploadApi, _ = uploader.New()

func UploadTestAsset(t *testing.T, publicID string) {
	params := uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      Tags,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func UploadTestVideoAsset(t *testing.T, publicID string) {
	params := uploader.UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      Tags,
	}

	resp, err := uploadApi.Upload(ctx, VideoFilePath, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func GetTestSuffix() string {
	testSuffix := os.Getenv("TRAVIS_JOB_ID")

	if testSuffix == "" {
		rand.Seed(time.Now().UnixNano())
		testSuffix = strconv.Itoa(rand.Intn(999999))
	}

	return testSuffix
}

func UniqueID(prefix string) string {
	return prefix + TestSuffix
}

func TestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))

	return filepath.Dir(d) + "/cldtest/testdata/"
}
