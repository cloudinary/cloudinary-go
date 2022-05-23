package admin_test

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/config"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"testing"
)

// Function that will be executed during the test
type AdminAPIRequestTest func(api *admin.API, ctx context.Context) (interface{}, error)

// Acceptance test case definition. See `TEST.md` for additional information.
type AdminAPIAcceptanceTestCase struct {
	Name              string                        // Name of the test case
	RequestTest       AdminAPIRequestTest           // Function which will be called as an API request. Put SDK calls here.
	ResponseTest      cldtest.APIResponseTest       // Function which will be called to test an API response.
	ExpectedRequest   cldtest.ExpectedRequestParams // Expected HTTP request to be sent to the server
	JsonResponse      string                        // Mock of the JSON response from server. This is used to check JSON parsing.
	ExpectedStatus    string                        // Expected HTTP status of the request. This status will be returned from the HTTP mock.
	ExpectedCallCount int                           // Expected call count to the server.
	Config            *config.Configuration         // Configuration
}

// Run acceptance tests by given test cases. See `TEST.md` for additional information.
func testAdminAPIByTestCases(cases []AdminAPIAcceptanceTestCase, t *testing.T) {
	for num, test := range cases {
		if test.Name == "" {
			t.Skipf("Test name should be set for test #%d. Skipping it.", num)
		}

		t.Run(test.Name, func(t *testing.T) {
			callCounter := 0
			srv := cldtest.GetServerMock(cldtest.GetTestHandler(test.JsonResponse, t, &callCounter, test.ExpectedRequest))

			res, _ := test.RequestTest(getTestableAdminAPI(srv.URL, test.Config, t), ctx)
			test.ResponseTest(res, t)

			if callCounter != test.ExpectedCallCount {
				t.Errorf("Expected %d call, %d given", test.ExpectedCallCount, callCounter)
			}

			srv.Close()
		})
	}
}

// Get configured API for test
func getTestableAdminAPI(mockServerUrl string, c *config.Configuration, t *testing.T) *admin.API {
	if c == nil {
		var err error
		c, err = config.NewFromParams("TEST", "key", "secret")
		if err != nil {
			t.Error(err)
		}
	}

	c.API.UploadPrefix = mockServerUrl

	api, err := admin.NewWithConfiguration(c)
	if err != nil {
		t.Error(err)
	}

	return api
}
