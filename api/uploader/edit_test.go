package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestUploader_Explicit(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	params := uploader.ExplicitParams{
		PublicID: cldtest.PublicID,
		Type:     api.Upload,
		Tags:     cldtest.Tags,
	}

	resp, err := uploadAPI.Explicit(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.Tags) != 2 {
		t.Error(resp)
	}
}

func TestUploader_Edit(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	params := uploader.RenameParams{
		FromPublicID: cldtest.PublicID,
		ToPublicID:   cldtest.PublicID2,
		Overwrite:    true,
	}

	resp, err := uploadAPI.Rename(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != cldtest.PublicID2 {
		t.Error(resp)
	}

	dParams := uploader.DestroyParams{
		PublicID: cldtest.PublicID2,
	}

	dResp, err := uploadAPI.Destroy(ctx, dParams)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || dResp.Result != "ok" {
		t.Error(resp)
	}
}
