package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
)

const (
	Folders api.EndPoint = "folders"
)

type RootFoldersParams struct {
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func (a *Api) RootFolders(ctx context.Context, params RootFoldersParams) (*FoldersResult, error) {
	res := &FoldersResult{}
	_, err := a.get(ctx, Folders, params, res)

	return res, err
}

type FoldersResult struct {
	Folders    []FolderResult `json:"folders"`
	TotalCount int            `json:"total_count"`
	NextCursor string         `json:"next_cursor"`
	Error      api.ErrorResp  `json:"error,omitempty"`
}

type FolderResult struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type SubFoldersParams struct {
	Folder     string `json:"-"`
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

func (a *Api) SubFolders(ctx context.Context, params SubFoldersParams) (*FoldersResult, error) {
	res := &FoldersResult{}
	_, err := a.get(ctx, api.BuildPath(Folders, params.Folder), params, res)

	return res, err
}

type CreateFolderParams struct {
	Folder string `json:"-"`
}

func (a *Api) CreateFolder(ctx context.Context, params CreateFolderParams) (*CreateFolderResult, error) {
	res := &CreateFolderResult{}
	_, err := a.post(ctx, api.BuildPath(Folders, params.Folder), params, res)

	return res, err
}

type CreateFolderResult struct {
	Success bool          `json:"success"`
	Path    string        `json:"path"`
	Name    string        `json:"name"`
	Error   api.ErrorResp `json:"error,omitempty"`
}

type DeleteFolderParams struct {
	Folder string `json:"-"`
}

func (a *Api) DeleteFolder(ctx context.Context, params DeleteFolderParams) (*DeleteFolderResult, error) {
	res := &DeleteFolderResult{}
	_, err := a.delete(ctx, api.BuildPath(Folders, params.Folder), params, res)

	return res, err
}

type DeleteFolderResult struct {
	Deleted []string      `json:"deleted"`
	Error   api.ErrorResp `json:"error,omitempty"`
}
