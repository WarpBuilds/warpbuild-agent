# \V1InsightsIntegrationsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GitHubCallback**](V1InsightsIntegrationsAPI.md#GitHubCallback) | **Post** /insights/integrations/github/callback | GitHub callback for insights



## GitHubCallback

> AuthUserResponse GitHubCallback(ctx).Body(body).Execute()

GitHub callback for insights



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
    body := *openapiclient.NewInsightsCallbackInput("Code_example", "SetupAction_example") // InsightsCallbackInput | Callback input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1InsightsIntegrationsAPI.GitHubCallback(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1InsightsIntegrationsAPI.GitHubCallback``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GitHubCallback`: AuthUserResponse
    fmt.Fprintf(os.Stdout, "Response from `V1InsightsIntegrationsAPI.GitHubCallback`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGitHubCallbackRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**InsightsCallbackInput**](InsightsCallbackInput.md) | Callback input | 

### Return type

[**AuthUserResponse**](AuthUserResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

