# \V1UiAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetBannerMessages**](V1UiAPI.md#GetBannerMessages) | **Get** /ui/banner-messages | Get specific banner messages for UI/Org or all



## GetBannerMessages

> []CommonsBannerMessage GetBannerMessages(ctx).Execute()

Get specific banner messages for UI/Org or all

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

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1UiAPI.GetBannerMessages(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1UiAPI.GetBannerMessages``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetBannerMessages`: []CommonsBannerMessage
    fmt.Fprintf(os.Stdout, "Response from `V1UiAPI.GetBannerMessages`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetBannerMessagesRequest struct via the builder pattern


### Return type

[**[]CommonsBannerMessage**](CommonsBannerMessage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

