package admin_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
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

// Acceptance test cases for `ping` method
func getUserAgentTestCases() []ApiAcceptanceTestCase {
	return []ApiAcceptanceTestCase{
		{
			Name: "Test User Agent",
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.PingResult)
				if !ok {
					t.Errorf("Response should be type of PingResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: expectedRequestParams{
				Method:  "GET",
				Uri:     "/ping",
				Headers: &map[string]string{"User-Agent": api.UserAgent},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
		{
			Name: "Test User Agent With User Platform",
			RequestTest: func(adminAPI *admin.API, ctx context.Context) (interface{}, error) {
				api.UserPlatform = "Test/1.2.3"
				return adminAPI.Ping(ctx)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.PingResult)
				if !ok {
					t.Errorf("Response should be type of PingResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: expectedRequestParams{
				Method:  "GET",
				Uri:     "/ping",
				Headers: &map[string]string{"User-Agent": fmt.Sprintf("Test/1.2.3 %s", api.UserAgent)},
			},
			JsonResponse:      "{\"status\": \"OK\"}",
			ExpectedCallCount: 1,
		},
	}
}

// Run tests
func TestAPI_Acceptance(t *testing.T) {
	t.Parallel()
	testApiByTestCases(getPingTestCases(), t)
	testApiByTestCases(getUserAgentTestCases(), t)
}
