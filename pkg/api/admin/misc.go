package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"encoding/json"
	"net/url"
)

type PingResult struct {
	Status string        `json:"status"`
	Error  api.ErrorResp `json:"error,omitempty"`
}

func (a *Api) Ping() (*PingResult, error) {
	resp := a.callApi("ping", url.Values{})

	res := &PingResult{}
	err := json.Unmarshal(resp, res)

	if err != nil {
		return nil, err
	}

	return res, nil
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
	Error api.ErrorResp `json:"error,omitempty"`
}

func (a *Api) Usage() (*UsageResult, error) {
	resp := a.callApi("usage", url.Values{})

	res := &UsageResult{}
	err := json.Unmarshal(resp, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

type TagsParams struct {
	AssetType  string `json:"-"`
	NextCursor string `json:"next_cursor,omitempty"`
	MaxResults int    `json:"max_results,omitempty"`
	Prefix     bool   `json:"prefix,omitempty"`
}

type TagsResult struct {
	Tags       []string `json:"tags"`
	NextCursor string   `json:"next_cursor"`
}

func (a *Api) Tags(tagsParams TagsParams) (*TagsResult, error) {
	if tagsParams.AssetType == "" {
		tagsParams.AssetType = "image"
	}
	params, err := api.StructToParams(tagsParams)
	if err != nil {
		return nil, err
	}
	resp := a.callApi("tags/" + tagsParams.AssetType, params)

	res := &TagsResult{}
	err = json.Unmarshal(resp, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
