// Package admin is used for accessing Cloudinary Admin API functionality.
//
// https://cloudinary.com/documentation/admin_api
package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/logger"
)

// API is the Admin API struct.
type API struct {
	Config config.Configuration
	Logger *logger.Logger
	Client http.Client
}

// Direction is the sorting direction.
type Direction string

const (
	// Ascending direction.
	Ascending Direction = "asc"
	// Descending direction.
	Descending = "desc"
)

// OrderByField is the field to order by.
type OrderByField string

// OrderFieldValue defines to order by value.
const OrderFieldValue OrderByField = "value"

// OrderFieldLabel defines to order by label.
const OrderFieldLabel OrderByField = "label"

// OrderFieldExternalID defines to order by external_id.
const OrderFieldExternalID OrderByField = "external_id"

// OrderFieldCreatedAt defines to order by created_at.
const OrderFieldCreatedAt OrderByField = "created_at"

// New creates a new Admin API instance from the environment variable (CLOUDINARY_URL).
func New() (*API, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}
	return NewWithConfiguration(c)
}

// NewWithConfiguration a new Admin API instance with the given Configuration
func NewWithConfiguration(c *config.Configuration) (*API, error) {
	return &API{
		Config: *c,
		Client: http.Client{},
		Logger: logger.New(),
	}, nil
}

func (a *API) get(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodGet, path, requestParams, result)
}

func (a *API) post(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodPost, path, requestParams, result)
}

func (a *API) put(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodPut, path, requestParams, result)
}

func (a *API) delete(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callAPI(ctx, http.MethodDelete, path, requestParams, result)
}

func (a *API) callAPI(ctx context.Context, method string, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	var body io.Reader = nil

	// Populate body for POST/PUT/DELETE.
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		jsonReq, err := json.Marshal(requestParams)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	}
	req, err := http.NewRequest(method,
		fmt.Sprintf("%v/%v/%v", api.BaseURL(a.Config.API.UploadPrefix), a.Config.Cloud.CloudName, api.BuildPath(path)),
		body,
	)
	if err != nil {
		a.Logger.Error(err)
		return nil, err
	}

	// Handle GET request query parameters
	if body == nil {
		params, err := api.StructToParams(requestParams)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = params.Encode()
	}

	req.Header.Set("User-Agent", api.GetUserAgent())
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	setAuth(a, req)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.Config.API.Timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	a.Logger.Debug(string(bodyBytes))

	err = json.Unmarshal(bodyBytes, result)
	if err != nil {
		return resp, err
	}

	err = api.HandleRawResponse(bodyBytes, result)

	return resp, err
}

func setAuth(a *API, req *http.Request) {
	if a.Config.Cloud.OAuthToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Config.Cloud.OAuthToken))
	} else {
		req.SetBasicAuth(a.Config.Cloud.APIKey, a.Config.Cloud.APISecret)
	}
}
