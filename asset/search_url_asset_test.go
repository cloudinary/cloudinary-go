package asset_test

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin/search"
	"github.com/cloudinary/cloudinary-go/v2/asset"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"github.com/stretchr/testify/assert"
	"testing"
)

var query = search.Query{
	Expression: "resource_type:image AND tags=kitten AND uploaded_at>1d AND bytes>1m",
	SortBy:     []search.SortByField{{"public_id": search.Descending}},
	MaxResults: 30,
}
var b64Query = "eyJleHByZXNzaW9uIjoicmVzb3VyY2VfdHlwZTppbWFnZSBBTkQgdGFncz1raXR0ZW4gQU5EIHVwbG9hZGVkX2F0" +
	"PjFkIEFORCBieXRlcz4xbSIsIm1heF9yZXN1bHRzIjozMCwic29ydF9ieSI6W3sicHVibGljX2lkIjoiZGVzYyJ9XX0="

var ttl300Sig = "431454b74cefa342e2f03e2d589b2e901babb8db6e6b149abf25bc0dd7ab20b7"
var ttl1000Sig = "25b91426a37d4f633a9b34383c63889ff8952e7ffecef29a17d600eeb3db0db7"

var c, _ = cloudinary.NewFromURL(cldtest.CldURL)

var searchEndpoint = fmt.Sprintf("%s://%s/%s/search", c.Config.URL.Protocol(), c.Config.URL.SharedHost, c.Config.Cloud.CloudName)

func getSearchURL(t *testing.T) *asset.SearchURLAsset {
	su, err := c.SearchURL(query)
	if err != nil {
		t.Fatal(err)
	}

	return su
}
func TestSearchURL_Default(t *testing.T) {
	su := getSearchURL(t)

	assert.Contains(t, getAssetUrl(t, su), fmt.Sprintf("%s/%s/%d/%s", searchEndpoint, ttl300Sig, 300, b64Query))
}

func TestSearchURL_WithNextCursor(t *testing.T) {
	su := getSearchURL(t)
	url, err := su.ToURLWithNextCursor(cldtest.NextCursor)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, url, fmt.Sprintf("%s/%s/%d/%s/%s", searchEndpoint, ttl300Sig, 300, b64Query, cldtest.NextCursor))
}

func TestSearchURL_WithCustomTTLAndNextCursor(t *testing.T) {
	su := getSearchURL(t)
	url, err := su.ToURL(1000, cldtest.NextCursor)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, url, fmt.Sprintf("%s/%s/%d/%s/%s", searchEndpoint, ttl1000Sig, 1000, b64Query, cldtest.NextCursor))
}

func TestSearchURL_TTLAndNextCursorSetOnStruct(t *testing.T) {
	su := getSearchURL(t)
	su.TTL = 1000
	su.SearchQuery.NextCursor = cldtest.NextCursor

	assert.Contains(t, getAssetUrl(t, su), fmt.Sprintf("%s/%s/%d/%s/%s", searchEndpoint, ttl1000Sig, 1000, b64Query, cldtest.NextCursor))
}

func TestSearchURL_PrivateCDN(t *testing.T) {
	su := getSearchURL(t)
	su.Config.URL.PrivateCDN = true

	privateSearchEndpoint := fmt.Sprintf("%s://%s-%s/search", c.Config.URL.Protocol(), c.Config.Cloud.CloudName, c.Config.URL.SharedHost)

	assert.Contains(t, getAssetUrl(t, su), fmt.Sprintf("%s/%s/%d/%s", privateSearchEndpoint, ttl300Sig, 300, b64Query))
}

func TestSearchURL_NoAnalytics(t *testing.T) {
	su := getSearchURL(t)

	assert.Contains(t, getAssetUrl(t, su), fmt.Sprintf("%s/%s/%d/%s", searchEndpoint, ttl300Sig, 300, b64Query))
	assert.Contains(t, getAssetUrl(t, su), "?_a=")

	su.Config.URL.Analytics = false

	assert.Contains(t, getAssetUrl(t, su), fmt.Sprintf("%s/%s/%d/%s", searchEndpoint, ttl300Sig, 300, b64Query))
	assert.NotContains(t, getAssetUrl(t, su), "?_a=")
}
