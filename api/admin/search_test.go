package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin/search"
)

func TestSearch_SearchQuery(t *testing.T) {
	sq := search.Query{
		Expression: "format:png",
		WithField:  []string{search.TagsField, search.ContextField, search.ImageMetadataField, search.ImageAnalysisField},
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 2,
	}

	resp, err := adminAPI.Search(ctx, sq)

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}
}
