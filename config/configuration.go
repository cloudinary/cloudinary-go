// Package config defines the Cloudinary configuration.
package config

import (
	"net/url"
	"os"

	"github.com/creasty/defaults"
)

// Configuration is the main configuration struct.
type Configuration struct {
	Cloud Cloud
	Api   Api
}

// New returns a new Configuration instance from the environment variable
func New() (*Configuration, error) {
	return NewFromUrl(os.Getenv("CLOUDINARY_URL"))
}

// NewFromUrl returns a new Configuration instance from a cloudinary url.
func NewFromUrl(cldUrlStr string) (*Configuration, error) {
	cldUrl, err := url.Parse(cldUrlStr)
	if err != nil {
		return nil, err
	}

	pass, _ := cldUrl.User.Password()

	return NewFromParams(cldUrl.Host, cldUrl.User.Username(), pass)
}

// NewFromParams returns a new Configuration instance from the provided parameters.
func NewFromParams(cloud string, key string, secret string) (*Configuration, error) {
	conf := &Configuration{
		Cloud: Cloud{
			CloudName: cloud,
			ApiKey:    key,
			ApiSecret: secret,
		},
		Api: Api{},
	}

	if err := defaults.Set(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
