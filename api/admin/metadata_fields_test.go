package admin

import (
	"log"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/api/admin/metadata"
)

var metadataField = metadata.Field{
	Type:         metadata.SetFieldType,
	ExternalID:   "go_color_id_" + testSuffix,
	Label:        "GoColors" + testSuffix,
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
	resp, err := adminApi.AddMetadataField(ctx, metadataField)

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
	metadataField.Label = "GoUpdatedColors" + testSuffix

	resp, err := adminApi.UpdateMetadataField(ctx, UpdateMetadataFieldParams{
		FieldExternalId: metadataField.ExternalID,
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
	resp, err := adminApi.ListMetadataFields(ctx)

	if err != nil || len(resp.MetadataFields) < 1 {
		t.Error(resp)
	}
}

func TestAdmin_MetadataFieldByFieldId(t *testing.T) {
	params := MetadataFieldByFieldIdParams{FieldExternalId: metadataField.ExternalID}
	resp, err := adminApi.MetadataFieldByFieldId(ctx, params)

	if err != nil || resp.ExternalID != metadataField.ExternalID {
		t.Error(err, resp)
	}
}

func TestAdmin_UpdateMetadataFieldDataSource(t *testing.T) {
	resp, err := adminApi.UpdateMetadataFieldDataSource(ctx, UpdateMetadataFieldDataSourceParams{
		FieldExternalId: metadataField.ExternalID,
		DataSource:      dataSource2,
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_DeleteDataSourceEntries(t *testing.T) {
	resp, err := adminApi.DeleteDataSourceEntries(ctx, DeleteDataSourceEntriesParams{
		FieldExternalId:    metadataField.ExternalID,
		EntriesExternalIDs: []string{"go_color3", "go_color4"},
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_RestoreMetadataFieldDataSource(t *testing.T) {
	resp, err := adminApi.RestoreDatasourceEntries(ctx, RestoreDatasourceEntriesParams{
		FieldExternalId:    metadataField.ExternalID,
		EntriesExternalIDs: []string{"go_color3", "go_color4"},
	})

	if err != nil || len(resp.Values) < 2 {
		t.Error(err, resp)
	}
}

func TestAdmin_DeleteMetadataField(t *testing.T) {
	resp, err := adminApi.DeleteMetadataField(ctx, DeleteMetadataFieldParams{FieldExternalId: metadataField.ExternalID})

	if err != nil || resp.Message != "ok" {
		t.Error(err, resp)
	}
}

var mdIDs = map[string]string{
	"enum": "go_distinct_color_id_" + testSuffix,
	"int":  "go_17_integer_id_" + testSuffix,
	"str":  "go_string_id_" + testSuffix,
	"date": "go_date_id_" + testSuffix,
}

func TestAdmin_AddMetadataFields(t *testing.T) {
	var integerMetadataField = metadata.Field{
		Type:       metadata.IntegerFieldType,
		ExternalID: mdIDs["int"],
		Label:      "Go17Integer" + testSuffix,
		Validation: metadata.AndValidation(
			[]interface{}{
				metadata.GreaterThanValidation(17, true),
				metadata.LessThanValidation(17, true),
			}),
	}

	var stringMetadataField = metadata.Field{
		Type:         metadata.StringFieldType,
		ExternalID:   mdIDs["str"],
		Label:        "GoString" + testSuffix,
		DefaultValue: "Gopher",
		Validation:   metadata.StringLengthValidation(2, 6),
	}

	var dateMetadataField = metadata.Field{
		Type:         metadata.DateFieldType,
		ExternalID:   mdIDs["date"],
		Label:        "GoDate" + testSuffix,
		DefaultValue: time.Now().Format("2006-01-02"),
		Validation:   metadata.GreaterThanValidation(time.Now().AddDate(0, 0, -1), false),
	}

	var enumMetadataField = metadata.Field{
		Type:       metadata.EnumFieldType,
		ExternalID: mdIDs["enum"],
		Label:      "GoDistinctColors" + testSuffix,
		DataSource: dataSource1,
	}

	for _, f := range []metadata.Field{
		integerMetadataField,
		stringMetadataField,
		dateMetadataField,
		enumMetadataField,
	} {
		resp, err := adminApi.AddMetadataField(ctx, f)

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
		resp, err := adminApi.DeleteMetadataField(ctx, DeleteMetadataFieldParams{FieldExternalId: extID})
		if err != nil || resp.Message != "ok" {
			log.Println(err, resp)
		}
	}
}
