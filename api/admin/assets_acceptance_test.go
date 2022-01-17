package admin_test

// Acceptance tests for Assets. See `TEST.md` for additional information.
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/internal/cldtest"
	"net/url"
	"reflect"
	"testing"
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
				Uri:    t.uri,
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
			Uri:    "/resources/image/upload/restore",
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
				Uri:    t.uri,
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
				KeepOriginal: false,
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\"}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: true,
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":true}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: true,
				Invalidate:   false,
			},
			uri:          "/resources/ASSET_TYPE/DELIVERY_TYPE",
			expectedBody: "{\"public_ids\":\"1,2\",\"keep_original\":true}",
		},
		{
			requestParams: admin.DeleteAssetsParams{
				AssetType:    "ASSET_TYPE",
				DeliveryType: "DELIVERY_TYPE",
				PublicIDs:    []string{"1", "2"},
				KeepOriginal: true,
				Invalidate:   true,
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
			Uri:    "/resources/image/upload",
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
			Uri:    "/resources/image/upload",
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
				Uri:    t.uri,
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
				Tags:       true,
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
				Tags:       false,
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
				Tags:       false,
				Context:    true,
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"context":     []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:  "ASSET_TYPE",
				NextCursor: "NEXT_CURSOR",
				MaxResults: 100,
				Tags:       false,
				Context:    false,
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        false,
				Context:     false,
				Moderations: true,
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"moderations": []string{"true"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        false,
				Context:     false,
				Moderations: false,
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsParams{
				AssetType:   "ASSET_TYPE",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        false,
				Context:     false,
				Moderations: false,
				Direction:   "ASC",
			},
			uri: "/resources/ASSET_TYPE",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
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
			Uri:    "/resources/image",
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

			if !reflect.DeepEqual(expectedResponse, *v) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			Uri:    "/resources/image",
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
				Uri:    t.uri,
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
				Tags:       true,
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
				Tags:       false,
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
				Tags:       false,
				Context:    true,
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
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
				Tags:       false,
				Context:    false,
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:   "ASSET_TYPE",
				Kind:        "KIND",
				Status:      "STATUS",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        false,
				Context:     false,
				Moderations: true,
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
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
				Tags:        false,
				Context:     false,
				Moderations: false,
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
			},
		},
		{
			requestParams: admin.AssetsByModerationParams{
				AssetType:   "ASSET_TYPE",
				Kind:        "KIND",
				Status:      "STATUS",
				NextCursor:  "NEXT_CURSOR",
				MaxResults:  100,
				Tags:        false,
				Context:     false,
				Moderations: false,
				Direction:   "ASC",
			},
			uri: "/resources/ASSET_TYPE/moderations/KIND/STATUS",
			expectedParams: &url.Values{
				"next_cursor": []string{"NEXT_CURSOR"},
				"max_results": []string{"100"},
				"direction":   []string{"ASC"},
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
			Uri:    "/resources/image/moderations",
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

			if !reflect.DeepEqual(expectedResponse, *v) {
				t.Errorf("Expected response to be %v, %v given", expectedResponse, v)
			}
		},
		ExpectedRequest: cldtest.ExpectedRequestParams{
			Method: "GET",
			Uri:    "/resources/image/moderations",
			Params: &url.Values{},
		},
		JsonResponse:      string(responseJson),
		ExpectedCallCount: 1,
	})

	return testCases
}

// Run tests
func TestAssets_Acceptance(t *testing.T) {
	t.Parallel()
	testAdminAPIByTestCases(getAssetsTestCases(), t)
	testAdminAPIByTestCases(getAssetsByModerationTestCases(), t)
	testAdminAPIByTestCases(getDeleteAssetsTestCases(), t)
	testAdminAPIByTestCases(getRestoreAssetsTestCases(), t)
}
