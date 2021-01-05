package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"net/http"
)

const (
	Sprite  api.EndPoint = "sprite"
	Multi   api.EndPoint = "multi"
	Explode api.EndPoint = "explode"
	Text    api.EndPoint = "text"
)

// GenerateSpriteParams struct
type GenerateSpriteParams struct {
	Tag             string `json:"tag,omitempty"`
	NotificationUrl string `json:"notification_url,omitempty"`
	Async           bool   `json:"async,omitempty"`
	Transformation  string `json:"transformation,omitempty"`
	ResourceType    string `json:"-"`
}

// GenerateSprite creates a sprite from all images that have been assigned a specified tag.
//
// The process produces two files:
// * A single image file containing all the images with the specified tag (PNG by default).
// * A CSS file that includes the style class names and the location of the individual images in the sprite.
//
// https://cloudinary.com/documentation/image_upload_api_reference#sprite_method
func (u *Api) GenerateSprite(ctx context.Context, params GenerateSpriteParams) (*GenerateSpriteResult, error) {
	res := &GenerateSpriteResult{}
	err := u.callUploadApi(ctx, Sprite, params, res)

	return res, err
}

type GenerateSpriteResult struct {
	CssUrl         string               `json:"css_url"`
	ImageUrl       string               `json:"image_url"`
	SecureCssUrl   string               `json:"secure_css_url"`
	SecureImageURL string               `json:"secure_image_url"`
	JsonUrl        string               `json:"json_url"`
	SecureJsonUrl  string               `json:"secure_json_url"`
	Version        int                  `json:"version"`
	PublicID       string               `json:"public_id"`
	ImageInfos     map[string]ImageInfo `json:"image_infos"`
	Error          api.ErrorResp        `json:"error,omitempty"`
	Response       http.Response
}

type ImageInfo struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// MultiParams struct
type MultiParams struct {
	Tag             string `json:"tag,omitempty"`
	Format          string `json:"format,omitempty"`
	NotificationUrl string `json:"notification_url,omitempty"`
	Async           bool   `json:"async,omitempty"`
	Transformation  string `json:"transformation,omitempty"`
	ResourceType    string `json:"-"`
}

// Multi Creates a single animated image, video or PDF from all image assets that have been assigned a specified tag.
//
// https://cloudinary.com/documentation/image_upload_api_reference#multi_method
func (u *Api) Multi(ctx context.Context, params MultiParams) (*MultiResult, error) {
	res := &MultiResult{}
	err := u.callUploadApi(ctx, Multi, params, res)

	return res, err
}

type MultiResult struct {
	Url       string        `json:"url"`
	SecureUrl string        `json:"secure_url"`
	AssetID   string        `json:"asset_id"`
	PublicID  string        `json:"public_id"`
	Version   int           `json:"version"`
	Error     api.ErrorResp `json:"error,omitempty"`
	Response  http.Response
}

// ExplodeParams struct
type ExplodeParams struct {
	PublicID        string `json:"public_id"`
	Format          string `json:"format,omitempty"`
	Type            string `json:"type,omitempty"`
	NotificationUrl string `json:"notification_url,omitempty"`
	Transformation  string `json:"transformation,omitempty"`
	ResourceType    string `json:"-"`
}

// Explode creates derived images for all of the individual pages in a multi-page file (PDF or animated GIF).
//
// Each derived image is stored with the same public ID as the original file, and can be accessed using the page
// parameter, in order to deliver a specific image.
//
//https://cloudinary.com/documentation/image_upload_api_reference#explode_method
func (u *Api) Explode(ctx context.Context, params ExplodeParams) (*ExplodeResult, error) {
	params.Transformation = "pg_all" // Transformation must contain exactly one "pg_all" transformation parameter
	res := &ExplodeResult{}
	err := u.callUploadApi(ctx, Explode, params, res)

	return res, err
}

type ExplodeResult struct {
	Status   string        `json:"status"`
	BatchID  string        `json:"batch_id"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response http.Response
}
