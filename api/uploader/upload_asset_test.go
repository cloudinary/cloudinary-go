package uploader_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
)

var ctx = context.Background()

const largeImagePublicID = "go_test_large_image"
const largeImageSize = 5880138
const largeChunkSize = 5243000
const largeImageWidth = 1400
const largeImageHeight = 1400

func TestUploader_UploadLocalPath(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:              cldtest.PublicID,
		QualityAnalysis:       api.Bool(true),
		AccessibilityAnalysis: api.Bool(true),
		CinemagraphAnalysis:   api.Bool(true),
		Overwrite:             api.Bool(true),
		Timestamp:             time.Now().Unix(),
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
		QualityAnalysis:       api.Bool(true),
		AccessibilityAnalysis: api.Bool(true),
		CinemagraphAnalysis:   api.Bool(true),
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
		Overwrite: api.Bool(true),
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
		PublicID:      cldtest.PublicID,
		ResourceType:  "video",
		Overwrite:     api.Bool(true),
		MediaMetadata: api.Bool(true),
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.VideoURL, params)

	if err != nil {
		t.Error(err)
	}
	if resp == nil || resp.PublicID != cldtest.PublicID || resp.Error.Message != "" {
		t.Error(resp)
	}
	if resp.PlaybackURL == "" {
		t.Error("PlaybackURL is empty")
	}
}

func TestUploader_UploadBase64Image(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: api.Bool(true),
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
		Overwrite: api.Bool(true),
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
	// buffered upload allocated about 10000000 bytes and piped allocated 100000, so assertion value is between them
	const memoryAllocationLimit uint64 = 1000000

	uploadAPI.Config.API.ChunkSize = largeChunkSize

	largeImage := populateLargeImage()

	params := uploader.UploadParams{
		PublicID:  largeImagePublicID,
		Overwrite: api.Bool(true),
	}

	largeImageHandle, err := os.Open(largeImage)
	if err != nil {
		t.Error(err)
	}

	defer api.DeferredClose(largeImageHandle)

	fi, err := largeImageHandle.Stat()
	if err != nil {
		t.Error(err)
	}

	sectionReader := io.NewSectionReader(largeImageHandle, 0, fi.Size())

	largeUploads := []interface{}{
		largeImage,
		largeImageHandle,
		sectionReader,
	}

	for _, largeFile := range largeUploads {
		stop := watchMemoryAllocation()
		resp, err := uploadAPI.Upload(ctx, largeFile, params)
		maxAllocation := stop()

		if err != nil {
			t.Error(err)
		}

		if resp == nil ||
			resp.PublicID != largeImagePublicID ||
			resp.Width != largeImageWidth ||
			resp.Height != largeImageHeight {
			t.Error(resp)
		}

		if maxAllocation > memoryAllocationLimit {
			t.Errorf("Uploading took too much memory (allocated %d bytes, expected %d at most)", maxAllocation, memoryAllocationLimit)
		}
	}

	defer func() {
		err := os.Remove(largeImage)
		if err != nil {
			t.Error(err)
		}

		_, _ = uploadAPI.Destroy(ctx, uploader.DestroyParams{PublicID: largeImagePublicID})
	}()
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
		Overwrite: api.Bool(true),
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

func TestUploader_UploadWithMetadata(t *testing.T) {
	externalID := cldtest.CreateStringMetadataField(t, "upload_metadata_field_")
	externalID2 := cldtest.CreateStringMetadataField(t, "upload_metadata_field_2_")
	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Overwrite: api.Bool(true),
		Metadata: api.Metadata{
			externalID:  cldtest.UniqueID("1")[:6],
			externalID2: cldtest.UniqueID("2")[:6],
		},
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.LogoURL, params)

	// FIXME: use setUp/tearDown
	cldtest.DeleteTestMetadataField(t, externalID)
	cldtest.DeleteTestMetadataField(t, externalID2)

	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error(resp)
	}

	assert.Equal(t, params.Metadata[externalID], resp.Metadata[externalID])
	assert.Equal(t, params.Metadata[externalID2], resp.Metadata[externalID2])
}

func TestUploader_UploadWithResponsiveBreakpoints(t *testing.T) {
	params := uploader.UploadParams{
		PublicID:              cldtest.PublicID,
		Overwrite:             api.Bool(true),
		ResponsiveBreakpoints: uploader.ResponsiveBreakpointsParams{{CreateDerived: api.Bool(false), Transformation: "a_90"}},
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
			{CreateDerived: api.Bool(false), Transformation: "a_90"},
			{CreateDerived: api.Bool(false), Transformation: "a_45"},
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

func TestUploader_UploadWithAccessControl(t *testing.T) {
	now := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 1)

	params := uploader.UploadParams{
		PublicID:  cldtest.PublicID,
		Type:      api.Authenticated,
		Overwrite: api.Bool(true),
		AccessControl: api.AccessControl{
			{
				AccessType: api.Anonymous,
				Start:      api.TimePtr(now),
				End:        api.TimePtr(tomorrow),
			},
			{
				AccessType: api.Token,
			},
		},
	}

	resp, err := uploadAPI.Upload(ctx, cldtest.LogoURL, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil {
		t.Error(resp)
	}

	assert.Len(t, resp.AccessControl, 2)

	assert.Equal(t, api.Anonymous, resp.AccessControl[0].AccessType)
	assert.Equal(t, now.Unix(), resp.AccessControl[0].Start.Unix())
	assert.Equal(t, tomorrow.Unix(), resp.AccessControl[0].End.Unix())

	assert.Equal(t, api.Token, resp.AccessControl[1].AccessType)
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

func readMemoryAlloc() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

func watchMemoryAllocation() func() uint64 {
	var maxAlloc uint64
	proceed := true
	initialRead := readMemoryAlloc()
	go func() {
		for i := 0; i < 10000; i++ {
			if !proceed {
				return
			}
			read := readMemoryAlloc()
			if read-initialRead > maxAlloc {
				maxAlloc = read - initialRead
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()

	return func() uint64 {
		proceed = false
		return maxAlloc
	}
}
