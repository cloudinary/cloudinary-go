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

// AddTagParams struct
type AddContextParams struct {
	Context      api.CldApiMap   `json:"context,omitempty"`
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// Adds context to the assets
func (u *Api) AddContext(ctx context.Context, params AddContextParams) (*AddContextResult, error) {
	res := &AddContextResult{}
	err := u.callContextApi(ctx, AddContext, params, res)

	return res, err
}

type AddContextResult struct {
	ContextResult
}

// RemoveAllContextParams struct
type RemoveAllContextParams struct {
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveAllTags from the assets
func (u *Api) RemoveAllContext(ctx context.Context, params RemoveAllContextParams) (*RemoveAllContextResult, error) {
	res := &RemoveAllContextResult{}
	err := u.callContextApi(ctx, RemoveAllContext, params, res)

	return res, err
}

type RemoveAllContextResult struct {
	ContextResult
}

type ContextResult struct {
	PublicIds []string      `json:"public_ids"`
	Error     api.ErrorResp `json:"error,omitempty"`
	Response  http.Response
}

func (u *Api) callContextApi(ctx context.Context, command ContextCommand, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	formParams.Add("command", fmt.Sprintf("%v", command))

	return u.callUploadApiWithParams(ctx, api.BuildPath(getAssetType(requestParams), Context), formParams, result)
}
