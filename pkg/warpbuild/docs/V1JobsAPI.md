# \V1JobsApi

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCostSummary**](V1JobsApi.md#GetCostSummary) | **Get** /jobs/cost-summary | GetCostSummary
[**GetDaywiseCosts**](V1JobsApi.md#GetDaywiseCosts) | **Get** /jobs/daywise-costs | GetDaywiseCosts



## GetCostSummary

> CommonsCostSummary GetCostSummary(ctx).StartDate(startDate).EndDate(endDate).Execute()

GetCostSummary



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
    resp, r, err := apiClient.V1JobsApi.GetCostSummary(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1JobsApi.GetCostSummary``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetCostSummary`: CommonsCostSummary
    fmt.Fprintf(os.Stdout, "Response from `V1JobsApi.GetCostSummary`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetCostSummaryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** | Date range start | 
 **endDate** | **string** | Date range end | 

### Return type

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

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

