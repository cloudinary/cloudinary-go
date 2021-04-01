package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestUploader_Creative(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	cldtest.UploadTestAsset(t, cldtest.PublicID2)

	params := uploader.GenerateSpriteParams{
		Tag: cldtest.Tag1,
	}

	resp, err := uploadApi.GenerateSprite(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.ImageInfos) < 2 {
		t.Error(resp)
	}

	mParams := uploader.MultiParams{
		Tag: cldtest.Tag1,
	}

	mResp, err := uploadApi.Multi(ctx, mParams)

	if err != nil {
		t.Error(err)
	}

	if mResp == nil || mResp.PublicID != cldtest.Tag1 {
		t.Error(mResp)
	}

	eParams := uploader.ExplodeParams{
		PublicID: cldtest.Tag1,
		Type:     "multi",
	}

	eResp, err := uploadApi.Explode(ctx, eParams)

	if err != nil {
		t.Error(err)
	}

	if eResp == nil || eResp.BatchID == "" {
		t.Error(eResp)
	}
}
