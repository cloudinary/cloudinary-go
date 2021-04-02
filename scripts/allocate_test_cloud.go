package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

const apiEndpoint = "https://sub-account-testing.cloudinary.com/create_sub_account"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please specify prefix")
		return
	}
	var resp, err = http.PostForm(apiEndpoint, url.Values{"prefix": {os.Args[1]}})

	if nil != err {
		log.Fatal("error happened getting the response", err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		log.Fatal("error happened reading the body", err)
	}

	res := &cloudAllocationResult{}
	err = json.Unmarshal(body, res)

	if err != nil {
		log.Fatal("error happened reading the body", err)
		return
	}

	if res.Status != "success" {
		log.Fatal("error happened: ", res.Status, res.ErrMsg)
	}

	c := res.Payload

	fmt.Printf("cloudinary://%v:%v@%v\n", c.CloudAPIKey, c.CloudAPISecret, c.CloudName)
}

type cloudAllocationResult struct {
	Payload struct {
		CloudAPIKey    string `json:"cloudApiKey"`
		CloudAPISecret string `json:"cloudApiSecret"`
		CloudName      string `json:"cloudName"`
		ID             string `json:"id"`
	} `json:"payload"`
	ErrMsg    string `json:"errMsg"`
	Operation string `json:"operation"`
	Status    string `json:"status"`
}
