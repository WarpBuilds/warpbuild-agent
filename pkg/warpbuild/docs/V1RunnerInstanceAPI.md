# \V1RunnerInstanceAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetRunnerInstanceAllocationDetails**](V1RunnerInstanceAPI.md#GetRunnerInstanceAllocationDetails) | **Get** /runner_instance/{id}/allocation_details | Get runner instance allocation details for the id



## GetRunnerInstanceAllocationDetails

> CommonsRunnerInstanceAllocationDetails GetRunnerInstanceAllocationDetails(ctx, id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()

Get runner instance allocation details for the id

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
    id := "id_example" // string | runner instance id
    xPOLLINGSECRET := "xPOLLINGSECRET_example" // string | polling secret for validation

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails(context.Background(), id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerInstanceAllocationDetails`: CommonsRunnerInstanceAllocationDetails
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerInstanceAllocationDetailsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xPOLLINGSECRET** | **string** | polling secret for validation | 

### Return type

[**CommonsRunnerInstanceAllocationDetails**](CommonsRunnerInstanceAllocationDetails.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

