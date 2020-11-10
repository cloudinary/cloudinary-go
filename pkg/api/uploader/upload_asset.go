package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"encoding/json"
	"time"
)

// UploadParams struct
// See http://cloudinary.com/documentation/image_upload_api_reference#api_example_1
type UploadParams struct {
	PublicID                string          `json:"public_id,omitempty"`
	PublicIds               api.CldApiArray `json:"public_ids,omitempty"`
	UseFilename             bool            `json:"use_filename,omitempty"`
	UniqueFilename          bool            `json:"unique_filename,omitempty"`
	Folder                  string          `json:"folder,omitempty"`
	Overwrite               bool            `json:"overwrite,omitempty"`
	ResourceType            string          `json:"resource_type,omitempty"`
	Type                    string          `json:"type,omitempty"`
	Tags                    api.CldApiArray `json:"tags,omitempty"`
	Context                 api.Context     `json:"context,omitempty"`
	Metadata                api.Metadata    `json:"metadata,omitempty"`
	Transformation          string          `json:"transformation,omitempty"`
	Format                  string          `json:"format,omitempty"`
	AllowedFormats          api.CldApiArray `json:"allowed_formats,omitempty"`
	Eager                   string          `json:"eager,omitempty"`
	Eval                    string          `json:"eval,omitempty"`
	Async                   bool            `json:"async,omitempty"`
	EagerAsync              bool            `json:"eager_async,omitempty"`
	Proxy                   string          `json:"proxy,omitempty"`
	Headers                 string          `json:"headers,omitempty"`
	Callback                string          `json:"callback,omitempty"`
	NotificationURL         string          `json:"notification_url,omitempty"`
	EagerNotificationURL    string          `json:"eager_notification_url,omitempty"`
	Faces                   bool            `json:"faces,omitempty"`
	ImageMetadata           bool            `json:"image_metadata,omitempty"`
	Exif                    bool            `json:"exif,omitempty"`
	Colors                  bool            `json:"colors,omitempty"`
	Phash                   bool            `json:"phash,omitempty"`
	FaceCoordinates         string          `json:"face_coordinates,omitempty"`
	CustomCoordinates       string          `json:"custom_coordinates,omitempty"`
	Backup                  bool            `json:"backup,omitempty"`
	ReturnDeleteToken       bool            `json:"return_delete_token,omitempty"`
	Invalidate              bool            `json:"invalidate,omitempty"`
	DiscardOriginalFilename bool            `json:"discard_original_filename,omitempty"`
	Moderation              string          `json:"moderation,omitempty"`
	UploadPreset            string          `json:"upload_preset,omitempty"`
	RawConvert              string          `json:"raw_convert,omitempty"`
	Categorization          string          `json:"categorization,omitempty"`
	AutoTagging             float64         `json:"auto_tagging,omitempty"`
	BackgroundRemoval       string          `json:"background_removal,omitempty"`
	Detection               string          `json:"detection,omitempty"`
	OCR                     string          `json:"ocr,omitempty"`
	Timestamp               time.Time       `json:"timestamp,omitempty"`
	QualityAnalysis         bool            `json:"quality_analysis,omitempty"`
	AccessibilityAnalysis   bool            `json:"accessibility_analysis,omitempty"`
	CinemagraphAnalysis     bool            `json:"cinemagraph_analysis,omitempty"`
}

// UploadResult image success response struct
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
	Tags             api.CldApiArray `json:"tags,omitempty"`
	Pages            int             `json:"pages,omitempty"`
	Bytes            int             `json:"bytes"`
	Type             string          `json:"type"`
	Etag             string          `json:"etag"`
	Placeholder      bool            `json:"placeholder,omitempty"`
	URL              string          `json:"url"`
	SecureURL        string          `json:"secure_url"`
	AccessMode       string          `json:"access_mode"`
	Context          api.Context     `json:"context,omitempty"`
	Metadata         api.Metadata    `json:"metadata,omitempty"`
	Overwritten      bool            `json:"overwritten"`
	OriginalFilename string          `json:"original_filename"`
	Error            api.ErrorResp   `json:"error,omitempty"`
}

// UploadResult is uploading an image
func (u *Api) Upload(file string, uploadParams ...UploadParams) (*UploadResult, error) {
	formParams, err := api.StructToParams(uploadParams[0])
	if err != nil {
		return nil, err
	}

	body := u.postFile(file, formParams)

	upload := &UploadResult{}
	err = json.Unmarshal(body, upload)

	if err != nil {
		return nil, err
	}

	return upload, nil
}
