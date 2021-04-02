// Package config defines the Cloudinary configuration.
package config

import (
	"errors"
	"net/url"
	"os"

	"github.com/creasty/defaults"
)

// Configuration is the main configuration struct.
type Configuration struct {
	Cloud Cloud
	API   API
}

// New returns a new Configuration instance from the environment variable
func New() (*Configuration, error) {
	return NewFromURL(os.Getenv("CLOUDINARY_URL"))
}

// NewFromURL returns a new Configuration instance from a cloudinary url.
func NewFromURL(cldURLStr string) (*Configuration, error) {
	if cldURLStr == "" {
		return nil, errors.New("must provide CLOUDINARY_URL")
	}

	cldURL, err := url.Parse(cldURLStr)
	if err != nil {
		return nil, err
	}

	pass, _ := cldURL.User.Password()

	return NewFromParams(cldURL.Host, cldURL.User.Username(), pass)
}

// NewFromParams returns a new Configuration instance from the provided parameters.
func NewFromParams(cloud string, key string, secret string) (*Configuration, error) {
	conf := &Configuration{
		Cloud: Cloud{
			CloudName: cloud,
			APIKey:    key,
			APISecret: secret,
		},
		API: API{},
	}

	if err := defaults.Set(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
