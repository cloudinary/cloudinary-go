package admin_test

// Acceptance tests for API. See `TEST.md` for additional information.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/internal/cldtest"
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
		{
			requestParams: admin.AssetParams{PublicID: cldtest.PublicID,
				RelatedComplementary: api.Bool(true), RelatedComplementaryNextCursor: "NEXT_CURSOR"},
			uri: "/resources/image/upload/" + cldtest.PublicID,
			expectedParams: &url.Values{
				"related_complementary":             []string{"true"},
				"related_complementary_next_cursor": []string{"NEXT_CURSOR"},
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

	// Test case for AssetContextResult with map format
	assetWithMapContext := admin.AssetResult{
		AssetID:  "1",
		PublicID: cldtest.PublicID,
		Context: admin.AssetContextResult{
			Custom: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	responseJsonMapContext, _ := json.Marshal(assetWithMapContext)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Asset response with context in map format",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Asset(ctx, admin.AssetParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetResult)
			if !ok {
				t.Errorf("Response should be type of AssetResult, %s given", reflect.TypeOf(response))
			}
			v.Response = nil // omit raw response comparison
			if !reflect.DeepEqual(*v, assetWithMapContext) {
				t.Errorf("Response asset should be %v, %v given", assetWithMapContext, *v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/upload/" + cldtest.PublicID,
		},
		JsonResponse:      string(responseJsonMapContext),
		ExpectedCallCount: 1,
	})

	// Test case for AssetContextResult with array format
	// The JSON response will have context as an array, but the parsed result should be the same
	assetWithArrayContext := admin.AssetResult{
		AssetID:  "1",
		PublicID: cldtest.PublicID,
		Context: admin.AssetContextResult{
			Custom: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	responseJsonArrayContext, _ := json.Marshal(map[string]interface{}{
		"asset_id":  "1",
		"public_id": cldtest.PublicID,
		"context": []map[string]string{
			{"key": "key1", "value": "value1"},
			{"key": "key2", "value": "value2"},
		},
	})

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Asset response with context in array format",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Asset(ctx, admin.AssetParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetResult)
			if !ok {
				t.Errorf("Response should be type of AssetResult, %s given", reflect.TypeOf(response))
			}
			v.Response = nil // omit raw response comparison
			if !reflect.DeepEqual(*v, assetWithArrayContext) {
				t.Errorf("Response asset should be %v, %v given", assetWithArrayContext, *v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/upload/" + cldtest.PublicID,
		},
		JsonResponse:      string(responseJsonArrayContext),
		ExpectedCallCount: 1,
	})

	assetWithAdminContext := admin.AssetResult{
		AssetID:  "1",
		PublicID: cldtest.PublicID,
		AdminContext: []admin.AssetAdminContextResult{
			{
				Name:  "key1",
				Value: []interface{}{"value1", "value2"},
			},
		},
	}
	responseJsonAdminContext, _ := json.Marshal(map[string]interface{}{
		"asset_id":  "1",
		"public_id": cldtest.PublicID,
		"admin_context": []map[string]interface{}{
			{"name": "key1", "value": []string{"value1", "value2"}},
		},
	})

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Asset response with admin context",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Asset(ctx, admin.AssetParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetResult)
			if !ok {
				t.Errorf("Response should be type of AssetResult, %s given", reflect.TypeOf(response))
			}
			v.Response = nil // omit raw response comparison
			if !reflect.DeepEqual(*v, assetWithAdminContext) {
				t.Errorf("Response asset should be %+v\n%+v given", assetWithAdminContext, *v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/upload/" + cldtest.PublicID,
		},
		JsonResponse:      string(responseJsonAdminContext),
		ExpectedCallCount: 1,
	})

	return testCases
}

// Run tests
func TestAsset_Acceptance(t *testing.T) {
	t.Parallel()
	testAdminAPIByTestCases(getAssetTestCases(), t)
}
