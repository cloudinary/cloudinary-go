package admin_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
	"net/url"
	"reflect"
	"testing"
)

// Acceptance test cases for `asset` method
func getAssetTestCases() []AdminAPIAcceptanceTestCase {
	type assetTestCase struct {
		requestParams  admin.AssetParams
		uri            string
		expectedParams *url.Values
	}

	getTestCase := func(num int, t assetTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("Asset #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Asset(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AssetResult)
				if !ok {
					t.Errorf("Response should be type of AssetResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "GET",
				URI:    t.uri,
				Params: t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	assetTestCases := []assetTestCase{
		{
			requestParams:  admin.AssetParams{PublicID: cldtest.PublicID},
			uri:            "/resources/image/upload/" + cldtest.PublicID,
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetParams{PublicID: cldtest.PublicID,
				Related: api.Bool(true), RelatedNextCursor: "NEXT_CURSOR"},
			uri: "/resources/image/upload/" + cldtest.PublicID,
			expectedParams: &url.Values{
				"related":             []string{"true"},
				"related_next_cursor": []string{"NEXT_CURSOR"},
			},
		},
	}

	for num, testCase := range assetTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	asset := admin.AssetResult{AssetID: "1", PublicID: cldtest.PublicID}
	responseJson, _ := json.Marshal(asset)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Asset response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Asset(ctx, admin.AssetParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetResult)
			if !ok {
				t.Errorf("Response should be type of AssetResult, %s given", reflect.TypeOf(response))
			}
			v.Response = nil // omit raw response comparison
			if !reflect.DeepEqual(*v, asset) {
				t.Errorf("Response asset should be %v, %v given", asset, *v)
			}
		},

		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/upload/" + cldtest.PublicID,
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	},
	)

	return testCases
}

// Run tests
func TestAsset_Acceptance(t *testing.T) {
	t.Parallel()
	testAdminAPIByTestCases(getAssetTestCases(), t)
}
