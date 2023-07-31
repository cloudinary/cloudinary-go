package admin_test

import (
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
)

func TestSearchFolders_SearchQuery(t *testing.T) {
	cfResp, err := adminAPI.CreateFolder(ctx, admin.CreateFolderParams{Folder: testFolder})

	if err != nil || cfResp.Success != true {
		t.Error(cfResp, err)
	}

	sq := search.Query{
		Expression: "path=" + testFolder,
		MaxResults: 1,
	}

	resp, err := adminAPI.SearchFolders(ctx, sq)

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}
}
