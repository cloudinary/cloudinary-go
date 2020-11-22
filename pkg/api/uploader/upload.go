package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Upload Api main struct
type Api struct {
	Config config.Configuration
}

// Create is creating a new Api instance from environment variable
func Create() (*Api, error) {
	c, err := config.Create()
	if err != nil {
		return nil, err
	}
	return &Api{
		Config: *c,
	}, nil
}

func (u *Api) postForm(url string, formParams url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%v/%v/%v", api.BaseUrl, u.Config.Account.CloudName, url),
		strings.NewReader(formParams.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", api.UserAgent)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer api.DeferredClose(resp.Body)

	return ioutil.ReadAll(resp.Body)
}

func (u *Api) postAndSignForm(url string, formParams url.Values) ([]byte, error) {
	if u.Config.Account.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	formParams.Add("signature", api.SignRequest(formParams, u.Config.Account.ApiSecret))
	formParams.Add("api_key", u.Config.Account.ApiKey)

	return u.postForm(url, formParams)
}

func (u *Api) postFile(file string, formParams url.Values) ([]byte, error) {
	if u.Config.Account.ApiSecret == "" {
		return nil, errors.New("must provide Api Secret")
	}

	formParams.Add("signature", api.SignRequest(formParams, u.Config.Account.ApiSecret))
	formParams.Add("api_key", u.Config.Account.ApiKey)
	formParams.Add("file", file)

	return u.postForm("auto/upload", formParams)
}
