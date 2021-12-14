package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestUploader_Tag(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	cldtest.UploadTestAsset(t, cldtest.PublicID2)

	params := uploader.AddTagParams{
		PublicIDs: []string{cldtest.PublicID, cldtest.PublicID2},
		Tag:       cldtest.Tag1,
	}

	resp, err := uploadAPI.AddTag(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.PublicIds) != 2 || resp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}

	rParams := uploader.RemoveTagParams{
		PublicIDs: []string{cldtest.PublicID, cldtest.PublicID2},
		Tag:       cldtest.Tag1,
	}

	rResp, err := uploadAPI.RemoveTag(ctx, rParams)

	if err != nil {
		t.Error(err)
	}

	if rResp == nil || len(rResp.PublicIds) != 2 || rResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}
	// FIXME: add some tags :) before removing
	raParams := uploader.RemoveAllTagsParams{
		PublicIDs: api.CldAPIArray{cldtest.PublicID, cldtest.PublicID2},
	}

	raResp, err := uploadAPI.RemoveAllTags(ctx, raParams)

	if err != nil {
		t.Error(err)
	}

	if raResp == nil || len(raResp.PublicIds) != 2 || raResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}

	reParams := uploader.ReplaceTagParams{
		PublicIDs: api.CldAPIArray{cldtest.PublicID, cldtest.PublicID2},
		Tag:       cldtest.Tag2,
	}

	reResp, err := uploadAPI.ReplaceTag(ctx, reParams)

	if err != nil {
		t.Error(err)
	}

	if reResp == nil || len(reResp.PublicIds) != 2 || reResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}
}
