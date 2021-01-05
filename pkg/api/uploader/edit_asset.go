package uploader

import (
	"cloudinary/cloudinary-go/pkg/api"
	"context"
	"net/http"
)

const (
	Upload   api.EndPoint = "upload"
	Destroy  api.EndPoint = "destroy"
	Rename   api.EndPoint = "rename"
	Explicit api.EndPoint = "explicit"
	Metadata api.EndPoint = "metadata"
)

// DestroyParams struct
type DestroyParams struct {
	PublicID     string `json:"public_id,omitempty"`
	Type         string `json:"type,omitempty"`
	ResourceType string `json:"-"`
	Invalidate   bool   `json:"invalidate,omitempty"`
}

// Destroy immediately and permanently deletes a single asset from your Cloudinary account.
//
// Backed up assets are not deleted, and any assets and transformed assets already downloaded by visitors to your
// website might still be accessible through cached copies on the CDN. You can invalidate any cached copies on the
// CDN with the `Invalidate` parameter.
//
// https://cloudinary.com/documentation/image_upload_api_reference#destroy_method
func (u *Api) Destroy(ctx context.Context, params DestroyParams) (*DestroyResult, error) {
	res := &DestroyResult{}
	err := u.callUploadApi(ctx, Destroy, params, res)

	return res, err
}

type DestroyResult struct {
	Result   string        `json:"result"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response http.Response
}

// RenameParams struct
type RenameParams struct {
	FromPublicID string `json:"from_public_id,omitempty"`
	ToPublicID   string `json:"to_public_id,omitempty"`
	Type         string `json:"type,omitempty"`
	ToType       string `json:"to_type,omitempty"`
	ResourceType string `json:"-"`
	Overwrite    bool   `json:"overwrite,omitempty"`
	Invalidate   bool   `json:"invalidate,omitempty"`
}

// Rename renames the specified asset in your Cloudinary account.
//
// The existing URLs of renamed assets and their associated derived resources are no longer valid, although any
// assets and transformed assets already downloaded by visitors to your website might still be accessible through
// cached copies on the CDN. You can invalidate any cached copies on the CDN with the `invalidate` parameter.
//
// https://cloudinary.com/documentation/image_upload_api_reference#rename_method
func (u *Api) Rename(ctx context.Context, params RenameParams) (*RenameResult, error) {
	res := &RenameResult{}
	err := u.callUploadApi(ctx, Rename, params, res)

	return res, err
}

type RenameResult struct {
	api.BriefAssetResult
	Error interface{} `json:"error,omitempty"`
}

type ExplicitParams = UploadParams

// Explicit applies actions to already uploaded assets.
//
// https://cloudinary.com/documentation/image_upload_api_reference#explicit_method
func (u *Api) Explicit(ctx context.Context, params ExplicitParams) (*ExplicitResult, error) {
	res := &ExplicitResult{}
	err := u.callUploadApi(ctx, Explicit, params, res)

	return res, err
}

type ExplicitResult struct {
	UploadResult
}

// UpdateMetadataParams struct
type UpdateMetadataParams struct {
	PublicIDs api.CldApiArray `json:"public_ids"`
	Metadata  api.Metadata    `json:"metadata"`
	Type      string          `json:"type,omitempty"`
}

// UpdateMetadata populates metadata fields with the given values. Existing values will be overwritten.
//
// Any metadata-value pairs given are merged with any existing metadata-value pairs
// (an empty value for an existing metadata field clears the value).
//
// https://cloudinary.com/documentation/image_upload_api_reference#metadata_method
func (u *Api) UpdateMetadata(ctx context.Context, params RenameParams) (*UpdateMetadataResult, error) {
	res := &UpdateMetadataResult{}
	err := u.callUploadApi(ctx, Metadata, params, res)

	return res, err
}

type UpdateMetadataResult struct {
	PublicIds []string    `json:"public_ids"`
	Error     interface{} `json:"error,omitempty"`
}
