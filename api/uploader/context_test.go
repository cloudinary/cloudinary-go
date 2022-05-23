package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
)

func TestUploader_Context(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)

	params := uploader.AddContextParams{
		PublicIDs: api.CldAPIArray{cldtest.PublicID},
		Context:   cldtest.CldContext,
	}

	resp, err := uploadAPI.AddContext(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.PublicIDs) != 1 || resp.PublicIDs[0] != cldtest.PublicID {
		t.Error(resp)
	}

	raParams := uploader.RemoveAllContextParams{
		PublicIDs: api.CldAPIArray{cldtest.PublicID},
	}

	raResp, err := uploadAPI.RemoveAllContext(ctx, raParams)

	if err != nil {
		t.Error(err)
	}

	if raResp == nil || len(raResp.PublicIDs) != 1 || raResp.PublicIDs[0] != cldtest.PublicID {
		t.Error(resp)
	}
}
