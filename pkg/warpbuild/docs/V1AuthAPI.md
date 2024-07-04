<<<<<<< HEAD
# \V1AuthAPI
=======
# \V1AuthApi
>>>>>>> prajjwal-warp-323

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
<<<<<<< HEAD
[**AuthTokensGet**](V1AuthAPI.md#AuthTokensGet) | **Get** /auth/tokens | List user tokens
[**AuthUser**](V1AuthAPI.md#AuthUser) | **Post** /auth | Auth user
[**AuthUsersGet**](V1AuthAPI.md#AuthUsersGet) | **Get** /auth/users | List users
[**GetAuthURL**](V1AuthAPI.md#GetAuthURL) | **Get** /auth/login/{provider} | Get auth url
[**GetMe**](V1AuthAPI.md#GetMe) | **Get** /auth/me | Auth user
[**Logout**](V1AuthAPI.md#Logout) | **Patch** /auth/logout | Logout
[**RefreshToken**](V1AuthAPI.md#RefreshToken) | **Patch** /auth/token/refresh | Refresh token
[**SwitchOrganization**](V1AuthAPI.md#SwitchOrganization) | **Patch** /auth/switch | Switch organization
=======
[**AuthTokensGet**](V1AuthApi.md#AuthTokensGet) | **Get** /auth/tokens | List user tokens
[**AuthUser**](V1AuthApi.md#AuthUser) | **Post** /auth | Auth user
[**AuthUsersGet**](V1AuthApi.md#AuthUsersGet) | **Get** /auth/users | List users
[**GetAuthURL**](V1AuthApi.md#GetAuthURL) | **Get** /auth/login/{provider} | Get auth url
[**GetMe**](V1AuthApi.md#GetMe) | **Get** /auth/me | Auth user
[**Logout**](V1AuthApi.md#Logout) | **Patch** /auth/logout | Logout
[**RefreshToken**](V1AuthApi.md#RefreshToken) | **Patch** /auth/token/refresh | Refresh token
[**SwitchOrganization**](V1AuthApi.md#SwitchOrganization) | **Patch** /auth/switch | Switch organization
>>>>>>> prajjwal-warp-323



## AuthTokensGet

> []CommonsUserToken AuthTokensGet(ctx).Body(body).Execute()

List user tokens



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
    body := *openapiclient.NewCommonsListTokensOptions() // CommonsListTokensOptions | ListTokenOptions

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1AuthAPI.AuthTokensGet(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.AuthTokensGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthTokensGet`: []CommonsUserToken
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.AuthTokensGet`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.AuthTokensGet(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.AuthTokensGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthTokensGet`: []CommonsUserToken
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.AuthTokensGet`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthTokensGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsListTokensOptions**](CommonsListTokensOptions.md) | ListTokenOptions | 

### Return type

[**[]CommonsUserToken**](CommonsUserToken.md)

### Authorization

[WarpBuildAdminSecretAuth](../README.md#WarpBuildAdminSecretAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthUser

> AuthUserResponse AuthUser(ctx).Body(body).Execute()

Auth user



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
    body := *openapiclient.NewAuthUserRequest("Code_example", "State_example") // AuthUserRequest | Auth request body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1AuthAPI.AuthUser(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.AuthUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthUser`: AuthUserResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.AuthUser`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.AuthUser(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.AuthUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthUser`: AuthUserResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.AuthUser`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**AuthUserRequest**](AuthUserRequest.md) | Auth request body | 

### Return type

[**AuthUserResponse**](AuthUserResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AuthUsersGet

> []CommonsListUsersResponse AuthUsersGet(ctx).Body(body).Execute()

List users



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
    body := *openapiclient.NewCommonsListUsersOptions() // CommonsListUsersOptions | ListUsersOptions

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1AuthAPI.AuthUsersGet(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.AuthUsersGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthUsersGet`: []CommonsListUsersResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.AuthUsersGet`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.AuthUsersGet(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.AuthUsersGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AuthUsersGet`: []CommonsListUsersResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.AuthUsersGet`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAuthUsersGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsListUsersOptions**](CommonsListUsersOptions.md) | ListUsersOptions | 

### Return type

[**[]CommonsListUsersResponse**](CommonsListUsersResponse.md)

### Authorization

[WarpBuildAdminSecretAuth](../README.md#WarpBuildAdminSecretAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAuthURL

<<<<<<< HEAD
> GetAuthURL(ctx, provider).InviteCode(inviteCode).Execute()
=======
> GetAuthURL(ctx, provider).Execute()
>>>>>>> prajjwal-warp-323

Get auth url



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
    provider := "provider_example" // string | Provider
<<<<<<< HEAD
    inviteCode := "inviteCode_example" // string | Invite code if any (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.V1AuthAPI.GetAuthURL(context.Background(), provider).InviteCode(inviteCode).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.GetAuthURL``: %v\n", err)
=======

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.V1AuthApi.GetAuthURL(context.Background(), provider).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.GetAuthURL``: %v\n", err)
>>>>>>> prajjwal-warp-323
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**provider** | **string** | Provider | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetAuthURLRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

<<<<<<< HEAD
 **inviteCode** | **string** | Invite code if any | 
=======
>>>>>>> prajjwal-warp-323

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetMe

> MeResponse GetMe(ctx).Execute()

Auth user



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
    resp, r, err := apiClient.V1AuthAPI.GetMe(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.GetMe``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetMe`: MeResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.GetMe`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.GetMe(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.GetMe``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetMe`: MeResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.GetMe`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetMeRequest struct via the builder pattern


### Return type

[**MeResponse**](MeResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## Logout

> Logout(ctx).Execute()

Logout



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
    r, err := apiClient.V1AuthAPI.Logout(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.Logout``: %v\n", err)
=======
    r, err := apiClient.V1AuthApi.Logout(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.Logout``: %v\n", err)
>>>>>>> prajjwal-warp-323
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiLogoutRequest struct via the builder pattern


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


## RefreshToken

> TokenRefreshResponse RefreshToken(ctx).Body(body).Execute()

Refresh token



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
    body := *openapiclient.NewTokenRefreshRequest("RefreshToken_example") // TokenRefreshRequest | Refresh token

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1AuthAPI.RefreshToken(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.RefreshToken``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RefreshToken`: TokenRefreshResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.RefreshToken`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.RefreshToken(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.RefreshToken``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RefreshToken`: TokenRefreshResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.RefreshToken`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRefreshTokenRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**TokenRefreshRequest**](TokenRefreshRequest.md) | Refresh token | 

### Return type

[**TokenRefreshResponse**](TokenRefreshResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SwitchOrganization

> SwitchOrganizationResponse SwitchOrganization(ctx).Body(body).Execute()

Switch organization

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
    body := *openapiclient.NewSwitchOrganizationRequest("OrganizationId_example", "RefreshToken_example") // SwitchOrganizationRequest | Switch organization

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
<<<<<<< HEAD
    resp, r, err := apiClient.V1AuthAPI.SwitchOrganization(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthAPI.SwitchOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SwitchOrganization`: SwitchOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthAPI.SwitchOrganization`: %v\n", resp)
=======
    resp, r, err := apiClient.V1AuthApi.SwitchOrganization(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1AuthApi.SwitchOrganization``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SwitchOrganization`: SwitchOrganizationResponse
    fmt.Fprintf(os.Stdout, "Response from `V1AuthApi.SwitchOrganization`: %v\n", resp)
>>>>>>> prajjwal-warp-323
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiSwitchOrganizationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**SwitchOrganizationRequest**](SwitchOrganizationRequest.md) | Switch organization | 

### Return type

[**SwitchOrganizationResponse**](SwitchOrganizationResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

