package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"time"
)

const SearchEndPoint api.EndPoint = "resources/search"

type SearchQuery struct {
	Expression string        `json:"expression,omitempty"`
	SortBy     []SortByField `json:"sort_by,omitempty"`
	Aggregate  []Aggregation `json:"aggregate,omitempty"`
	WithField  []WithField   `json:"with_field,omitempty"`
	MaxResults int           `json:"max_results,omitempty"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

func (a *Api) Search(ctx context.Context, searchQuery SearchQuery) (*SearchResult, error) {
	res := &SearchResult{}
	_, err := a.post(ctx, api.BuildPath(SearchEndPoint), searchQuery, res)

	return res, err
}

type SearchResult struct {
	TotalCount int           `json:"total_count"`
	Time       int           `json:"time"`
	Assets     []SearchAsset `json:"resources"`
	Error      api.ErrorResp `json:"error,omitempty"`
}

type SearchAsset struct {
	PublicID      string              `json:"public_id"`
	Folder        string              `json:"folder"`
	Filename      string              `json:"filename"`
	Format        string              `json:"format"`
	Version       int                 `json:"version"`
	ResourceType  string              `json:"resource_type"`
	Type          string              `json:"type"`
	CreatedAt     time.Time           `json:"created_at"`
	UploadedAt    time.Time           `json:"uploaded_at"`
	Bytes         int                 `json:"bytes"`
	BackupBytes   int                 `json:"backup_bytes"`
	Width         int                 `json:"width"`
	Height        int                 `json:"height"`
	AspectRatio   float64             `json:"aspect_ratio"`
	Pixels        int                 `json:"pixels"`
	Tags          []string            `json:"tags"`
	ImageMetadata ImageMetadataResult `json:"image_metadata"`
	ImageAnalysis ImageAnalysis       `json:"image_analysis"`
	URL           string              `json:"url"`
	SecureURL     string              `json:"secure_url"`
	Status        string              `json:"status"`
	AccessMode    string              `json:"access_mode"`
	AccessControl interface{}         `json:"access_control"`
	Etag          string              `json:"etag"`
	CreatedBy     SearchUser          `json:"created_by"`
	UploadedBy    SearchUser          `json:"uploaded_by"`
}

type ImageAnalysis struct {
	FaceCount         int               `json:"face_count"`
	Faces             [][]int           `json:"faces"`
	Grayscale         bool              `json:"grayscale"`
	IllustrationScore int               `json:"illustration_score"`
	Transparent       bool              `json:"transparent"`
	Etag              string            `json:"etag"`
	Colors            map[string]string `json:"colors"`
}

type SearchUser struct {
	AccessKey string `json:"access_key"`
}

type Search struct {
	query SearchQuery
}

func (s *Search) Expression(expression string) *Search {
	s.query.Expression = expression
	return s
}

func (s *Search) MaxResults(maxResults int) *Search {
	s.query.MaxResults = maxResults
	return s
}

func (s *Search) NextCursor(nextCursor string) *Search {
	s.query.NextCursor = nextCursor
	return s
}

func (s *Search) SortBy(fieldName string, direction Direction) *Search {
	s.query.SortBy = append(s.query.SortBy, SortByField{fieldName: direction})
	return s
}

func (s *Search) Aggregate(aggregation Aggregation) *Search {
	s.query.Aggregate = append(s.query.Aggregate, aggregation)
	return s
}

func (s *Search) WithField(field WithField) *Search {
	s.query.WithField = append(s.query.WithField, field)
	return s
}

func (s *Search) GetQuery() SearchQuery {
	return s.query
}

type Direction string

const (
	Ascending  Direction = "asc"
	Descending           = "desc"
)

type Aggregation string

const (
	AssetType    Aggregation = "resource_type"
	DeliveryType             = "type"
	Pixels                   = "pixels"
	Duration                 = "duration"
	Format                   = "format"
	Bytes                    = "bytes"
)

type WithField string

const (
	ContextField       WithField = "context"
	TagsField                    = "tags"
	ImageMetadataField           = "image_metadata"
	ImageAnalysisField           = "image_analysis"
)

type SortByField map[string]Direction
