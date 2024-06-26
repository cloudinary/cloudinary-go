package admin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
)

var ctx = context.Background()
var adminAPI, _ = admin.New()

func TestAPI_Timeout(t *testing.T) {
	var originalTimeout = adminAPI.Config.API.Timeout

	adminAPI.Config.API.Timeout = 0 // should time out immediately

	_, err := adminAPI.Ping(ctx)

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	adminAPI.Config.API.Timeout = originalTimeout
}
