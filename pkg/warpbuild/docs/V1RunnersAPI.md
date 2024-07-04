<<<<<<< HEAD
# \V1RunnersAPI
=======
# \V1RunnersApi
>>>>>>> prajjwal-warp-323

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
<<<<<<< HEAD
[**DeleteRunner**](V1RunnersAPI.md#DeleteRunner) | **Delete** /runners/{id} | delete runner for the id. Current organization is figured from the authorization token
[**GetRunner**](V1RunnersAPI.md#GetRunner) | **Get** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token
[**GetRuntimes**](V1RunnersAPI.md#GetRuntimes) | **Get** /runners/runtimes | Get runtimes for runners of the organisation
[**ListRunners**](V1RunnersAPI.md#ListRunners) | **Get** /runners | ListRunners lists all the runners for an org.
[**SetupRunner**](V1RunnersAPI.md#SetupRunner) | **Post** /runners | Adds a new runner for a current organization
[**UpdateRunner**](V1RunnersAPI.md#UpdateRunner) | **Patch** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token



## DeleteRunner

> CommonsRunner DeleteRunner(ctx).Execute()
=======
[**ComputeCustomRunnerRate**](V1RunnersApi.md#ComputeCustomRunnerRate) | **Post** /runners/cost/calculator | Get ComputeCustomRunnerRate details
[**DeleteRunner**](V1RunnersApi.md#DeleteRunner) | **Delete** /runners/{id} | delete runner for the id. Current organization is figured from the authorization token
[**GetRunner**](V1RunnersApi.md#GetRunner) | **Get** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token
[**GetRunnerSetDefaultGroup**](V1RunnersApi.md#GetRunnerSetDefaultGroup) | **Get** /runners/default-group | Get default group for runner set
[**GetRunnersUsage**](V1RunnersApi.md#GetRunnersUsage) | **Get** /runners/usage | Get runtimes for runners of the organisation
[**ListRunners**](V1RunnersApi.md#ListRunners) | **Get** /runners | ListRunners lists all the runners for an org.
[**SetRunnerSetDefaultGroup**](V1RunnersApi.md#SetRunnerSetDefaultGroup) | **Patch** /runners/default-group | Set default group for runner set
[**SetupRunner**](V1RunnersApi.md#SetupRunner) | **Post** /runners | Adds a new runner for a current organization
[**UpdateRunner**](V1RunnersApi.md#UpdateRunner) | **Patch** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token



## ComputeCustomRunnerRate

> CommonsRateCalculationOutput ComputeCustomRunnerRate(ctx).Body(body).Execute()

Get ComputeCustomRunnerRate details

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
    body := *openapiclient.NewCommonsRateCalculationInput() // CommonsRateCalculationInput | Compute custom runner rate

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersApi.ComputeCustomRunnerRate(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.ComputeCustomRunnerRate``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ComputeCustomRunnerRate`: CommonsRateCalculationOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.ComputeCustomRunnerRate`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiComputeCustomRunnerRateRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsRateCalculationInput**](CommonsRateCalculationInput.md) | Compute custom runner rate | 

### Return type

[**CommonsRateCalculationOutput**](CommonsRateCalculationOutput.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteRunner

> CommonsRunner DeleteRunner(ctx, id).Execute()
>>>>>>> prajjwal-warp-323

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
<<<<<<< HEAD

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersAPI.DeleteRunner(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.DeleteRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.DeleteRunner`: %v\n", resp)
=======
    id := "id_example" // string | Runner ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersApi.DeleteRunner(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.DeleteRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.DeleteRunner`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters

<<<<<<< HEAD
This endpoint does not need any parameter.
=======

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner ID | 
>>>>>>> prajjwal-warp-323

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerRequest struct via the builder pattern


<<<<<<< HEAD
=======
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


>>>>>>> prajjwal-warp-323
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

<<<<<<< HEAD
> CommonsRunner GetRunner(ctx).Execute()
=======
> CommonsRunner GetRunner(ctx, id).Execute()
>>>>>>> prajjwal-warp-323

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
<<<<<<< HEAD

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersAPI.GetRunner(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.GetRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.GetRunner`: %v\n", resp)
=======
    id := "id_example" // string | Runner ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersApi.GetRunner(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.GetRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.GetRunner`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters

<<<<<<< HEAD
This endpoint does not need any parameter.
=======

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner ID | 
>>>>>>> prajjwal-warp-323

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerRequest struct via the builder pattern


<<<<<<< HEAD
=======
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


>>>>>>> prajjwal-warp-323
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


<<<<<<< HEAD
## GetRuntimes

> []CommonsRuntime GetRuntimes(ctx).Execute()

Get runtimes for runners of the organisation
=======
## GetRunnerSetDefaultGroup

> CommonsRunnerSetDefaultGroup GetRunnerSetDefaultGroup(ctx).Execute()

Get default group for runner set
>>>>>>> prajjwal-warp-323

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
<<<<<<< HEAD
    resp, r, err := apiClient.V1RunnersAPI.GetRuntimes(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.GetRuntimes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRuntimes`: []CommonsRuntime
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.GetRuntimes`: %v\n", resp)
=======
    resp, r, err := apiClient.V1RunnersApi.GetRunnerSetDefaultGroup(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.GetRunnerSetDefaultGroup``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerSetDefaultGroup`: CommonsRunnerSetDefaultGroup
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.GetRunnerSetDefaultGroup`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

<<<<<<< HEAD
Other parameters are passed through a pointer to a apiGetRuntimesRequest struct via the builder pattern
=======
Other parameters are passed through a pointer to a apiGetRunnerSetDefaultGroupRequest struct via the builder pattern
>>>>>>> prajjwal-warp-323


### Return type

<<<<<<< HEAD
[**[]CommonsRuntime**](CommonsRuntime.md)
=======
[**CommonsRunnerSetDefaultGroup**](CommonsRunnerSetDefaultGroup.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnersUsage

> CommonsRunnersUsage GetRunnersUsage(ctx).StartDate(startDate).EndDate(endDate).CapacityTypes(capacityTypes).Archs(archs).Images(images).Cores(cores).RunnerTypes(runnerTypes).SearchTerm(searchTerm).Execute()

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
    startDate := "startDate_example" // string | Date range start
    endDate := "endDate_example" // string | Date range end
    capacityTypes := []string{"CapacityTypes_example"} // []string | Capacity types (optional)
    archs := []string{"Archs_example"} // []string | Architectures (optional)
    images := []string{"Images_example"} // []string | Images (optional)
    cores := []int32{int32(123)} // []int32 | Cores (optional)
    runnerTypes := []string{"RunnerTypes_example"} // []string | Runner types (optional)
    searchTerm := "searchTerm_example" // string | Search term (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersApi.GetRunnersUsage(context.Background()).StartDate(startDate).EndDate(endDate).CapacityTypes(capacityTypes).Archs(archs).Images(images).Cores(cores).RunnerTypes(runnerTypes).SearchTerm(searchTerm).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.GetRunnersUsage``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnersUsage`: CommonsRunnersUsage
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.GetRunnersUsage`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnersUsageRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **startDate** | **string** | Date range start | 
 **endDate** | **string** | Date range end | 
 **capacityTypes** | **[]string** | Capacity types | 
 **archs** | **[]string** | Architectures | 
 **images** | **[]string** | Images | 
 **cores** | **[]int32** | Cores | 
 **runnerTypes** | **[]string** | Runner types | 
 **searchTerm** | **string** | Search term | 

### Return type

[**CommonsRunnersUsage**](CommonsRunnersUsage.md)
>>>>>>> prajjwal-warp-323

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
<<<<<<< HEAD
    resp, r, err := apiClient.V1RunnersAPI.ListRunners(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.ListRunners``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunners`: []CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.ListRunners`: %v\n", resp)
=======
    resp, r, err := apiClient.V1RunnersApi.ListRunners(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.ListRunners``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListRunners`: []CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.ListRunners`: %v\n", resp)
>>>>>>> prajjwal-warp-323
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


<<<<<<< HEAD
=======
## SetRunnerSetDefaultGroup

> CommonsRunnerSetDefaultGroup SetRunnerSetDefaultGroup(ctx).Body(body).Execute()

Set default group for runner set

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
    body := *openapiclient.NewCommonsSetRunnerSetDefaultGroupInput() // CommonsSetRunnerSetDefaultGroupInput | Set default group for runner set body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnersApi.SetRunnerSetDefaultGroup(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.SetRunnerSetDefaultGroup``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetRunnerSetDefaultGroup`: CommonsRunnerSetDefaultGroup
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.SetRunnerSetDefaultGroup`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSetRunnerSetDefaultGroupRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsSetRunnerSetDefaultGroupInput**](CommonsSetRunnerSetDefaultGroupInput.md) | Set default group for runner set body | 

### Return type

[**CommonsRunnerSetDefaultGroup**](CommonsRunnerSetDefaultGroup.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


>>>>>>> prajjwal-warp-323
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
<<<<<<< HEAD
    resp, r, err := apiClient.V1RunnersAPI.SetupRunner(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.SetupRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetupRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.SetupRunner`: %v\n", resp)
=======
    resp, r, err := apiClient.V1RunnersApi.SetupRunner(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.SetupRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SetupRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.SetupRunner`: %v\n", resp)
>>>>>>> prajjwal-warp-323
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

<<<<<<< HEAD
> CommonsRunner UpdateRunner(ctx).Body(body).Execute()
=======
> CommonsRunner UpdateRunner(ctx, id).Body(body).Execute()
>>>>>>> prajjwal-warp-323

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
<<<<<<< HEAD
=======
    id := "id_example" // string | Runner ID
>>>>>>> prajjwal-warp-323
    body := *openapiclient.NewCommonsUpdateRunnerInput() // CommonsUpdateRunnerInput | update runner body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1RunnersAPI.UpdateRunner(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersAPI.UpdateRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersAPI.UpdateRunner`: %v\n", resp)
=======
    resp, r, err := apiClient.V1RunnersApi.UpdateRunner(context.Background(), id).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnersApi.UpdateRunner``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateRunner`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnersApi.UpdateRunner`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters


<<<<<<< HEAD
=======
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner ID | 
>>>>>>> prajjwal-warp-323

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateRunnerRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
<<<<<<< HEAD
=======

>>>>>>> prajjwal-warp-323
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

