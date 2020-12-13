package admin

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var ctx = context.Background()
var adminApi, _ = Create()

var testSuffix = getTestSuffix()

func getTestSuffix() string {
	testSuffix := os.Getenv("TRAVIS_JOB_ID")

	if testSuffix == "" {
		rand.Seed(time.Now().UnixNano())
		testSuffix = strconv.Itoa(rand.Intn(999999))
	}

	return testSuffix
}
