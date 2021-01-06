package config

// Api defines the configuration when making requests to the Cloudinary API.
type Api struct {
	UploadPrefix string `default:"https://api.cloudinary.com"`
	Timeout      int64  `default:"60"` // seconds
}
