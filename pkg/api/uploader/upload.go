package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
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
func Create() *Api {
	return &Api{
		Config: *config.Create(),
	}
}

func (u *Api) postForm(url string, formParams url.Values) []byte {
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%v/%v/%v", api.ApiEndpoint, u.Config.Account.CloudName, url),
		strings.NewReader(formParams.Encode()),
	)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", api.UserAgent)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

func (u *Api) postAndSignForm(url string, formParams url.Values) []byte {
	if u.Config.Account.ApiSecret == "" {
		panic("Must provide Api Secret")
	}
	formParams.Add("signature", api.SignRequest(formParams, u.Config.Account.ApiSecret))
	formParams.Add("api_key", u.Config.Account.ApiKey)

	return u.postForm(url, formParams)
}

func (u *Api) postFile(file string, formParams url.Values) []byte {
	if u.Config.Account.ApiSecret == "" {
		panic("Must provide Api Secret")
	}

	formParams.Add("signature", api.SignRequest(formParams, u.Config.Account.ApiSecret))
	formParams.Add("api_key", u.Config.Account.ApiKey)
	formParams.Add("file", file)

	return u.postForm("auto/upload", formParams)
}
