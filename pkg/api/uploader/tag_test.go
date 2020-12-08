package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"testing"
)

func TestUploader_Tag(t *testing.T) {
	UploadTestAsset(t, publicID)

	params := AddTagParams{
		PublicIDs: api.CldApiArray{publicID},
		Tag:       tag1,
	}

	resp, err := uploadApi.AddTag(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.PublicIds) != 1 || resp.PublicIds[0] != publicID {
		t.Error(resp)
	}

	rParams := RemoveTagParams{
		PublicIDs: api.CldApiArray{publicID},
		Tag:       tag1,
	}

	rResp, err := uploadApi.RemoveTag(ctx, rParams)

	if err != nil {
		t.Error(err)
	}

	if rResp == nil || len(rResp.PublicIds) != 1 || rResp.PublicIds[0] != publicID {
		t.Error(resp)
	}
	// FIXME: add some tags :) before removing
	raParams := RemoveAllTagsParams{
		PublicIDs: api.CldApiArray{publicID},
	}

	raResp, err := uploadApi.RemoveAllTags(ctx, raParams)

	if err != nil {
		t.Error(err)
	}

	if raResp == nil || len(raResp.PublicIds) != 1 || raResp.PublicIds[0] != publicID {
		t.Error(resp)
	}

	reParams := ReplaceTagParams{
		PublicIDs: api.CldApiArray{publicID},
		Tag:       tag2,
	}

	reResp, err := uploadApi.ReplaceTag(ctx, reParams)

	if err != nil {
		t.Error(err)
	}

	if reResp == nil || len(reResp.PublicIds) != 1 || reResp.PublicIds[0] != publicID {
		t.Error(resp)
	}
}
