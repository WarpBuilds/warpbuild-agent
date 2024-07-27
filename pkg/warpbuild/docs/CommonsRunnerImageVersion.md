# CommonsRunnerImageVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Arch** | Pointer to **string** |  | [optional] 
**ContainerRunnerImage** | Pointer to [**CommonsContainerRunnerImageVersion**](CommonsContainerRunnerImageVersion.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** | ExternalID is the ID of the runner image version in the external system. | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Os** | Pointer to **string** |  | [optional] 
**RunnerImageId** | Pointer to **string** |  | [optional] 
**RunnerImagePullSecretId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**VersionTimeId** | Pointer to **int32** |  | [optional] 

## Methods

### NewCommonsRunnerImageVersion

`func NewCommonsRunnerImageVersion() *CommonsRunnerImageVersion`

NewCommonsRunnerImageVersion instantiates a new CommonsRunnerImageVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImageVersionWithDefaults

`func NewCommonsRunnerImageVersionWithDefaults() *CommonsRunnerImageVersion`

NewCommonsRunnerImageVersionWithDefaults instantiates a new CommonsRunnerImageVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArch

`func (o *CommonsRunnerImageVersion) GetArch() string`

GetArch returns the Arch field if non-nil, zero value otherwise.

### GetArchOk

`func (o *CommonsRunnerImageVersion) GetArchOk() (*string, bool)`

GetArchOk returns a tuple with the Arch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArch

`func (o *CommonsRunnerImageVersion) SetArch(v string)`

SetArch sets Arch field to given value.

### HasArch

`func (o *CommonsRunnerImageVersion) HasArch() bool`

HasArch returns a boolean if a field has been set.

### GetContainerRunnerImage

`func (o *CommonsRunnerImageVersion) GetContainerRunnerImage() CommonsContainerRunnerImageVersion`

GetContainerRunnerImage returns the ContainerRunnerImage field if non-nil, zero value otherwise.

### GetContainerRunnerImageOk

`func (o *CommonsRunnerImageVersion) GetContainerRunnerImageOk() (*CommonsContainerRunnerImageVersion, bool)`

GetContainerRunnerImageOk returns a tuple with the ContainerRunnerImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerRunnerImage

`func (o *CommonsRunnerImageVersion) SetContainerRunnerImage(v CommonsContainerRunnerImageVersion)`

SetContainerRunnerImage sets ContainerRunnerImage field to given value.

### HasContainerRunnerImage

`func (o *CommonsRunnerImageVersion) HasContainerRunnerImage() bool`

HasContainerRunnerImage returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsRunnerImageVersion) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsRunnerImageVersion) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsRunnerImageVersion) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CommonsRunnerImageVersion) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetExternalId

`func (o *CommonsRunnerImageVersion) GetExternalId() string`

GetExternalId returns the ExternalId field if non-nil, zero value otherwise.

### GetExternalIdOk

`func (o *CommonsRunnerImageVersion) GetExternalIdOk() (*string, bool)`

GetExternalIdOk returns a tuple with the ExternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalId

`func (o *CommonsRunnerImageVersion) SetExternalId(v string)`

SetExternalId sets ExternalId field to given value.

### HasExternalId

`func (o *CommonsRunnerImageVersion) HasExternalId() bool`

HasExternalId returns a boolean if a field has been set.

### GetId

`func (o *CommonsRunnerImageVersion) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsRunnerImageVersion) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsRunnerImageVersion) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsRunnerImageVersion) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsRunnerImageVersion) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsRunnerImageVersion) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsRunnerImageVersion) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsRunnerImageVersion) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetOs

`func (o *CommonsRunnerImageVersion) GetOs() string`

GetOs returns the Os field if non-nil, zero value otherwise.

### GetOsOk

`func (o *CommonsRunnerImageVersion) GetOsOk() (*string, bool)`

GetOsOk returns a tuple with the Os field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOs

`func (o *CommonsRunnerImageVersion) SetOs(v string)`

SetOs sets Os field to given value.

### HasOs

`func (o *CommonsRunnerImageVersion) HasOs() bool`

HasOs returns a boolean if a field has been set.

### GetRunnerImageId

`func (o *CommonsRunnerImageVersion) GetRunnerImageId() string`

GetRunnerImageId returns the RunnerImageId field if non-nil, zero value otherwise.

### GetRunnerImageIdOk

`func (o *CommonsRunnerImageVersion) GetRunnerImageIdOk() (*string, bool)`

GetRunnerImageIdOk returns a tuple with the RunnerImageId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerImageId

`func (o *CommonsRunnerImageVersion) SetRunnerImageId(v string)`

SetRunnerImageId sets RunnerImageId field to given value.

### HasRunnerImageId

`func (o *CommonsRunnerImageVersion) HasRunnerImageId() bool`

HasRunnerImageId returns a boolean if a field has been set.

### GetRunnerImagePullSecretId

`func (o *CommonsRunnerImageVersion) GetRunnerImagePullSecretId() string`

GetRunnerImagePullSecretId returns the RunnerImagePullSecretId field if non-nil, zero value otherwise.

### GetRunnerImagePullSecretIdOk

`func (o *CommonsRunnerImageVersion) GetRunnerImagePullSecretIdOk() (*string, bool)`

GetRunnerImagePullSecretIdOk returns a tuple with the RunnerImagePullSecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerImagePullSecretId

`func (o *CommonsRunnerImageVersion) SetRunnerImagePullSecretId(v string)`

SetRunnerImagePullSecretId sets RunnerImagePullSecretId field to given value.

### HasRunnerImagePullSecretId

`func (o *CommonsRunnerImageVersion) HasRunnerImagePullSecretId() bool`

HasRunnerImagePullSecretId returns a boolean if a field has been set.

### GetStatus

`func (o *CommonsRunnerImageVersion) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CommonsRunnerImageVersion) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CommonsRunnerImageVersion) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CommonsRunnerImageVersion) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetType

`func (o *CommonsRunnerImageVersion) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsRunnerImageVersion) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsRunnerImageVersion) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *CommonsRunnerImageVersion) HasType() bool`

HasType returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *CommonsRunnerImageVersion) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsRunnerImageVersion) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsRunnerImageVersion) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *CommonsRunnerImageVersion) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetVersionTimeId

`func (o *CommonsRunnerImageVersion) GetVersionTimeId() int32`

GetVersionTimeId returns the VersionTimeId field if non-nil, zero value otherwise.

### GetVersionTimeIdOk

`func (o *CommonsRunnerImageVersion) GetVersionTimeIdOk() (*int32, bool)`

GetVersionTimeIdOk returns a tuple with the VersionTimeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersionTimeId

`func (o *CommonsRunnerImageVersion) SetVersionTimeId(v int32)`

SetVersionTimeId sets VersionTimeId field to given value.

### HasVersionTimeId

`func (o *CommonsRunnerImageVersion) HasVersionTimeId() bool`

HasVersionTimeId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


