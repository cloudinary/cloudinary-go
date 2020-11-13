package admin

import (
	"context"
	"testing"
	"time"
)

var ctx = context.Background()
var adminApi = Create()

func TestAdmin_Ping(t *testing.T) {
	resp, err := adminApi.Ping(ctx)

	if err != nil || resp.Status != "ok" {
		t.Error(resp)
	}
}

func TestAdmin_Usage(t *testing.T) {
	resp, err := adminApi.Usage(ctx, UsageParams{})

	if err != nil || len(resp.Plan) < 1 {
		t.Error(err, resp)
	}
}

func TestAdmin_UsageYesterday(t *testing.T) {
	resp, err := adminApi.Usage(ctx, UsageParams{Date: time.Now().AddDate(0, 0, -1)})

	if err != nil || len(resp.Plan) < 1 {
		t.Error(err, resp)
	}
}

func TestAdmin_Tags(t *testing.T) {
	resp, err := adminApi.Tags(ctx, TagsParams{})

	if err != nil || len(resp.Tags) < 1 {
		t.Error(err, resp)
	}

	if resp.NextCursor != "" {
		resp2, err := adminApi.Tags(ctx, TagsParams{NextCursor: resp.NextCursor})

		if err != nil || len(resp2.Tags) < 1 {
			t.Error(err, resp)
		}
	}
}
