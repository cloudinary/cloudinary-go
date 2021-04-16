package config

import "github.com/cloudinary/cloudinary-go/internal/signature"

// URL defines the configuration applied when generating Cloudinary URLs.
//
// https://cloudinary.com/documentation/how_to_integrate_cloudinary#get_familiar_with_the_cloudinary_console
type URL struct {
	Domain             string `default:"cloudinary.com"`
	SubDomain          string `default:"res"`
	SharedHost         string `default:"res.cloudinary.com"`
	CName              string
	SecureCName        string
	Secure             bool `default:"true"`
	CDNSubDomain       bool
	SecureCDNSubDomain bool
	PrivateCDN         bool
	SignURL            bool
	LongURLSignature   bool
	Shorten            bool
	UseRootPath        bool
	ForceVersion       bool `default:"true"`
	Analytics          bool `default:"true"`
}

// Protocol returns URL protocol (http or https).
func (uc URL) Protocol() string {
	if uc.Secure {
		return "https"
	}

	return "http"
}

// GetSignatureLength returns the length of the URL signature.
func (uc URL) GetSignatureLength() uint8 {
	if uc.LongURLSignature {
		return signature.Long
	}

	return signature.Short
}
