package uploader

import (
	"bufio"
	"bytes"
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
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

func (u *Api) callUploadApi(ctx context.Context, path interface{}, requestParams interface{}, result interface{}) error {
	formParams, err := api.StructToParams(requestParams)
	if err != nil {
		return err
	}

	return u.callUploadApiWithParams(ctx, api.BuildPath(getAssetType(requestParams), path), formParams, result)
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

func (u *Api) postAndSignForm(ctx context.Context, urlPath string, formParams url.Values) ([]byte, error) {
	formParams, err := u.signRequest(formParams)
	if err != nil {
		return nil, err
	}

	return u.postForm(ctx, urlPath, formParams)
}

func (u *Api) signRequest(requestParams url.Values) (url.Values, error) {
	if u.Config.Account.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	signature, err := api.SignParameters(requestParams, u.Config.Account.ApiSecret)
	if err != nil {
		return nil, err
	}
	requestParams.Add("signature", signature)
	requestParams.Add("api_key", u.Config.Account.ApiKey)

	return requestParams, nil
}

func (u *Api) postForm(ctx context.Context, urlPath interface{}, formParams url.Values) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	_, err := bodyBuf.Write([]byte(formParams.Encode()))
	if err != nil {
		return nil, err
	}

	return u.postBody(ctx, urlPath, bodyBuf, nil)
}

func (u *Api) postFile(ctx context.Context, file interface{}, formParams url.Values) ([]byte, error) {
	formParams, err := u.signRequest(formParams)
	if err != nil {
		return nil, err
	}
	uploadEndpoint := api.BuildPath(api.Auto, Upload)
	switch fileValue := file.(type) {
	case string:
		if !api.IsLocalFilePath(file) {
			formParams.Add("file", fileValue)

			return u.postForm(ctx, uploadEndpoint, formParams)
		} else {
			return u.postLocalFile(ctx, uploadEndpoint, fileValue, formParams)
		}
	case io.Reader:
		return u.postIOReader(ctx, uploadEndpoint, fileValue, "file", formParams)
	default:
		return nil, errors.New("unsupported file type")
	}
}

// Creates a new file upload http request with optional extra params
func (u *Api) postLocalFile(ctx context.Context, urlPath string, filePath string, formParams url.Values) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	bufferedReader := bufio.NewReader(file)

	defer api.DeferredClose(file)

	return u.postIOReader(ctx, urlPath, bufferedReader, fi.Name(), formParams)
}

func (u *Api) postIOReader(ctx context.Context, urlPath string, reader io.Reader, name string, formParams url.Values) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	formWriter := multipart.NewWriter(bodyBuf)

	headers := map[string]string{
		"Content-Type": formWriter.FormDataContentType(),
	}

	for key, val := range formParams {
		_ = formWriter.WriteField(key, val[0])
	}

	partWriter, err := formWriter.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(partWriter, reader)
	if err != nil {
		return nil, err
	}

	err = formWriter.Close()
	if err != nil {
		return nil, err
	}

	return u.postBody(ctx, urlPath, bodyBuf, headers)
}

func (u *Api) postBody(ctx context.Context, urlPath interface{}, bodyBuf *bytes.Buffer, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost,
		u.getUploadURL(urlPath),
		bodyBuf,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)
	for key, val := range headers {
		req.Header.Add(key, val)
	}

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

func (u *Api) getUploadURL(urlPath interface{}) string {
	return fmt.Sprintf("%v/%v/%v", api.BaseUrl, u.Config.Account.CloudName, api.BuildPath(urlPath))
}

func getAssetType(requestParams interface{}) string {
	// FIXME: define interface or something to just access the field, and/or have a default value ("image") in the struct
	assetType := fmt.Sprintf("%v", reflect.ValueOf(requestParams).FieldByName("ResourceType"))
	if assetType == "" {
		assetType = fmt.Sprintf("%v", api.Image)
	}

	return assetType
}
