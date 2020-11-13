package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"net/http"
	"time"
)

const (
	Ping  api.EndPoint = "ping"
	Usage api.EndPoint = "usage"
)

func (a *Api) Ping(ctx context.Context) (*PingResult, error) {
	res := &PingResult{}
	_, err := a.get(ctx, Ping, nil, res)

	return res, err
}

type PingResult struct {
	Status string        `json:"status"`
	Error  api.ErrorResp `json:"error,omitempty"`
	Response http.Response
}

type UsageParams struct {
	Date time.Time `json:"-"`
}

func (a *Api) Usage(ctx context.Context, params UsageParams) (*UsageResult, error) {
	date := ""
	if !params.Date.IsZero() {
		date = params.Date.Format("02-01-2006")
	}
	res := &UsageResult{}
	_, err := a.get(ctx, api.BuildPath(Usage, date), params, res)

	return res, err
}

type UsageResult struct {
	Plan            string `json:"plan"`
	LastUpdated     string `json:"last_updated"`
	Transformations struct {
		Usage       int     `json:"usage"`
		Limit       int     `json:"limit"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"transformations"`
	Objects struct {
		Usage       int     `json:"usage"`
		Limit       int     `json:"limit"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"objects"`
	Bandwidth struct {
		Usage       int64   `json:"usage"`
		Limit       int64   `json:"limit"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"bandwidth"`
	Storage struct {
		Usage       int64   `json:"usage"`
		Limit       int64   `json:"limit"`
		UsedPercent float64 `json:"used_percent"`
	} `json:"storage"`
	Requests         int64 `json:"requests"`
	Resources        int   `json:"resources"`
	DerivedResources int   `json:"derived_resources"`
	MediaLimits      struct {
		ImageMaxSizeBytes int `json:"image_max_size_bytes"`
		VideoMaxSizeBytes int `json:"video_max_size_bytes"`
		RawMaxSizeBytes   int `json:"raw_max_size_bytes"`
		ImageMaxPx        int `json:"image_max_px"`
		AssetMaxTotalPx   int `json:"asset_max_total_px"`
	} `json:"media_limits"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

type TagsParams struct {
	AssetType  api.AssetType `json:"-"`
	NextCursor string        `json:"next_cursor,omitempty"`
	MaxResults int           `json:"max_results,omitempty"`
	Prefix     string        `json:"prefix,omitempty"`
}

func (a *Api) Tags(ctx context.Context, params TagsParams) (*TagsResult, error) {
	res := &TagsResult{}
	_, err := a.get(ctx, api.BuildPath(Tags, params.AssetType.ToString()), params, res)

	return res, err
}

type TagsResult struct {
	Tags       []string      `json:"tags"`
	NextCursor string        `json:"next_cursor"`
	Error      api.ErrorResp `json:"error,omitempty"`
}
