package config

// Api Configuration struct
type Api struct {
	UploadPrefix string `default:"https://api.cloudinary.com"`
	Timeout      int    `default:"60"`
	// TODO: add the rest
}
