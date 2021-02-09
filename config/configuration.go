// Package config defines the Cloudinary configuration.
package config

import (
	"log"
	"net/url"
	"os"

	"github.com/creasty/defaults"
)

// Logger cloudinary logger
type Logger func(...interface{})

// Configuration is the main configuration struct.
type Configuration struct {
	Account  Account
	Api      Api
	DebugLog Logger
	ErrorLog Logger
}

// Create returns a new Configuration instance from the environment variable
func Create() (*Configuration, error) {
	return CreateFromUrl(os.Getenv("CLOUDINARY_URL"))
}

// CreateFromUrl returns a new Configuration instance from a cloudinary url.
func CreateFromUrl(cldUrlStr string) (*Configuration, error) {
	cldUrl, err := url.Parse(cldUrlStr)
	if err != nil {
		return nil, err
	}

	pass, _ := cldUrl.User.Password()

	return CreateFromParams(cldUrl.Host, cldUrl.User.Username(), pass)
}

// CreateFromParams returns a new Configuration instance from the provided parameters.
func CreateFromParams(cloud string, key string, secret string) (*Configuration, error) {
	conf := &Configuration{
		Account: Account{
			CloudName: cloud,
			ApiKey:    key,
			ApiSecret: secret,
		},
		Api: Api{},
	}

	if err := defaults.Set(conf); err != nil {
		return nil, err
	}

	conf.DebugLog = log.Println
	conf.ErrorLog = log.Println

	return conf, nil
}
