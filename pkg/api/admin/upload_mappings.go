package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
)

const (
	UploadMappings api.EndPoint = "upload_mappings"
)

type ListUploadMappingsParams struct {
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func (a *Api) ListUploadMappings(ctx context.Context, params ListUploadMappingsParams) (*ListUploadMappingsResult, error) {
	res := &ListUploadMappingsResult{}
	_, err := a.get(ctx, UploadMappings, nil, res)

	return res, err
}

type ListUploadMappingsResult struct {
	Mappings []UploadMapping `json:"mappings"`
	Error    api.ErrorResp   `json:"error,omitempty"`
}

type UploadMapping struct {
	Folder   string `json:"folder"`
	Template string `json:"template"`
}

type GetUploadMappingParams struct {
	Folder string `json:"folder"`
}

func (a *Api) GetUploadMapping(ctx context.Context, params GetUploadMappingParams) (*GetUploadMappingResult, error) {
	res := &GetUploadMappingResult{}
	_, err := a.get(ctx, api.BuildPath(UploadMappings), params, res)

	return res, err
}

type GetUploadMappingResult struct {
	Folder   string        `json:"folder"`
	Template string        `json:"template"`
	Error    api.ErrorResp `json:"error,omitempty"`
}

type CreateUploadMappingParams struct {
	Folder   string `json:"folder"`
	Template string `json:"template"`
}

func (a *Api) CreateUploadMapping(ctx context.Context, params CreateUploadMappingParams) (*CreateUploadMappingResult, error) {
	res := &CreateUploadMappingResult{}
	_, err := a.post(ctx, api.BuildPath(UploadMappings), params, res)

	return res, err
}

type CreateUploadMappingResult struct {
	Message string        `json:"message"`
	Folder  string        `json:"folder"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type UploadMappingResult struct {
	Message string        `json:"message"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type UpdateUploadMappingParams struct {
	Folder   string `json:"folder"`
	Template string `json:"template"`
}

func (a *Api) UpdateUploadMapping(ctx context.Context, params UpdateUploadMappingParams) (*UploadMappingResult, error) {
	res := &UploadMappingResult{}
	_, err := a.put(ctx, api.BuildPath(UploadMappings), params, res)

	return res, err
}

type DeleteUploadMappingParams struct {
	Folder string `json:"folder"`
}

func (a *Api) DeleteUploadMapping(ctx context.Context, params DeleteUploadMappingParams) (*UploadMappingResult, error) {
	res := &UploadMappingResult{}
	_, err := a.delete(ctx, api.BuildPath(UploadMappings), params, res)

	return res, err
}
