package uploader_test

import (
	"github.com/stretchr/testify/assert"
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


func TestUploader_UpdateMetadata(t *testing.T) {
	externalID := cldtest.CreateStringMetadataField(t, "update_metadata_field_")
	externalID2 := cldtest.CreateStringMetadataField(t, "update_metadata_field2_")
	pID1 := cldtest.UniqueID(cldtest.PublicID)
	pID2 := cldtest.UniqueID(cldtest.PublicID2)
	cldtest.UploadTestAsset(t, pID1)
	cldtest.UploadTestAsset(t, pID2)

	params := uploader.UpdateMetadataParams{
		PublicIDs: []string{pID1, pID2},
		Metadata: api.CldAPIMap{externalID: "upd1", externalID2: "upd2"},
	}

	resp, err := uploadAPI.UpdateMetadata(ctx, params)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(resp.PublicIds))
}
