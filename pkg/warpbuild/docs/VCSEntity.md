# VCSEntity

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | Pointer to **string** | AccountID is the github account  TODO: AccountID should be group id in case of gitlab | [optional] 
**CommitId** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**DisplayName** | **string** |  | 
**EntityType** | **string** |  | 
**ExternalId** | **string** |  | 
**IntegrationId** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 
**Name** | **string** |  | 
**Provider** | **string** |  | 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**WebUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewVCSEntity

`func NewVCSEntity(displayName string, entityType string, externalId string, name string, provider string, ) *VCSEntity`

NewVCSEntity instantiates a new VCSEntity object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVCSEntityWithDefaults

`func NewVCSEntityWithDefaults() *VCSEntity`

NewVCSEntityWithDefaults instantiates a new VCSEntity object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *VCSEntity) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *VCSEntity) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *VCSEntity) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *VCSEntity) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### GetCommitId

`func (o *VCSEntity) GetCommitId() string`

GetCommitId returns the CommitId field if non-nil, zero value otherwise.

### GetCommitIdOk

`func (o *VCSEntity) GetCommitIdOk() (*string, bool)`

GetCommitIdOk returns a tuple with the CommitId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitId

`func (o *VCSEntity) SetCommitId(v string)`

SetCommitId sets CommitId field to given value.

### HasCommitId

`func (o *VCSEntity) HasCommitId() bool`

HasCommitId returns a boolean if a field has been set.

### GetCreatedAt

`func (o *VCSEntity) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *VCSEntity) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *VCSEntity) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *VCSEntity) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetDisplayName

`func (o *VCSEntity) GetDisplayName() string`

GetDisplayName returns the DisplayName field if non-nil, zero value otherwise.

### GetDisplayNameOk

`func (o *VCSEntity) GetDisplayNameOk() (*string, bool)`

GetDisplayNameOk returns a tuple with the DisplayName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisplayName

`func (o *VCSEntity) SetDisplayName(v string)`

SetDisplayName sets DisplayName field to given value.


### GetEntityType

`func (o *VCSEntity) GetEntityType() string`

GetEntityType returns the EntityType field if non-nil, zero value otherwise.

### GetEntityTypeOk

`func (o *VCSEntity) GetEntityTypeOk() (*string, bool)`

GetEntityTypeOk returns a tuple with the EntityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEntityType

`func (o *VCSEntity) SetEntityType(v string)`

SetEntityType sets EntityType field to given value.


### GetExternalId

`func (o *VCSEntity) GetExternalId() string`

GetExternalId returns the ExternalId field if non-nil, zero value otherwise.

### GetExternalIdOk

`func (o *VCSEntity) GetExternalIdOk() (*string, bool)`

GetExternalIdOk returns a tuple with the ExternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalId

`func (o *VCSEntity) SetExternalId(v string)`

SetExternalId sets ExternalId field to given value.


### GetIntegrationId

`func (o *VCSEntity) GetIntegrationId() string`

GetIntegrationId returns the IntegrationId field if non-nil, zero value otherwise.

### GetIntegrationIdOk

`func (o *VCSEntity) GetIntegrationIdOk() (*string, bool)`

GetIntegrationIdOk returns a tuple with the IntegrationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrationId

`func (o *VCSEntity) SetIntegrationId(v string)`

SetIntegrationId sets IntegrationId field to given value.

### HasIntegrationId

`func (o *VCSEntity) HasIntegrationId() bool`

HasIntegrationId returns a boolean if a field has been set.

### GetInternalId

`func (o *VCSEntity) GetInternalId() string`

GetInternalId returns the InternalId field if non-nil, zero value otherwise.

### GetInternalIdOk

`func (o *VCSEntity) GetInternalIdOk() (*string, bool)`

GetInternalIdOk returns a tuple with the InternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInternalId

`func (o *VCSEntity) SetInternalId(v string)`

SetInternalId sets InternalId field to given value.

### HasInternalId

`func (o *VCSEntity) HasInternalId() bool`

HasInternalId returns a boolean if a field has been set.

### GetName

`func (o *VCSEntity) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VCSEntity) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VCSEntity) SetName(v string)`

SetName sets Name field to given value.


### GetProvider

`func (o *VCSEntity) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *VCSEntity) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *VCSEntity) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetUpdatedAt

`func (o *VCSEntity) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *VCSEntity) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *VCSEntity) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *VCSEntity) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetWebUrl

`func (o *VCSEntity) GetWebUrl() string`

GetWebUrl returns the WebUrl field if non-nil, zero value otherwise.

### GetWebUrlOk

`func (o *VCSEntity) GetWebUrlOk() (*string, bool)`

GetWebUrlOk returns a tuple with the WebUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWebUrl

`func (o *VCSEntity) SetWebUrl(v string)`

SetWebUrl sets WebUrl field to given value.

### HasWebUrl

`func (o *VCSEntity) HasWebUrl() bool`

HasWebUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


