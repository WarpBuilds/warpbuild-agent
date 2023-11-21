# CreateVCSIntegrationRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AdditionalState** | Pointer to **map[string]interface{}** |  | [optional] 
**ApplicationCallbackUrl** | Pointer to **string** |  | [optional] 
**ApplicationId** | Pointer to **string** |  | [optional] 
**ApplicationSecret** | Pointer to **string** |  | [optional] 
**Provider** | **string** |  | 
**VcsUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewCreateVCSIntegrationRequest

`func NewCreateVCSIntegrationRequest(provider string, ) *CreateVCSIntegrationRequest`

NewCreateVCSIntegrationRequest instantiates a new CreateVCSIntegrationRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCreateVCSIntegrationRequestWithDefaults

`func NewCreateVCSIntegrationRequestWithDefaults() *CreateVCSIntegrationRequest`

NewCreateVCSIntegrationRequestWithDefaults instantiates a new CreateVCSIntegrationRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAdditionalState

`func (o *CreateVCSIntegrationRequest) GetAdditionalState() map[string]interface{}`

GetAdditionalState returns the AdditionalState field if non-nil, zero value otherwise.

### GetAdditionalStateOk

`func (o *CreateVCSIntegrationRequest) GetAdditionalStateOk() (*map[string]interface{}, bool)`

GetAdditionalStateOk returns a tuple with the AdditionalState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdditionalState

`func (o *CreateVCSIntegrationRequest) SetAdditionalState(v map[string]interface{})`

SetAdditionalState sets AdditionalState field to given value.

### HasAdditionalState

`func (o *CreateVCSIntegrationRequest) HasAdditionalState() bool`

HasAdditionalState returns a boolean if a field has been set.

### GetApplicationCallbackUrl

`func (o *CreateVCSIntegrationRequest) GetApplicationCallbackUrl() string`

GetApplicationCallbackUrl returns the ApplicationCallbackUrl field if non-nil, zero value otherwise.

### GetApplicationCallbackUrlOk

`func (o *CreateVCSIntegrationRequest) GetApplicationCallbackUrlOk() (*string, bool)`

GetApplicationCallbackUrlOk returns a tuple with the ApplicationCallbackUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplicationCallbackUrl

`func (o *CreateVCSIntegrationRequest) SetApplicationCallbackUrl(v string)`

SetApplicationCallbackUrl sets ApplicationCallbackUrl field to given value.

### HasApplicationCallbackUrl

`func (o *CreateVCSIntegrationRequest) HasApplicationCallbackUrl() bool`

HasApplicationCallbackUrl returns a boolean if a field has been set.

### GetApplicationId

`func (o *CreateVCSIntegrationRequest) GetApplicationId() string`

GetApplicationId returns the ApplicationId field if non-nil, zero value otherwise.

### GetApplicationIdOk

`func (o *CreateVCSIntegrationRequest) GetApplicationIdOk() (*string, bool)`

GetApplicationIdOk returns a tuple with the ApplicationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplicationId

`func (o *CreateVCSIntegrationRequest) SetApplicationId(v string)`

SetApplicationId sets ApplicationId field to given value.

### HasApplicationId

`func (o *CreateVCSIntegrationRequest) HasApplicationId() bool`

HasApplicationId returns a boolean if a field has been set.

### GetApplicationSecret

`func (o *CreateVCSIntegrationRequest) GetApplicationSecret() string`

GetApplicationSecret returns the ApplicationSecret field if non-nil, zero value otherwise.

### GetApplicationSecretOk

`func (o *CreateVCSIntegrationRequest) GetApplicationSecretOk() (*string, bool)`

GetApplicationSecretOk returns a tuple with the ApplicationSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplicationSecret

`func (o *CreateVCSIntegrationRequest) SetApplicationSecret(v string)`

SetApplicationSecret sets ApplicationSecret field to given value.

### HasApplicationSecret

`func (o *CreateVCSIntegrationRequest) HasApplicationSecret() bool`

HasApplicationSecret returns a boolean if a field has been set.

### GetProvider

`func (o *CreateVCSIntegrationRequest) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *CreateVCSIntegrationRequest) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *CreateVCSIntegrationRequest) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetVcsUrl

`func (o *CreateVCSIntegrationRequest) GetVcsUrl() string`

GetVcsUrl returns the VcsUrl field if non-nil, zero value otherwise.

### GetVcsUrlOk

`func (o *CreateVCSIntegrationRequest) GetVcsUrlOk() (*string, bool)`

GetVcsUrlOk returns a tuple with the VcsUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsUrl

`func (o *CreateVCSIntegrationRequest) SetVcsUrl(v string)`

SetVcsUrl sets VcsUrl field to given value.

### HasVcsUrl

`func (o *CreateVCSIntegrationRequest) HasVcsUrl() bool`

HasVcsUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


