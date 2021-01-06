package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/config"
)

// Cloudinary main struct
type Cloudinary struct {
	Config config.Configuration
	Admin  admin.Api
	Upload uploader.Api
}

// Create returns a new Cloudinary instance from environment variable.
func Create() (*Cloudinary, error) {
	c, err := config.Create()
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromUrl returns a new Cloudinary instance from a cloudinary url.
func CreateFromUrl(cloudinaryUrl string) (*Cloudinary, error) {
	c, err := config.CreateFromUrl(cloudinaryUrl)
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromParams returns a new Cloudinary instance from the provided parameters.
func CreateFromParams(cloud string, key string, secret string) (*Cloudinary, error) {
	c, err := config.CreateFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromConfiguration returns a new Cloudinary instance from the provided configuration.
func CreateFromConfiguration(configuration config.Configuration) (*Cloudinary, error) {
	return &Cloudinary{
		Config: configuration,
		Admin: admin.Api{
			Config: configuration,
		},
		Upload: uploader.Api{
			Config: configuration,
		},
	}, nil
}
