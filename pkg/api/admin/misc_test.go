package admin

import (
	"testing"
)

func TestAdmin_Ping(t *testing.T) {
	a := Create()

	resp, err := a.Ping()

	if err != nil {
		t.Error(err)
	}

	if resp == nil || resp.Status != "ok" {
		t.Error(resp)
	}
}

func TestAdmin_Usage(t *testing.T) {
	a := Create()

	resp, err := a.Usage()

	if err != nil || len(resp.Plan) < 1{
		t.Error(err, resp)
	}
}

func TestAdmin_Tags(t *testing.T) {
	a := Create()

	resp, err := a.Tags(TagsParams{})

	if err != nil || len(resp.Tags) < 1{
		t.Error(err, resp)
	}

	if resp.NextCursor != "" {
		resp2, err := a.Tags(TagsParams{NextCursor: resp.NextCursor})

		if err != nil || len(resp2.Tags) < 1{
			t.Error(err, resp)
		}
	}
}
