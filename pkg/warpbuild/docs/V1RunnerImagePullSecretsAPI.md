# \V1RunnerImagePullSecretsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateRunnerImagePullSecret**](V1RunnerImagePullSecretsAPI.md#CreateRunnerImagePullSecret) | **Post** /runner-image-pull-secrets | Create a new runner image pull secret.
[**DeleteRunnerImagePullSecret**](V1RunnerImagePullSecretsAPI.md#DeleteRunnerImagePullSecret) | **Delete** /runner-image-pull-secrets/{id} | Delete runner image pull secret details for the id.
[**GetRunnerImagePullSecret**](V1RunnerImagePullSecretsAPI.md#GetRunnerImagePullSecret) | **Get** /runner-image-pull-secrets/{id} | Get runner image pull secret details for the id.
[**ListRunnerImagePullSecrets**](V1RunnerImagePullSecretsAPI.md#ListRunnerImagePullSecrets) | **Get** /runner-image-pull-secrets | List all runner image pull secrets.
[**UpdateRunnerImagePullSecret**](V1RunnerImagePullSecretsAPI.md#UpdateRunnerImagePullSecret) | **Put** /runner-image-pull-secrets/{id} | Update runner image pull secret details for the id.



## CreateRunnerImagePullSecret

> CommonsRunnerImagePullSecret CreateRunnerImagePullSecret(ctx).Body(body).Execute()

Create a new runner image pull secret.

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
    body := *openapiclient.NewCommonsCreateRunnerImagePullSecretInput("Alias_example", "Type_example") // CommonsCreateRunnerImagePullSecretInput | Runner Image Pull Secret

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagePullSecretsAPI.CreateRunnerImagePullSecret(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagePullSecretsAPI.CreateRunnerImagePullSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateRunnerImagePullSecret`: CommonsRunnerImagePullSecret
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagePullSecretsAPI.CreateRunnerImagePullSecret`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateRunnerImagePullSecretRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsCreateRunnerImagePullSecretInput**](CommonsCreateRunnerImagePullSecretInput.md) | Runner Image Pull Secret | 

### Return type

[**CommonsRunnerImagePullSecret**](CommonsRunnerImagePullSecret.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteRunnerImagePullSecret

> TypesGenericSuccessMessage DeleteRunnerImagePullSecret(ctx, id).Execute()

Delete runner image pull secret details for the id.

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
    id := "id_example" // string | Runner Image Pull Secret ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagePullSecretsAPI.DeleteRunnerImagePullSecret(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagePullSecretsAPI.DeleteRunnerImagePullSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunnerImagePullSecret`: TypesGenericSuccessMessage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagePullSecretsAPI.DeleteRunnerImagePullSecret`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Pull Secret ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerImagePullSecretRequest struct via the builder pattern


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


## GetRunnerImagePullSecret

> CommonsRunnerImagePullSecret GetRunnerImagePullSecret(ctx, id).Execute()

Get runner image pull secret details for the id.

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
    id := "id_example" // string | Runner Image Pull Secret ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagePullSecretsAPI.GetRunnerImagePullSecret(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagePullSecretsAPI.GetRunnerImagePullSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerImagePullSecret`: CommonsRunnerImagePullSecret
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagePullSecretsAPI.GetRunnerImagePullSecret`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Pull Secret ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerImagePullSecretRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunnerImagePullSecret**](CommonsRunnerImagePullSecret.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListRunnerImagePullSecrets

> CommonsListRunnerImagePullSecretsOutput ListRunnerImagePullSecrets(ctx).Execute()

List all runner image pull secrets.

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
    resp, r, err := apiClient.V1RunnerImagePullSecretsAPI.ListRunnerImagePullSecrets(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagePullSecretsAPI.ListRunnerImagePullSecrets``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunnerImagePullSecrets`: CommonsListRunnerImagePullSecretsOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagePullSecretsAPI.ListRunnerImagePullSecrets`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListRunnerImagePullSecretsRequest struct via the builder pattern


### Return type

[**CommonsListRunnerImagePullSecretsOutput**](CommonsListRunnerImagePullSecretsOutput.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateRunnerImagePullSecret

> CommonsRunnerImagePullSecret UpdateRunnerImagePullSecret(ctx, id).Body(body).Execute()

Update runner image pull secret details for the id.

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
    id := "id_example" // string | Runner Image Pull Secret ID
    body := *openapiclient.NewCommonsUpdateRunnerImagePullSecretInput() // CommonsUpdateRunnerImagePullSecretInput | Runner Image Pull Secret

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerImagePullSecretsAPI.UpdateRunnerImagePullSecret(context.Background(), id).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerImagePullSecretsAPI.UpdateRunnerImagePullSecret``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunnerImagePullSecret`: CommonsRunnerImagePullSecret
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerImagePullSecretsAPI.UpdateRunnerImagePullSecret`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Image Pull Secret ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateRunnerImagePullSecretRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**CommonsUpdateRunnerImagePullSecretInput**](CommonsUpdateRunnerImagePullSecretInput.md) | Runner Image Pull Secret | 

### Return type

[**CommonsRunnerImagePullSecret**](CommonsRunnerImagePullSecret.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

