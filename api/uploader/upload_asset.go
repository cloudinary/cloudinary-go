package uploader

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloudinary/cloudinary-go/api"
)

// UploadParams struct allows to customize upload behaviour.
// For additional information about each one of the parameters please use the link below.
//
// https://cloudinary.com/documentation/image_upload_api_reference#upload_optional_parameters
type UploadParams struct {
	PublicID                string           `json:"public_id,omitempty"`
	PublicIds               api.CldAPIArray  `json:"public_ids,omitempty"`
	UseFilename             bool             `json:"use_filename,omitempty"`
	UniqueFilename          bool             `json:"unique_filename,omitempty"`
	FilenameOverride        string           `json:"filename_override,omitempty"`
	Folder                  string           `json:"folder,omitempty"`
	Overwrite               bool             `json:"overwrite,omitempty"`
	ResourceType            string           `json:"resource_type,omitempty"`
	Type                    api.DeliveryType `json:"type,omitempty"`
	Tags                    api.CldAPIArray  `json:"tags,omitempty"`
	Context                 api.CldAPIMap    `json:"context,omitempty"`
	Metadata                api.Metadata     `json:"metadata,omitempty"`
	Transformation          string           `json:"transformation,omitempty"`
	Format                  string           `json:"format,omitempty"`
	AllowedFormats          api.CldAPIArray  `json:"allowed_formats,omitempty"`
	Eager                   string           `json:"eager,omitempty"`
	Eval                    string           `json:"eval,omitempty"`
	Async                   bool             `json:"async,omitempty"`
	EagerAsync              bool             `json:"eager_async,omitempty"`
	Unsigned                bool             `json:"unsigned,omitempty"`
	Proxy                   string           `json:"proxy,omitempty"`
	Headers                 string           `json:"headers,omitempty"`
	Callback                string           `json:"callback,omitempty"`
	NotificationURL         string           `json:"notification_url,omitempty"`
	EagerNotificationURL    string           `json:"eager_notification_url,omitempty"`
	Faces                   bool             `json:"faces,omitempty"`
	ImageMetadata           bool             `json:"image_metadata,omitempty"`
	Exif                    bool             `json:"exif,omitempty"`
	Colors                  bool             `json:"colors,omitempty"`
	Phash                   bool             `json:"phash,omitempty"`
	FaceCoordinates         api.Coordinates  `json:"face_coordinates,omitempty"`
	CustomCoordinates       api.Coordinates  `json:"custom_coordinates,omitempty"`
	Backup                  bool             `json:"backup,omitempty"`
	ReturnDeleteToken       bool             `json:"return_delete_token,omitempty"`
	Invalidate              bool             `json:"invalidate,omitempty"`
	DiscardOriginalFilename bool             `json:"discard_original_filename,omitempty"`
	Moderation              string           `json:"moderation,omitempty"`
	UploadPreset            string           `json:"upload_preset,omitempty"`
	RawConvert              string           `json:"raw_convert,omitempty"`
	Categorization          string           `json:"categorization,omitempty"`
	AutoTagging             float64          `json:"auto_tagging,omitempty"`
	BackgroundRemoval       string           `json:"background_removal,omitempty"`
	Detection               string           `json:"detection,omitempty"`
	OCR                     string           `json:"ocr,omitempty"`
	Timestamp               time.Time        `json:"timestamp,omitempty"`
	QualityAnalysis         bool             `json:"quality_analysis,omitempty"`
	AccessibilityAnalysis   bool             `json:"accessibility_analysis,omitempty"`
	CinemagraphAnalysis     bool             `json:"cinemagraph_analysis,omitempty"`
}

// Upload uploads an asset to a Cloudinary account.
//
// The asset can be:
//   * a local file path
//   * the actual data (io.Reader)
//   * the Data URI (Base64 encoded), max ~60 MB (62,910,000 chars)
//   * the remote FTP, HTTP or HTTPS URL address of an existing file
//   * a private storage bucket (S3 or Google Storage) URL of a whitelisted bucket
//
// https://cloudinary.com/documentation/image_upload_api_reference#upload_method
func (u *API) Upload(ctx context.Context, file interface{}, uploadParams UploadParams) (*UploadResult, error) {
	formParams, err := api.StructToParams(uploadParams)
	if err != nil {
		return nil, err
	}

	body, err := u.postFile(ctx, file, formParams)

	if err != nil {
		return nil, err
	}
	upload := &UploadResult{}
	err = json.Unmarshal(body, upload)

	if err != nil {
		return nil, err
	}

	return upload, nil
}

// Eager contains information about eagerly transformed derived assets.
type Eager struct {
	Transformation string `json:"transformation"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	Bytes          int    `json:"bytes"`
	Format         string `json:"format"`
	URL            string `json:"url"`
	SecureURL      string `json:"secure_url"`
}

// UploadResult image success response struct.
type UploadResult struct {
	AssetID          string          `json:"asset_id"`
	PublicID         string          `json:"public_id"`
	Version          int             `json:"version"`
	VersionID        string          `json:"version_id"`
	Signature        string          `json:"signature"`
	Width            int             `json:"width,omitempty"`
	Height           int             `json:"height,omitempty"`
	Format           string          `json:"format"`
	ResourceType     string          `json:"resource_type"`
	CreatedAt        time.Time       `json:"created_at"`
	Tags             api.CldAPIArray `json:"tags,omitempty"`
	Pages            int             `json:"pages,omitempty"`
	Bytes            int             `json:"bytes"`
	Type             string          `json:"type"`
	Etag             string          `json:"etag"`
	Placeholder      bool            `json:"placeholder,omitempty"`
	URL              string          `json:"url"`
	SecureURL        string          `json:"secure_url"`
	AccessMode       string          `json:"access_mode"`
	Context          api.Metadata    `json:"context,omitempty"`
	Metadata         api.Metadata    `json:"metadata,omitempty"`
	Overwritten      bool            `json:"overwritten"`
	OriginalFilename string          `json:"original_filename"`
	Eager            []Eager         `json:"eager"`
	Error            api.ErrorResp   `json:"error,omitempty"`
}

// UnsignedUpload uploads an asset to a Cloudinary account.
//
// The upload is not signed so an upload preset is required.
//
// https://cloudinary.com/documentation/image_upload_api_reference#unsigned_upload_syntax
func (u *API) UnsignedUpload(ctx context.Context, file string, uploadPreset string, uploadParams UploadParams) (*UploadResult, error) {
	uploadParams.Unsigned = true
	uploadParams.UploadPreset = uploadPreset

	return u.Upload(ctx, file, uploadParams)
}
