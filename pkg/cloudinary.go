package cloudinary

import (
	"cloudinary-labs/cloudinary-go/pkg/api/admin"
	"cloudinary-labs/cloudinary-go/pkg/api/uploader"
	"cloudinary-labs/cloudinary-go/pkg/config"
)

// Cloudinary main struct
type Cloudinary struct {
	Config config.Configuration
	Admin  admin.Api
	Upload uploader.Api
}

// Create is creating a new Cloudinary instance from environment variable
func Create() (*Cloudinary, error) {
	c, err := config.Create()
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromUrl is creating a new Cloudinary instance from a cloudinary url
func CreateFromUrl(cloudinaryUrl string) (*Cloudinary, error) {
	c, err := config.CreateFromUrl(cloudinaryUrl)
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromParams is creating a new Cloudinary instance from provided parameters
func CreateFromParams(cloud string, key string, secret string) (*Cloudinary, error) {
	c, err := config.CreateFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}
	return CreateFromConfiguration(*c)
}

// CreateFromParams is creating a new Cloudinary instance from provided configuration
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
