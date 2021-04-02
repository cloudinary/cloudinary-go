package admin

// Enables you to manage stored transformations.
//
// https://cloudinary.com/documentation/admin_api#transformations

import (
	"context"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/transformation"
)

const (
	transformations api.EndPoint = "transformations"
)

type ListTransformationsParams struct {
	Named      bool   `json:"named,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// ListTransformations lists stored transformations.
//
// https://cloudinary.com/documentation/admin_api#get_transformations
func (a *Api) ListTransformations(ctx context.Context, params ListTransformationsParams) (*ListTransformationsResult, error) {
	res := &ListTransformationsResult{}
	_, err := a.get(ctx, transformations, params, res)

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
	Transformation transformation.RawTransformation `json:"transformation"` // The transformation string.
	MaxResults     int                              `json:"max_results,omitempty"`
	NextCursor     string                           `json:"next_cursor,omitempty"`
}

// GetTransformation returns the details of a single transformation.
//
// https://cloudinary.com/documentation/admin_api#get_transformation_details
func (a *Api) GetTransformation(ctx context.Context, params GetTransformationParams) (*GetTransformationResult, error) {
	res := &GetTransformationResult{}
	_, err := a.get(ctx, api.BuildPath(transformations), params, res)

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

// CreateTransformation creates a named transformation.
//
// https://cloudinary.com/documentation/admin_api#create_named_transformation
func (a *Api) CreateTransformation(ctx context.Context, params CreateTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.post(ctx, api.BuildPath(transformations), params, res)

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

// UpdateTransformation updates the specified transformation.
//
// https://cloudinary.com/documentation/admin_api#update_transformation
func (a *Api) UpdateTransformation(ctx context.Context, params UpdateTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.put(ctx, api.BuildPath(transformations), params, res)

	return res, err
}

type DeleteTransformationParams struct {
	Transformation transformation.RawTransformation `json:"transformation"`
}

// DeleteTransformation deletes the specified stored transformation.
//
// Deleting a transformation also deletes all the stored derived resources based on this transformation (up to 1000).
// The method returns an error if there are more than 1000 derived resources based on this transformation.
//
// https://cloudinary.com/documentation/admin_api#delete_transformation
func (a *Api) DeleteTransformation(ctx context.Context, params DeleteTransformationParams) (*TransformationResult, error) {
	res := &TransformationResult{}
	_, err := a.delete(ctx, api.BuildPath(transformations), params, res)

	return res, err
}
