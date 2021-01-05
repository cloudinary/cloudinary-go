package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"fmt"
	"net/http"
)

const (
	Context api.EndPoint = "context"
)

type ContextCommand string

const (
	AddContext       ContextCommand = "add"
	RemoveAllContext ContextCommand = "remove_all"
)

// AddContextParams struct
type AddContextParams struct {
	Context      api.CldApiMap   `json:"context,omitempty"`
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// AddContext adds context metadata as key-value pairs to the the specified assets.
//
// https://cloudinary.com/documentation/image_upload_api_reference#context_method
func (u *Api) AddContext(ctx context.Context, params AddContextParams) (*AddContextResult, error) {
	res := &AddContextResult{}
	err := u.callContextApi(ctx, AddContext, params, res)

	return res, err
}

type AddContextResult struct {
	ContextResult
}

type ContextResult struct {
	PublicIds []string      `json:"public_ids"`
	Error     api.ErrorResp `json:"error,omitempty"`
	Response  http.Response
}

// RemoveAllContextParams struct
type RemoveAllContextParams struct {
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveAllContext removes all context metadata from the specified assets.
//
// https://cloudinary.com/documentation/image_upload_api_reference#context_method
func (u *Api) RemoveAllContext(ctx context.Context, params RemoveAllContextParams) (*RemoveAllContextResult, error) {
	res := &RemoveAllContextResult{}
	err := u.callContextApi(ctx, RemoveAllContext, params, res)

	return res, err
}

type RemoveAllContextResult struct {
	ContextResult
}

// callContextApi is an internal method that is used to call to the context API.
func (u *Api) callContextApi(ctx context.Context, command ContextCommand, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	formParams.Add("command", fmt.Sprintf("%v", command))

	return u.callUploadApiWithParams(ctx, api.BuildPath(getAssetType(requestParams), Context), formParams, result)
}
