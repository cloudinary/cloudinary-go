package admin_test

import (
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/stretchr/testify/assert"
)

const testTriggerURI = "https://example.com/000-go-trigger-test"

var testTriggerID string

func TestTriggers_CreateTrigger(t *testing.T) {
	resp, err := adminAPI.CreateTrigger(ctx, admin.CreateTriggerParams{
		URI:       testTriggerURI,
		EventType: "upload",
	})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
		return
	}

	testTriggerID = resp.ID
	assert.NotEmpty(t, testTriggerID)
}

func TestTriggers_ListTriggers(t *testing.T) {
	resp, err := adminAPI.ListTriggers(ctx, admin.ListTriggersParams{})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
	}
}

func TestTriggers_UpdateTrigger(t *testing.T) {
	if testTriggerID == "" {
		t.Skip("create trigger test did not run or failed")
	}

	resp, err := adminAPI.UpdateTrigger(ctx, admin.UpdateTriggerParams{
		TriggerID:  testTriggerID,
		URI:        testTriggerURI,
		EventType:  "upload",
		AuthScheme: "default",
	})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
	}
}

func TestTriggers_DeleteTrigger(t *testing.T) {
	if testTriggerID == "" {
		t.Skip("create trigger test did not run or failed")
	}

	resp, err := adminAPI.DeleteTrigger(ctx, admin.DeleteTriggerParams{
		TriggerID: testTriggerID,
	})

	if err != nil || resp.Error.Message != "" {
		t.Error(resp, err)
	}
}
