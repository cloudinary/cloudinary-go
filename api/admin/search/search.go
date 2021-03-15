package search

// Query struct includes the search query.
type Query struct {
	// Expression is the (Lucene-like) string expression specifying the search query.
	// If not provided then all resources are listed (up to MaxResults).
	Expression string `json:"expression,omitempty"`
	// SortBy is the the field to sort by. You can specify more than one SortBy parameter; results will be sorted
	// according to the order of the fields provided.
	SortBy []SortByField `json:"sort_by,omitempty"`
	// Aggregate is the the name of a field (attribute) for which an aggregation count should be calculated and returned in the response.
	// (Tier 2 only)
	// You can specify more than one aggregate parameter.
	// For aggregation fields without discrete values, the results are divided into categories.
	Aggregate []Aggregation `json:"aggregate,omitempty"`
	//WithField contains names of additional asset attributes to include for each asset in the response.
	WithField []WithField `json:"with_field,omitempty"`
	// MaxResults is the maximum number of results to return. Default 50. Maximum 500.
	MaxResults int `json:"max_results,omitempty"`
	// NextCursor value is returned as part of the response when a search request has more results to return than MaxResults.
	// You can then specify this value as the NextCursor parameter of the following request.
	NextCursor string `json:"next_cursor,omitempty"`
}

type Aggregation = string

const (
	AssetType    Aggregation = "resource_type"
	DeliveryType             = "type"
	Pixels                   = "pixels"   // Pixels only the image assets in the response are aggregated.
	Duration                 = "duration" // Duration only the video assets in the response are aggregated.
	Format                   = "format"
	Bytes                    = "bytes"
)

type WithField = string

const (
	ContextField       WithField = "context"
	TagsField                    = "tags"
	ImageMetadataField           = "image_metadata"
	ImageAnalysisField           = "image_analysis"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending           = "desc"
)

type SortByField map[string]Direction
