# CommonsCreateRunnerImageInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alias** | Pointer to **string** |  | [optional] 
**Arch** | Pointer to **string** |  | [optional] 
**ContainerRunnerImage** | Pointer to [**CommonsContainerRunnerImage**](CommonsContainerRunnerImage.md) |  | [optional] 
**Hooks** | Pointer to [**[]CommonsRunnerImageHook**](CommonsRunnerImageHook.md) |  | [optional] 
**Os** | Pointer to **string** |  | [optional] 
**RunnerImagePullSecretId** | Pointer to **string** |  | [optional] 
**Settings** | Pointer to [**CommonsRunnerImageSettings**](CommonsRunnerImageSettings.md) |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsCreateRunnerImageInput

`func NewCommonsCreateRunnerImageInput() *CommonsCreateRunnerImageInput`

NewCommonsCreateRunnerImageInput instantiates a new CommonsCreateRunnerImageInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsCreateRunnerImageInputWithDefaults

`func NewCommonsCreateRunnerImageInputWithDefaults() *CommonsCreateRunnerImageInput`

NewCommonsCreateRunnerImageInputWithDefaults instantiates a new CommonsCreateRunnerImageInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlias

`func (o *CommonsCreateRunnerImageInput) GetAlias() string`

GetAlias returns the Alias field if non-nil, zero value otherwise.

### GetAliasOk

`func (o *CommonsCreateRunnerImageInput) GetAliasOk() (*string, bool)`

GetAliasOk returns a tuple with the Alias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlias

`func (o *CommonsCreateRunnerImageInput) SetAlias(v string)`

SetAlias sets Alias field to given value.

### HasAlias

`func (o *CommonsCreateRunnerImageInput) HasAlias() bool`

HasAlias returns a boolean if a field has been set.

### GetArch

`func (o *CommonsCreateRunnerImageInput) GetArch() string`

GetArch returns the Arch field if non-nil, zero value otherwise.

### GetArchOk

`func (o *CommonsCreateRunnerImageInput) GetArchOk() (*string, bool)`

GetArchOk returns a tuple with the Arch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArch

`func (o *CommonsCreateRunnerImageInput) SetArch(v string)`

SetArch sets Arch field to given value.

### HasArch

`func (o *CommonsCreateRunnerImageInput) HasArch() bool`

HasArch returns a boolean if a field has been set.

### GetContainerRunnerImage

`func (o *CommonsCreateRunnerImageInput) GetContainerRunnerImage() CommonsContainerRunnerImage`

GetContainerRunnerImage returns the ContainerRunnerImage field if non-nil, zero value otherwise.

### GetContainerRunnerImageOk

`func (o *CommonsCreateRunnerImageInput) GetContainerRunnerImageOk() (*CommonsContainerRunnerImage, bool)`

GetContainerRunnerImageOk returns a tuple with the ContainerRunnerImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerRunnerImage

`func (o *CommonsCreateRunnerImageInput) SetContainerRunnerImage(v CommonsContainerRunnerImage)`

SetContainerRunnerImage sets ContainerRunnerImage field to given value.

### HasContainerRunnerImage

`func (o *CommonsCreateRunnerImageInput) HasContainerRunnerImage() bool`

HasContainerRunnerImage returns a boolean if a field has been set.

### GetHooks

`func (o *CommonsCreateRunnerImageInput) GetHooks() []CommonsRunnerImageHook`

GetHooks returns the Hooks field if non-nil, zero value otherwise.

### GetHooksOk

`func (o *CommonsCreateRunnerImageInput) GetHooksOk() (*[]CommonsRunnerImageHook, bool)`

GetHooksOk returns a tuple with the Hooks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHooks

`func (o *CommonsCreateRunnerImageInput) SetHooks(v []CommonsRunnerImageHook)`

SetHooks sets Hooks field to given value.

### HasHooks

`func (o *CommonsCreateRunnerImageInput) HasHooks() bool`

HasHooks returns a boolean if a field has been set.

### GetOs

`func (o *CommonsCreateRunnerImageInput) GetOs() string`

GetOs returns the Os field if non-nil, zero value otherwise.

### GetOsOk

`func (o *CommonsCreateRunnerImageInput) GetOsOk() (*string, bool)`

GetOsOk returns a tuple with the Os field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOs

`func (o *CommonsCreateRunnerImageInput) SetOs(v string)`

SetOs sets Os field to given value.

### HasOs

`func (o *CommonsCreateRunnerImageInput) HasOs() bool`

HasOs returns a boolean if a field has been set.

### GetRunnerImagePullSecretId

`func (o *CommonsCreateRunnerImageInput) GetRunnerImagePullSecretId() string`

GetRunnerImagePullSecretId returns the RunnerImagePullSecretId field if non-nil, zero value otherwise.

### GetRunnerImagePullSecretIdOk

`func (o *CommonsCreateRunnerImageInput) GetRunnerImagePullSecretIdOk() (*string, bool)`

GetRunnerImagePullSecretIdOk returns a tuple with the RunnerImagePullSecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerImagePullSecretId

`func (o *CommonsCreateRunnerImageInput) SetRunnerImagePullSecretId(v string)`

SetRunnerImagePullSecretId sets RunnerImagePullSecretId field to given value.

### HasRunnerImagePullSecretId

`func (o *CommonsCreateRunnerImageInput) HasRunnerImagePullSecretId() bool`

HasRunnerImagePullSecretId returns a boolean if a field has been set.

### GetSettings

`func (o *CommonsCreateRunnerImageInput) GetSettings() CommonsRunnerImageSettings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *CommonsCreateRunnerImageInput) GetSettingsOk() (*CommonsRunnerImageSettings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *CommonsCreateRunnerImageInput) SetSettings(v CommonsRunnerImageSettings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *CommonsCreateRunnerImageInput) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

### GetType

`func (o *CommonsCreateRunnerImageInput) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsCreateRunnerImageInput) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsCreateRunnerImageInput) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *CommonsCreateRunnerImageInput) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


