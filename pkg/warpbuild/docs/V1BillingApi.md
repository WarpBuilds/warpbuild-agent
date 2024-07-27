# \V1BillingAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostUsageForInternalService**](V1BillingAPI.md#PostUsageForInternalService) | **Post** /billing/usage/internal | Post Usage for internal service



## PostUsageForInternalService

> CommonsInternalPostUsageOutput PostUsageForInternalService(ctx).Body(body).Execute()

Post Usage for internal service

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
    body := *openapiclient.NewCommonsInternalPostUsageInput() // CommonsInternalPostUsageInput | post usage input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1BillingAPI.PostUsageForInternalService(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1BillingAPI.PostUsageForInternalService``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PostUsageForInternalService`: CommonsInternalPostUsageOutput
    fmt.Fprintf(os.Stdout, "Response from `V1BillingAPI.PostUsageForInternalService`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostUsageForInternalServiceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsInternalPostUsageInput**](CommonsInternalPostUsageInput.md) | post usage input | 

### Return type

[**CommonsInternalPostUsageOutput**](CommonsInternalPostUsageOutput.md)

### Authorization

[WarpBuildServiceSecretAuth](../README.md#WarpBuildServiceSecretAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

