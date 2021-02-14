package admin

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/transformation"
)

const TName = "go_transformation"
const TTransformation = "c_fill,h_500,w_500"
const TTransformationUpdated = "c_fill,h_501,w_501"

func TestTransformations_Create(t *testing.T) {
	resp, err := adminApi.CreateTransformation(ctx, CreateTransformationParams{
		Name:           TName,
		Transformation: TTransformation,
	})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
	}
}

func TestTransformations_List(t *testing.T) {
	resp, err := adminApi.ListTransformations(ctx, ListTransformationsParams{})

	if err != nil || len(resp.Transformations) < 1 {
		t.Error(resp, err)
	}
}

func TestTransformations_Get(t *testing.T) {
	lResp, err := adminApi.ListTransformations(ctx, ListTransformationsParams{})

	if err != nil || lResp.Error.Message != "" {
		t.Error(lResp, err)
	}

	resp, err := adminApi.GetTransformation(ctx, GetTransformationParams{Transformation: transformation.RawTransformation(lResp.Transformations[0].Name)})

	if err != nil || len(resp.Info) < 1 {
		t.Error(resp, err)
	}
}

func TestTransformations_Update(t *testing.T) {
	resp, err := adminApi.UpdateTransformation(ctx, UpdateTransformationParams{
		Transformation: TName,
		UnsafeUpdate:   TTransformationUpdated,
	})

	if err != nil || resp.Message != "updated" {
		t.Error(resp, err)
	}
}

func TestTransformations_Delete(t *testing.T) {
	resp, err := adminApi.DeleteTransformation(ctx, DeleteTransformationParams{Transformation: TName})

	if err != nil || resp.Message != "deleted" {
		t.Error(resp, err)
	}
}
