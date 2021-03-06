package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

const UPName = "go-upload-preset"

func TestUploadPresets_Create(t *testing.T) {

	params := admin.CreateUploadPresetParams{
		Name:         UPName,
		Unsigned:     true,
		Live:         true,
		UploadParams: uploader.UploadParams{Tags: api.CldAPIArray{"go-tag1", "go-tag2"}},
	}

	resp, err := adminAPI.CreateUploadPreset(ctx, params)

	if err != nil || resp.Message != "created" {
		t.Error(resp, err)
	}
}

func TestUploadPresets_List(t *testing.T) {
	resp, err := adminAPI.ListUploadPresets(ctx, admin.ListUploadPresetsParams{})

	if err != nil || len(resp.Presets) < 1 {
		t.Error(resp, err)
	}
}

func TestUploadPresets_Get(t *testing.T) {
	resp, err := adminAPI.GetUploadPreset(ctx, admin.GetUploadPresetParams{Name: UPName})

	if err != nil {
		t.Error(resp, err)
	}
}

func TestUploadPresets_Update(t *testing.T) {
	updateUPParams := admin.UpdateUploadPresetParams{
		Name:         UPName,
		Unsigned:     false,
		Live:         false,
		UploadParams: uploader.UploadParams{Tags: api.CldAPIArray{"go-tag3", "go-tag4"}},
	}

	resp, err := adminAPI.UpdateUploadPreset(ctx, updateUPParams)

	if err != nil || resp.Message != "updated" {
		t.Error(resp, err)
	}

	gResp, err := adminAPI.GetUploadPreset(ctx, admin.GetUploadPresetParams{Name: UPName})

	if err != nil {
		t.Error(gResp, err)
	}
}

func TestUploadPresets_Delete(t *testing.T) {
	resp, err := adminAPI.DeleteUploadPreset(ctx, admin.DeleteUploadPresetParams{Name: UPName})

	if err != nil || resp.Message != "deleted" {
		t.Error(resp, err)
	}
}
