package admin

import (
	"testing"
)

func TestSearch_Search(t *testing.T) {
	s := Search{}

	s.Expression("format:jpg").
		WithField(TagsField).WithField(ContextField).WithField(ImageMetadataField).WithField(ImageAnalysisField).
		SortBy("created_at", Descending).
		MaxResults(2)

	resp, err := adminApi.Search(ctx, s.GetQuery())

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp)
	}
}

func TestSearch_SearchQuery(t *testing.T) {
	sq := SearchQuery{
		Expression: "format:jpg",
		WithField:  []WithField{TagsField, ContextField, ImageMetadataField, ImageAnalysisField},
		SortBy:     []SortByField{{"created_at": Descending}},
		MaxResults: 2,
	}

	resp, err := adminApi.Search(ctx, sq)

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp)
	}
}
