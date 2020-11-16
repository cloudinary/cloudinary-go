package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/transformation"
	"context"
)

const (
	Transformations api.EndPoint = "transformations"
)

type ListTransformationsParams struct {
	Named      bool   `json:"named,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func (a *Api) ListTransformations(ctx context.Context, params ListTransformationsParams) (*ListTransformationsResult, error) {
	res := &ListTransformationsResult{}
	_, err := a.get(ctx, Transformations, params, res)

	return res, err
}

type ListTransformationsResult struct {
	Transformations []TransformationListItem `json:"transformations"`
	Error           api.ErrorResp            `json:"error,omitempty"`
}

type TransformationListItem struct {
	Name             string `json:"name"`
	AllowedForStrict bool   `json:"allowed_for_strict"`
	Used             bool   `json:"used"`
	Named            bool   `json:"named"`
}

type GetTransformationParams struct {
	Transformation transformation.RawTransformation `json:"transformation"`
	MaxResults     int                              `json:"max_results,omitempty"`
	NextCursor     string                           `json:"next_cursor,omitempty"`
}

func (a *Api) GetTransformation(ctx context.Context, params GetTransformationParams) (*GetTransformationResult, error) {
	res := &GetTransformationResult{}
	_, err := a.get(ctx, api.BuildPath(Transformations), params, res)

	return res, err
}

type GetTransformationResult struct {
	Name             string                        `json:"name"`
	AllowedForStrict bool                          `json:"allowed_for_strict"`
	Used             bool                          `json:"used"`
	Named            bool                          `json:"named"`
	Info             transformation.Transformation `json:"info"`
	Derived          []DerivedAsset                `json:"derived"`
	Error            api.ErrorResp                 `json:"error,omitempty"`
}

type DerivedAsset struct {
	PublicID     string `json:"public_id"`
	ResourceType string `json:"resource_type"`
	Type         string `json:"type"`
	Format       string `json:"format"`
	URL          string `json:"url"`
	SecureURL    string `json:"secure_url"`
	Bytes        int    `json:"bytes"`
	ID           string `json:"id"`
}

type CreateTransformationParams struct {
	Name           string                           `json:"name"`
	Transformation transformation.RawTransformation `json:"transformation"`
}

func (a *Api) CreateTransformation(ctx context.Context, params CreateTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.post(ctx, api.BuildPath(Transformations), params, res)

	return res, err
}

type TransformationResult struct {
	Message string        `json:"message"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type UpdateTransformationParams struct {
	Transformation   transformation.RawTransformation `json:"transformation"`
	AllowedForStrict bool                             `json:"allowed_for_strict,omitempty"`
	UnsafeUpdate     transformation.RawTransformation `json:"unsafe_update,omitempty"`
}

func (a *Api) UpdateTransformation(ctx context.Context, params UpdateTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.put(ctx, api.BuildPath(Transformations), params, res)

	return res, err
}

type DeleteTransformationParams struct {
	Transformation transformation.RawTransformation `json:"transformation"`
}

func (a *Api) DeleteTransformation(ctx context.Context, params DeleteTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.delete(ctx, api.BuildPath(Transformations), params, res)

	return res, err
}
