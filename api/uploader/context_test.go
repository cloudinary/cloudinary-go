package uploader

import (
	"github.com/cloudinary/cloudinary-go/api"
	"testing"
)

func TestUploader_Context(t *testing.T) {
	UploadTestAsset(t, publicID)

	params := AddContextParams{
		PublicIDs: api.CldApiArray{publicID},
		Context:   cContext,
	}

	resp, err := uploadApi.AddContext(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.PublicIds) != 1 || resp.PublicIds[0] != publicID {
		t.Error(resp)
	}

	raParams := RemoveAllContextParams{
		PublicIDs: api.CldApiArray{publicID},
	}

	raResp, err := uploadApi.RemoveAllContext(ctx, raParams)

	if err != nil {
		t.Error(err)
	}

	if raResp == nil || len(raResp.PublicIds) != 1 || raResp.PublicIds[0] != publicID {
		t.Error(resp)
	}
}
