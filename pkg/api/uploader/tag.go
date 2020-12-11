package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"fmt"
	"net/http"
)

const (
	Tags api.EndPoint = "tags"
)

type TagCommand string

const (
	AddTag        TagCommand = "add"
	RemoveTag     TagCommand = "remove"
	ReplaceTag    TagCommand = "replace"
	RemoveAllTags TagCommand = "remove_all"
)

// AddTagParams struct
type AddTagParams struct {
	Tag          string          `json:"tag,omitempty"`
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// AddTag to the assets
func (u *Api) AddTag(ctx context.Context, params AddTagParams) (*AddTagResult, error) {
	res := &AddTagResult{}
	err := u.callTagsApi(ctx, AddTag, params, res)

	return res, err
}

type AddTagResult struct {
	TagResult
}

// RemoveTagParams struct
type RemoveTagParams struct {
	Tag          string          `json:"tag,omitempty"`
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveTag from the assets
func (u *Api) RemoveTag(ctx context.Context, params RemoveTagParams) (*RemoveTagResult, error) {
	res := &RemoveTagResult{}
	err := u.callTagsApi(ctx, RemoveTag, params, res)

	return res, err
}

type RemoveTagResult struct {
	TagResult
}

// RemoveTagParams struct
type RemoveAllTagsParams struct {
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveAllTags from the assets
func (u *Api) RemoveAllTags(ctx context.Context, params RemoveAllTagsParams) (*RemoveAllTagsResult, error) {
	res := &RemoveAllTagsResult{}
	err := u.callTagsApi(ctx, RemoveAllTags, params, res)

	return res, err
}

type RemoveAllTagsResult struct {
	TagResult
}

// ReplaceTagParams struct
type ReplaceTagParams struct {
	Tag          string          `json:"tag,omitempty"`
	PublicIDs    api.CldApiArray `json:"public_ids"`
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// Replaces all existing tags on the assets specified with the tag specified.
func (u *Api) ReplaceTag(ctx context.Context, params ReplaceTagParams) (*ReplaceTagResult, error) {
	res := &ReplaceTagResult{}
	err := u.callTagsApi(ctx, ReplaceTag, params, res)

	return res, err
}

type ReplaceTagResult struct {
	TagResult
}

type TagResult struct {
	PublicIds []string      `json:"public_ids"`
	Error     api.ErrorResp `json:"error,omitempty"`
	Response  http.Response
}

func (u *Api) callTagsApi(ctx context.Context, command TagCommand, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	formParams.Add("command", fmt.Sprintf("%v", command))

	return u.callUploadApiWithParams(ctx, api.BuildPath(getAssetType(requestParams), Tags), formParams, result)
}
