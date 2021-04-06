package config

// API defines the configuration for making requests to the Cloudinary API.
type API struct {
	UploadPrefix  string `default:"https://api.cloudinary.com"`
	Timeout       int64  `default:"60"` // seconds
	UploadTimeout int64
	ChunkSize     int64 `default:"20000000"` //bytes
}
