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

const Version = "0.0.1-Alpha0"
const UserAgent = "CloudinaryGo/" + Version

var ApiEndpoint = "https://api.cloudinary.com/v1_1"

// Option is the optional parameters custom struct
type Option map[string]interface{}

type CldApiArray []string

type Context map[string]string
type Metadata map[string]interface{}

// MarshalJSON writes a quoted string in the custom format
func (cdlApiArr CldApiArray) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strings.Join(cdlApiArr[:], ","))), nil
}

// UnmarshalJSON Parses the json string in the custom format
func (cdlApiArr *CldApiArray) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	*cdlApiArr = strings.Split(s, ",")
	return
}

// ErrorResp is the failed api request main struct
type ErrorResp struct {
	Message string `json:"message"`
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
		params.Add(paramName, fmt.Sprintf("%v", value))
	}

	return params, nil
}
