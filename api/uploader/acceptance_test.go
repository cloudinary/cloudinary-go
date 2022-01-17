package uploader_test

import (
	"context"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"testing"
)

// UploadAPIRequestTest is a function that will be executed during the test.
type UploadAPIRequestTest func(api *uploader.API, ctx context.Context) (interface{}, error)

// UploadAPIAcceptanceTestCase is the acceptance test case definition. See `TEST.md` for an additional information.
type UploadAPIAcceptanceTestCase struct {
	Name              string                        // Name of the test case
	RequestTest       UploadAPIRequestTest          // Function which will be called as a API request. Put SDK calls here.
	ResponseTest      cldtest.ApiResponseTest       // Function which will be called to test an API response.
	ExpectedRequest   cldtest.ExpectedRequestParams // Expected HTTP request to be sent to the server
	JsonResponse      string                        // Mock of the JSON response from server. This is used to check JSON parsing.
	ExpectedStatus    string                        // Expected HTTP status of the request. This status will be returned from the HTTP mock.
	ExpectedCallCount int                           // Expected call count to the server.
	Config            *config.Configuration         // Configuration
}

// testUploadAPIByTestCases run acceptance tests by the given test cases. See `TEST.md` for an additional information.
func testUploadAPIByTestCases(cases []UploadAPIAcceptanceTestCase, t *testing.T) {
	for num, test := range cases {
		if test.Name == "" {
			t.Skipf("Test name should be set for test #%d. Skipping it.", num)
		}

		t.Run(test.Name, func(t *testing.T) {
			callCounter := 0
			srv := cldtest.GetServerMock(cldtest.GetTestHandler(test.JsonResponse, t, &callCounter, test.ExpectedRequest))

			res, _ := test.RequestTest(getTestableUploadAPI(srv.URL, test.Config, t), ctx)
			test.ResponseTest(res, t)

			if callCounter != test.ExpectedCallCount {
				t.Errorf("Expected %d call, %d given", test.ExpectedCallCount, callCounter)
			}

			srv.Close()
		})
	}
}

// Get configured API for test
func getTestableUploadAPI(mockServerUrl string, c *config.Configuration, t *testing.T) *uploader.API {
	if c == nil {
		var err error
		c, err = config.NewFromParams("TEST", "key", "secret")
		if err != nil {
			t.Error(err)
		}
	}

	c.API.UploadPrefix = mockServerUrl

	api, err := uploader.NewWithConfiguration(c)
	if err != nil {
		t.Error(err)
	}

	return api
}
