package admin

// Enables you to manage the folders in your account or cloud.
//
// https://cloudinary.com/documentation/admin_api#folders
import (
	"context"
	"github.com/cloudinary/cloudinary-go/api"
)

const (
	Folders api.EndPoint = "folders"
)

type RootFoldersParams struct {
	MaxResults int    `json:"max_results,omitempty"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// RootFolders lists all root folders.
//
// https://cloudinary.com/documentation/admin_api#get_root_folders
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

// SubFolders lists sub-folders.
//
// Returns the name and path of all the sub-folders of a specified parent folder. Limited to 2000 results.
//
// https://cloudinary.com/documentation/admin_api#get_subfolders
func (a *Api) SubFolders(ctx context.Context, params SubFoldersParams) (*FoldersResult, error) {
	res := &FoldersResult{}
	_, err := a.get(ctx, api.BuildPath(Folders, params.Folder), params, res)

	return res, err
}

type CreateFolderParams struct {
	Folder string `json:"-"` // The full path of the new folder to create.
}

// CreateFolder creates a new empty folder.
//
// https://cloudinary.com/documentation/admin_api#create_folder
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
	Folder string `json:"-"` // The full path of the empty folder to delete.
}

// DeleteFolder deletes an empty folder.
//
// The specified folder cannot contain any assets, but can have empty descendant sub-folders.
//
// https://cloudinary.com/documentation/admin_api#delete_folder
func (a *Api) DeleteFolder(ctx context.Context, params DeleteFolderParams) (*DeleteFolderResult, error) {
	res := &DeleteFolderResult{}
	_, err := a.delete(ctx, api.BuildPath(Folders, params.Folder), params, res)

	return res, err
}

type DeleteFolderResult struct {
	Deleted []string      `json:"deleted"`
	Error   api.ErrorResp `json:"error,omitempty"`
}
