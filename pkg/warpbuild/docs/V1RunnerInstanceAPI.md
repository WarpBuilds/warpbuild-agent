# \V1RunnerInstanceAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddRunnerInstance**](V1RunnerInstanceAPI.md#AddRunnerInstance) | **Post** /runner_instance | Add a new runner instance
[**AllocateRunnerInstance**](V1RunnerInstanceAPI.md#AllocateRunnerInstance) | **Post** /runner_instance/allocate | Allocate a runner instance
[**DeleteRunnerInstance**](V1RunnerInstanceAPI.md#DeleteRunnerInstance) | **Delete** /runner_instance/{id} | delete runner instance for the id. Current organization is figured from the authorization token
[**GetRunnerInstance**](V1RunnerInstanceAPI.md#GetRunnerInstance) | **Get** /runner_instance/{id} | Get runner instance by ID for the organization. Current organization is figured from the authorization token
[**GetRunnerInstanceAllocationDetails**](V1RunnerInstanceAPI.md#GetRunnerInstanceAllocationDetails) | **Get** /runners_instance/{id}/allocation_details | Get runner instance allocation details for the id
[**GetRunnerInstancePresignedLogUploadURL**](V1RunnerInstanceAPI.md#GetRunnerInstancePresignedLogUploadURL) | **Get** /runners_instance/{id}/presigned_log_upload_url | Gets a presigned url for uploading logs for a runner instance
[**GetRunnerInstanceTelemetry**](V1RunnerInstanceAPI.md#GetRunnerInstanceTelemetry) | **Get** /runner_instance/{id}/telemetry | Gets metrics log files for a runner instance with presigned URLs
[**GetRunnerLastJobProcessedMeta**](V1RunnerInstanceAPI.md#GetRunnerLastJobProcessedMeta) | **Get** /runner_instance/internal/{id}/last_job_processed_meta | Get runner last used job meta
[**RunnerInstanceCleanupHook**](V1RunnerInstanceAPI.md#RunnerInstanceCleanupHook) | **Post** /runners_instance/{id}/cleanup_hook | Get runner instance allocation details for the id



## AddRunnerInstance

> CommonsRunnerInstance AddRunnerInstance(ctx).Body(body).Execute()

Add a new runner instance



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
    body := *openapiclient.NewCommonsAddRunnerInstanceInput() // CommonsAddRunnerInstanceInput | Add runner instance body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.AddRunnerInstance(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.AddRunnerInstance``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AddRunnerInstance`: CommonsRunnerInstance
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.AddRunnerInstance`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAddRunnerInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsAddRunnerInstanceInput**](CommonsAddRunnerInstanceInput.md) | Add runner instance body | 

### Return type

[**CommonsRunnerInstance**](CommonsRunnerInstance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AllocateRunnerInstance

> CommonsRunnerInstance AllocateRunnerInstance(ctx).Body(body).Execute()

Allocate a runner instance



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
    body := *openapiclient.NewCommonsAllocateRunnerInstanceInput() // CommonsAllocateRunnerInstanceInput | Allocate runner instance body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.AllocateRunnerInstance(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.AllocateRunnerInstance``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AllocateRunnerInstance`: CommonsRunnerInstance
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.AllocateRunnerInstance`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAllocateRunnerInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsAllocateRunnerInstanceInput**](CommonsAllocateRunnerInstanceInput.md) | Allocate runner instance body | 

### Return type

[**CommonsRunnerInstance**](CommonsRunnerInstance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteRunnerInstance

> CommonsRunner DeleteRunnerInstance(ctx, id).Execute()

delete runner instance for the id. Current organization is figured from the authorization token

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
    id := "id_example" // string | Runner Instance ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.DeleteRunnerInstance(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.DeleteRunnerInstance``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteRunnerInstance`: CommonsRunner
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.DeleteRunnerInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Instance ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteRunnerInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunner**](CommonsRunner.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerInstance

> CommonsRunnerInstance GetRunnerInstance(ctx, id).Execute()

Get runner instance by ID for the organization. Current organization is figured from the authorization token

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
    id := "id_example" // string | Runner Instance ID

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerInstance(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerInstance``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerInstance`: CommonsRunnerInstance
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerInstance`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Runner Instance ID | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerInstanceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRunnerInstance**](CommonsRunnerInstance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerInstanceAllocationDetails

> CommonsRunnerInstanceAllocationDetails GetRunnerInstanceAllocationDetails(ctx, id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()

Get runner instance allocation details for the id

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
    id := "id_example" // string | runner instance id
    xPOLLINGSECRET := "xPOLLINGSECRET_example" // string | polling secret for validation

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails(context.Background(), id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerInstanceAllocationDetails`: CommonsRunnerInstanceAllocationDetails
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerInstanceAllocationDetails`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerInstanceAllocationDetailsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xPOLLINGSECRET** | **string** | polling secret for validation | 

### Return type

[**CommonsRunnerInstanceAllocationDetails**](CommonsRunnerInstanceAllocationDetails.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerInstancePresignedLogUploadURL

> CommonsGetPresignedLogUploadURLOutput GetRunnerInstancePresignedLogUploadURL(ctx, id).XPOLLINGSECRET(xPOLLINGSECRET).LogFileName(logFileName).Execute()

Gets a presigned url for uploading logs for a runner instance

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
    id := "id_example" // string | runner instance id
    xPOLLINGSECRET := "xPOLLINGSECRET_example" // string | polling secret for validation
    logFileName := "logFileName_example" // string | Log file name (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerInstancePresignedLogUploadURL(context.Background(), id).XPOLLINGSECRET(xPOLLINGSECRET).LogFileName(logFileName).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerInstancePresignedLogUploadURL``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerInstancePresignedLogUploadURL`: CommonsGetPresignedLogUploadURLOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerInstancePresignedLogUploadURL`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerInstancePresignedLogUploadURLRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xPOLLINGSECRET** | **string** | polling secret for validation | 
 **logFileName** | **string** | Log file name | 

### Return type

[**CommonsGetPresignedLogUploadURLOutput**](CommonsGetPresignedLogUploadURLOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerInstanceTelemetry

> CommonsGetRunnerInstanceTelemetryOutput GetRunnerInstanceTelemetry(ctx, id).MetricType(metricType).Execute()

Gets metrics log files for a runner instance with presigned URLs

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
    id := "id_example" // string | runner instance id
    metricType := "metricType_example" // string | metric type (logs, metrics, or gha-logs) (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerInstanceTelemetry(context.Background(), id).MetricType(metricType).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerInstanceTelemetry``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerInstanceTelemetry`: CommonsGetRunnerInstanceTelemetryOutput
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerInstanceTelemetry`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerInstanceTelemetryRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **metricType** | **string** | metric type (logs, metrics, or gha-logs) | 

### Return type

[**CommonsGetRunnerInstanceTelemetryOutput**](CommonsGetRunnerInstanceTelemetryOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetRunnerLastJobProcessedMeta

> CommonsLastJobProcessedMeta GetRunnerLastJobProcessedMeta(ctx, id).Execute()

Get runner last used job meta

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
    id := "id_example" // string | runner instance id

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.GetRunnerLastJobProcessedMeta(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.GetRunnerLastJobProcessedMeta``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetRunnerLastJobProcessedMeta`: CommonsLastJobProcessedMeta
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.GetRunnerLastJobProcessedMeta`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetRunnerLastJobProcessedMetaRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsLastJobProcessedMeta**](CommonsLastJobProcessedMeta.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RunnerInstanceCleanupHook

> map[string]interface{} RunnerInstanceCleanupHook(ctx, id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()

Get runner instance allocation details for the id

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
    id := "id_example" // string | runner instance id
    xPOLLINGSECRET := "xPOLLINGSECRET_example" // string | polling secret for validation

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1RunnerInstanceAPI.RunnerInstanceCleanupHook(context.Background(), id).XPOLLINGSECRET(xPOLLINGSECRET).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1RunnerInstanceAPI.RunnerInstanceCleanupHook``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RunnerInstanceCleanupHook`: map[string]interface{}
    fmt.Fprintf(os.Stdout, "Response from `V1RunnerInstanceAPI.RunnerInstanceCleanupHook`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | runner instance id | 

### Other Parameters

Other parameters are passed through a pointer to a apiRunnerInstanceCleanupHookRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xPOLLINGSECRET** | **string** | polling secret for validation | 

### Return type

**map[string]interface{}**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

