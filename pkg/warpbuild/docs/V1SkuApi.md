# \V1SkuApi

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetSku**](V1SkuApi.md#GetSku) | **Get** /sku/{id} | Get default group for runner set
[**ListSku**](V1SkuApi.md#ListSku) | **Get** /sku | ListAllSku lists all the runners sku for an org.



## GetSku

> CommonsInstanceSku GetSku(ctx).Execute()

Get default group for runner set

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
    resp, r, err := apiClient.V1SkuApi.GetSku(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SkuApi.GetSku``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetSku`: CommonsInstanceSku
    fmt.Fprintf(os.Stdout, "Response from `V1SkuApi.GetSku`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetSkuRequest struct via the builder pattern


### Return type

[**CommonsInstanceSku**](CommonsInstanceSku.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListSku

> []CommonsInstanceSku ListSku(ctx).Ids(ids).Cores(cores).Memory(memory).Arch(arch).Os(os).Manufacturer(manufacturer).PerformanceCategory(performanceCategory).HasGpu(hasGpu).Burstable(burstable).IncludeInternal(includeInternal).Execute()

ListAllSku lists all the runners sku for an org.

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
    ids := []string{"Inner_example"} // []string | list of sku id's (optional)
    cores := int32(56) // int32 | cores (optional)
    memory := int32(56) // int32 | memory (optional)
    arch := "arch_example" // string | Architectures (optional)
    os := []string{"Os_example"} // []string | operating system (optional)
    manufacturer := "manufacturer_example" // string | Images (optional)
    performanceCategory := "performanceCategory_example" // string | performance category (optional)
    hasGpu := true // bool | has gpu (optional)
    burstable := true // bool | burstable (optional)
    includeInternal := true // bool | include internal skus (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SkuApi.ListSku(context.Background()).Ids(ids).Cores(cores).Memory(memory).Arch(arch).Os(os).Manufacturer(manufacturer).PerformanceCategory(performanceCategory).HasGpu(hasGpu).Burstable(burstable).IncludeInternal(includeInternal).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SkuApi.ListSku``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListSku`: []CommonsInstanceSku
    fmt.Fprintf(os.Stdout, "Response from `V1SkuApi.ListSku`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListSkuRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ids** | **[]string** | list of sku id&#39;s | 
 **cores** | **int32** | cores | 
 **memory** | **int32** | memory | 
 **arch** | **string** | Architectures | 
 **os** | **[]string** | operating system | 
 **manufacturer** | **string** | Images | 
 **performanceCategory** | **string** | performance category | 
 **hasGpu** | **bool** | has gpu | 
 **burstable** | **bool** | burstable | 
 **includeInternal** | **bool** | include internal skus | 

### Return type

[**[]CommonsInstanceSku**](CommonsInstanceSku.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

