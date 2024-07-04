<<<<<<< HEAD
# \V1JobsAPI
=======
# \V1JobsApi
>>>>>>> prajjwal-warp-323

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
<<<<<<< HEAD
[**GetEstimatedCosts**](V1JobsAPI.md#GetEstimatedCosts) | **Get** /jobs/estimated-costs | GetEstimatedCosts



## GetEstimatedCosts

> []EstimatedCost GetEstimatedCosts(ctx).StartDate(startDate).EndDate(endDate).Execute()

GetEstimatedCosts
=======
[**GetCostSummary**](V1JobsApi.md#GetCostSummary) | **Get** /jobs/cost-summary | GetCostSummary
[**GetDaywiseCosts**](V1JobsApi.md#GetDaywiseCosts) | **Get** /jobs/daywise-costs | GetDaywiseCosts



## GetCostSummary

> CommonsCostSummary GetCostSummary(ctx).StartDate(startDate).EndDate(endDate).Execute()

GetCostSummary
>>>>>>> prajjwal-warp-323



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    startDate := "startDate_example" // string | Date range start
    endDate := "endDate_example" // string | Date range end

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1JobsAPI.GetEstimatedCosts(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1JobsAPI.GetEstimatedCosts``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetEstimatedCosts`: []EstimatedCost
    fmt.Fprintf(os.Stdout, "Response from `V1JobsAPI.GetEstimatedCosts`: %v\n", resp)
=======
    resp, r, err := apiClient.V1JobsApi.GetCostSummary(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1JobsApi.GetCostSummary``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetCostSummary`: CommonsCostSummary
    fmt.Fprintf(os.Stdout, "Response from `V1JobsApi.GetCostSummary`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

<<<<<<< HEAD
Other parameters are passed through a pointer to a apiGetEstimatedCostsRequest struct via the builder pattern
=======
Other parameters are passed through a pointer to a apiGetCostSummaryRequest struct via the builder pattern
>>>>>>> prajjwal-warp-323


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** | Date range start | 
 **endDate** | **string** | Date range end | 

### Return type

<<<<<<< HEAD
[**[]EstimatedCost**](EstimatedCost.md)
=======
[**CommonsCostSummary**](CommonsCostSummary.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetDaywiseCosts

> []CommonsDaywiseCost GetDaywiseCosts(ctx).StartDate(startDate).EndDate(endDate).Execute()

GetDaywiseCosts



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    startDate := "startDate_example" // string | Date range start
    endDate := "endDate_example" // string | Date range end

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1JobsApi.GetDaywiseCosts(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1JobsApi.GetDaywiseCosts``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetDaywiseCosts`: []CommonsDaywiseCost
    fmt.Fprintf(os.Stdout, "Response from `V1JobsApi.GetDaywiseCosts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetDaywiseCostsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** | Date range start | 
 **endDate** | **string** | Date range end | 

### Return type

[**[]CommonsDaywiseCost**](CommonsDaywiseCost.md)
>>>>>>> prajjwal-warp-323

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

