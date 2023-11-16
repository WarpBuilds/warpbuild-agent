# \V1JobsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetEstimatedCosts**](V1JobsAPI.md#GetEstimatedCosts) | **Get** /jobs/estimated-costs | GetEstimatedCosts



## GetEstimatedCosts

> []EstimatedCost GetEstimatedCosts(ctx).StartDate(startDate).EndDate(endDate).Execute()

GetEstimatedCosts



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
    resp, r, err := apiClient.V1JobsAPI.GetEstimatedCosts(context.Background()).StartDate(startDate).EndDate(endDate).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1JobsAPI.GetEstimatedCosts``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetEstimatedCosts`: []EstimatedCost
    fmt.Fprintf(os.Stdout, "Response from `V1JobsAPI.GetEstimatedCosts`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetEstimatedCostsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** | Date range start | 
 **endDate** | **string** | Date range end | 

### Return type

[**[]EstimatedCost**](EstimatedCost.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

