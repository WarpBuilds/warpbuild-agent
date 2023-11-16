# \V1DebuggerAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DebugPublishEvent**](V1DebuggerAPI.md#DebugPublishEvent) | **Post** /debugger/events/publish | Publish an event to the event bus



## DebugPublishEvent

> TypesGenericSuccessMessage DebugPublishEvent(ctx).Body(body).Execute()

Publish an event to the event bus

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
    body := *openapiclient.NewDebuggerPublishEventInput() // DebuggerPublishEventInput | Event to publish

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1DebuggerAPI.DebugPublishEvent(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1DebuggerAPI.DebugPublishEvent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DebugPublishEvent`: TypesGenericSuccessMessage
    fmt.Fprintf(os.Stdout, "Response from `V1DebuggerAPI.DebugPublishEvent`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDebugPublishEventRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**DebuggerPublishEventInput**](DebuggerPublishEventInput.md) | Event to publish | 

### Return type

[**TypesGenericSuccessMessage**](TypesGenericSuccessMessage.md)

### Authorization

[WarpBuildAdminSecretAuth](../README.md#WarpBuildAdminSecretAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

