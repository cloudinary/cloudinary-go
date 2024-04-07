package admin_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"reflect"
	"testing"
)

// Acceptance test cases for `analyze` method
func getAnalysisTestCases() []AdminAPIAcceptanceTestCase {
	type analyzeTestCase struct {
		requestParams admin.AnalyzeParams
		uri           string
		expectedBody  string
	}

	getTestCase := func(num int, t analyzeTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("Analyze #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Analyze(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AnalyzeResult)
				if !ok {
					t.Errorf("Response should be type of AnalyzeResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method:     "POST",
				URI:        t.uri,
				APIVersion: "v2",
				Body:       &t.expectedBody,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	analysisTestCases := []analyzeTestCase{
		{
			requestParams: admin.AnalyzeParams{
				Uri:          cldtest.LogoURL,
				AnalysisType: "captioning",
			},
			uri:          "/analysis/analyze/uri",
			expectedBody: `{"uri":"` + cldtest.LogoURL + `","analysis_type":"captioning"}`,
		},
	}

	for num, testCase := range analysisTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	analyzeRes := admin.AnalyzeResult{
		Data: admin.AnalysisPayload{
			Entity:   cldtest.LogoURL,
			Analysis: map[string]interface{}{"data": "data"},
		},
		RequestId: cldtest.AssetID,
	}
	responseJson, _ := json.Marshal(analyzeRes)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Analyze response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Analyze(ctx, admin.AnalyzeParams{
				Uri:          cldtest.LogoURL,
				AnalysisType: "captioning",
			})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AnalyzeResult)
			if !ok {
				t.Errorf("Response should be type of AnalyzeResult, %s given", reflect.TypeOf(response))
			}
			v.Response = nil // omit raw response comparison
			if !reflect.DeepEqual(*v, analyzeRes) {
				t.Errorf("Response analyzeRes should be %v, %v given", analyzeRes, *v)
			}
		},

		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method:     "POST",
			APIVersion: "v2",
			URI:        "/analysis/analyze/uri",
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	},
	)

	return testCases
}

// Run tests
func TestAnalysis_Acceptance(t *testing.T) {
	t.Parallel()
	testAdminAPIByTestCases(getAnalysisTestCases(), t)
}
