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

type EndPoint string

const Version = "0.2.0"

const UserAgent = "CloudinaryGo/" + Version

var apiVersion = "1_1"

func BaseUrl(uploadPrefix string) string {
	return fmt.Sprintf("%s/v%s", uploadPrefix, apiVersion)
}

var base64DataRegex = regexp.MustCompile("^data:([\\w-]+/[\\w\\-+.]+)?(;[\\w-]+=[\\w-]+)*;base64,([a-zA-Z0-9/+\\n=]+)$")

type AssetType string

func (a AssetType) String() string {
	if a == "" {
		a = Image
	}
	return string(a)
}

func (e EndPoint) String() string {
	return string(e)
}

const (
	Image AssetType = "image"
	Video           = "video"
	File            = "raw"
	Auto            = "auto"
	All             = "all"
)

type DeliveryType string

func (d DeliveryType) String() string {
	if d == "" {
		d = Upload
	}
	return string(d)
}

const (
	Upload          DeliveryType = "upload"
	Private                      = "private"
	Public                       = "public"
	Authenticated                = "authenticated"
	Fetch                        = "fetch"
	Sprite                       = "sprite"
	Text                         = "text"
	Multi                        = "multi"
	Facebook                     = "facebook"
	Twitter                      = "twitter"
	TwitterName                  = "twitter_name"
	Gravatar                     = "gravatar"
	Youtube                      = "youtube"
	Hulu                         = "hulu"
	Vimeo                        = "vimeo"
	Animoto                      = "animoto"
	Worldstarhiphop              = "worldstarhiphop"
	Dailymotion                  = "dailymotion"
)

type ModerationStatus string

const (
	Pending  ModerationStatus = "pending"
	Approved                  = "approved"
	Rejected                  = "rejected"
)

// Option is the optional parameters custom struct
type Option map[string]interface{}

type Coordinates [][]int
type CldApiArray []string

type CldApiMap map[string]string
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
	Context     CldApiMap `json:"context,omitempty"`
	Metadata    Metadata  `json:"metadata,omitempty"`
	Placeholder bool      `json:"placeholder,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// MarshalJSON writes a quoted string in the custom format.
func (cldApiMap CldApiMap) MarshalJSON() ([]byte, error) {
	// FIXME: handle escaping
	var params []string
	for name, value := range cldApiMap {
		params = append(params, strings.Join([]string{name, value}, "="))
	}

	return []byte(strconv.Quote(strings.Join(params, "|"))), nil
}

// MarshalJSON writes a quoted string in the custom format.
func (cldApiArr CldApiArray) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(cldApiArr[:], ","))), nil
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
	paramsJsonObj, _ := json.Marshal(inputStruct)
	err := json.Unmarshal(paramsJsonObj, &paramsMap)
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

// IsValidUrl checks whether urlCandidate string is a valid URL.
func IsValidUrl(urlCandidate string) bool {
	_, err := url.ParseRequestURI(urlCandidate)
	if err != nil {
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
		return !(IsValidUrl(pathV) || IsBase64Data(pathV))
	default:
		return false
	}
}
