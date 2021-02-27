package admin_test

import (
	"log"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/admin/metadata"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
)

var metadataField = metadata.Field{
	Type:         metadata.SetFieldType,
	ExternalID:   cldtest.UniqueID("go_color_id_"),
	Label:        cldtest.UniqueID("GoColors"),
	Mandatory:    true,
	DefaultValue: []string{"go_color1", "go_color2"},
	DataSource:   dataSource1,
}

var dataSource1 = metadata.DataSource{
	Values: []metadata.DataSourceValue{
		{
			ExternalID: "go_color1",
			Value:      "red",
			State:      "active",
		},
		{
			ExternalID: "go_color2",
			Value:      "green",
			State:      "active",
		},
	},
}

var dataSource2 = metadata.DataSource{
	Values: []metadata.DataSourceValue{
		{
			ExternalID: "go_color3",
			Value:      "blue",
			State:      "active",
		},
		{
			ExternalID: "go_color4",
			Value:      "yellow",
			State:      "active",
		},
	},
}

func TestAdmin_AddMetadataField(t *testing.T) {
	resp, err := adminAPI.AddMetadataField(ctx, metadataField)

	if err != nil {
		t.Error(err)
	}

	if resp.Error.Message == "external id "+metadataField.ExternalID+" already exists" {
		t.Skip(resp.Error.Message)
	}

	if resp.ExternalID != metadataField.ExternalID {
		t.Error(resp)
	}
}

func TestAdmin_UpdateMetadataField(t *testing.T) {
	metadataField.Label = cldtest.UniqueID("GoUpdatedColors")

	resp, err := adminAPI.UpdateMetadataField(ctx, admin.UpdateMetadataFieldParams{
		FieldExternalID: metadataField.ExternalID,
		Field:           metadataField,
	})

	if err != nil {
		t.Error(err)
	}

	if resp.Label != metadataField.Label {
		t.Error(resp)
	}
}

func TestAdmin_ListMetadataFields(t *testing.T) {
	resp, err := adminAPI.ListMetadataFields(ctx)

	if err != nil || len(resp.MetadataFields) < 1 {
		t.Error(resp)
	}
}

func TestAdmin_MetadataFieldByFieldID(t *testing.T) {
	params := admin.MetadataFieldByFieldIDParams{FieldExternalID: metadataField.ExternalID}
	resp, err := adminAPI.MetadataFieldByFieldID(ctx, params)

	if err != nil || resp.ExternalID != metadataField.ExternalID {
		t.Error(err, resp)
	}
}

func TestAdmin_UpdateMetadataFieldDataSource(t *testing.T) {
	resp, err := adminAPI.UpdateMetadataFieldDataSource(ctx, admin.UpdateMetadataFieldDataSourceParams{
		FieldExternalID: metadataField.ExternalID,
		DataSource:      dataSource2,
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_DeleteDataSourceEntries(t *testing.T) {
	resp, err := adminAPI.DeleteDataSourceEntries(ctx, admin.DeleteDataSourceEntriesParams{
		FieldExternalID:    metadataField.ExternalID,
		EntriesExternalIDs: []string{"go_color3", "go_color4"},
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_RestoreMetadataFieldDataSource(t *testing.T) {
	resp, err := adminAPI.RestoreDatasourceEntries(ctx, admin.RestoreDatasourceEntriesParams{
		FieldExternalID:    metadataField.ExternalID,
		EntriesExternalIDs: []string{"go_color3", "go_color4"},
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_SortMetadataFieldsDatasource(t *testing.T) {
	resp, err := adminAPI.SortMetadataFieldDatasource(ctx, admin.SortMetadataFieldDatasourceParams{FieldExternalId: metadataField.ExternalID, FieldSortBy: "value", FieldDirection: admin.Ascending})

	if err != nil {
		t.Error(err, resp)
	}

	if resp.Values[0].Value != dataSource2.Values[0].Value {
		t.Error("Wrong response. Metadata fields should be sorted in ascending order")
	}

	resp, err = adminAPI.SortMetadataFieldDatasource(ctx, admin.SortMetadataFieldDatasourceParams{FieldExternalId: metadataField.ExternalID, FieldSortBy: "value", FieldDirection: admin.Descending})

	if err != nil {
		t.Error(err, resp)
	}

	if resp.Values[0].Value != dataSource2.Values[1].Value {
		t.Error("Wrong response. Metadata fields should be sorted in descending order")
	}
}

func TestAdmin_DeleteMetadataField(t *testing.T) {
	resp, err := adminAPI.DeleteMetadataField(ctx, admin.DeleteMetadataFieldParams{FieldExternalID: metadataField.ExternalID})

	if err != nil || resp.Message != "ok" {
		t.Error(err, resp)
	}
}

var mdIDs = map[string]string{
	"enum": cldtest.UniqueID("go_distinct_color_id_"),
	"int":  cldtest.UniqueID("go_17_integer_id_"),
	"str":  cldtest.UniqueID("go_string_id_"),
	"date": cldtest.UniqueID("go_date_id_"),
}

func TestAdmin_AddMetadataFields(t *testing.T) {
	var integerMetadataField = metadata.Field{
		Type:       metadata.IntegerFieldType,
		ExternalID: mdIDs["int"],
		Label:      cldtest.UniqueID("Go17Integer"),
		Validation: metadata.AndValidation(
			[]interface{}{
				metadata.GreaterThanValidation(17, true),
				metadata.LessThanValidation(17, true),
			}),
	}

	var stringMetadataField = metadata.Field{
		Type:         metadata.StringFieldType,
		ExternalID:   mdIDs["str"],
		Label:        cldtest.UniqueID("GoString"),
		DefaultValue: "Gopher",
		Validation:   metadata.StringLengthValidation(2, 6),
	}

	var dateMetadataField = metadata.Field{
		Type:         metadata.DateFieldType,
		ExternalID:   mdIDs["date"],
		Label:        cldtest.UniqueID("GoDate"),
		DefaultValue: time.Now().Format("2006-01-02"),
		Validation:   metadata.GreaterThanValidation(time.Now().AddDate(0, 0, -1), false),
	}

	var enumMetadataField = metadata.Field{
		Type:       metadata.EnumFieldType,
		ExternalID: mdIDs["enum"],
		Label:      cldtest.UniqueID("GoDistinctColors"),
		DataSource: dataSource1,
	}

	for _, f := range []metadata.Field{
		integerMetadataField,
		stringMetadataField,
		dateMetadataField,
		enumMetadataField,
	} {
		resp, err := adminAPI.AddMetadataField(ctx, f)

		if err != nil {
			t.Error(err)
		}

		if resp.Error.Message == "external id "+f.ExternalID+" already exists" {
			t.Skip(resp.Error.Message)
		}

		if resp.ExternalID != f.ExternalID {
			t.Error(resp)
		}
	}
}

//FIXME; find a good library with a proper TearDown method
func TestAdmin_MetadataFieldsCleanup(t *testing.T) {
	for _, extID := range mdIDs {
		resp, err := adminAPI.DeleteMetadataField(ctx, admin.DeleteMetadataFieldParams{FieldExternalID: extID})
		if err != nil || resp.Message != "ok" {
			log.Println(err, resp)
		}
	}
}
