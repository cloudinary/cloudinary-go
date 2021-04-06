// Package api contains packages used for accessing Cloudinary API functionality.
//
// https://cloudinary.com/documentation/cloudinary_references
package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// EndPoint represents the API endpoint.
type EndPoint = string

// Version is the Cloudinary Go package version.
const Version = "0.2.0"

// UserAgent contains information about the SDK user agent. Passed to the Cloudinary servers.
const UserAgent = "CloudinaryGo/" + Version

// apiVersion is the current Cloudinary API version.
var apiVersion = "1_1"

// BaseURL is the base API url.
func BaseURL(uploadPrefix string) string {
	return fmt.Sprintf("%s/v%s", uploadPrefix, apiVersion)
}

// base64DataRegex is the regular expression for detecting base64 encoded strings.
var base64DataRegex = regexp.MustCompile("^data:([\\w-]+/[\\w\\-+.]+)?(;[\\w-]+=[\\w-]+)*;base64,([a-zA-Z0-9/+\\n=]+)$")

// AssetType is the type of the asset.
type AssetType string

// String serializes AssetType to string.
func (a AssetType) String() string {
	if a == "" {
		return string(Image)
	}

	return string(a)
}

const (
	// Image is the image asset type.
	Image AssetType = "image"
	// Video is the video asset type.
	Video = "video"
	// File is the raw asset type.
	File = "raw"
	// Auto is the auto asset type. Tells Cloudinary to automatically detect the type of the uploaded asset.
	Auto = "auto"
	// All is the all asset type. Used for downloading folders with all assets inside.
	All = "all"
)

// DeliveryType is the delivery type of the asset.
type DeliveryType string

// String serializes DeliveryType to string.
func (d DeliveryType) String() string {
	if d == "" {
		return string(Upload)
	}

	return string(d)
}

const (
	// Upload is the upload delivery type.
	Upload DeliveryType = "upload"
	// Private is the private delivery type.
	Private = "private"
	// Public is the  delivery type.
	Public = "public"
	// Authenticated is the  delivery type.
	Authenticated = "authenticated"
	// Fetch is the fetch delivery type.
	Fetch = "fetch"
	// Sprite is the sprite delivery type.
	Sprite = "sprite"
	// Text is the text delivery type.
	Text = "text"
	// Multi is the multi delivery type.
	Multi = "multi"
	// Facebook is the facebook delivery type.
	Facebook = "facebook"
	// Twitter is the twitter delivery type.
	Twitter = "twitter"
	// TwitterName is the twitter name delivery type.
	TwitterName = "twitter_name"
	// Gravatar is the gravatar delivery type.
	Gravatar = "gravatar"
	// Youtube is the youtube delivery type.
	Youtube = "youtube"
	// Hulu is the hulu delivery type.
	Hulu = "hulu"
	// Vimeo is the vimeo delivery type.
	Vimeo = "vimeo"
	// Animoto is the animoto delivery type.
	Animoto = "animoto"
	// Worldstarhiphop is the world star hip hop delivery type.
	Worldstarhiphop = "worldstarhiphop"
	// Dailymotion is the daily motion delivery type.
	Dailymotion = "dailymotion"
)

// ModerationStatus is the moderation status of the asset.
type ModerationStatus string

const (
	// Pending is the pending moderation status.
	Pending ModerationStatus = "pending"
	// Approved is the approved moderation status.
	Approved = "approved"
	// Rejected is the rejected moderation status.
	Rejected = "rejected"
)

// Option is the optional parameters custom struct.
type Option map[string]interface{}

// Coordinates represents coordinates on the asset.
type Coordinates [][]int

// CldAPIArray is not just an alias, in addition it has a custom MarshalJSON() for serialisation purposes.
type CldAPIArray []string

// CldAPIMap is not just an alias, in addition it has a custom MarshalJSON() for serialisation purposes.
type CldAPIMap map[string]string

// Metadata is the Cloudinary structured metadata.
type Metadata map[string]interface{}

// BriefAssetResult represents a partial asset result that is returned when assets are listed.
type BriefAssetResult struct {
	AssetID     string    `json:"asset_id"`
	PublicID    string    `json:"public_id"`
	Format      string    `json:"format"`
	Version     int       `json:"version"`
	AssetType   string    `json:"resource_type"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	Bytes       int       `json:"bytes"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Backup      bool      `json:"backup"`
	AccessMode  string    `json:"access_mode"`
	URL         string    `json:"url"`
	SecureURL   string    `json:"secure_url"`
	Tags        []string  `json:"tags,omitempty"`
	Context     CldAPIMap `json:"context,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	Placeholder bool      `json:"placeholder,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// MarshalJSON writes a quoted string in the custom format.
func (cldAPIMap CldAPIMap) MarshalJSON() ([]byte, error) {
	// FIXME: handle escaping
	var params []string
	for name, value := range cldAPIMap {
		params = append(params, strings.Join([]string{name, value}, "="))
	}

	return []byte(strconv.Quote(strings.Join(params, "|"))), nil
}

// MarshalJSON writes a quoted string in the custom format.
func (cldAPIArr CldAPIArray) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(cldAPIArr[:], ","))), nil
}

// ErrorResp is the failed api request main struct.
type ErrorResp struct {
	Message string `json:"message"`
}

// BuildPath builds (joins) the URL path from the provided parts.
func BuildPath(parts ...interface{}) string {
	var partsSlice []string
	//TODO: make it more elegant (?)
	for _, part := range parts {
		partRes := ""
		switch partVal := part.(type) {
		case string:
			partRes = partVal
		case fmt.Stringer:
			partRes = partVal.String()
		default:
			partRes = fmt.Sprintf("%v", partVal)
		}
		if len(partRes) > 0 {
			partsSlice = append(partsSlice, partRes)
		}
	}

	return strings.Join(partsSlice, "/")
}

// SignParameters signs parameters using the provided secret.
func SignParameters(params url.Values, secret string) (string, error) {
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	encodedUnescapedParams, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	hash.Write([]byte(encodedUnescapedParams + secret))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// StructToParams serializes struct to url.Values, which can be further sent to the http client.
func StructToParams(inputStruct interface{}) (url.Values, error) {
	var paramsMap map[string]interface{}
	paramsJSONObj, _ := json.Marshal(inputStruct)
	err := json.Unmarshal(paramsJSONObj, &paramsMap)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	for paramName, value := range paramsMap {
		resBytes, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		res := string(resBytes)
		if strings.HasPrefix(res, "\"") { // FIXME: Fix this dirty hack that prevents double quoting of strings
			res, _ = strconv.Unquote(res)
		}

		params.Add(paramName, res)
	}

	return params, nil
}

// DeferredClose is a wrapper around io.Closer.Close method.
// Logs error if occurred.
func DeferredClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

// IsValidURL checks whether urlCandidate string is a valid URL.
func IsValidURL(urlCandidate string) bool {
	urlStruct, err := url.Parse(urlCandidate)
	if err != nil || urlStruct.Scheme == "" {
		return false
	}

	return true
}

// IsBase64Data checks whether base64Candidate represents a valid base64 encoded string.
func IsBase64Data(base64Candidate string) bool {
	return base64DataRegex.MatchString(base64Candidate)
}

// IsLocalFilePath determines whether the provided path can be a local file.
//
// Since a unix file path can include almost any characters, the way to distinguish between file path and non-file path
// is to check if it can be URL or Base64 encoded data.
func IsLocalFilePath(path interface{}) bool {
	switch pathV := path.(type) {
	case string:
		return !(IsValidURL(pathV) || IsBase64Data(pathV))
	default:
		return false
	}
}
