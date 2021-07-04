package admin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

var ctx = context.Background()
var adminAPI, _ = admin.New()

const apiVersion = "v1_1"

func TestAPI_Timeout(t *testing.T) {
	var originalTimeout = adminAPI.Config.API.Timeout

	adminAPI.Config.API.Timeout = 0 // should timeout immediately

	_, err := adminAPI.Ping(ctx)

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	adminAPI.Config.API.Timeout = originalTimeout
}
