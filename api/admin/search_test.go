package admin_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
)

func TestSearch_SearchQuery(t *testing.T) {
	sq := search.Query{
		Expression: "format:png",
		WithField:  []string{search.TagsField, search.ContextField, search.ImageMetadataField, search.ImageAnalysisField},
		SortBy:     []search.SortByField{{"created_at": search.Descending}},
		Fields:     []string{"tags", "secure_url"},
		MaxResults: 2,
	}

	resp, err := adminAPI.Search(ctx, sq)

	if err != nil || resp.TotalCount < 1 {
		t.Error(resp, err)
	}

	assert.NotEmpty(t, resp.Assets[0].AssetID)
	assert.NotEmpty(t, resp.Assets[0].SecureURL)
	assert.Empty(t, resp.Assets[0].URL)
}
