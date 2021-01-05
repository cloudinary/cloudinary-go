package admin

// Enables you to manage the assets in your account or cloud.
//
// https://cloudinary.com/documentation/admin_api#resources

import (
	"cloudinary/cloudinary-go/pkg/api"
	"context"
	"time"
)

const (
	Assets        api.EndPoint = "resources"
	DerivedAssets api.EndPoint = "derived_resources"
	Tags          api.EndPoint = "tags"
	Context       api.EndPoint = "context"
	Moderations   api.EndPoint = "moderations"
	Restore       api.EndPoint = "restore"
)

// AssetTypes lists available asset types.
func (a *Api) AssetTypes(ctx context.Context) (*AssetTypesResult, error) {
	res := &AssetTypesResult{}
	_, err := a.get(ctx, Assets, nil, res)

	return res, err
}

type AssetTypesResult struct {
	AssetTypes []string      `json:"resource_types"`
	Error      api.ErrorResp `json:"error,omitempty"`
}

type AssetsParams struct {
	AssetType   api.AssetType `json:"-"`
	Prefix      string        `json:"prefix,omitempty"`
	StartAt     *time.Time    `json:"start_at,omitempty"`
	NextCursor  string        `json:"next_cursor,omitempty"`
	MaxResults  int           `json:"max_results,omitempty"`
	Tags        bool          `json:"tags,omitempty"`
	Context     bool          `json:"context,omitempty"`
	Moderations bool          `json:"moderations,omitempty"`
	Direction   string        `json:"direction,omitempty"`
}

// Assets lists all uploaded assets filtered by any specified AssetsParams.
//
//https://cloudinary.com/documentation/admin_api#get_resources
func (a *Api) Assets(ctx context.Context, params AssetsParams) (*AssetsResult, error) {
	res := &AssetsResult{}
	_, err := a.get(ctx, api.BuildPath(Assets, params.AssetType), params, res)

	return res, err
}

type AssetsResult struct {
	Assets     []api.BriefAssetResult `json:"resources"`
	NextCursor string                 `json:"next_cursor"`
	Error      api.ErrorResp          `json:"error,omitempty"`
}

type AssetsByTagParams struct {
	AssetType   api.AssetType `json:"-"`
	Tag         string        `json:"-"`
	NextCursor  string        `json:"next_cursor,omitempty"`
	MaxResults  int           `json:"max_results,omitempty"`
	Tags        bool          `json:"tags,omitempty"`
	Context     bool          `json:"context,omitempty"`
	Moderations bool          `json:"moderations,omitempty"`
	Direction   string        `json:"direction,omitempty"`
}

// AssetsByTag lists assets with the specified tag.
//
// This method does not return matching deleted assets, even if they have been backed up.
//
// https://cloudinary.com/documentation/admin_api#get_resources_by_tag
func (a *Api) AssetsByTag(ctx context.Context, params AssetsByTagParams) (*AssetsResult, error) {
	res := &AssetsResult{}
	_, err := a.get(ctx, api.BuildPath(Assets, params.AssetType, Tags, params.Tag), params, res)

	return res, err
}

type AssetsByContextParams struct {
	AssetType   api.AssetType `json:"-"`
	Key         string        `json:"key"`
	Value       string        `json:"value,omitempty"`
	NextCursor  string        `json:"next_cursor,omitempty"`
	MaxResults  int           `json:"max_results,omitempty"`
	Tags        bool          `json:"tags,omitempty"`
	Context     bool          `json:"context,omitempty"`
	Moderations bool          `json:"moderations,omitempty"`
	Direction   string        `json:"direction,omitempty"`
}

// AssetsByContext sists assets with the specified contextual metadata.
//
// This method does not return matching deleted assets, even if they have been backed up.
//
// https://cloudinary.com/documentation/admin_api#get_resources_by_context
func (a *Api) AssetsByContext(ctx context.Context, params AssetsByContextParams) (*AssetsResult, error) {
	res := &AssetsResult{}
	_, err := a.get(ctx, api.BuildPath(Assets, params.AssetType, Context), params, res)

	return res, err
}

type AssetsByModerationParams struct {
	AssetType   api.AssetType `json:"-"`
	Kind        string        `json:"-"`
	Status      string        `json:"-"`
	NextCursor  string        `json:"next_cursor,omitempty"`
	MaxResults  int           `json:"max_results,omitempty"`
	Tags        bool          `json:"tags,omitempty"`
	Context     bool          `json:"context,omitempty"`
	Moderations bool          `json:"moderations,omitempty"`
	Direction   string        `json:"direction,omitempty"`
}

// AssetsByModeration lists assets currently in the specified moderation queue and status.
//
// https://cloudinary.com/documentation/admin_api#get_resources_in_moderation_queues
func (a *Api) AssetsByModeration(ctx context.Context, params AssetsByModerationParams) (*AssetsResult, error) {
	res := &AssetsResult{}
	_, err := a.get(ctx, api.BuildPath(Assets, params.AssetType, Moderations, params.Kind, params.Status), params, res)

	return res, err
}

type AssetsByIDsParams struct {
	AssetType    api.AssetType    `json:"-"`
	DeliveryType api.DeliveryType `json:"-"`
	PublicIDs    api.CldApiArray  `json:"public_ids"`
	Tags         bool             `json:"tags,omitempty"`
	Context      bool             `json:"context,omitempty"`
	Moderations  bool             `json:"moderations,omitempty"`
}

// AssetsByIDs lists assets with the specified public IDs.
//
// https://cloudinary.com/documentation/admin_api#get_resources
func (a *Api) AssetsByIDs(ctx context.Context, params AssetsByIDsParams) (*AssetsResult, error) {
	res := &AssetsResult{}
	_, err := a.get(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType), params, res)

	return res, err
}

type RestoreAssetsParams struct {
	AssetType    api.AssetType    `json:"-"`
	DeliveryType api.DeliveryType `json:"-"`
	PublicIDs    api.CldApiArray  `json:"public_ids"`
	Versions     api.CldApiArray  `json:"versions"`
}

// RestoreAssets reverts to the latest backed up version of the specified deleted assets.
//
// https://cloudinary.com/documentation/admin_api#restore_resources
func (a *Api) RestoreAssets(ctx context.Context, params RestoreAssetsParams) (*RestoreAssetsResult, error) {
	res := &RestoreAssetsResult{}
	_, err := a.post(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType, Restore), params, res)

	return res, err
}

type RestoreAssetsResult map[string]api.BriefAssetResult

type DeleteAssetsParams struct {
	AssetType       api.AssetType    `json:"-"`
	DeliveryType    api.DeliveryType `json:"-"`
	PublicIDs       api.CldApiArray  `json:"public_ids"` // The public IDs of the assets to delete (up to 100).
	KeepOriginal    bool             `json:"keep_original,omitempty"`
	Invalidate      bool             `json:"invalidate,omitempty"`
	Transformations string           `json:"transformations,omitempty"`
	NextCursor      string           `json:"next_cursor,omitempty"`
}

// DeleteAssets deletes the specified assets.
//
// https://cloudinary.com/documentation/admin_api#delete_resources
func (a *Api) DeleteAssets(ctx context.Context, params DeleteAssetsParams) (*DeleteAssetsResult, error) {
	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType), params, res)

	return res, err
}

type DeleteAssetsResult struct {
	Deleted       map[string]string      `json:"deleted"`
	DeletedCounts map[string]interface{} `json:"deleted_counts"`
	Partial       bool                   `json:"partial"`
	Error         api.ErrorResp          `json:"error,omitempty"`
}

type DeleteAssetsByPrefixParams struct {
	AssetType       api.AssetType    `json:"-"`
	DeliveryType    api.DeliveryType `json:"-"`
	Prefix          api.CldApiArray  `json:"prefix"`
	KeepOriginal    bool             `json:"keep_original,omitempty"`
	Invalidate      bool             `json:"invalidate,omitempty"`
	Transformations string           `json:"transformations,omitempty"`
	NextCursor      string           `json:"next_cursor,omitempty"`
}

// DeleteAssetsByPrefix deletes assets by prefix.
//
// https://cloudinary.com/documentation/admin_api#delete_resources
func (a *Api) DeleteAssetsByPrefix(ctx context.Context, params DeleteAssetsByPrefixParams) (*DeleteAssetsResult, error) {
	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType), params, res)

	return res, err
}

type DeleteAssetsByTagParams struct {
	AssetType       api.AssetType `json:"-"`
	Tag             string        `json:"-"`
	KeepOriginal    bool          `json:"keep_original,omitempty"`
	Invalidate      bool          `json:"invalidate,omitempty"`
	Transformations string        `json:"transformations,omitempty"`
	NextCursor      string        `json:"next_cursor,omitempty"`
}

// DeleteAssetsByTag deletes assets with the specified tag, including their derived resources.
//
// Supports deleting up to a maximum of 1000 original assets in a single call.
//
// https://cloudinary.com/documentation/admin_api#delete_resources_by_tags
func (a *Api) DeleteAssetsByTag(ctx context.Context, params DeleteAssetsByTagParams) (*DeleteAssetsResult, error) {
	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(Assets, params.AssetType, Tags, params.Tag), params, res)

	return res, err
}

type DeleteAllAssetsParams struct {
	AssetType       api.AssetType    `json:"-"`
	DeliveryType    api.DeliveryType `json:"-"`
	All             bool             `json:"all"`
	KeepOriginal    bool             `json:"keep_original,omitempty"`
	Invalidate      bool             `json:"invalidate,omitempty"`
	Transformations string           `json:"transformations,omitempty"`
	NextCursor      string           `json:"next_cursor,omitempty"`
}

// DeleteAllAssets deletes all assets of the specified asset and delivery type, including their derived resources.
//
// Supports deleting up to a maximum of 1000 original assets in a single call.
//
// https://cloudinary.com/documentation/admin_api#delete_resources
func (a *Api) DeleteAllAssets(ctx context.Context, params DeleteAllAssetsParams) (*DeleteAssetsResult, error) {
	params.All = true

	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType), params, res)

	return res, err
}

type DeleteDerivedAssetsParams struct {
	DerivedAssetIDs api.CldApiArray `json:"derived_resource_ids"`
}

// DeleteDerivedAssets deletes the specified derived resources by derived resource ID.
//
// The derived resource IDs for a particular original asset are returned when calling the `resource` method to return
// the details of a single asset.
//
// https://cloudinary.com/documentation/admin_api##delete_resources
func (a *Api) DeleteDerivedAssets(ctx context.Context, params DeleteDerivedAssetsParams) (*DeleteAssetsResult, error) {
	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(DerivedAssets), params, res)

	return res, err
}

type DeleteDerivedAssetsByTransformationParams struct {
	AssetType       api.AssetType    `json:"-"`
	DeliveryType    api.DeliveryType `json:"-"`
	PublicIDs       api.CldApiArray  `json:"public_ids"`      // The public IDs for which you want to delete derived resources.
	Transformations string           `json:"transformations"` // The transformation(s) associated with the derived resources to delete.
	KeepOriginal    bool             `json:"keep_original"`
	Invalidate      bool             `json:"invalidate,omitempty"`
}

// DeleteDerivedAssetsByTransformation deletes derived resources identified by transformation and public_ids.
func (a *Api) DeleteDerivedAssetsByTransformation(ctx context.Context, params DeleteDerivedAssetsByTransformationParams) (*DeleteAssetsResult, error) {
	params.KeepOriginal = true

	res := &DeleteAssetsResult{}
	_, err := a.delete(ctx, api.BuildPath(Assets, params.AssetType, params.DeliveryType), params, res)

	return res, err
}
