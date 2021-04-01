package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

const SPName = "00-go-sp"

func TestStreamingProfiles_List(t *testing.T) {
	resp, err := adminApi.ListStreamingProfiles(ctx)

	if err != nil || len(resp.Data) < 1 {
		t.Error(resp, err)
	}
}

func TestStreamingProfiles_Get(t *testing.T) {
	lResp, err := adminApi.ListStreamingProfiles(ctx)

	if err != nil || lResp.Error.Message != "" {
		t.Error(lResp, err)
	}

	resp, err := adminApi.GetStreamingProfile(ctx, admin.GetStreamingProfileParams{Name: lResp.Data[0].Name})

	if err != nil {
		t.Error(resp, err)
	}
}

func TestStreamingProfiles_Create(t *testing.T) {
	resp, err := adminApi.CreateStreamingProfile(ctx, admin.CreateStreamingProfileParams{
		Name:            SPName,
		DisplayName:     "Go SP",
		Representations: admin.StreamingProfileRepresentations{{Transformation: "c_fill,w_1000,h_1000"}},
	})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
	}
}

func TestStreamingProfiles_Update(t *testing.T) {
	resp, err := adminApi.UpdateStreamingProfile(ctx, admin.UpdateStreamingProfileParams{
		Name:            SPName,
		DisplayName:     "Go SP Updated",
		Representations: admin.StreamingProfileRepresentations{{"c_fill,w_1001,h_1001"}},
	})

	if err != nil || resp.Data.DisplayName != "Go SP Updated" {
		t.Error(resp, err)
	}
}

func TestStreamingProfiles_Delete(t *testing.T) {
	resp, err := adminApi.DeleteStreamingProfile(ctx, admin.DeleteStreamingProfileParams{Name: SPName})

	if err != nil || resp.Message != "deleted" {
		t.Error(resp, err)
	}
}
