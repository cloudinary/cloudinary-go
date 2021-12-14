package uploader_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

var ctx = context.Background()
var uploadAPI, _ = uploader.New()

const largeImagePublicID = "go_test_large_image"
const largeImageSize = 5880138
const largeChunkSize = 5243000
const largeImageWidth = 1400
const largeImageHeight = 1400

func TestUploader_UploadLocalPath(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:              cldtest.PublicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
		Overwrite:             true,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.ImageFilePath, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID {
		t.Error(resp)
	}
}

func TestUploader_UploadIOReader(t *testing.T) {
	file, err := os.Open(cldtest.ImageFilePath)
	if err != nil {
		t.Error(fmt.Printf("unable to open a file: %v\n", err))
	}

	defer api.DeferredClose(file)

	params := uploader.UploadParams{
		PublicID:              cldtest.PublicID,
		QualityAnalysis:       true,
		AccessibilityAnalysis: true,
		CinemagraphAnalysis:   true,
	}

	resp, err := uploadAPI.Upload(ctx, file, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID {
		t.Error(resp)
	}
}

func TestUploader_UploadURL(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: true,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.LogoURL, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID {
		t.Error(resp)
	}
}

func TestUploader_UploadVideoURL(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:     cldtest.PublicID,
		ResourceType: "video",
		Overwrite:    true,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.VideoURL, params)

	if err != nil {
		t.Error(err)
	}
	if resp == nil || resp.PublicID != cldtest.PublicID || resp.Error.Message != "" {
		t.Error(resp)
	}
}

func TestUploader_UploadBase64Image(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: true,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.Base64Image, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID {
		t.Error(resp)
	}
}

func TestUploader_UploadAuthenticated(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: true,
		Type:      api.Authenticated,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.Base64Image, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID {
		t.Error(resp)
	}
}

func TestUploader_UploadLargeFile(t *testing.T) {
	uploadAPI.Config.API.ChunkSize = largeChunkSize

	largeImage := populateLargeImage()

	defer func() {
		err := os.Remove(largeImage)
		if err != nil {
			t.Error(err)
		}
	}()

	params := uploader.UploadParams{
		PublicID:  largeImagePublicID,
		Overwrite: true,
	}

	resp, err := uploadAPI.Upload(ctx, largeImage, params)

	if err != nil {
		t.Error(err)
	}

	// FIXME: destroy in teardown when available
	_, _ = uploadAPI.Destroy(ctx, uploader.DestroyParams{PublicID: largeImagePublicID})

	if resp == nil ||
		resp.PublicID != largeImagePublicID ||
		resp.Width != largeImageWidth ||
		resp.Height != largeImageHeight {
		t.Error(resp)
	}

}

func TestUploader_Timeout(t *testing.T) {
	var originalTimeout = uploadAPI.Config.API.Timeout

	uploadAPI.Config.API.Timeout = 0 // should timeout immediately

	_, err := uploadAPI.Upload(ctx, cldtest.LogoURL, uploader.UploadParams{})

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	uploadAPI.Config.API.Timeout = originalTimeout
}

func TestUploader_UploadWithContext(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: true,
		Context:   cldtest.CldContext,
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.LogoURL, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error(resp)
	}

	assert.Equal(t, fmt.Sprintf("%v", cldtest.CldContext), fmt.Sprintf("%v", resp.Context["custom"]))
}

func TestUploader_UploadWithResponsiveBreakpoints(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:              cldtest.PublicID,
		Overwrite:             true,
		ResponsiveBreakpoints: uploader.ResponsiveBreakpointsParams{{CreateDerived: false, Transformation: "a_90"}},
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.LogoURL, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error(resp)
	}

	assert.Len(t, resp.ResponsiveBreakpoints, 1)
	assert.Equal(t, "a_90", resp.ResponsiveBreakpoints[0].Transformation)

	eParams := uploader.ExplicitParams{
		PublicID: resp.PublicID,
		Type:     api.Upload,
		ResponsiveBreakpoints: uploader.ResponsiveBreakpointsParams{
			{CreateDerived: false, Transformation: "a_90"},
			{CreateDerived: false, Transformation: "a_45"},
		}}

	eResp, err := uploadAPI.Explicit(ctx, eParams)

	if err != nil {
		t.Error(err)
	}

	if eResp == nil {
		t.Error(resp)
	}

	assert.Len(t, eResp.ResponsiveBreakpoints, 2)
	assert.Equal(t, "a_90", eResp.ResponsiveBreakpoints[0].Transformation)
	assert.Equal(t, "a_45", eResp.ResponsiveBreakpoints[1].Transformation)
}

func populateLargeImage() string {
	head := "BMJ\xB9Y\x00\x00\x00\x00\x00\x8A\x00\x00\x00|\x00\x00\x00x\x05\x00\x00x\x05\x00\x00\x01\x00\x18\x00" +
		"\x00\x00\x00\x00\xC0\xB8Y\x00a\x0F\x00\x00a\x0F\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\xFF" +
		"\x00\x00\xFF\x00\x00\xFF\x00\x00\x00\x00\x00\x00\xFFBGRs\x00\x00\x00\x00\x00\x00\x00\x00T\xB8\x1E" +
		"\xFC\x00\x00\x00\x00\x00\x00\x00\x00fff\xFC\x00\x00\x00\x00\x00\x00\x00\x00\xC4\xF5(\xFF\x00\x00\x00" +
		"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	tmpFile, err := ioutil.TempFile(cldtest.TestDataDir(), largeImagePublicID+".*.bmp")
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
