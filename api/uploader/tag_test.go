package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestUploader_Tag(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	params := uploader.AddTagParams{
		PublicIDs: []string{cldtest.PublicID},
		Tag:       cldtest.Tag1,
	}

	resp, err := uploadApi.AddTag(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.PublicIds) != 1 || resp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}

	rParams := uploader.RemoveTagParams{
		PublicIDs: []string{cldtest.PublicID},
		Tag:       cldtest.Tag1,
	}

	rResp, err := uploadApi.RemoveTag(ctx, rParams)

	if err != nil {
		t.Error(err)
	}

	if rResp == nil || len(rResp.PublicIds) != 1 || rResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}
	// FIXME: add some tags :) before removing
	raParams := uploader.RemoveAllTagsParams{
		PublicIDs: api.CldApiArray{cldtest.PublicID},
	}

	raResp, err := uploadApi.RemoveAllTags(ctx, raParams)

	if err != nil {
		t.Error(err)
	}

	if raResp == nil || len(raResp.PublicIds) != 1 || raResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}

	reParams := uploader.ReplaceTagParams{
		PublicIDs: api.CldApiArray{cldtest.PublicID},
		Tag:       cldtest.Tag2,
	}

	reResp, err := uploadApi.ReplaceTag(ctx, reParams)

	if err != nil {
		t.Error(err)
	}

	if reResp == nil || len(reResp.PublicIds) != 1 || reResp.PublicIds[0] != cldtest.PublicID {
		t.Error(resp)
	}
}
