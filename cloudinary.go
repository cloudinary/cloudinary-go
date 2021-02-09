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

// New returns a new Cloudinary instance from environment variable.
func New() (*Cloudinary, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}
	return NewFromConfiguration(*c)
}

// NewFromUrl returns a new Cloudinary instance from a cloudinary url.
func NewFromUrl(cloudinaryUrl string) (*Cloudinary, error) {
	c, err := config.NewFromUrl(cloudinaryUrl)
	if err != nil {
		return nil, err
	}
	return NewFromConfiguration(*c)
}

// NewFromParams returns a new Cloudinary instance from the provided parameters.
func NewFromParams(cloud string, key string, secret string) (*Cloudinary, error) {
	c, err := config.NewFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}
	return NewFromConfiguration(*c)
}

// NewFromConfiguration returns a new Cloudinary instance from the provided configuration.
func NewFromConfiguration(configuration config.Configuration) (*Cloudinary, error) {
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
