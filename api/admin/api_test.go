package admin

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var ctx = context.Background()
var adminApi, _ = New()

var testSuffix = getTestSuffix()

func getTestSuffix() string {
	testSuffix := os.Getenv("TRAVIS_JOB_ID")

	if testSuffix == "" {
		rand.Seed(time.Now().UnixNano())
		testSuffix = strconv.Itoa(rand.Intn(999999))
	}

	return testSuffix
}

func TestApi_Timeout(t *testing.T) {
	var originalTimeout = adminApi.Config.Api.Timeout

	adminApi.Config.Api.Timeout = 0 // should timeout immediately

	_, err := adminApi.Ping(ctx)

	if err == nil || !strings.HasSuffix(err.Error(), "context deadline exceeded") {
		t.Error("Expected context timeout did not happen")
	}

	adminApi.Config.Api.Timeout = originalTimeout
}
