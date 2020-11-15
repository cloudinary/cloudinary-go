package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type EndPoint string

const Version = "0.0.1-Alpha0"
const UserAgent = "CloudinaryGo/" + Version

var BaseUrl = "https://api.cloudinary.com/v1_1"

type AssetType string

func (a AssetType) ToString() string {
	if a == "" {
		a = Image
	}
	return string(a)
}

const (
	Image AssetType = "image"
	Video           = "video"
	File            = "raw"
	Auto            = "auto"
	All             = "all"
)

type DeliveryType string

func (d DeliveryType) ToString() string {
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

type Coordinates   [][]int
type CldApiArray []string

type Context map[string]string
type Metadata map[string]interface{}

// MarshalJSON writes a quoted string in the custom format
func (cdlApiArr CldApiArray) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(cdlApiArr[:], ","))), nil
}

// ErrorResp is the failed api request main struct
type ErrorResp struct {
	Message string `json:"message"`
}

func BuildPath(parts ...interface{}) string {
	var partsSlice []string
	for _, part := range parts {
		if part != "" {
			partsSlice = append(partsSlice, fmt.Sprintf("%v", part))
		}
	}

	return strings.Join(partsSlice, "/")
}

func SignRequest(params url.Values, secret string) string {
	params.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	encodedParams := params.Encode()

	hash := sha1.New()
	hash.Write([]byte(encodedParams + secret))

	return hex.EncodeToString(hash.Sum(nil))
}

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
			res, _ = strconv.Unquote(string(res))
		}

		params.Add(paramName, res)
	}

	return params, nil
}
