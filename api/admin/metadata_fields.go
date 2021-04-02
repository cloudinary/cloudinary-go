package admin

import (
	"context"
	"net/http"

	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin/metadata"
)

const (
	metadataFields    api.EndPoint = "metadata_fields"
	dataSource        api.EndPoint = "datasource"
	dataSourceRestore api.EndPoint = "datasource_restore"
)

// ListMetadataFields lists all metadata field definitions.
//
// https://cloudinary.com/documentation/admin_api#get_metadata_fields
func (a *Api) ListMetadataFields(ctx context.Context) (*ListMetadataFieldsResult, error) {
	res := &ListMetadataFieldsResult{}
	_, err := a.get(ctx, metadataFields, nil, res)

	return res, err
}

type ListMetadataFieldsResult struct {
	MetadataFields []metadata.Field `json:"metadata_fields"`
	Error          api.ErrorResp    `json:"error,omitempty"`
	Response       http.Response
}

type MetadataFieldByFieldIdParams struct {
	FieldExternalId string `json:"-"`
}

// MetadataFieldByFieldId gets a single metadata field definition by external ID.
//
// https://cloudinary.com/documentation/admin_api#get_a_metadata_field_by_external_id
func (a *Api) MetadataFieldByFieldId(ctx context.Context, params MetadataFieldByFieldIdParams) (*MetadataFieldByFieldIdResult, error) {
	res := &MetadataFieldByFieldIdResult{}
	_, err := a.get(ctx, api.BuildPath(metadataFields, params.FieldExternalId), params, res)

	return res, err
}

type MetadataFieldByFieldIdResult struct {
	metadata.Field
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

// AddMetadataField creates a new metadata field definition.
//
// https://cloudinary.com/documentation/admin_api#create_a_metadata_field
func (a *Api) AddMetadataField(ctx context.Context, params metadata.Field) (*AddMetadataFieldResult, error) {
	res := &AddMetadataFieldResult{}
	_, err := a.post(ctx, metadataFields, params, res)

	return res, err
}

type AddMetadataFieldResult struct {
	metadata.Field
	Error    api.ErrorResp `json:"error,omitempty"`
	Response http.Response
}

type UpdateMetadataFieldParams struct {
	metadata.Field
	FieldExternalId string `json:"-"`
}

// UpdateMetadataField updates a metadata field by external ID.
//
// Updates a metadata field definition (partially, no need to pass the entire object) passed as JSON data.
//
// https://cloudinary.com/documentation/admin_api#update_a_metadata_field_by_external_id
func (a *Api) UpdateMetadataField(ctx context.Context, params UpdateMetadataFieldParams) (*UpdateMetadataFieldResult, error) {
	res := &UpdateMetadataFieldResult{}
	_, err := a.put(ctx, api.BuildPath(metadataFields, params.FieldExternalId), params, res)

	return res, err
}

type UpdateMetadataFieldResult struct {
	metadata.Field
	Error    api.ErrorResp `json:"error,omitempty"`
	Response http.Response
}

type DeleteMetadataFieldParams struct {
	FieldExternalId string `json:"-"`
}

// DeleteMetadataField deletes a metadata field definition by external ID.
//
// The external ID is immutable. Therefore, once deleted, the field's external ID can no longer be used for future purposes.
//
// https://cloudinary.com/documentation/admin_api#delete_a_metadata_field_by_external_id
func (a *Api) DeleteMetadataField(ctx context.Context, params DeleteMetadataFieldParams) (*DeleteMetadataFieldResult, error) {
	res := &DeleteMetadataFieldResult{}
	_, err := a.delete(ctx, api.BuildPath(metadataFields, params.FieldExternalId), params, res)

	return res, err
}

type DeleteMetadataFieldResult struct {
	Message  string        `json:"message"`
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

type DeleteDataSourceEntriesParams struct {
	FieldExternalId    string   `json:"-"`
	EntriesExternalIDs []string `json:"external_ids"`
}

// DeleteDataSourceEntries deletes entries in a metadata single or multi-select field's datasource.
//
// Deletes (blocks) the datasource (list) entries from the specified metadata field definition. Sets the state of
// the entries to inactive. This is a soft delete. The entries still exist in the database and can be reactivated
// using the RestoreMetadataFieldDataSource method.
//
// https://cloudinary.com/documentation/admin_api#delete_entries_in_a_metadata_field_datasource
func (a *Api) DeleteDataSourceEntries(ctx context.Context, params DeleteDataSourceEntriesParams) (*DeleteDataSourceEntriesResult, error) {
	res := &DeleteDataSourceEntriesResult{}
	_, err := a.delete(ctx, api.BuildPath(metadataFields, params.FieldExternalId, dataSource), params, res)

	return res, err
}

type DeleteDataSourceEntriesResult struct {
	metadata.DataSource
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

type UpdateMetadataFieldDataSourceParams struct {
	metadata.DataSource
	FieldExternalId string `json:"-"`
}

// UpdateMetadataFieldDataSource updates a metadata field datasource.
//
// Updates the datasource of a supported field type (currently enum or set), passed as JSON data. The
// update is partial: datasource entries with an existing external_id will be updated and entries with new
// external_id’s (or without external_id’s) will be appended.
//
// https://cloudinary.com/documentation/admin_api#update_a_metadata_field_datasource
func (a *Api) UpdateMetadataFieldDataSource(ctx context.Context, params UpdateMetadataFieldDataSourceParams) (*UpdateMetadataFieldDataSourceResult, error) {
	res := &UpdateMetadataFieldDataSourceResult{}
	_, err := a.put(ctx, api.BuildPath(metadataFields, params.FieldExternalId, dataSource), params, res)

	return res, err
}

type UpdateMetadataFieldDataSourceResult struct {
	metadata.DataSource
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}

type RestoreDatasourceEntriesParams struct {
	FieldExternalId    string   `json:"-"`
	EntriesExternalIDs []string `json:"external_ids"`
}

// RestoreDatasourceEntries restores entries in a metadata field datasource.
//
// Restores (unblocks) any previously deleted datasource entries for a specified metadata field definition.
// Sets the state of the entries to active.
//
// https://cloudinary.com/documentation/admin_api#restore_entries_in_a_metadata_field_datasource
func (a *Api) RestoreDatasourceEntries(ctx context.Context, params RestoreDatasourceEntriesParams) (*RestoreDatasourceEntriesResult, error) {
	res := &RestoreDatasourceEntriesResult{}
	_, err := a.post(ctx, api.BuildPath(metadataFields, params.FieldExternalId, dataSourceRestore), params, res)

	return res, err
}

type RestoreDatasourceEntriesResult struct {
	metadata.DataSource
	Error    api.ErrorResp `json:"error,omitempty"`
	Response interface{}
}
