package admin_test

import (
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

const testFolder = "000-go-folder"

func TestFolders_CreateFolder(t *testing.T) {
	resp, err := adminApi.CreateFolder(ctx, admin.CreateFolderParams{Folder: testFolder})

	if err != nil || resp.Success != true {
		t.Error(resp, err)
	}
}

func TestFolders_RootFolders(t *testing.T) {
	resp, err := adminApi.RootFolders(ctx, admin.RootFoldersParams{MaxResults: 5})

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}
}

func TestFolders_DeleteFolder(t *testing.T) {
	resp, err := adminApi.DeleteFolder(ctx, admin.DeleteFolderParams{Folder: testFolder})

	if err != nil || len(resp.Deleted) < 1 {
		t.Error(resp, err)
	}
}

func TestFolders_SubFolders(t *testing.T) {
	cfResp, err := adminApi.CreateFolder(ctx, admin.CreateFolderParams{Folder: testFolder})
	if err != nil || cfResp.Success != true {
		t.Error(cfResp, err)
	}

	cfResp, err = adminApi.CreateFolder(ctx, admin.CreateFolderParams{Folder: testFolder + "/" + testFolder})
	if err != nil || cfResp.Success != true {
		t.Error(cfResp, err)
	}

	time.Sleep(1 * time.Second)

	resp, err := adminApi.RootFolders(ctx, admin.RootFoldersParams{MaxResults: 1})
	if err != nil || resp == nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}

	resp, err = adminApi.SubFolders(ctx, admin.SubFoldersParams{Folder: resp.Folders[0].Path, MaxResults: 2})
	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}
}
