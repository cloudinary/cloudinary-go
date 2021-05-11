package config

// AuthToken defines the configuration for delivering token-based authenticated media assets.
//
// https://cloudinary.com/documentation/control_access_to_media#delivering_token_based_authenticated_media_assets
type AuthToken struct {
	Key        string
	IP         string
	ACL        string
	StartTime  int64
	Expiration int64
	Duration   int64
}
