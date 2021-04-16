package config

// AuthToken defines the configuration for delivering token-based authenticated media assets.
//
// https://cloudinary.com/documentation/control_access_to_media#delivering_token_based_authenticated_media_assets
type AuthToken struct {
	Key        string
	IP         string
	ACL        string
	StartTime  uint64
	Expiration uint64
	Duration   uint64
}

func (a AuthToken) isEnabled() bool {
	return a.Key != ""
}
