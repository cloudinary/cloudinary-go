package uploader_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

func TestUploader_GenerateSprite(t *testing.T) {
	cldtest.UploadTestAsset(t, cldtest.PublicID)
	cldtest.UploadTestAsset(t, cldtest.PublicID2)

	params := uploader.GenerateSpriteParams{
		Tag: cldtest.Tag1,
	}

	resp, err := uploadAPI.GenerateSprite(ctx, params)

	if err != nil {
		t.Error(err)
	}

	if resp == nil || len(resp.ImageInfos) < 2 {
		t.Error(resp)
	}
}
func TestUploader_Multi(t *testing.T) {
	mParams := uploader.MultiParams{
		Tag: cldtest.Tag1,
	}

	mResp, err := uploadAPI.Multi(ctx, mParams)

	if err != nil {
		t.Error(err)
	}

	if mResp == nil || mResp.PublicID != cldtest.Tag1 {
		t.Error(mResp)
	}
}
func TestUploader_Explode(t *testing.T) {
	eParams := uploader.ExplodeParams{
		PublicID: cldtest.Tag1,
		Type:     "multi",
	}

	eResp, err := uploadAPI.Explode(ctx, eParams)

	if err != nil {
		t.Error(err)
	}

	if eResp == nil || eResp.BatchID == "" {
		t.Error(eResp)
	}
}

func TestUploader_Text(t *testing.T) {
	tParams := uploader.TextParams{
		Text:       "HelloGo",
		PublicID:   cldtest.PublicID,
		FontFamily: "Arial",
		FontSize:   20,
	}

	tResp, err := uploadAPI.Text(ctx, tParams)

	if err != nil {
		t.Error(err)
	}

	if tResp == nil || tResp.PublicID == "" {
		t.Error(tResp)
	}
}
