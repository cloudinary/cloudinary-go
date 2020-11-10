package admin

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"cloudinary-labs/cloudinary-go/pkg/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Api main struct
type Api struct {
	Config config.Configuration
}

// Create is creating a new Api instance from environment variable
func Create() *Api {
	return &Api{
		Config: *config.Create(),
	}
}

func (a *Api) callApi(url string, params url.Values) []byte {

	req, err := http.NewRequest("GET",
		fmt.Sprintf("%v/%v/%v", api.ApiEndpoint, a.Config.Account.CloudName, url),
		strings.NewReader(params.Encode()),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", api.UserAgent)
	req.SetBasicAuth(a.Config.Account.ApiKey, a.Config.Account.ApiSecret)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}
