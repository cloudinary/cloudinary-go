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
func Create() *Cloudinary {
	return CreateFromConfiguration(*config.Create())
}

// CreateFromUrl is creating a new Cloudinary instance from a cloudinary url
func CreateFromUrl(cloudinaryUrl string) *Cloudinary {
	return CreateFromConfiguration(*config.CreateFromUrl(cloudinaryUrl))
}

// CreateFromParams is creating a new Cloudinary instance from provided parameters
func CreateFromParams(cloud string, key string, secret string) *Cloudinary {
	return CreateFromConfiguration(*config.CreateFromParams(cloud, key, secret))
}

// CreateFromParams is creating a new Cloudinary instance from provided configuration
func CreateFromConfiguration(configuration config.Configuration) *Cloudinary {
	return &Cloudinary{
		Config: configuration,
		Admin: admin.Api{
			Config: configuration,
		},
		Upload: uploader.Api{
			Config: configuration,
		},
	}
}
