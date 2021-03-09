package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const ApiEndpoint = "https://sub-account-testing.cloudinary.com/create_sub_account"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify prefix")
		return
	}
	var resp, err = http.PostForm(ApiEndpoint,
		url.Values{"prefix": {os.Args[1]}})

	if nil != err {
		fmt.Println("errorination happened getting the response", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		fmt.Println("errorination happened reading the body", err)
		return
	}

	fmt.Println(string(body[:]))
}
