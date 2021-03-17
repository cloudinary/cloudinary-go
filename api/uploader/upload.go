// Package uploader is used for accessing Cloudinary Upload API functionality.
//
//https://cloudinary.com/documentation/image_upload_api_reference
package uploader

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Upload Api main struct
type Api struct {
	Config config.Configuration
	client http.Client
}

// New creates a new Admin Api instance from the environment variable.
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

	//log.Println(string(resp)) FIXME: find a good logger

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
	if u.Config.Cloud.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	signature, err := api.SignParameters(requestParams, u.Config.Cloud.ApiSecret)
	if err != nil {
		return nil, err
	}
	requestParams.Add("signature", signature)
	requestParams.Add("api_key", u.Config.Cloud.ApiKey)

	return requestParams, nil
}

func (u *Api) postForm(ctx context.Context, urlPath interface{}, formParams url.Values) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	_, err := bodyBuf.Write([]byte(formParams.Encode()))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(u.Config.Api.Timeout)*time.Second)
	defer cancel()

	return u.postBody(ctx, urlPath, bodyBuf, nil)
}

func (u *Api) postFile(ctx context.Context, file interface{}, formParams url.Values) ([]byte, error) {
	unsigned, _ := strconv.ParseBool(formParams.Get("unsigned"))

	if !unsigned {
		var err error
		formParams, err = u.signRequest(formParams)
		if err != nil {
			return nil, err
		}
	}

	uploadEndpoint := api.BuildPath(api.Auto, Upload)
	switch fileValue := file.(type) {
	case string:
		if !api.IsLocalFilePath(file) {
			// Can be URL, Base64 encoded string, etc.
			formParams.Add("file", fileValue)

			return u.postForm(ctx, uploadEndpoint, formParams)
		} else {
			return u.postLocalFile(ctx, uploadEndpoint, fileValue, formParams)
		}
	case io.Reader:
		return u.postIOReader(ctx, uploadEndpoint, fileValue, "file", formParams, map[string]string{}, 0)
	default:
		return nil, errors.New("unsupported file type")
	}
}

// postLocalFile creates a new file upload http request with optional extra params.
func (u *Api) postLocalFile(ctx context.Context, urlPath string, filePath string, formParams url.Values) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(file)

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fi.Size() > u.Config.Api.ChunkSize {
		return u.postLargeFile(ctx, urlPath, file, formParams)
	}

	return u.postIOReader(ctx, urlPath, file, fi.Name(), formParams, map[string]string{}, 0)
}

func (u *Api) postLargeFile(ctx context.Context, urlPath string, file *os.File, formParams url.Values) ([]byte, error) {
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-Unique-Upload-Id": randomPublicId(),
	}

	var res []byte

	fileSize := fi.Size()
	var currPos int64 = 0
	for currPos < fileSize {
		currChunkSize := min(fileSize-currPos, u.Config.Api.ChunkSize)

		headers["Content-Range"] = fmt.Sprintf("bytes %v-%v/%v", currPos, currPos+currChunkSize-1, fileSize)

		res, err = u.postIOReader(ctx, urlPath, file, fi.Name(), formParams, headers, currChunkSize)
		if err != nil {
			return nil, err
		}

		currPos += currChunkSize
	}

	return res, nil
}

func (u *Api) postIOReader(ctx context.Context, urlPath string, reader io.Reader, name string, formParams url.Values, headers map[string]string, chunkSize int64) ([]byte, error) {
	bodyBuf := new(bytes.Buffer)
	formWriter := multipart.NewWriter(bodyBuf)

	headers["Content-Type"] = formWriter.FormDataContentType()

	for key, val := range formParams {
		_ = formWriter.WriteField(key, val[0])
	}

	partWriter, err := formWriter.CreateFormFile("file", name)
	if err != nil {
		return nil, err
	}

	if chunkSize != 0 {
		_, err = io.CopyN(partWriter, reader, chunkSize)
	} else {
		_, err = io.Copy(partWriter, reader)
	}
	if err != nil {
		return nil, err
	}

	err = formWriter.Close()
	if err != nil {
		return nil, err
	}

	if u.Config.Api.UploadTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(u.Config.Api.UploadTimeout)*time.Second)
		defer cancel()
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
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

func (u *Api) getUploadURL(urlPath interface{}) string {
	return fmt.Sprintf("%v/%v/%v", api.BaseUrl(u.Config.Api.UploadPrefix), u.Config.Cloud.CloudName, api.BuildPath(urlPath))
}

func getAssetType(requestParams interface{}) string {
	// FIXME: define interface or something to just access the field, and/or have a default value ("image") in the struct
	assetType := fmt.Sprintf("%v", reflect.ValueOf(requestParams).FieldByName("ResourceType"))
	if assetType == "" {
		assetType = api.Image.String()
	}

	return assetType
}

// randomPublicId generates a random public ID string, which is the first 16 characters of sha1 of uuid.
func randomPublicId() string {
	hash := sha1.New()
	hash.Write([]byte(uuid.NewString()))

	return hex.EncodeToString(hash.Sum(nil))[0:16]
}

// min returns minimum of two integers
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
