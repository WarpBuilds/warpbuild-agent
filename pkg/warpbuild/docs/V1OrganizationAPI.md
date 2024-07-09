# \V1OrganizationApi

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateOrganization**](V1OrganizationApi.md#CreateOrganization) | **Post** /organization | Adds a new organisation for a current user
[**GetOrganization**](V1OrganizationApi.md#GetOrganization) | **Get** /organization | Get organization details for the current organization. Current organization is figured from the authorization token
[**ListOrgUsers**](V1OrganizationApi.md#ListOrgUsers) | **Get** /organization/users | ListOrgUsers list the users for the current organization
[**ListUserOrganizations**](V1OrganizationApi.md#ListUserOrganizations) | **Get** /organizations | ListUserOrganizations lists all the organization user has access to.
[**UpdateOrganization**](V1OrganizationApi.md#UpdateOrganization) | **Patch** /organization | Updates existing organization based on the fields provided.



## CreateOrganization

> SwitchOrganizationResponse CreateOrganization(ctx).Execute()

Adds a new organisation for a current user



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
    resp, r, err := apiClient.V1OrganizationApi.CreateOrganization(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1OrganizationApi.CreateOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateOrganization`: SwitchOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `V1OrganizationApi.CreateOrganization`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiCreateOrganizationRequest struct via the builder pattern


### Return type

[**SwitchOrganizationResponse**](SwitchOrganizationResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetOrganization

> CommonsOrganization GetOrganization(ctx).Execute()

Get organization details for the current organization. Current organization is figured from the authorization token

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
    resp, r, err := apiClient.V1OrganizationApi.GetOrganization(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1OrganizationApi.GetOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetOrganization`: CommonsOrganization
    fmt.Fprintf(os.Stdout, "Response from `V1OrganizationApi.GetOrganization`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetOrganizationRequest struct via the builder pattern


### Return type

[**CommonsOrganization**](CommonsOrganization.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrgUsers

> []V1ListUsersForOrganizationResult ListOrgUsers(ctx).Execute()

ListOrgUsers list the users for the current organization

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
    resp, r, err := apiClient.V1OrganizationApi.ListOrgUsers(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1OrganizationApi.ListOrgUsers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrgUsers`: []V1ListUsersForOrganizationResult
    fmt.Fprintf(os.Stdout, "Response from `V1OrganizationApi.ListOrgUsers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOrgUsersRequest struct via the builder pattern


### Return type

[**[]V1ListUsersForOrganizationResult**](V1ListUsersForOrganizationResult.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListUserOrganizations

> []V1Organization ListUserOrganizations(ctx).Execute()

ListUserOrganizations lists all the organization user has access to.

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
    resp, r, err := apiClient.V1OrganizationApi.ListUserOrganizations(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1OrganizationApi.ListUserOrganizations``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListUserOrganizations`: []V1Organization
    fmt.Fprintf(os.Stdout, "Response from `V1OrganizationApi.ListUserOrganizations`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListUserOrganizationsRequest struct via the builder pattern


### Return type

[**[]V1Organization**](V1Organization.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateOrganization

> CommonsOrganization UpdateOrganization(ctx).Body(body).Execute()

Updates existing organization based on the fields provided.



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
    body := *openapiclient.NewUpdateOrganizationRequest() // UpdateOrganizationRequest | Update existing organization body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1OrganizationApi.UpdateOrganization(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1OrganizationApi.UpdateOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateOrganization`: CommonsOrganization
    fmt.Fprintf(os.Stdout, "Response from `V1OrganizationApi.UpdateOrganization`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**UpdateOrganizationRequest**](UpdateOrganizationRequest.md) | Update existing organization body | 

### Return type

[**CommonsOrganization**](CommonsOrganization.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

