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

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/config"
)

type Api struct {
	Config config.Configuration
	client http.Client
}

// New creates a new Admin Api instance from the environment variable (CLOUDINARY_URL).
func New() (*Api, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}
	return &Api{
		Config: *c,
		client: http.Client{},
	}, nil
}

func (a *Api) get(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callApi(ctx, http.MethodGet, path, requestParams, result)
}

func (a *Api) post(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callApi(ctx, http.MethodPost, path, requestParams, result)
}

func (a *Api) put(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callApi(ctx, http.MethodPut, path, requestParams, result)
}

func (a *Api) delete(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	return a.callApi(ctx, http.MethodDelete, path, requestParams, result)
}

func (a *Api) callApi(ctx context.Context, method string, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error) {
	var body io.Reader = nil

	// Populate body for POST/PUT/DELETE
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		jsonReq, err := json.Marshal(requestParams)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonReq)
	}
	req, err := http.NewRequest(method,
		fmt.Sprintf("%v/%v/%v", api.BaseUrl(a.Config.Api.UploadPrefix), a.Config.Cloud.CloudName, api.BuildPath(path)),
		body,
	)
	if err != nil {
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

	req.Header.Set("User-Agent", api.UserAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(a.Config.Cloud.ApiKey, a.Config.Cloud.ApiSecret)

	ctx, cancel := context.WithTimeout(ctx, time.Duration(a.Config.Api.Timeout)*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//log.println(string(bodyBytes)) FIXME: find a good logger

	err = json.Unmarshal(bodyBytes, result)

	return resp, err
}
