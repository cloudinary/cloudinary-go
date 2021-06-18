package admin_test

import (
	"context"
	"github.com/cloudinary/cloudinary-go/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/cloudinary/cloudinary-go/api/admin"
)

var ctx = context.Background()
var adminAPI, _ = admin.New()

type ApiTestRequest func(api *admin.API, ctx context.Context) (interface{}, error)
type ApiTestResponse func(response interface{}, t *testing.T)
type ApiAcceptanceTestCase struct {
	Name              string
	TestRequest       ApiTestRequest
	TestResponse      ApiTestResponse
	ExpectedRequest   expectedRequestParams
	JsonResponse      string
	ExpectedStatus    string
	ExpectedCallCount int
}

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

func getPingTestCases() []ApiAcceptanceTestCase {
	return []ApiAcceptanceTestCase{
		{
			Name: "Ping",
			TestRequest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			TestResponse: func(response interface{}, t *testing.T) {
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
			TestRequest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			TestResponse: func(response interface{}, t *testing.T) {
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
			TestRequest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Ping(ctx)
			},
			TestResponse: func(response interface{}, t *testing.T) {
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

func TestAPI_Acceptance(t *testing.T) {
	t.Parallel()
	testApiByTestCases(getPingTestCases(), t)
}

func testApiByTestCases(cases []ApiAcceptanceTestCase, t *testing.T) {
	for num, test := range cases {
		if test.Name == "" {
			t.Skipf("Test name should be set for test #%d. Skipping it.", num)
		}

		t.Run(test.Name, func(t *testing.T) {
			callCounter := 0
			srv := getServerMock(getTestHandler(test.JsonResponse, t, &callCounter, test.ExpectedRequest))

			res, _ := test.TestRequest(getTestableApi(srv.URL, t), ctx)
			test.TestResponse(res, t)

			if callCounter != test.ExpectedCallCount {
				t.Errorf("Expected %d call, %d given", test.ExpectedCallCount, callCounter)
			}

			srv.Close()
		})
	}
}

type testFunction func(w http.ResponseWriter, r *http.Request)
type expectedRequestParams struct {
	Method  string
	Uri     string
	Params  *url.Values
	Body    *string
	Headers *map[string]string
}

func getTestHandler(response string, t *testing.T, callCounter *int, ep expectedRequestParams) testFunction {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != ep.Method {
			t.Errorf("HTTP method should be %s", ep.Method)
		}

		if ep.Params != nil && ep.Params.Encode() != r.URL.Query().Encode() {
			t.Errorf(
				"Expected query string: %s, got: %s\n",
				ep.Params.Encode(),
				r.URL.Query().Encode(),
			)
		}

		expectedURI := "/" + apiVersion + "/TEST" + ep.Uri
		if expectedURI != r.URL.Path {
			t.Errorf(
				"Expected request URI: %s, got: %s\n",
				expectedURI,
				r.URL.Path,
			)
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			if r.Body != nil && ep.Body != nil {
				bodyString, err := ioutil.ReadAll(r.Body)

				if err != nil {
					t.Error(err)
				}

				if string(bodyString) != *ep.Body {
					t.Errorf("Wrong request body. Expected: %s, given: %s", *ep.Body, string(bodyString))
				}
			}
		}

		*callCounter++
		io.WriteString(w, response)
	}
}
func getTestableApi(mockServerUrl string, t *testing.T) *admin.API {
	c, err := config.NewFromParams("TEST", "", "")
	if err != nil {
		t.Error(err)
	}

	c.API.UploadPrefix = mockServerUrl

	api, err := admin.NewWithConfiguration(c)
	if err != nil {
		t.Error(err)
	}

	return api
}
func getServerMock(fn testFunction) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	})

	srv := httptest.NewServer(handler)

	return srv
}
