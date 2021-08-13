package admin_test

import (
	"context"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// Function that will be executed during the test
type ApiRequestTest func(api *admin.API, ctx context.Context) (interface{}, error)

// Test function for the response from the API.
type ApiResponseTest func(response interface{}, t *testing.T)

// Acceptance test case definition. See `TEST.md` for an additional information.
type ApiAcceptanceTestCase struct {
	Name              string                // Name of the test case
	RequestTest       ApiRequestTest        // Function which will be called as a API request. Put SDK calls here.
	ResponseTest      ApiResponseTest       // Function which will be called to test an API response.
	ExpectedRequest   expectedRequestParams // Expected HTTP request to be sent to the server
	JsonResponse      string                // Mock of the JSON response from server. This is used to check JSON parsing.
	ExpectedStatus    string                // Expected HTTP status of the request. This status will be returned from the HTTP mock.
	ExpectedCallCount int                   // Expected call count to the server.
}

// Run acceptance tests by given test cases. See `TEST.md` for an additional information.
func testApiByTestCases(cases []ApiAcceptanceTestCase, t *testing.T) {
	for num, test := range cases {
		if test.Name == "" {
			t.Skipf("Test name should be set for test #%d. Skipping it.", num)
		}

		t.Run(test.Name, func(t *testing.T) {
			callCounter := 0
			srv := getServerMock(getTestHandler(test.JsonResponse, t, &callCounter, test.ExpectedRequest))

			res, _ := test.RequestTest(getTestableApi(srv.URL, t), ctx)
			test.ResponseTest(res, t)

			if callCounter != test.ExpectedCallCount {
				t.Errorf("Expected %d call, %d given", test.ExpectedCallCount, callCounter)
			}

			srv.Close()
		})
	}
}

type testFunction func(w http.ResponseWriter, r *http.Request)
type expectedRequestParams struct {
	Method  string             // Expected HTTP method of the request
	Uri     string             // Expected URI
	Params  *url.Values        // Expected URI params
	Body    *string            // Expected HTTP body (for POST / PUT requests)
	Headers *map[string]string // Expected HTTP request headers
}

// Get test handler for HTTP server. Contains basic checks by expected request params.
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
		_, err := io.WriteString(w, response)
		if err != nil {
			t.Error(err)
		}
	}
}

// Get configured API for test
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

// Get HTTP server mock
func getServerMock(fn testFunction) *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	})

	srv := httptest.NewServer(handler)

	return srv
}
