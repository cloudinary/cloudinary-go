package config

import "github.com/cloudinary/cloudinary-go/internal/signature"

// Cloud defines the cloud configuration required to connect your application to Cloudinary.
//
// https://cloudinary.com/documentation/how_to_integrate_cloudinary#get_familiar_with_the_cloudinary_console
type Cloud struct {
	CloudName          string
	APIKey             string
	APISecret          string
	SignatureAlgorithm string
}
// GetSignatureAlgorithm returns the signature algorithm.
func (c Cloud) GetSignatureAlgorithm() string {
	if c.SignatureAlgorithm == "" {
		return signature.SHA1
	}

	return c.SignatureAlgorithm
}
