package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"testing"
)

func TestUploader_Explicit(t *testing.T) {
	UploadTestAsset(t, publicID)

	params := ExplicitParams{
		UploadParams{
			PublicID: publicID,
			Type:     api.Upload,
			Tags:     tags,
		},
	}

	resp, err := uploadApi.Explicit(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.Tags) != 2 {
		t.Error(resp)
	}
}

func TestUploader_Edit(t *testing.T) {
	UploadTestAsset(t, publicID)

	params := RenameParams{
		FromPublicID: publicID,
		ToPublicID:   publicID2,
		Overwrite:    true,
	}

	resp, err := uploadApi.Rename(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.PublicID != publicID2 {
		t.Error(resp)
	}

	dParams := DestroyParams{
		PublicID: publicID2,
	}

	dResp, err := uploadApi.Destroy(ctx, dParams)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || dResp.Result != "ok" {
		t.Error(resp)
	}
}
