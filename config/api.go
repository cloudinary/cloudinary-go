package config

// Api defines the configuration for making requests to the Cloudinary API.
type Api struct {
	UploadPrefix string `default:"https://api.cloudinary.com"`
	Timeout      int64  `default:"60"` // seconds
}