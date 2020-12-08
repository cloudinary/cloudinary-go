package uploader

import (
	"testing"
)

func TestUploader_Creative(t *testing.T) {
	UploadTestAsset(t, publicID)
	UploadTestAsset(t, publicID2)

	params := GenerateSpriteParams{
		Tag: tag1,
	}

	resp, err := uploadApi.GenerateSprite(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.ImageInfos) != 2 {
		t.Error(resp)
	}

	mParams := MultiParams{
		Tag: tag1,
	}

	mResp, err := uploadApi.Multi(ctx, mParams)

	if err != nil {
		t.Error(err)
	}

	if mResp == nil || mResp.PublicID != tag1 {
		t.Error(mResp)
	}

	eParams := ExplodeParams{
		PublicID: tag1,
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
