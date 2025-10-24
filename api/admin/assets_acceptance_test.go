package admin_test

// Acceptance tests for Assets. See `TEST.md` for additional information.
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

// Acceptance test cases for `restore` method
func getRestoreAssetsTestCases() []AdminAPIAcceptanceTestCase {
	type restoreAssetsTestCase struct {
		requestParams admin.RestoreAssetsParams
		uri           string
		expectedBody  string
	}

	getTestCase := func(num int, t restoreAssetsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("RestoreAssets #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.RestoreAssets(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.RestoreAssetsResult)
				if !ok {
					t.Errorf("Response should be type of RestoreAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    t.uri,
				Body:   &t.expectedBody,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	restoreAssetsTestCases := []restoreAssetsTestCase{
		{
			requestParams: admin.RestoreAssetsParams{},
			uri:           "/resources/image/upload/restore",
			expectedBody:  "{\"public_ids\":\"\",\"versions\":\"\"}",
		},
		{
			requestParams: admin.RestoreAssetsParams{
				AssetType: "ASSET_TYPE",
			},
			uri:          "/resources/ASSET_TYPE/upload/restore",
			expectedBody: "{\"public_ids\":\"\",\"versions\":\"\"}",
		},
		{
			requestParams: admin.RestoreAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE/restore",
			expectedBody: "{\"public_ids\":\"\",\"versions\":\"\"}",
		},
		{
			requestParams: admin.RestoreAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE/restore",
			expectedBody: "{\"public_ids\":\"1,2\",\"versions\":\"\"}",
		},
		{
			requestParams: admin.RestoreAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				Versions:     []string{"3", "4"},
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE/restore",
			expectedBody: "{\"public_ids\":\"1,2\",\"versions\":\"3,4\"}",
		},
	}

	for num, testCase := range restoreAssetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	asset := api.BriefAssetResult{AssetID: "1"}
	response := map[string]api.BriefAssetResult{"1": asset}
	responseJson, _ := json.Marshal(response)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "RestoreAssets response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.RestoreAssets(ctx, admin.RestoreAssetsParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.RestoreAssetsResult)
			if !ok {
				t.Errorf("Response should be type of RestoreAssetsResult, %s given", reflect.TypeOf(response))
			}

			directMap := *v
			if responseAsset, ok := directMap["1"]; ok {
				if !reflect.DeepEqual(directMap["1"], asset) {
					t.Errorf("Response asset should be %v, %v given", asset, responseAsset)
				}
			} else {
				t.Errorf("Asset #1 is not found in response %v", v)
			}
		},

		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "POST",
			URI:    "/resources/image/upload/restore",
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	},
	)

	return testCases
}

// Acceptance test cases for `delete` method
func getDeleteAssetsTestCases() []AdminAPIAcceptanceTestCase {
	type deleteAssetsTestCase struct {
		requestParams admin.DeleteAssetsParams
		uri           string
		expectedBody  string
	}

	getTestCase := func(num int, t deleteAssetsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("DeleteAssets #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.DeleteAssets(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.DeleteAssetsResult)
				if !ok {
					t.Errorf("Response should be type of DeleteAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "DELETE",
				URI:    t.uri,
				Body:   &t.expectedBody,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	restoreAssetsTestCases := []deleteAssetsTestCase{
		{
			requestParams: admin.DeleteAssetsParams{},
			uri:           "/resources/image/upload",
			expectedBody:  "{\"public_ids\":\"\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType: "ASSET_TYPE",
			},
			uri:          "/resources/ASSET_TYPE/upload",
			expectedBody: "{\"public_ids\":\"\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: api.Bool(false),
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":false}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: api.Bool(true),
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":true}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: api.Bool(true),
				Invalidate:   api.Bool(false),
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":true,\"invalidate\":false}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: api.Bool(true),
				Invalidate:   api.Bool(true),
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":true,\"invalidate\":true}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:       "ASSET_TYPE",
				DeliveryType:    "DELIVERY_TYPE",
				PublicIDs:       []string{"1", "2"},
				Transformations: "TEST_TRANSFORMATIONS",
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"transformations\":\"TEST_TRANSFORMATIONS\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				NextCursor:   "NEXT_CURSOR",
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"next_cursor\":\"NEXT_CURSOR\"}",
		},
	}

	for num, testCase := range restoreAssetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "DeleteAssets error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.DeleteAssets(ctx, admin.DeleteAssetsParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.DeleteAssetsResult)
			if !ok {
				t.Errorf("Response should be type of DeleteAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Error.Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Error.Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "DELETE",
			URI:    "/resources/image/upload",
		},
		JsonResponse:      "{\"error\":{\"message\":\"TEST ERROR\"}}",
		ExpectedCallCount: 1,
	})

	response := map[string]interface{}{
		"deleted":        map[string]string{"1": "TEST", "2": "TEST_2"},
		"deleted_counts": map[string]interface{}{"1": 1, "2": "2"},
		"partial":        true,
	}
	responseJson, _ := json.Marshal(response)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "DeleteAssets response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.DeleteAssets(ctx, admin.DeleteAssetsParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.DeleteAssetsResult)
			if !ok {
				t.Errorf("Response should be type of DeleteAssetsResult, %s given", reflect.TypeOf(response))
			}

			expectedResponse := &admin.DeleteAssetsResult{
				Deleted:       map[string]string{"1": "TEST", "2": "TEST_2"},
				DeletedCounts: map[string]interface{}{"1": 1, "2": "2"},
				Partial:       true,
			}

			// ugly solution for map[string]interface{} below. deepequal does not work for this case :(
			if !reflect.DeepEqual(expectedResponse.Deleted, v.Deleted) {
				t.Errorf("Response.Deleted expected to be %v, %v given", expectedResponse.Deleted, v.Deleted)
			}

			if expectedResponse.Partial != v.Partial {
				t.Errorf("Response.Partial expected to be %v, %v given", expectedResponse.Partial, v.Partial)
			}

			if fmt.Sprintf("%v", expectedResponse.DeletedCounts) != fmt.Sprintf("%v", v.DeletedCounts) {
				t.Errorf("Response.DeletedCounts expected to be %v, %v given", expectedResponse.DeletedCounts, v.DeletedCounts)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "DELETE",
			URI:    "/resources/image/upload",
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `Assets` method
func getAssetsTestCases() []AdminAPIAcceptanceTestCase {
	type assetsTestCase struct {
		requestParams  admin.AssetsParams
		uri            string
		expectedParams *url.Values
	}

	getTestCase := func(num int, t assetsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("Assets #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.Assets(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AssetsResult)
				if !ok {
					t.Errorf("Response should be type of AssetsResult, %s given", reflect.TypeOf(response))
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

	assetsTestCases := []assetsTestCase{
		{
			requestParams:  admin.AssetsParams{},
			uri:            "/resources/image",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType: "ASSET_TYPE",
			},
			uri:            "/resources/ASSET_TYPE",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
			},
			uri:            "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
				Context:    api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
				Context:    api.Bool(false),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(false),
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"false"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(false),
				Direction:   "ASC",
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"false"},
				"direction":   []string{"ASC"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Prefix:     "PREFIX",
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"prefix":      []string{"PREFIX"},
			},
		}, {
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Fields:     []string{"tags", "secure_url"},
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"fields":      []string{"tags,secure_url"},
			},
		},
	}

	for num, testCase := range assetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Assets error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Assets(ctx, admin.AssetsParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetsResult)
			if !ok {
				t.Errorf("Response should be type of AssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Error.Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Error.Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image",
			Params: &url.Values{},
		},
		JsonResponse:      "{\"error\":{\"message\": \"TEST ERROR\"}}",
		ExpectedCallCount: 1,
	})

	asset := api.BriefAssetResult{AssetID: "1"}
	response := map[string]interface{}{
		"resources":   []api.BriefAssetResult{asset},
		"next_cursor": "NEXT_CURSOR",
	}
	responseJson, _ := json.Marshal(response)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Assets response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.Assets(ctx, admin.AssetsParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			expectedResponse := admin.AssetsResult{
				Assets:     []api.BriefAssetResult{asset},
				NextCursor: "NEXT_CURSOR",
			}

			v, ok := response.(*admin.AssetsResult)
			if !ok {
				t.Errorf("Response should be type of %s, %s given", reflect.TypeOf(expectedResponse), reflect.TypeOf(response))
			}

			if !reflect.DeepEqual(expectedResponse.Assets, v.Assets) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}
			if !reflect.DeepEqual(expectedResponse.NextCursor, v.NextCursor) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image",
			Params: &url.Values{},
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `assetsByModeration` method
func getAssetsByModerationTestCases() []AdminAPIAcceptanceTestCase {
	type assetByModerationTestCase struct {
		requestParams  admin.AssetsByModerationParams
		uri            string
		expectedParams *url.Values
	}

	getTestCase := func(num int, t assetByModerationTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("AssetsByModeration #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.AssetsByModeration(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AssetsResult)
				if !ok {
					t.Errorf("Response should be type of AssetsResult, %s given", reflect.TypeOf(response))
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

	assetByModerationTestCases := []assetByModerationTestCase{
		{
			requestParams:  admin.AssetsByModerationParams{},
			uri:            "/resources/image/moderations",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType: "ASSET_TYPE",
			},
			uri:            "/resources/ASSET_TYPE/moderations",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType: "ASSET_TYPE",
				Kind:      "KIND",
			},
			uri:            "/resources/ASSET_TYPE/moderations/KIND",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType: "ASSET_TYPE",
				Kind:      "KIND",
				Status:    "STATUS",
			},
			uri:            "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
				Context:    api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:  "ASSET_TYPE",
				Kind:       "KIND",
				Status:     "STATUS",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       api.Bool(false),
				Context:    api.Bool(false),
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:   "ASSET_TYPE",
				Kind:        "KIND",
				Status:      "STATUS",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(true),
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:   "ASSET_TYPE",
				Kind:        "KIND",
				Status:      "STATUS",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(false),
				Fields:      []string{"tags", "secure_url"},
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"false"},
				"fields":      []string{"tags,secure_url"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:   "ASSET_TYPE",
				Kind:        "KIND",
				Status:      "STATUS",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        api.Bool(false),
				Context:     api.Bool(false),
				Moderations: api.Bool(false),
				Direction:   "ASC",
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"direction":   []string{"ASC"},
				"tags":        []string{"false"},
				"context":     []string{"false"},
				"moderations": []string{"false"},
			},
		},
	}

	for num, testCase := range assetByModerationTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "AssetsByModeration error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.AssetsByModeration(ctx, admin.AssetsByModerationParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AssetsResult)
			if !ok {
				t.Errorf("Response should be type of AssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Error.Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Error.Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/moderations",
			Params: &url.Values{},
		},
		JsonResponse:      "{\"error\":{\"message\": \"TEST ERROR\"}}",
		ExpectedCallCount: 1,
	})

	asset := api.BriefAssetResult{AssetID: "1"}
	response := map[string]interface{}{
		"resources":   []api.BriefAssetResult{asset},
		"next_cursor": "NEXT_CURSOR",
	}
	responseJson, _ := json.Marshal(response)

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "AssetsByModeration response parsing case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.AssetsByModeration(ctx, admin.AssetsByModerationParams{})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			expectedResponse := admin.AssetsResult{
				Assets:     []api.BriefAssetResult{asset},
				NextCursor: "NEXT_CURSOR",
			}

			v, ok := response.(*admin.AssetsResult)
			if !ok {
				t.Errorf("Response should be type of %s, %s given", reflect.TypeOf(expectedResponse), reflect.TypeOf(response))
			}

			if !reflect.DeepEqual(expectedResponse.Assets, v.Assets) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}

			if !reflect.DeepEqual(expectedResponse.NextCursor, v.NextCursor) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			URI:    "/resources/image/moderations",
			Params: &url.Values{},
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `AddRelatedAssets` method
func getAddRelatedAssetsTestCases() []AdminAPIAcceptanceTestCase {
	type addRelatedAssetsTestCase struct {
		requestParams  admin.AddRelatedAssetsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t addRelatedAssetsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("AddRelatedAssets #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.AddRelatedAssets(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AddRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of AddRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase
	params := fmt.Sprintf("{\"assets_to_relate\":[\"%s\",\"%s\"]}", cldtest.FQPublicID2, cldtest.FQPublicID3)
	addRelatedAssetsTestCases := []addRelatedAssetsTestCase{
		{
			requestParams: admin.AddRelatedAssetsParams{
				PublicID:       cldtest.PublicID,
				AssetsToRelate: []string{cldtest.FQPublicID2, cldtest.FQPublicID3},
			},
			uri:            "/resources/related_assets/image/upload/" + cldtest.PublicID,
			expectedParams: &params,
		},
		{
			requestParams: admin.AddRelatedAssetsParams{
				AssetType:      "ASSET_TYPE",
				DeliveryType:   "DELIVERY_TYPE",
				PublicID:       cldtest.PublicID,
				AssetsToRelate: []string{cldtest.FQPublicID2, cldtest.FQPublicID3},
			},
			uri:            "/resources/related_assets/ASSET_TYPE/DELIVERY_TYPE/" + cldtest.PublicID,
			expectedParams: &params,
		},
	}

	for num, testCase := range addRelatedAssetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Related Assets error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.AddRelatedAssets(ctx, admin.AddRelatedAssetsParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AddRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of AddRelatedAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "POST",
			URI:    "/resources/related_assets/image/upload/" + cldtest.PublicID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `DeleteRelatedAssets` method
func getDeleteRelatedAssetsTestCases() []AdminAPIAcceptanceTestCase {
	type deleteRelatedAssetsTestCase struct {
		requestParams  admin.DeleteRelatedAssetsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t deleteRelatedAssetsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("DeleteRelatedAssets #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.DeleteRelatedAssets(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.DeleteRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "DELETE",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase
	params := fmt.Sprintf("{\"assets_to_unrelate\":[\"%s\",\"%s\"]}", cldtest.FQPublicID2, cldtest.FQPublicID3)
	addRelatedAssetsTestCases := []deleteRelatedAssetsTestCase{
		{
			requestParams: admin.DeleteRelatedAssetsParams{
				PublicID:         cldtest.PublicID,
				AssetsToUnrelate: []string{cldtest.FQPublicID2, cldtest.FQPublicID3},
			},
			uri:            "/resources/related_assets/image/upload/" + cldtest.PublicID,
			expectedParams: &params,
		},
		{
			requestParams: admin.DeleteRelatedAssetsParams{
				AssetType:        "ASSET_TYPE",
				DeliveryType:     "DELIVERY_TYPE",
				PublicID:         cldtest.PublicID,
				AssetsToUnrelate: []string{cldtest.FQPublicID2, cldtest.FQPublicID3},
			},
			uri:            "/resources/related_assets/ASSET_TYPE/DELIVERY_TYPE/" + cldtest.PublicID,
			expectedParams: &params,
		},
	}

	for num, testCase := range addRelatedAssetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Delete Related Assets Error Case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.DeleteRelatedAssets(ctx, admin.DeleteRelatedAssetsParams{PublicID: cldtest.PublicID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.DeleteRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "DELETE",
			URI:    "/resources/related_assets/image/upload/" + cldtest.PublicID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `AddRelatedAssetsByAssetIDs` method
func getAddRelatedAssetsByAssetIDsTestCases() []AdminAPIAcceptanceTestCase {
	type addRelatedAssetsByAssetIDsTestCase struct {
		requestParams  admin.AddRelatedAssetsByAssetIDsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t addRelatedAssetsByAssetIDsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("AddRelatedAssetsByAssetIDs #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.AddRelatedAssetsByAssetIDs(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AddRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of AddRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	params := fmt.Sprintf("{\"assets_to_relate\":[\"%s\",\"%s\"]}", cldtest.AssetID2, cldtest.AssetID3)
	addRelatedAssetsByAssetIDsTestCases := []addRelatedAssetsByAssetIDsTestCase{
		{
			requestParams: admin.AddRelatedAssetsByAssetIDsParams{
				AssetID:        cldtest.AssetID,
				AssetsToRelate: []string{cldtest.AssetID2, cldtest.AssetID3},
			},
			uri:            "/resources/related_assets/" + cldtest.AssetID,
			expectedParams: &params,
		},
	}

	for num, testCase := range addRelatedAssetsByAssetIDsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Related Assets By Asset IDs error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.AddRelatedAssetsByAssetIDs(ctx, admin.AddRelatedAssetsByAssetIDsParams{AssetID: cldtest.AssetID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AddRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of AddRelatedAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "POST",
			URI:    "/resources/related_assets/" + cldtest.AssetID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `DeleteRelatedAssetsByAssetIDs` method
func getDeleteRelatedAssetsByAssetIDsTestCases() []AdminAPIAcceptanceTestCase {
	type addRelatedAssetsByAssetIDsTestCase struct {
		requestParams  admin.DeleteRelatedAssetsByAssetIDsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t addRelatedAssetsByAssetIDsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("AddRelatedAssetsByAssetIDs #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.DeleteRelatedAssetsByAssetIDs(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.DeleteRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "DELETE",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	params := fmt.Sprintf("{\"assets_to_unrelate\":[\"%s\",\"%s\"]}", cldtest.AssetID2, cldtest.AssetID3)
	addRelatedAssetsByAssetIDsTestCases := []addRelatedAssetsByAssetIDsTestCase{
		{
			requestParams: admin.DeleteRelatedAssetsByAssetIDsParams{
				AssetID:          cldtest.AssetID,
				AssetsToUnrelate: []string{cldtest.AssetID2, cldtest.AssetID3},
			},
			uri:            "/resources/related_assets/" + cldtest.AssetID,
			expectedParams: &params,
		},
	}

	for num, testCase := range addRelatedAssetsByAssetIDsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Related Assets By Asset IDs error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.DeleteRelatedAssetsByAssetIDs(ctx, admin.DeleteRelatedAssetsByAssetIDsParams{AssetID: cldtest.AssetID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.DeleteRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "DELETE",
			URI:    "/resources/related_assets/" + cldtest.AssetID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `AddRelatedAssetsByAssetIDs` method
func getAddRelatedComplementaryAssetsByAssetIDsTestCases() []AdminAPIAcceptanceTestCase {
	type addRelatedComplementaryAssetsByAssetIDsTestCase struct {
		requestParams  admin.AddRelatedComplementaryAssetsByAssetIDsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t addRelatedComplementaryAssetsByAssetIDsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("AddRelatedComplementaryAssetsByAssetIDs #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.AddRelatedComplementaryAssetsByAssetIDs(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.AddRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of AddRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "POST",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	params := fmt.Sprintf("{\"complementary_assets_to_relate\":[\"%s\",\"%s\"],\"subkind\":\"subkind\"}", cldtest.AssetID2, cldtest.AssetID3)
	addRelatedComplementaryAssetsByAssetIDsTestCases := []addRelatedComplementaryAssetsByAssetIDsTestCase{
		{
			requestParams: admin.AddRelatedComplementaryAssetsByAssetIDsParams{
				AssetID:                     cldtest.AssetID,
				ComplementaryAssetsToRelate: []string{cldtest.AssetID2, cldtest.AssetID3},
				Subkind:                     "subkind",
			},
			uri:            "/resources/related_complementary_assets/" + cldtest.AssetID,
			expectedParams: &params,
		},
	}

	for num, testCase := range addRelatedComplementaryAssetsByAssetIDsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Related Complementary Assets By Asset IDs error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.AddRelatedComplementaryAssetsByAssetIDs(ctx, admin.AddRelatedComplementaryAssetsByAssetIDsParams{AssetID: cldtest.AssetID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.AddRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of AddRelatedComplementaryAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "POST",
			URI:    "/resources/related_complementary_assets/" + cldtest.AssetID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `DeleteRelatedAssetsByAssetIDs` method
func getDeleteRelatedComplementaryAssetsByAssetIDsTestCases() []AdminAPIAcceptanceTestCase {
	type deleteRelatedComplementaryAssetsByAssetIDsTestCase struct {
		requestParams  admin.DeleteRelatedComplementaryAssetsByAssetIDsParams
		uri            string
		expectedParams *string
	}

	getTestCase := func(num int, t deleteRelatedComplementaryAssetsByAssetIDsTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("DeleteRelatedComplementaryAssetsByAssetIDs #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.DeleteRelatedComplementaryAssetsByAssetIDs(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.DeleteRelatedAssetsResult)
				if !ok {
					t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
				}
			},
			ExpectedRequest: cldtest.ExpectedRequestParams{
				Method: "DELETE",
				URI:    t.uri,
				Body:   t.expectedParams,
			},
			JsonResponse:      "{}",
			ExpectedCallCount: 1,
		}
	}

	var testCases []AdminAPIAcceptanceTestCase

	params := fmt.Sprintf("{\"complementary_assets_to_unrelate\":[\"%s\",\"%s\"],\"subkind\":\"subkind\"}", cldtest.AssetID2, cldtest.AssetID3)
	deleteRelatedComplementaryAssetsByAssetIDsTestCases := []deleteRelatedComplementaryAssetsByAssetIDsTestCase{
		{
			requestParams: admin.DeleteRelatedComplementaryAssetsByAssetIDsParams{
				AssetID:                       cldtest.AssetID,
				ComplementaryAssetsToUnrelate: []string{cldtest.AssetID2, cldtest.AssetID3},
				Subkind:                       "subkind",
			},
			uri:            "/resources/related_complementary_assets/" + cldtest.AssetID,
			expectedParams: &params,
		},
	}

	for num, testCase := range deleteRelatedComplementaryAssetsByAssetIDsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	testCases = append(testCases, AdminAPIAcceptanceTestCase{
		Name: "Related Complementary Assets By Asset IDs error case",
		RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
			return api.DeleteRelatedComplementaryAssetsByAssetIDs(ctx, admin.DeleteRelatedComplementaryAssetsByAssetIDsParams{AssetID: cldtest.AssetID})
		},
		ResponseTest: func(response interface{}, t *testing.T) {
			v, ok := response.(*admin.DeleteRelatedAssetsResult)
			if !ok {
				t.Errorf("Response should be type of DeleteRelatedAssetsResult, %s given", reflect.TypeOf(response))
			}

			if v.Failed[0].Message != "TEST ERROR" {
				t.Errorf("Error message should be %s, %s given", "TEST ERROR", v.Failed[0].Message)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "DELETE",
			URI:    "/resources/related_complementary_assets/" + cldtest.AssetID,
			Params: &url.Values{},
		},
		JsonResponse:      "{\"failed\":[{\"message\": \"TEST ERROR\"}]}",
		ExpectedCallCount: 1,
	})

	return testCases
}

// Acceptance test cases for `VisualSearch` method
func getVisualSearchTestCases() []AdminAPIAcceptanceTestCase {
	type visualSearchTestCase struct {
		requestParams  admin.VisualSearchParams
		uri            string
		expectedParams *url.Values
	}

	getTestCase := func(num int, t visualSearchTestCase) AdminAPIAcceptanceTestCase {
		return AdminAPIAcceptanceTestCase{
			Name: fmt.Sprintf("VisualSearch #%d", num),
			RequestTest: func(api *admin.API, ctx context.Context) (interface{}, error) {
				return api.VisualSearch(ctx, t.requestParams)
			},
			ResponseTest: func(response interface{}, t *testing.T) {
				_, ok := response.(*admin.VisualSearchResult)
				if !ok {
					t.Errorf("Response should be type of VisualSearchResult, %s given", reflect.TypeOf(response))
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
	addRelatedAssetsTestCases := []visualSearchTestCase{
		{
			requestParams: admin.VisualSearchParams{
				Text: "TEST",
			},
			uri: "/resources/visual_search",
			expectedParams: &url.Values{
				"text": []string{"TEST"},
			},
		},
		{
			requestParams: admin.VisualSearchParams{
				ImageAssetID: cldtest.AssetID,
			},
			uri: "/resources/visual_search",
			expectedParams: &url.Values{
				"image_asset_id": []string{cldtest.AssetID},
			},
		},
		{
			requestParams: admin.VisualSearchParams{
				ImageURL: cldtest.LogoURL,
			},
			uri: "/resources/visual_search",
			expectedParams: &url.Values{
				"image_url": []string{cldtest.LogoURL},
			},
		},
	}

	for num, testCase := range addRelatedAssetsTestCases {
		testCases = append(testCases, getTestCase(num, testCase))
	}

	return testCases
}

// Run tests
func TestAssets_Acceptance(t *testing.T) {
	t.Parallel()
	testAdminAPIByTestCases(getAssetsTestCases(), t)
	testAdminAPIByTestCases(getAddRelatedAssetsTestCases(), t)
	testAdminAPIByTestCases(getDeleteRelatedAssetsTestCases(), t)
	testAdminAPIByTestCases(getAddRelatedAssetsByAssetIDsTestCases(), t)
	testAdminAPIByTestCases(getDeleteRelatedAssetsByAssetIDsTestCases(), t)
	testAdminAPIByTestCases(getAddRelatedComplementaryAssetsByAssetIDsTestCases(), t)
	testAdminAPIByTestCases(getDeleteRelatedComplementaryAssetsByAssetIDsTestCases(), t)
	testAdminAPIByTestCases(getAssetsByModerationTestCases(), t)
	testAdminAPIByTestCases(getVisualSearchTestCases(), t)
	testAdminAPIByTestCases(getDeleteAssetsTestCases(), t)
	testAdminAPIByTestCases(getRestoreAssetsTestCases(), t)
}
