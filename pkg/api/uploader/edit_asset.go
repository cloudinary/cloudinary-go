package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"net/http"
)

const (
	Upload   api.EndPoint = "upload"
	Destroy  api.EndPoint = "destroy"
	Rename   api.EndPoint = "rename"
	Explicit api.EndPoint = "explicit"
)

// DestroyParams struct
type DestroyParams struct {
	PublicID     string `json:"public_id,omitempty"`
	Type         string `json:"type,omitempty"`
	ResourceType string `json:"-"`
	Invalidate   bool   `json:"invalidate,omitempty"`
}

// Destroy the asset
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

// DestroyParams struct
type RenameParams struct {
	FromPublicID string `json:"from_public_id,omitempty"`
	ToPublicID   string `json:"to_public_id,omitempty"`
	Type         string `json:"type,omitempty"`
	ToType       string `json:"to_type,omitempty"`
	ResourceType string `json:"-"`
	Overwrite    bool   `json:"overwrite,omitempty"`
	Invalidate   bool   `json:"invalidate,omitempty"`
}

// Destroy the asset
func (u *Api) Rename(ctx context.Context, params RenameParams) (*RenameResult, error) {
	res := &RenameResult{}
	err := u.callUploadApi(ctx, Rename, params, res)

	return res, err
}

type RenameResult struct {
	api.BriefAssetResult
	Error interface{} `json:"error,omitempty"`
}

type ExplicitParams struct {
	UploadParams
}

// Destroy the asset
func (u *Api) Explicit(ctx context.Context, params ExplicitParams) (*ExplicitResult, error) {
	res := &ExplicitResult{}
	err := u.callUploadApi(ctx, Explicit, params, res)

	return res, err
}

type ExplicitResult struct {
	UploadResult
}
