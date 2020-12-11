package uploader

import (
	"cloudinary-labs/cloudinary-go/pkg/api"
	"context"
	"net/http"
	"net/url"
	"time"
)

const (
	GenerateArchive api.EndPoint = "generate_archive"
)

type ArchiveFormat string

const (
	Zip ArchiveFormat = "zip"
	Tgz ArchiveFormat = "tgz"
)

type ArchiveMode string

const (
	CreateArchive   ArchiveMode = "create"
	DownloadArchive ArchiveMode = "download"
)

// GenerateSpriteParams struct
type CreateArchiveParams struct {
	AllowMissing            bool             `json:"allow_mising,omitempty"`
	Async                   bool             `json:"async,omitempty"`
	ExpiresAt               *time.Time       `json:"expires_at,omitempty"`
	FlattenFolders          bool             `json:"flatten_folders,omitempty"`
	FlattenTransformations  bool             `json:"flatten_transformations,omitempty"`
	FullyQualifiedPublicIds api.CldApiArray  `json:"fully_qualified_public_ids,omitempty"`
	KeepDerived             bool             `json:"keep_derived,omitempty"`
	Mode                    ArchiveMode      `json:"mode,omitempty"`
	NotificationUrl         string           `json:"notification_url,omitempty"`
	Phash                   string           `json:"phash,omitempty"`
	Prefixes                api.CldApiArray  `json:"prefixes,omitempty"`
	PublicIds               api.CldApiArray  `json:"public_ids,omitempty"`
	ResourceType            api.AssetType    `json:"-"`
	SkipTransformationName  bool             `json:"skip_transformation_name,omitempty"`
	TargetFormat            ArchiveFormat    `json:"target_format,omitempty"`
	TargetPublicId          string           `json:"target_public_id,omitempty"`
	TargetTags              api.CldApiArray  `json:"target_tags,omitempty"`
	Tags                    api.CldApiArray  `json:"tags,omitempty"`
	Transformations         string           `json:"transformations,omitempty"`
	Type                    api.DeliveryType `json:"type,omitempty"`
	UseOriginalFilename     bool             `json:"use_original_filename,omitempty"`
}

// Creates a new archive in the server and returns information in JSON format.
func (u *Api) CreateArchive(ctx context.Context, params CreateArchiveParams) (*CreateArchiveResult, error) {
	res := &CreateArchiveResult{}
	err := u.callUploadApi(ctx, GenerateArchive, params, res)

	return res, err
}

type CreateArchiveResult struct {
	AssetID       string        `json:"asset_id"`
	PublicID      string        `json:"public_id"`
	Version       int           `json:"version"`
	VersionID     string        `json:"version_id"`
	Signature     string        `json:"signature"`
	ResourceType  string        `json:"resource_type"`
	CreatedAt     time.Time     `json:"created_at"`
	Tags          []string      `json:"tags"`
	Bytes         int           `json:"bytes"`
	Type          string        `json:"type"`
	Etag          string        `json:"etag"`
	Placeholder   bool          `json:"placeholder"`
	URL           string        `json:"url"`
	SecureURL     string        `json:"secure_url"`
	AccessMode    string        `json:"access_mode"`
	ResourceCount int           `json:"resource_count"`
	FileCount     int           `json:"file_count"`
	Error         api.ErrorResp `json:"error,omitempty"`
	Response      http.Response
}

// Creates a new zip archive in the server and returns information in JSON format.
func (u *Api) CreateZip(ctx context.Context, params CreateArchiveParams) (*CreateArchiveResult, error) {
	params.TargetFormat = Zip

	return u.CreateArchive(ctx, params)
}

// Returns a URL that when invoked creates an archive and returns it.
func (u *Api) DownloadArchiveUrl(params CreateArchiveParams) (string, error) {
	params.Mode = DownloadArchive

	queryParams, err := api.StructToParams(params)
	if err != nil {
		return "", err
	}
	queryParams, err = u.signRequest(queryParams)
	if err != nil {
		return "", err
	}

	assetType := getAssetType(params)

	archiveEndpointURL := u.getUploadURL(api.BuildPath(assetType, GenerateArchive))

	urlStruct, err := url.Parse(archiveEndpointURL)
	if err != nil {
		return "", err
	}

	urlStruct.RawQuery = queryParams.Encode()

	return urlStruct.String(), nil
}

// Returns a URL that when invokes creates a zip archive and returns it.
func (u *Api) DownloadZipUrl(params CreateArchiveParams) (string, error) {
	params.TargetFormat = Zip

	return u.DownloadArchiveUrl(params)
}

// Creates and returns a URL that when invoked creates an archive of a folder.
func (u *Api) DownloadFolder(folderPath string, params CreateArchiveParams) (string, error) {
	params.Prefixes = api.CldApiArray{folderPath}
	if len(params.ResourceType) == 0 {
		params.ResourceType = api.All
	}

	return u.DownloadArchiveUrl(params)
}
