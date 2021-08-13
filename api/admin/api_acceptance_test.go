package admin_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"reflect"
	"testing"
)

// Acceptance test cases for `ping` method
func getPingTestCases() []ApiAcceptanceTestCase {
	return []ApiAcceptanceTestCase{
		{
			Name: "Ping",
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.PingResult)
				if !ok {
					t.Errorf("Response should be type of PingResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest:   expectedRequestParams{Method: "GET", Uri: "/ping"},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
		{
			Name: "Ping error check",
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				v, ok := response.(*admin.PingResult)
				if !ok {
					t.Errorf("Response should be type %s, %s given", reflect.TypeOf(admin.PingResult{}), reflect.TypeOf(response))
				} else {
					if v.Status == "OK" {
						t.Error("Response status should not be OK")
					}

					if v.Error.Message != "ERROR MESSAGE" {
						t.Errorf("Error message should be %s, %s given", "ERROR MESSAGE", v.Error.Message)
					}
				}
			},
			ExpectedRequest:   expectedRequestParams{Method: "GET", Uri: "/ping"},
			JsonResponse:      "{\"error\":{\"message\": \"ERROR MESSAGE\"}}",
			ExpectedCallCount: 1,
		},
		{
			Name: "Ping result struct check",
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				v, ok := response.(*admin.PingResult)
				if ok {
					if v.Status != "OK" {
						t.Errorf("Status should be %s, %s given\n", "OK", v.Status)
					}
				} else {
					t.Errorf("Response should be type %s, %s given", reflect.TypeOf(admin.PingResult{}), reflect.TypeOf(response))
				}
			},
			ExpectedRequest:   expectedRequestParams{Method: "GET", Uri: "/ping"},
			JsonResponse:      "{\"status\":\"OK\"}",
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestAPI_Acceptance(t *testing.T) {
	t.Parallel()
	testApiByTestCases(getPingTestCases(), t)
}
