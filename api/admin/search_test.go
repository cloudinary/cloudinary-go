package admin

import (
	"github.com/cloudinary/cloudinary-go/api/admin/search"
	"testing"
)

func TestSearch_SearchQuery(t *testing.T) {
	sq := search.Query{
		Expression: "format:jpg",
		WithField:  []string{search.TagsField, search.ContextField, search.ImageMetadataField, search.ImageAnalysisField},
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		MaxResults: 2,
	}

	resp, err := adminApi.Search(ctx, sq)

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}
}
