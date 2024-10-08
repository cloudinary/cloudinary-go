package admin_test

import (
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

func TestAdmin_Ping(t *testing.T) {
	resp, err := adminAPI.Ping(ctx)

	if err != nil || resp.Status != "ok" {
		t.Error(resp)
	}

	rawResponse := resp.Response.(*map[string]interface{})

	if (*rawResponse)["status"] != "ok" {
		t.Error(resp)
	}
}

func TestAdmin_GetConfig(t *testing.T) {
	resp, err := adminAPI.GetConfig(ctx, admin.GetConfigParams{Settings: api.Bool(true)})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp)
	}

	assert.Equal(t, adminAPI.Config.Cloud.CloudName, resp.CloudName)
	assert.NotEmpty(t, resp.Settings.FolderMode)
}

func TestAdmin_Usage(t *testing.T) {
	resp, err := adminAPI.Usage(ctx, admin.UsageParams{})

	if err != nil || len(resp.Plan) < 1 {
		t.Error(err, resp)
	}
}

func TestAdmin_UsageYesterday(t *testing.T) {
	resp, err := adminAPI.Usage(ctx, admin.UsageParams{Date: time.Now().AddDate(0, 0, -2)})

	if err != nil || len(resp.Plan) < 1 {
		t.Error(err, resp)
	}
}

func TestAdmin_Tags(t *testing.T) {
	resp, err := adminAPI.Tags(ctx, admin.TagsParams{})

	if err != nil || len(resp.Tags) < 1 {
		t.Error(err, resp)
	}

	if resp.NextCursor != "" {
		resp2, err := adminAPI.Tags(ctx, admin.TagsParams{NextCursor: resp.NextCursor})

		if err != nil || len(resp2.Tags) < 1 {
			t.Error(err, resp)
		}
	}
}
