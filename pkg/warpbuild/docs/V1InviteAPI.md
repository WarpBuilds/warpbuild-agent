# \V1InviteAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateInvite**](V1InviteAPI.md#CreateInvite) | **Post** /invite | Adds a new organisation for a current user
[**ListOrgInvites**](V1InviteAPI.md#ListOrgInvites) | **Get** /invites | ListOrgInvites lists all the invite user has access to.
[**SendInvite**](V1InviteAPI.md#SendInvite) | **Post** /invite/send/email | Adds a new organisation for a current user
[**UpdateInvite**](V1InviteAPI.md#UpdateInvite) | **Patch** /invite | Updates existing invite based on the fields provided.



## CreateInvite

> V1Invite CreateInvite(ctx).Body(body).Execute()

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
    body := *openapiclient.NewCreateInviteRequest("RoleKind_example") // CreateInviteRequest | Create new invite body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1InviteAPI.CreateInvite(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1InviteAPI.CreateInvite``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateInvite`: V1Invite
    fmt.Fprintf(os.Stdout, "Response from `V1InviteAPI.CreateInvite`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateInviteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CreateInviteRequest**](CreateInviteRequest.md) | Create new invite body | 

### Return type

[**V1Invite**](V1Invite.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListOrgInvites

> []V1Invite ListOrgInvites(ctx).Execute()

ListOrgInvites lists all the invite user has access to.

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
    resp, r, err := apiClient.V1InviteAPI.ListOrgInvites(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1InviteAPI.ListOrgInvites``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListOrgInvites`: []V1Invite
    fmt.Fprintf(os.Stdout, "Response from `V1InviteAPI.ListOrgInvites`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListOrgInvitesRequest struct via the builder pattern


### Return type

[**[]V1Invite**](V1Invite.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SendInvite

> SendInvite(ctx).Body(body).Execute()

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
    body := *openapiclient.NewV1SendInviteRequest("Email_example", "InviteId_example") // V1SendInviteRequest | Send invite body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.V1InviteAPI.SendInvite(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1InviteAPI.SendInvite``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSendInviteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**V1SendInviteRequest**](V1SendInviteRequest.md) | Send invite body | 

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


## UpdateInvite

> V1Invite UpdateInvite(ctx).Body(body).Execute()

Updates existing invite based on the fields provided.



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
    body := *openapiclient.NewUpdateInviteRequest() // UpdateInviteRequest | Update existing invite body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1InviteAPI.UpdateInvite(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1InviteAPI.UpdateInvite``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateInvite`: V1Invite
    fmt.Fprintf(os.Stdout, "Response from `V1InviteAPI.UpdateInvite`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateInviteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**UpdateInviteRequest**](UpdateInviteRequest.md) | Update existing invite body | 

### Return type

[**V1Invite**](V1Invite.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

