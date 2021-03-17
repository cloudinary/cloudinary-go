package uploader

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
)

var ctx = context.Background()
var uploadApi, _ = New()

const LogoUrl = "https://cloudinary-res.cloudinary.com/image/upload/cloudinary_logo.png"
const LogoFilePath = "testdata/cloudinary_logo.png"
const Base64Image = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
const publicID = "go_test_image"
const publicID2 = "go_test_image_2"
const largeImagePublicID = "go_test_large_image"
const largeImageSize = 5880138
const largeChunkSize = 5243000
const largeImageWidth = 1400
const largeImageHeight = 1400
const tag1 = "go_tag1"
const tag2 = "go_tag2"

var tags = api.CldApiArray{tag1, tag2}
var cContext = api.CldApiMap{"go-context-key": "go-context-value"}

func TestUploader_UploadLocalPath(t *testing.T) {
	params := UploadParams{
		PublicID:              publicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
		Overwrite:             true,
	}

	resp, err := uploadApi.Upload(ctx, LogoFilePath, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadIOReader(t *testing.T) {
	file, err := os.Open(LogoFilePath)
	if err != nil {
		t.Error(fmt.Printf("unable to open a file: %v\n", err))
	}

	defer api.DeferredClose(file)

	params := UploadParams{
		PublicID:              publicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
	}

	resp, err := uploadApi.Upload(ctx, file, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadUrl(t *testing.T) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadBase64Image(t *testing.T) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
	}

	resp, err := uploadApi.Upload(ctx, Base64Image, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func TestUploader_UploadLargeFile(t *testing.T) {
	uploadApi.Config.Api.ChunkSize = largeChunkSize

	largeImage := populateLargeImage()

	defer func() {
		err := os.Remove(largeImage)
		if err != nil {
			t.Error(err)
		}
	}()

	params := UploadParams{
		PublicID:  largeImagePublicID,
		Overwrite: true,
	}

	resp, err := uploadApi.Upload(ctx, largeImage, params)

	if err != nil {
		t.Error(err)
	}

	// FIXME: destroy in teardown when available
	_, _ = uploadApi.Destroy(ctx, DestroyParams{PublicID: largeImagePublicID})

	if resp == nil ||
		resp.PublicID != largeImagePublicID ||
		resp.Width != largeImageWidth ||
		resp.Height != largeImageHeight {
		t.Error(resp)
	}

}

func TestUploader_Timeout(t *testing.T) {
	var originalTimeout = uploadApi.Config.Api.Timeout

	uploadApi.Config.Api.Timeout = 0 // should timeout immediately

	_, err := uploadApi.Upload(ctx, LogoUrl, UploadParams{})

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	uploadApi.Config.Api.Timeout = originalTimeout
}

func UploadTestAsset(t *testing.T, publicID string) {
	params := UploadParams{
		PublicID:  publicID,
		Overwrite: true,
		Tags:      tags,
	}

	resp, err := uploadApi.Upload(ctx, LogoUrl, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID {
		t.Error(resp)
	}
}

func populateLargeImage() string {
	head := "BMJ\xB9Y\x00\x00\x00\x00\x00\x8A\x00\x00\x00|\x00\x00\x00x\x05\x00\x00x\x05\x00\x00\x01\x00\x18\x00" +
		"\x00\x00\x00\x00\xC0\xB8Y\x00a\x0F\x00\x00a\x0F\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF" +
		"\x00\x00\xFF\x00\x00\xFF\x00\x00\x00\x00\x00\x00\xFFBGRs\x00\x00\x00\x00\x00\x00\x00\x00T\xB8\x1E" +
		"\xFC\x00\x00\x00\x00\x00\x00\x00\x00fff\xFC\x00\x00\x00\x00\x00\x00\x00\x00\xC4\xF5(\xFF\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	tmpFile, err := ioutil.TempFile("testdata", largeImagePublicID+".*.bmp")
	if err != nil {
		log.Fatal(err)
	}

	content := head + strings.Repeat("\xFF", largeImageSize-len(head))

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		_ = tmpFile.Close()
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpFile.Name()
}
