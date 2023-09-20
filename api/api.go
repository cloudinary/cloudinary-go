// Package api contains packages used for accessing Cloudinary API functionality.
//
// https://cloudinary.com/documentation/cloudinary_references
package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/internal/signature"
	"io"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// EndPoint represents the API endpoint.
type EndPoint = string

// Version is the Cloudinary Go package version.
const Version = "2.5.1"

// UserAgent contains information about the SDK user agent. Passed to the Cloudinary servers.
var UserAgent = fmt.Sprintf("CloudinaryGo/%s (Go %s)", Version, strings.TrimPrefix(runtime.Version(), "go"))

// UserPlatform provides additional information to be passed with the UserAgent, e.g. "CloudinaryIntegration/1.2.3".
//
// This value is set in platform-specific implementations that use cloudinary-go.
//
// The format of the value should be <ProductName>/Version[ (comment)].
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.43
//
// **Do not set this value in application code!**
var UserPlatform = ""

// GetUserAgent provides the UserAgent string that is passed to the Cloudinary servers.
//
// Prepends UserPlatform if it is defined.
func GetUserAgent() string {
	if UserPlatform == "" {
		return UserAgent
	}

	return fmt.Sprintf("%s %s", UserPlatform, UserAgent)
}

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

// HookExecution is the result of a hook execution.
type HookExecution map[string]interface{}

// BriefAssetResult represents a partial asset result that is returned when assets are listed.
type BriefAssetResult struct {
	AssetID     string    `json:"asset_id"`
	PublicID    string    `json:"public_id"`
	AssetFolder string    `json:"asset_folder"`
	DisplayName string    `json:"display_name"`
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
	Context     Metadata  `json:"context,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	Placeholder bool      `json:"placeholder,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// LastUpdated represents the details of the asset last updated time.
type LastUpdated struct {
	AccessControlUpdatedAt time.Time `json:"access_control_updated_at,omitempty"`
	ContextUpdatedAt       time.Time `json:"context_updated_at,omitempty"`
	MetadataUpdatedAt      time.Time `json:"metadata_updated_at,omitempty"`
	PublicIDUpdatedAt      time.Time `json:"public_id_updated_at,omitempty"`
	TagsUpdatedAt          time.Time `json:"tags_updated_at,omitempty"`
	UpdatedAt              time.Time `json:"updated_at,omitempty"`
}

// MarshalJSON writes a quoted string in the custom format.
func (cldAPIMap Metadata) MarshalJSON() ([]byte, error) {
	// FIXME: handle escaping
	var params []string
	for name, value := range cldAPIMap {
		val, err := encodeParamValue(value)
		if err != nil {
			return nil, err
		}
		params = append(params, strings.Join([]string{name, val}, "="))
	}

	return []byte(strconv.Quote(strings.Join(params, "|"))), nil
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
	if _, withTimestamp := params["timestamp"]; !withTimestamp || params["timestamp"][0] == "0" {
		params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	}

	encodedUnescapedParams, err := url.QueryUnescape(params.Encode())
	if err != nil {
		return "", err
	}

	rawSignature, err := signature.Sign(encodedUnescapedParams, secret, signature.SHA1)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawSignature), nil
}

// MarshalJSONRaw marshals JSON without HTML characters escaping, which is enabled in the standard library.
// In addition, it removes the last newline character from the resulting bytes since there is no reason to keep it.
func MarshalJSONRaw(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	bufferBytes := buffer.Bytes()
	if len(bufferBytes) > 0 {
		bufferBytes = bufferBytes[:len(bufferBytes)-1]
	}
	return bufferBytes, err
}

// ReMarshalJSON unmarshals and then marshals data - as the result the data is sorted by key.
func ReMarshalJSON(bytes []byte) ([]byte, error) {
	var data interface{}
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return MarshalJSONRaw(data)
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
		kind := reflect.ValueOf(value).Kind()

		if kind == reflect.Slice || kind == reflect.Array {
			rVal := reflect.ValueOf(value)
			for i := 0; i < rVal.Len(); i++ {
				item := rVal.Index(i)
				val, err := encodeParamValue(item.Interface())
				if err != nil {
					return nil, err
				}

				arrParamName := fmt.Sprintf("%s[%d]", paramName, i)
				params.Add(arrParamName, val)
			}

			continue
		}

		val, err := encodeParamValue(value)
		if err != nil {
			return nil, err
		}

		params.Add(paramName, val)
	}

	return params, nil
}

func encodeParamValue(value interface{}) (string, error) {
	resBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	res := string(resBytes)
	if strings.HasPrefix(res, "\"") { // FIXME: Fix this dirty hack that prevents double quoting of strings
		res, _ = strconv.Unquote(res)
	}

	return res, nil
}

// DeferredClose is a wrapper around io.Closer.Close method.
// Logs error if occurred.
func DeferredClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

// HandleRawResponse sets a raw Response field value (JSON) in the Result structs that support it, for future proofing
func HandleRawResponse(bodyBytes []byte, result interface{}) error {
	resultMetaValue := reflect.ValueOf(result).Elem()

	if resultMetaValue.Kind() != reflect.Struct {
		return nil
	}

	responseField := resultMetaValue.FieldByName("Response")
	if responseField == (reflect.Value{}) {
		// no 'Response' field
		return nil
	}

	var rawResponse interface{}

	err := json.Unmarshal(bodyBytes, &rawResponse)
	if err != nil {
		return err
	}

	rawResponseValue := reflect.New(reflect.TypeOf(rawResponse))
	rawResponseValue.Elem().Set(reflect.ValueOf(rawResponse))

	responseField.Set(rawResponseValue)

	return nil
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

// Bool returns a pointer for the provided boolean.
func Bool(b bool) *bool {
	return &b
}
