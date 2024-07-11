# \V1WorkflowsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetPullRequestAuthURL**](V1WorkflowsAPI.md#GetPullRequestAuthURL) | **Get** /workflows/pr-auth-url | Get auth url required for GH PR
[**ListWorkflows**](V1WorkflowsAPI.md#ListWorkflows) | **Get** /workflows | Lists all workflows (workflows) for organization according to repo
[**PullWorkflows**](V1WorkflowsAPI.md#PullWorkflows) | **Patch** /workflows/pull | Pulls all workflows from the provider to the database
[**WarpWorkflows**](V1WorkflowsAPI.md#WarpWorkflows) | **Patch** /workflows/warp | Warps workflows for organization according to given internal workflow ids



## GetPullRequestAuthURL

> map[string]string GetPullRequestAuthURL(ctx).Execute()

Get auth url required for GH PR



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
    resp, r, err := apiClient.V1WorkflowsAPI.GetPullRequestAuthURL(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1WorkflowsAPI.GetPullRequestAuthURL``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetPullRequestAuthURL`: map[string]string
    fmt.Fprintf(os.Stdout, "Response from `V1WorkflowsAPI.GetPullRequestAuthURL`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetPullRequestAuthURLRequest struct via the builder pattern


### Return type

**map[string]string**

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListWorkflows

> ListWorkflowsResponse ListWorkflows(ctx).Execute()

Lists all workflows (workflows) for organization according to repo



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
    resp, r, err := apiClient.V1WorkflowsAPI.ListWorkflows(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1WorkflowsAPI.ListWorkflows``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListWorkflows`: ListWorkflowsResponse
    fmt.Fprintf(os.Stdout, "Response from `V1WorkflowsAPI.ListWorkflows`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListWorkflowsRequest struct via the builder pattern


### Return type

[**ListWorkflowsResponse**](ListWorkflowsResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PullWorkflows

> PullWorkflows(ctx).Execute()

Pulls all workflows from the provider to the database



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
    r, err := apiClient.V1WorkflowsAPI.PullWorkflows(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1WorkflowsAPI.PullWorkflows``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiPullWorkflowsRequest struct via the builder pattern


### Return type

 (empty response body)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## WarpWorkflows

> WarpWorkflowsResponse WarpWorkflows(ctx).Body(body).Execute()

Warps workflows for organization according to given internal workflow ids



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
    body := *openapiclient.NewWarpWorkflowsRequest("RunnerId_example", []string{"WorkflowIds_example"}) // WarpWorkflowsRequest | Warp workflows options

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1WorkflowsAPI.WarpWorkflows(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1WorkflowsAPI.WarpWorkflows``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `WarpWorkflows`: WarpWorkflowsResponse
    fmt.Fprintf(os.Stdout, "Response from `V1WorkflowsAPI.WarpWorkflows`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiWarpWorkflowsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**WarpWorkflowsRequest**](WarpWorkflowsRequest.md) | Warp workflows options | 

### Return type

[**WarpWorkflowsResponse**](WarpWorkflowsResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

