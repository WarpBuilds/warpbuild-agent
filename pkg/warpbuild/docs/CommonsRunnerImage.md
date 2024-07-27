# CommonsRunnerImage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alias** | Pointer to **string** |  | [optional] 
**Arch** | Pointer to **string** |  | [optional] 
**ContainerRunnerImage** | Pointer to [**CommonsContainerRunnerImage**](CommonsContainerRunnerImage.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Hooks** | Pointer to [**[]CommonsRunnerImageHook**](CommonsRunnerImageHook.md) |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Os** | Pointer to **string** |  | [optional] 
**RunnerImagePullSecretId** | Pointer to **string** |  | [optional] 
**Settings** | Pointer to [**CommonsRunnerImageSettings**](CommonsRunnerImageSettings.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**WarpbuildImage** | Pointer to [**CommonsWarpbuildImage**](CommonsWarpbuildImage.md) |  | [optional] 

## Methods

### NewCommonsRunnerImage

`func NewCommonsRunnerImage() *CommonsRunnerImage`

NewCommonsRunnerImage instantiates a new CommonsRunnerImage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImageWithDefaults

`func NewCommonsRunnerImageWithDefaults() *CommonsRunnerImage`

NewCommonsRunnerImageWithDefaults instantiates a new CommonsRunnerImage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlias

`func (o *CommonsRunnerImage) GetAlias() string`

GetAlias returns the Alias field if non-nil, zero value otherwise.

### GetAliasOk

`func (o *CommonsRunnerImage) GetAliasOk() (*string, bool)`

GetAliasOk returns a tuple with the Alias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlias

`func (o *CommonsRunnerImage) SetAlias(v string)`

SetAlias sets Alias field to given value.

### HasAlias

`func (o *CommonsRunnerImage) HasAlias() bool`

HasAlias returns a boolean if a field has been set.

### GetArch

`func (o *CommonsRunnerImage) GetArch() string`

GetArch returns the Arch field if non-nil, zero value otherwise.

### GetArchOk

`func (o *CommonsRunnerImage) GetArchOk() (*string, bool)`

GetArchOk returns a tuple with the Arch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArch

`func (o *CommonsRunnerImage) SetArch(v string)`

SetArch sets Arch field to given value.

### HasArch

`func (o *CommonsRunnerImage) HasArch() bool`

HasArch returns a boolean if a field has been set.

### GetContainerRunnerImage

`func (o *CommonsRunnerImage) GetContainerRunnerImage() CommonsContainerRunnerImage`

GetContainerRunnerImage returns the ContainerRunnerImage field if non-nil, zero value otherwise.

### GetContainerRunnerImageOk

`func (o *CommonsRunnerImage) GetContainerRunnerImageOk() (*CommonsContainerRunnerImage, bool)`

GetContainerRunnerImageOk returns a tuple with the ContainerRunnerImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerRunnerImage

`func (o *CommonsRunnerImage) SetContainerRunnerImage(v CommonsContainerRunnerImage)`

SetContainerRunnerImage sets ContainerRunnerImage field to given value.

### HasContainerRunnerImage

`func (o *CommonsRunnerImage) HasContainerRunnerImage() bool`

HasContainerRunnerImage returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsRunnerImage) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsRunnerImage) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsRunnerImage) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CommonsRunnerImage) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetHooks

`func (o *CommonsRunnerImage) GetHooks() []CommonsRunnerImageHook`

GetHooks returns the Hooks field if non-nil, zero value otherwise.

### GetHooksOk

`func (o *CommonsRunnerImage) GetHooksOk() (*[]CommonsRunnerImageHook, bool)`

GetHooksOk returns a tuple with the Hooks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHooks

`func (o *CommonsRunnerImage) SetHooks(v []CommonsRunnerImageHook)`

SetHooks sets Hooks field to given value.

### HasHooks

`func (o *CommonsRunnerImage) HasHooks() bool`

HasHooks returns a boolean if a field has been set.

### GetId

`func (o *CommonsRunnerImage) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsRunnerImage) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsRunnerImage) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsRunnerImage) HasId() bool`

HasId returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsRunnerImage) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsRunnerImage) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsRunnerImage) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsRunnerImage) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetOs

`func (o *CommonsRunnerImage) GetOs() string`

GetOs returns the Os field if non-nil, zero value otherwise.

### GetOsOk

`func (o *CommonsRunnerImage) GetOsOk() (*string, bool)`

GetOsOk returns a tuple with the Os field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOs

`func (o *CommonsRunnerImage) SetOs(v string)`

SetOs sets Os field to given value.

### HasOs

`func (o *CommonsRunnerImage) HasOs() bool`

HasOs returns a boolean if a field has been set.

### GetRunnerImagePullSecretId

`func (o *CommonsRunnerImage) GetRunnerImagePullSecretId() string`

GetRunnerImagePullSecretId returns the RunnerImagePullSecretId field if non-nil, zero value otherwise.

### GetRunnerImagePullSecretIdOk

`func (o *CommonsRunnerImage) GetRunnerImagePullSecretIdOk() (*string, bool)`

GetRunnerImagePullSecretIdOk returns a tuple with the RunnerImagePullSecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerImagePullSecretId

`func (o *CommonsRunnerImage) SetRunnerImagePullSecretId(v string)`

SetRunnerImagePullSecretId sets RunnerImagePullSecretId field to given value.

### HasRunnerImagePullSecretId

`func (o *CommonsRunnerImage) HasRunnerImagePullSecretId() bool`

HasRunnerImagePullSecretId returns a boolean if a field has been set.

### GetSettings

`func (o *CommonsRunnerImage) GetSettings() CommonsRunnerImageSettings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *CommonsRunnerImage) GetSettingsOk() (*CommonsRunnerImageSettings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *CommonsRunnerImage) SetSettings(v CommonsRunnerImageSettings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *CommonsRunnerImage) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

### GetStatus

`func (o *CommonsRunnerImage) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CommonsRunnerImage) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CommonsRunnerImage) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CommonsRunnerImage) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetType

`func (o *CommonsRunnerImage) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsRunnerImage) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsRunnerImage) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *CommonsRunnerImage) HasType() bool`

HasType returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *CommonsRunnerImage) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsRunnerImage) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsRunnerImage) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *CommonsRunnerImage) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetWarpbuildImage

`func (o *CommonsRunnerImage) GetWarpbuildImage() CommonsWarpbuildImage`

GetWarpbuildImage returns the WarpbuildImage field if non-nil, zero value otherwise.

### GetWarpbuildImageOk

`func (o *CommonsRunnerImage) GetWarpbuildImageOk() (*CommonsWarpbuildImage, bool)`

GetWarpbuildImageOk returns a tuple with the WarpbuildImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarpbuildImage

`func (o *CommonsRunnerImage) SetWarpbuildImage(v CommonsWarpbuildImage)`

SetWarpbuildImage sets WarpbuildImage field to given value.

### HasWarpbuildImage

`func (o *CommonsRunnerImage) HasWarpbuildImage() bool`

HasWarpbuildImage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


