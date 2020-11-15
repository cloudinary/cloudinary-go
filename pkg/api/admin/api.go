package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Api main struct
type Api struct {
	Config config.Configuration
	client http.Client
}

// Create is creating a new Api instance from environment variable
func Create() *Api {
	return &Api{
		Config: *config.Create(),
		client: http.Client{},
	}
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

func (a *Api) callApi(ctx context.Context, method string, path interface{}, requestParams interface{}, result interface{}) (*http.Response, error){
	params, err := api.StructToParams(requestParams)
	if err != nil {
		return nil, err
	}

	var body io.Reader = nil

	if method == http.MethodPost || method == http.MethodPut {
		decodedValue, err := url.QueryUnescape(params.Encode())
		if err != nil {
			return nil, err
		}
		body = strings.NewReader(decodedValue)
	}
	req, err := http.NewRequest(method,
		fmt.Sprintf("%v/%v/%v", api.BaseUrl, a.Config.Account.CloudName, api.BuildPath(path)),
		body,
	)
	if err != nil {
		return nil, err
	}

	if body == nil {
		req.URL.RawQuery = params.Encode()
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.SetBasicAuth(a.Config.Account.ApiKey, a.Config.Account.ApiSecret)

	req = req.WithContext(ctx)

	resp, err := a.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer deferredClose(resp.Body)

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	println(string(bodyBytes))

	err = json.Unmarshal(bodyBytes, result)

	return resp, err
}

func deferredClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
