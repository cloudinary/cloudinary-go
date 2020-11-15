package admin

import (
	"testing"
)

const UMFolder = "wiki"
const UMTemplate = "https://u.wiki.example.com/wiki-images/"
const UMTemplateUpdated = "https://images.example.com/product_assets/images/"

func TestUploadMappings_Create(t *testing.T) {
	resp, err := adminApi.CreateUploadMapping(ctx, CreateUploadMappingParams{
		Folder:   UMFolder,
		Template: UMTemplate,
	})

	if err != nil || resp.Message != "created" {
		t.Error(resp, err)
	}
}

func TestUploadMappings_List(t *testing.T) {
	resp, err := adminApi.ListUploadMappings(ctx, ListUploadMappingsParams{})

	if err != nil || len(resp.Mappings) < 1 {
		t.Error(resp, err)
	}
}

func TestUploadMappings_Get(t *testing.T) {
	lResp, err := adminApi.ListUploadMappings(ctx, ListUploadMappingsParams{})

	if err != nil || lResp.Error.Message != "" {
		t.Error(lResp, err)
	}

	resp, err := adminApi.GetUploadMapping(ctx, GetUploadMappingParams{Folder: lResp.Mappings[0].Folder})

	if err != nil {
		t.Error(resp, err)
	}
}

func TestUploadMappings_Update(t *testing.T) {
	resp, err := adminApi.UpdateUploadMapping(ctx, UpdateUploadMappingParams{
		Folder:   UMFolder,
		Template: UMTemplateUpdated,
	})

	if err != nil || resp.Message != "updated" {
		t.Error(resp, err)
	}
}

func TestUploadMappings_Delete(t *testing.T) {
	resp, err := adminApi.DeleteUploadMapping(ctx, DeleteUploadMappingParams{Folder: UMFolder})

	if err != nil || resp.Message != "deleted" {
		t.Error(resp, err)
	}
}
