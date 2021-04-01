package admin_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

var ctx = context.Background()
var adminApi, _ = admin.New()

func TestApi_Timeout(t *testing.T) {
	var originalTimeout = adminApi.Config.Api.Timeout

	adminApi.Config.Api.Timeout = 0 // should timeout immediately

	_, err := adminApi.Ping(ctx)

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	adminApi.Config.Api.Timeout = originalTimeout
}
