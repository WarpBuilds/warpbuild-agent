# \V1RunnerImagesAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateRunnerImage**](V1RunnerImagesAPI.md#CreateRunnerImage) | **Post** /runner-images | Create a new runner image.
[**DeleteRunnerImage**](V1RunnerImagesAPI.md#DeleteRunnerImage) | **Delete** /runner-images/{id} | Delete runner image details for the id.
[**GetRunnerImage**](V1RunnerImagesAPI.md#GetRunnerImage) | **Get** /runner-images/{id} | Get runner image details for the id.
[**ListRunnerImages**](V1RunnerImagesAPI.md#ListRunnerImages) | **Get** /runner-images | List all runner images.
[**UpdateRunnerImage**](V1RunnerImagesAPI.md#UpdateRunnerImage) | **Put** /runner-images/{id} | Update runner image details for the id.



## CreateRunnerImage

> CommonsRunnerImage CreateRunnerImage(ctx).Body(body).Execute()

Create a new runner image.

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
    body := *openapiclient.NewCommonsCreateRunnerImageInput() // CommonsCreateRunnerImageInput | Runner Image

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagesAPI.CreateRunnerImage(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagesAPI.CreateRunnerImage``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateRunnerImage`: CommonsRunnerImage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagesAPI.CreateRunnerImage`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateRunnerImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsCreateRunnerImageInput**](CommonsCreateRunnerImageInput.md) | Runner Image | 

### Return type

[**CommonsRunnerImage**](CommonsRunnerImage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteRunnerImage

> CommonsRunnerImage DeleteRunnerImage(ctx, id).Execute()

Delete runner image details for the id.

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
    id := "id_example" // string | Runner Image ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagesAPI.DeleteRunnerImage(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagesAPI.DeleteRunnerImage``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunnerImage`: CommonsRunnerImage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagesAPI.DeleteRunnerImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunnerImage**](CommonsRunnerImage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerImage

> CommonsRunnerImage GetRunnerImage(ctx, id).Execute()

Get runner image details for the id.

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
    id := "id_example" // string | Runner Image ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagesAPI.GetRunnerImage(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagesAPI.GetRunnerImage``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerImage`: CommonsRunnerImage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagesAPI.GetRunnerImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunnerImage**](CommonsRunnerImage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRunnerImages

> CommonsListRunnerImagesOutput ListRunnerImages(ctx).Execute()

List all runner images.

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
    resp, r, err := apiClient.V1RunnerImagesAPI.ListRunnerImages(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagesAPI.ListRunnerImages``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunnerImages`: CommonsListRunnerImagesOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagesAPI.ListRunnerImages`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListRunnerImagesRequest struct via the builder pattern


### Return type

[**CommonsListRunnerImagesOutput**](CommonsListRunnerImagesOutput.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateRunnerImage

> CommonsRunnerImage UpdateRunnerImage(ctx, id).Body(body).Execute()

Update runner image details for the id.

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
    id := "id_example" // string | Runner Image ID
    body := *openapiclient.NewCommonsUpdateRunnerImageInput() // CommonsUpdateRunnerImageInput | Runner Image

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagesAPI.UpdateRunnerImage(context.Background(), id).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagesAPI.UpdateRunnerImage``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunnerImage`: CommonsRunnerImage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagesAPI.UpdateRunnerImage`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateRunnerImageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**CommonsUpdateRunnerImageInput**](CommonsUpdateRunnerImageInput.md) | Runner Image | 

### Return type

[**CommonsRunnerImage**](CommonsRunnerImage.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

