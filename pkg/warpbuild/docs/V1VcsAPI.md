# \V1VcsApi

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApproveVCSIntegration**](V1VcsApi.md#ApproveVCSIntegration) | **Put** /vcs/approve-integration | This handles the callback for approving an installation
[**CreateVCSGitRepo**](V1VcsApi.md#CreateVCSGitRepo) | **Post** /vcs/repos | create vcs repo based on repo internal id
[**CreateVCSIntegration**](V1VcsApi.md#CreateVCSIntegration) | **Post** /vcs/integrations | Create a new vcs integration
[**DeleteVCSIntegration**](V1VcsApi.md#DeleteVCSIntegration) | **Delete** /vcs/integrations/{integration_id} | Delete an existing vcs integration
[**GetVCSGitRepo**](V1VcsApi.md#GetVCSGitRepo) | **Get** /vcs/repos/{id} | get vcs repo based on repo internal id
[**ListVCSEntites**](V1VcsApi.md#ListVCSEntites) | **Get** /vcs/entities | Lists all vcs entities for vcs integration
[**ListVCSIntegration**](V1VcsApi.md#ListVCSIntegration) | **Get** /vcs/integrations | Lists all vcs integration for provider
[**ListVCSRepos**](V1VcsApi.md#ListVCSRepos) | **Get** /vcs/repos | Lists all vcs repos for vcs integration
[**ListVCSRunnerGroups**](V1VcsApi.md#ListVCSRunnerGroups) | **Post** /vcs/list-runner-groups | Lists all vcs runner groups
[**UpdateVCSIntegration**](V1VcsApi.md#UpdateVCSIntegration) | **Put** /vcs/integrations/{integration_id} | Update an existing vcs integration



## ApproveVCSIntegration

> AuthUserResponse ApproveVCSIntegration(ctx).Body(body).Execute()

This handles the callback for approving an installation

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
    body := *openapiclient.NewApproveVCSIntegrationRequest() // ApproveVCSIntegrationRequest | Approve vcs integration app installation

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.ApproveVCSIntegration(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.ApproveVCSIntegration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApproveVCSIntegration`: AuthUserResponse
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.ApproveVCSIntegration`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApproveVCSIntegrationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ApproveVCSIntegrationRequest**](ApproveVCSIntegrationRequest.md) | Approve vcs integration app installation | 

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


## CreateVCSGitRepo

> CommonsRepo CreateVCSGitRepo(ctx).Body(body).Execute()

create vcs repo based on repo internal id

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
    body := *openapiclient.NewCommonsCreateRepoOptions("IntegrationId_example", "Name_example", "Owner_example") // CommonsCreateRepoOptions | create repo options

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.CreateVCSGitRepo(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.CreateVCSGitRepo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateVCSGitRepo`: CommonsRepo
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.CreateVCSGitRepo`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateVCSGitRepoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsCreateRepoOptions**](CommonsCreateRepoOptions.md) | create repo options | 

### Return type

[**CommonsRepo**](CommonsRepo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## CreateVCSIntegration

> VCSIntegration CreateVCSIntegration(ctx).Body(body).Execute()

Create a new vcs integration

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
    body := *openapiclient.NewCreateVCSIntegrationRequest("Provider_example") // CreateVCSIntegrationRequest | Create new vcs integration body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.CreateVCSIntegration(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.CreateVCSIntegration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `CreateVCSIntegration`: VCSIntegration
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.CreateVCSIntegration`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiCreateVCSIntegrationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CreateVCSIntegrationRequest**](CreateVCSIntegrationRequest.md) | Create new vcs integration body | 

### Return type

[**VCSIntegration**](VCSIntegration.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteVCSIntegration

> TypesGenericSuccessMessage DeleteVCSIntegration(ctx, integrationId).Provider(provider).Execute()

Delete an existing vcs integration

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
    integrationId := "integrationId_example" // string | ID for the vcs integration
    provider := "provider_example" // string | ID for the vcs integration

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.DeleteVCSIntegration(context.Background(), integrationId).Provider(provider).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.DeleteVCSIntegration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteVCSIntegration`: TypesGenericSuccessMessage
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.DeleteVCSIntegration`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**integrationId** | **string** | ID for the vcs integration | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteVCSIntegrationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **provider** | **string** | ID for the vcs integration | 

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


## GetVCSGitRepo

> CommonsRepo GetVCSGitRepo(ctx, id).Execute()

get vcs repo based on repo internal id

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
    id := "id_example" // string | internal id for the vcs git repo

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.GetVCSGitRepo(context.Background(), id).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.GetVCSGitRepo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetVCSGitRepo`: CommonsRepo
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.GetVCSGitRepo`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | internal id for the vcs git repo | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetVCSGitRepoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsRepo**](CommonsRepo.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVCSEntites

> []VCSEntity ListVCSEntites(ctx).EntityType(entityType).Provider(provider).Name(name).IntegrationId(integrationId).ParentId(parentId).Execute()

Lists all vcs entities for vcs integration



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
    entityType := "entityType_example" // string | VCS Entity that should be returned
    provider := "provider_example" // string | Git provider (optional)
    name := "name_example" // string | Filter using organization name (optional)
    integrationId := "integrationId_example" // string | IntegrationID used by the git provider (optional)
    parentId := "parentId_example" // string | VCS Entity Parent ID (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.ListVCSEntites(context.Background()).EntityType(entityType).Provider(provider).Name(name).IntegrationId(integrationId).ParentId(parentId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.ListVCSEntites``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListVCSEntites`: []VCSEntity
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.ListVCSEntites`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVCSEntitesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **entityType** | **string** | VCS Entity that should be returned | 
 **provider** | **string** | Git provider | 
 **name** | **string** | Filter using organization name | 
 **integrationId** | **string** | IntegrationID used by the git provider | 
 **parentId** | **string** | VCS Entity Parent ID | 

### Return type

[**[]VCSEntity**](VCSEntity.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVCSIntegration

> []VCSIntegration ListVCSIntegration(ctx).Provider(provider).Status(status).Execute()

Lists all vcs integration for provider

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
    provider := "provider_example" // string | vcs integration provider filter (optional)
    status := "status_example" // string | vcs integration provider filter (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.ListVCSIntegration(context.Background()).Provider(provider).Status(status).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.ListVCSIntegration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListVCSIntegration`: []VCSIntegration
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.ListVCSIntegration`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVCSIntegrationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **provider** | **string** | vcs integration provider filter | 
 **status** | **string** | vcs integration provider filter | 

### Return type

[**[]VCSIntegration**](VCSIntegration.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVCSRepos

> []CommonsRepo ListVCSRepos(ctx).Execute()

Lists all vcs repos for vcs integration

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
    resp, r, err := apiClient.V1VcsApi.ListVCSRepos(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.ListVCSRepos``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListVCSRepos`: []CommonsRepo
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.ListVCSRepos`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListVCSReposRequest struct via the builder pattern


### Return type

[**[]CommonsRepo**](CommonsRepo.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListVCSRunnerGroups

> CommonsListVCSRunnerGroupsResponse ListVCSRunnerGroups(ctx).Body(body).Execute()

Lists all vcs runner groups

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
    body := *openapiclient.NewCommonsListVCSRunnerGroupsInput() // CommonsListVCSRunnerGroupsInput | List runner groups input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.ListVCSRunnerGroups(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.ListVCSRunnerGroups``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListVCSRunnerGroups`: CommonsListVCSRunnerGroupsResponse
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.ListVCSRunnerGroups`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiListVCSRunnerGroupsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsListVCSRunnerGroupsInput**](CommonsListVCSRunnerGroupsInput.md) | List runner groups input | 

### Return type

[**CommonsListVCSRunnerGroupsResponse**](CommonsListVCSRunnerGroupsResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateVCSIntegration

> UpdateVCSIntegrationResponse UpdateVCSIntegration(ctx, integrationId).Body(body).Execute()

Update an existing vcs integration

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
    integrationId := "integrationId_example" // string | ID for the vcs integration
    body := *openapiclient.NewUpdateVCSIntegrationRequest("Id_example") // UpdateVCSIntegrationRequest | Update vcs integration body

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1VcsApi.UpdateVCSIntegration(context.Background(), integrationId).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1VcsApi.UpdateVCSIntegration``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateVCSIntegration`: UpdateVCSIntegrationResponse
    fmt.Fprintf(os.Stdout, "Response from `V1VcsApi.UpdateVCSIntegration`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**integrationId** | **string** | ID for the vcs integration | 

### Other Parameters

Other parameters are passed through a pointer to a apiUpdateVCSIntegrationRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **body** | [**UpdateVCSIntegrationRequest**](UpdateVCSIntegrationRequest.md) | Update vcs integration body | 

### Return type

[**UpdateVCSIntegrationResponse**](UpdateVCSIntegrationResponse.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

