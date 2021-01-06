package uploader

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
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
	Tag          string          `json:"tag,omitempty"` // The name of the tag to add.
	PublicIDs    api.CldApiArray `json:"public_ids"`    // The public IDs of the assets to add the tag to.
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// AddTag adds a tag to the assets specified.
//
// https://cloudinary.com/documentation/image_upload_api_reference#tags_method
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
	Tag          string          `json:"tag,omitempty"` // The name of the tag to remove.
	PublicIDs    api.CldApiArray `json:"public_ids"`    // The public IDs of the assets to remove the tags from.
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveTag removes a tag from the assets specified.
//
// https://cloudinary.com/documentation/image_upload_api_reference#tags_method
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
	PublicIDs    api.CldApiArray `json:"public_ids"` // The public IDs of the assets to remove all tags from.
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// RemoveAllTags removes all tags from the assets specified.
//
// https://cloudinary.com/documentation/image_upload_api_reference#tags_method
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
	Tag          string          `json:"tag"`        // The new tag with which to replace the existing tags.
	PublicIDs    api.CldApiArray `json:"public_ids"` // The public IDs of the assets to replace the tags of.
	Type         string          `json:"type,omitempty"`
	ResourceType string          `json:"-"`
}

// Replaces all existing tags on the assets specified with the tag specified.
//
// https://cloudinary.com/documentation/image_upload_api_reference#tags_method
func (u *Api) ReplaceTag(ctx context.Context, params ReplaceTagParams) (*ReplaceTagResult, error) {
	res := &ReplaceTagResult{}
	err := u.callTagsApi(ctx, ReplaceTag, params, res)

	return res, err
}

type ReplaceTagResult struct {
	TagResult
}

type TagResult struct {
	PublicIds []string      `json:"public_ids"` // The public IDs of the assets that were affected.
	Error     api.ErrorResp `json:"error,omitempty"`
	Response  http.Response
}

// callTagsApi is an internal method that is used to call to the tags API.
func (u *Api) callTagsApi(ctx context.Context, command TagCommand, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	formParams.Add("command", fmt.Sprintf("%v", command))

	return u.callUploadApiWithParams(ctx, api.BuildPath(getAssetType(requestParams), Tags), formParams, result)
}
