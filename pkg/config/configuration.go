package config

import (
	"net/url"
	"os"
)

// Cloudinary Configuration
type Configuration struct {
	Account Account
	Api     Api
}

// Create is creating a new Configuration instance from environment variable
func Create() (*Configuration, error) {
	return CreateFromUrl(os.Getenv("CLOUDINARY_URL"))
}

// CreateFromUrl is creating a new Configuration instance from a cloudinary url
func CreateFromUrl(cldUrlStr string) (*Configuration, error)  {
	cldUrl, err := url.Parse(cldUrlStr)
	if err != nil {
		return nil, err
	}

	pass, _ := cldUrl.User.Password()

	return CreateFromParams(cldUrl.Host, cldUrl.User.Username(), pass)
}

// CreateFromParams is creating a new Configuration instance from provided parameters
func CreateFromParams(cloud string, key string, secret string) (*Configuration, error)  {
	return &Configuration{
		Account: Account{
			CloudName: cloud,
			ApiKey:    key,
			ApiSecret: secret,
		},
	}, nil
}
