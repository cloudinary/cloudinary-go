package config

import "time"

// Api defines the configuration when making requests to the Cloudinary API.
type Api struct {
	UploadPrefix string        `default:"https://api.cloudinary.com"`
	Timeout      time.Duration `default:"60"` // seconds
}
