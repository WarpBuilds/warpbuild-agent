# \V1RunnerImageVersionsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteRunnerImageVersion**](V1RunnerImageVersionsAPI.md#DeleteRunnerImageVersion) | **Delete** /runner-image-versions/{id} | Delete runner image version details for the id.
[**GetRunnerImageVersion**](V1RunnerImageVersionsAPI.md#GetRunnerImageVersion) | **Get** /runner-image-versions/{id} | Get runner image version details for the id.
[**ListRunnerImageVersions**](V1RunnerImageVersionsAPI.md#ListRunnerImageVersions) | **Get** /runner-image-versions | List all runner image versions.
[**UpdateRunnerImageVersion**](V1RunnerImageVersionsAPI.md#UpdateRunnerImageVersion) | **Patch** /runner-image-versions/{id} | Update runner image version details for the id.



## DeleteRunnerImageVersion

> TypesGenericSuccessMessage DeleteRunnerImageVersion(ctx, id).Execute()

Delete runner image version details for the id.

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
    id := "id_example" // string | Runner Image Version ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImageVersionsAPI.DeleteRunnerImageVersion(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImageVersionsAPI.DeleteRunnerImageVersion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunnerImageVersion`: TypesGenericSuccessMessage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImageVersionsAPI.DeleteRunnerImageVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Version ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerImageVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**TypesGenericSuccessMessage**](TypesGenericSuccessMessage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerImageVersion

> CommonsRunnerImageVersion GetRunnerImageVersion(ctx, id).Execute()

Get runner image version details for the id.

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
    id := "id_example" // string | Runner Image Version ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImageVersionsAPI.GetRunnerImageVersion(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImageVersionsAPI.GetRunnerImageVersion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerImageVersion`: CommonsRunnerImageVersion
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImageVersionsAPI.GetRunnerImageVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Version ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerImageVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunnerImageVersion**](CommonsRunnerImageVersion.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRunnerImageVersions

> CommonsListRunnerImageVersionsOutput ListRunnerImageVersions(ctx).RunnerImageId(runnerImageId).Execute()

List all runner image versions.

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
    runnerImageId := "runnerImageId_example" // string | Runner Image ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImageVersionsAPI.ListRunnerImageVersions(context.Background()).RunnerImageId(runnerImageId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImageVersionsAPI.ListRunnerImageVersions``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunnerImageVersions`: CommonsListRunnerImageVersionsOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImageVersionsAPI.ListRunnerImageVersions`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListRunnerImageVersionsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **runnerImageId** | **string** | Runner Image ID | 

### Return type

[**CommonsListRunnerImageVersionsOutput**](CommonsListRunnerImageVersionsOutput.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateRunnerImageVersion

> CommonsRunnerImageVersion UpdateRunnerImageVersion(ctx, id).Body(body).Execute()

Update runner image version details for the id.

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
    id := "id_example" // string | Runner Image Version ID
    body := *openapiclient.NewCommonsUpdateRunnerImageVersionInput() // CommonsUpdateRunnerImageVersionInput | Runner Image Version

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImageVersionsAPI.UpdateRunnerImageVersion(context.Background(), id).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImageVersionsAPI.UpdateRunnerImageVersion``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunnerImageVersion`: CommonsRunnerImageVersion
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImageVersionsAPI.UpdateRunnerImageVersion`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Version ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateRunnerImageVersionRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**CommonsUpdateRunnerImageVersionInput**](CommonsUpdateRunnerImageVersionInput.md) | Runner Image Version | 

### Return type

[**CommonsRunnerImageVersion**](CommonsRunnerImageVersion.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

