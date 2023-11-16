# \V1RunnersAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteRunner**](V1RunnersAPI.md#DeleteRunner) | **Delete** /runners/{id} | delete runner for the id. Current organization is figured from the authorization token
[**GetRunner**](V1RunnersAPI.md#GetRunner) | **Get** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token
[**GetRuntimes**](V1RunnersAPI.md#GetRuntimes) | **Get** /runners/runtimes | Get runtimes for runners of the organisation
[**ListRunners**](V1RunnersAPI.md#ListRunners) | **Get** /runners | ListRunners lists all the runners for an org.
[**SetupRunner**](V1RunnersAPI.md#SetupRunner) | **Post** /runners | Adds a new runner for a current organization
[**UpdateRunner**](V1RunnersAPI.md#UpdateRunner) | **Patch** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token



## DeleteRunner

> CommonsRunner DeleteRunner(ctx).Execute()

delete runner for the id. Current organization is figured from the authorization token

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
    resp, r, err := apiClient.V1RunnersAPI.DeleteRunner(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.DeleteRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.DeleteRunner`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerRequest struct via the builder pattern


### Return type

[**CommonsRunner**](CommonsRunner.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunner

> CommonsRunner GetRunner(ctx).Execute()

Get runner details for the id. Current organization is figured from the authorization token

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
    resp, r, err := apiClient.V1RunnersAPI.GetRunner(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.GetRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.GetRunner`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerRequest struct via the builder pattern


### Return type

[**CommonsRunner**](CommonsRunner.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRuntimes

> []CommonsRuntime GetRuntimes(ctx).Execute()

Get runtimes for runners of the organisation

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
    resp, r, err := apiClient.V1RunnersAPI.GetRuntimes(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.GetRuntimes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRuntimes`: []CommonsRuntime
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.GetRuntimes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetRuntimesRequest struct via the builder pattern


### Return type

[**[]CommonsRuntime**](CommonsRuntime.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRunners

> []CommonsRunner ListRunners(ctx).Execute()

ListRunners lists all the runners for an org.

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
    resp, r, err := apiClient.V1RunnersAPI.ListRunners(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.ListRunners``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunners`: []CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.ListRunners`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListRunnersRequest struct via the builder pattern


### Return type

[**[]CommonsRunner**](CommonsRunner.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SetupRunner

> CommonsRunner SetupRunner(ctx).Body(body).Execute()

Adds a new runner for a current organization



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
    body := *openapiclient.NewCommonsSetupRunnerInput() // CommonsSetupRunnerInput | Create new runner body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersAPI.SetupRunner(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.SetupRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetupRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.SetupRunner`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSetupRunnerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsSetupRunnerInput**](CommonsSetupRunnerInput.md) | Create new runner body | 

### Return type

[**CommonsRunner**](CommonsRunner.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateRunner

> CommonsRunner UpdateRunner(ctx).Body(body).Execute()

Get runner details for the id. Current organization is figured from the authorization token

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
    body := *openapiclient.NewCommonsUpdateRunnerInput() // CommonsUpdateRunnerInput | update runner body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersAPI.UpdateRunner(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.UpdateRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.UpdateRunner`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateRunnerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsUpdateRunnerInput**](CommonsUpdateRunnerInput.md) | update runner body | 

### Return type

[**CommonsRunner**](CommonsRunner.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

