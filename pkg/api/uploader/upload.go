package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// Upload Api main struct
type Api struct {
	Config config.Configuration
	client http.Client
}

// Create is creating a new Api instance from environment variable
func Create() (*Api, error) {
	c, err := config.Create()
	if err != nil {
		return nil, err
	}
	return &Api{
		Config: *c,
		client: http.Client{},
	}, nil
}

func (u *Api) postForm(ctx context.Context, path interface{}, formParams url.Values) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%v/%v/%v", api.BaseUrl, u.Config.Account.CloudName, api.BuildPath(path)),
		strings.NewReader(formParams.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)

	req = req.WithContext(ctx)

	resp, err := u.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

func (u *Api) postAndSignForm(ctx context.Context, path string, formParams url.Values) ([]byte, error) {
	if u.Config.Account.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	signature, err := api.SignRequest(formParams, u.Config.Account.ApiSecret)
	if err != nil {
		return nil, err
	}
	formParams.Add("signature", signature)
	formParams.Add("api_key", u.Config.Account.ApiKey)

	return u.postForm(ctx, path, formParams)
}

func (u *Api) postFile(ctx context.Context, file string, formParams url.Values) ([]byte, error) {
	if u.Config.Account.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	signature, err := api.SignRequest(formParams, u.Config.Account.ApiSecret)
	if err != nil {
		return nil, err
	}
	formParams.Add("signature", signature)
	formParams.Add("api_key", u.Config.Account.ApiKey)
	formParams.Add("file", file)

	return u.postForm(ctx, api.BuildPath("auto", Upload), formParams)
}

func (u *Api) callUploadApi(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	return u.callUploadApiWithParams(ctx, api.BuildPath(getResourceType(requestParams), path), formParams, result)
}

func (u *Api) callUploadApiWithParams(ctx context.Context, path string, formParams url.Values, result interface{}) error {
	resp, err := u.postAndSignForm(ctx, path, formParams)
	if err != nil {
		return err
	}

	println(string(resp))

	err = json.Unmarshal(resp, result)

	return err

}

func getResourceType(requestParams interface{}) string {
	// FIXME: define interface or something to just access the field, and/or have a default value ("image") in the struct
	resourceType := fmt.Sprintf("%v", reflect.ValueOf(requestParams).FieldByName("ResourceType"))
	if resourceType == "" {
		resourceType = fmt.Sprintf("%v", api.Image)
	}

	return resourceType
}
