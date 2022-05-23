package config

import "github.com/cloudinary/cloudinary-go/v2/internal/signature"

// Cloud defines the cloud configuration required to connect your application to Cloudinary.
//
// https://cloudinary.com/documentation/how_to_integrate_cloudinary#get_familiar_with_the_cloudinary_console
type Cloud struct {
	CloudName          string `schema:"-"`
	APIKey             string `schema:"-"`
	APISecret          string `schema:"-"`
	OAuthToken         string `schema:"oauth_token"`
	SignatureAlgorithm string `schema:"signature_algorithm"`
}

// GetSignatureAlgorithm returns the signature algorithm.
func (c Cloud) GetSignatureAlgorithm() string {
	if c.SignatureAlgorithm == "" {
		return signature.SHA1
	}

	return c.SignatureAlgorithm
}
