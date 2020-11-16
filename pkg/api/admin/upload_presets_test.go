package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"testing"
)

const UPName = "go-upload-preset"

func TestUploadPresets_Create(t *testing.T) {

	params := CreateUploadPresetParams{
		Name:     UPName,
		Unsigned: true,
		Live:     true,
	}

	params.Tags = api.CldApiArray{"go-tag1", "go-tag2"}

	resp, err := adminApi.CreateUploadPreset(ctx, params)

	if err != nil || resp.Message != "created" {
		t.Error(resp, err)
	}
}

func TestUploadPresets_List(t *testing.T) {
	resp, err := adminApi.ListUploadPresets(ctx, ListUploadPresetsParams{})

	if err != nil || len(resp.Presets) < 1 {
		t.Error(resp, err)
	}
}

func TestUploadPresets_Get(t *testing.T) {
	resp, err := adminApi.GetUploadPreset(ctx, GetUploadPresetParams{Name: UPName})

	if err != nil {
		t.Error(resp, err)
	}
}

func TestUploadPresets_Update(t *testing.T) {
	updateUPParams := UpdateUploadPresetParams{
		Name:     UPName,
		Unsigned: false,
		Live:     false,
	}
	updateUPParams.Tags = api.CldApiArray{"go-tag3", "go-tag4"}

	resp, err := adminApi.UpdateUploadPreset(ctx, updateUPParams)

	if err != nil || resp.Message != "updated" {
		t.Error(resp, err)
	}

	gResp, err := adminApi.GetUploadPreset(ctx, GetUploadPresetParams{Name: UPName})

	if err != nil {
		t.Error(gResp, err)
	}
}

func TestUploadPresets_Delete(t *testing.T) {
	resp, err := adminApi.DeleteUploadPreset(ctx, DeleteUploadPresetParams{Name: UPName})

	if err != nil || resp.Message != "deleted" {
		t.Error(resp, err)
	}
}
