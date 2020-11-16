package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/api/uploader"
	"context"
)

const (
	UploadPresets api.EndPoint = "upload_presets"
)

type ListUploadPresetsParams struct {
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func (a *Api) ListUploadPresets(ctx context.Context, params ListUploadPresetsParams) (*ListUploadPresetsResult, error) {
	res := &ListUploadPresetsResult{}
	_, err := a.get(ctx, UploadPresets, params, res)

	return res, err
}

type ListUploadPresetsResult struct {
	Presets []UploadPreset `json:"presets"`
	Error   api.ErrorResp  `json:"error,omitempty"`
}

type UploadPreset struct {
	Name     string      `json:"name"`
	Unsigned bool        `json:"unsigned"`
	Settings interface{} `json:"settings"`
}

type GetUploadPresetParams struct {
	Name       string `json:"-"`
	MaxResults int    `json:"max_results,omitempty"`
}

func (a *Api) GetUploadPreset(ctx context.Context, params GetUploadPresetParams) (*GetUploadPresetResult, error) {
	res := &GetUploadPresetResult{}
	_, err := a.get(ctx, api.BuildPath(UploadPresets, params.Name), params, res)

	return res, err
}

type GetUploadPresetResult struct {
	Name     string        `json:"name"`
	Unsigned bool          `json:"unsigned"`
	Settings interface{}   `json:"settings"`
	Error    api.ErrorResp `json:"error,omitempty"`
}

type CreateUploadPresetParams struct {
	Name             string `json:"name,omitempty"`
	Unsigned         bool   `json:"unsigned,omitempty"`
	DisallowPublicId bool   `json:"disallow_public_id,omitempty"`
	Live             bool   `json:"live,omitempty"`
	uploader.UploadParams
}

func (a *Api) CreateUploadPreset(ctx context.Context, params CreateUploadPresetParams) (*CreateUploadPresetResult, error) {
	res := &CreateUploadPresetResult{}
	_, err := a.post(ctx, api.BuildPath(UploadPresets), params, res)

	return res, err
}

type CreateUploadPresetResult struct {
	Message string        `json:"message"`
	Name    string        `json:"name"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type UpdateUploadPresetParams struct {
	Name             string `json:"name,omitempty"`
	Unsigned         bool   `json:"unsigned,omitempty"`
	DisallowPublicId bool   `json:"disallow_public_id,omitempty"`
	Live             bool   `json:"live,omitempty"`
	uploader.UploadParams
}

func (a *Api) UpdateUploadPreset(ctx context.Context, params UpdateUploadPresetParams) (*UploadPresetResult, error) {
	res := &UploadPresetResult{}
	_, err := a.put(ctx, api.BuildPath(UploadPresets, params.Name), params, res)

	return res, err
}

type UploadPresetResult struct {
	Message string        `json:"message"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type DeleteUploadPresetParams struct {
	Name string `json:"-"`
}

func (a *Api) DeleteUploadPreset(ctx context.Context, params DeleteUploadPresetParams) (*UploadPresetResult, error) {
	res := &UploadPresetResult{}
	_, err := a.delete(ctx, api.BuildPath(UploadPresets, params.Name), params, res)

	return res, err
}
