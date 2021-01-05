package admin

// The Cloudinary API search method allows you fine control on filtering and retrieving information on all the assets
// in your account with the help of query expressions in a Lucene-like query language. A few examples of what you can
// accomplish using the search method include:
//
//  * Searching by descriptive attributes such as public ID, filename, folders, tags, context, etc.
//  * Searching by file details such as type, format, file size, dimensions, etc.
//  * Searching by embedded data such as Exif, XMP, etc.
//  * Searching by analyzed data such as the number of faces, predominant colors, auto-tags, etc.
//  * Requesting aggregation counts on specified parameters, for example the number of assets found broken down by file
// format.
//
// https://cloudinary.com/documentation/search_api
import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"time"
)

const SearchEndPoint api.EndPoint = "resources/search"

// SearchQuery struct includes the search query.
type SearchQuery struct {
	// Expression is the (Lucene-like) string expression specifying the search query.
	// If not provided then all resources are listed (up to MaxResults).
	Expression string        `json:"expression,omitempty"`
	// SortBy is the the field to sort by. You can specify more than one SortBy parameter; results will be sorted
	// according to the order of the fields provided.
	SortBy     []SortByField `json:"sort_by,omitempty"`
	// Aggregate is the the name of a field (attribute) for which an aggregation count should be calculated and returned in the response.
	// (Tier 2 only)
	// You can specify more than one aggregate parameter.
	// For aggregation fields without discrete values, the results are divided into categories.
	Aggregate  []Aggregation `json:"aggregate,omitempty"`
	//WithField contains names of additional asset attributes to include for each asset in the response.
	WithField  []WithField   `json:"with_field,omitempty"`
	// MaxResults is the maximum number of results to return. Default 50. Maximum 500.
	MaxResults int           `json:"max_results,omitempty"`
	// NextCursor value is returned as part of the response when a search request has more results to return than MaxResults.
	// You can then specify this value as the NextCursor parameter of the following request.
	NextCursor string        `json:"next_cursor,omitempty"`
}

// Search executes the search API request.
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
	FaceCount         int                `json:"face_count"`
	Faces             [][]int            `json:"faces"`
	Grayscale         bool               `json:"grayscale"`
	IllustrationScore int                `json:"illustration_score"`
	Transparent       bool               `json:"transparent"`
	Etag              string             `json:"etag"`
	Colors            map[string]float64 `json:"colors"`
}

type SearchUser struct {
	AccessKey string `json:"access_key"`
}

type Search struct {
	query SearchQuery
}

// Expression sets the query string for filtering the assets in your account.
func (s *Search) Expression(expression string) *Search {
	s.query.Expression = expression
	return s
}

// MaxResults sets the maximum number of results to return.
func (s *Search) MaxResults(maxResults int) *Search {
	s.query.MaxResults = maxResults
	return s
}

// NextCursor sets the next cursor.
func (s *Search) NextCursor(nextCursor string) *Search {
	s.query.NextCursor = nextCursor
	return s
}

// SortBy sets the field to sort by.
func (s *Search) SortBy(fieldName string, direction Direction) *Search {
	s.query.SortBy = append(s.query.SortBy, SortByField{fieldName: direction})
	return s
}

// Aggregate sets the name of a field (attribute) for which an aggregation count should be calculated and returned in the response.
func (s *Search) Aggregate(aggregation Aggregation) *Search {
	s.query.Aggregate = append(s.query.Aggregate, aggregation)
	return s
}

// WithField sets the name of an additional asset attribute to include for each asset in the response.
func (s *Search) WithField(field WithField) *Search {
	s.query.WithField = append(s.query.WithField, field)
	return s
}

// GetQuery returns the SearchQuery
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
	// Pixels only the image assets in the response are aggregated.
	Pixels                   = "pixels"
	// Duration only the video assets in the response are aggregated.
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
