package cloudinary

import (
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/asset"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/logger"
)

// Cloudinary main struct
type Cloudinary struct {
	Config config.Configuration
	Admin  admin.API
	Upload uploader.API
	Logger *logger.Logger
}

// New returns a new Cloudinary instance from environment variable.
func New() (*Cloudinary, error) {
	c, err := config.New()
	if err != nil {
		return nil, err
	}

	return NewFromConfiguration(*c)
}

// NewFromURL returns a new Cloudinary instance from a cloudinary url.
func NewFromURL(cloudinaryURL string) (*Cloudinary, error) {
	c, err := config.NewFromURL(cloudinaryURL)
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
	logger := logger.New()

	return &Cloudinary{
		Config: configuration,
		Admin: admin.API{
			Config: configuration,
			Logger: logger,
		},
		Upload: uploader.API{
			Config: configuration,
			Logger: logger,
		},
		Logger: logger,
	}, nil
}

// Image creates a new asset.Image instance.
func (c Cloudinary) Image(publicID string) (*asset.Asset, error) {
	return asset.Image(publicID, &c.Config)
}

// Video creates a new asset.Video instance.
func (c Cloudinary) Video(publicID string) (*asset.Asset, error) {
	return asset.Video(publicID, &c.Config)
}

// File creates a new asset.File instance.
func (c Cloudinary) File(publicID string) (*asset.Asset, error) {
	return asset.File(publicID, &c.Config)
}

// Media creates a new asset.Media instance.
func (c Cloudinary) Media(publicID string) (*asset.Asset, error) {
	return asset.Media(publicID, &c.Config)
}
