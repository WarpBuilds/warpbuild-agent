# CommonsUpdateRunnerImageInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ContainerRunnerImage** | Pointer to [**CommonsContainerRunnerImageUpdate**](CommonsContainerRunnerImageUpdate.md) |  | [optional] 
**Hooks** | Pointer to [**[]CommonsRunnerImageHook**](CommonsRunnerImageHook.md) |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**RunnerImagePullSecretId** | Pointer to **string** |  | [optional] 
**Settings** | Pointer to [**CommonsRunnerImageSettings**](CommonsRunnerImageSettings.md) |  | [optional] 

## Methods

### NewCommonsUpdateRunnerImageInput

`func NewCommonsUpdateRunnerImageInput() *CommonsUpdateRunnerImageInput`

NewCommonsUpdateRunnerImageInput instantiates a new CommonsUpdateRunnerImageInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsUpdateRunnerImageInputWithDefaults

`func NewCommonsUpdateRunnerImageInputWithDefaults() *CommonsUpdateRunnerImageInput`

NewCommonsUpdateRunnerImageInputWithDefaults instantiates a new CommonsUpdateRunnerImageInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContainerRunnerImage

`func (o *CommonsUpdateRunnerImageInput) GetContainerRunnerImage() CommonsContainerRunnerImageUpdate`

GetContainerRunnerImage returns the ContainerRunnerImage field if non-nil, zero value otherwise.

### GetContainerRunnerImageOk

`func (o *CommonsUpdateRunnerImageInput) GetContainerRunnerImageOk() (*CommonsContainerRunnerImageUpdate, bool)`

GetContainerRunnerImageOk returns a tuple with the ContainerRunnerImage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainerRunnerImage

`func (o *CommonsUpdateRunnerImageInput) SetContainerRunnerImage(v CommonsContainerRunnerImageUpdate)`

SetContainerRunnerImage sets ContainerRunnerImage field to given value.

### HasContainerRunnerImage

`func (o *CommonsUpdateRunnerImageInput) HasContainerRunnerImage() bool`

HasContainerRunnerImage returns a boolean if a field has been set.

### GetHooks

`func (o *CommonsUpdateRunnerImageInput) GetHooks() []CommonsRunnerImageHook`

GetHooks returns the Hooks field if non-nil, zero value otherwise.

### GetHooksOk

`func (o *CommonsUpdateRunnerImageInput) GetHooksOk() (*[]CommonsRunnerImageHook, bool)`

GetHooksOk returns a tuple with the Hooks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHooks

`func (o *CommonsUpdateRunnerImageInput) SetHooks(v []CommonsRunnerImageHook)`

SetHooks sets Hooks field to given value.

### HasHooks

`func (o *CommonsUpdateRunnerImageInput) HasHooks() bool`

HasHooks returns a boolean if a field has been set.

### GetId

`func (o *CommonsUpdateRunnerImageInput) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsUpdateRunnerImageInput) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsUpdateRunnerImageInput) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsUpdateRunnerImageInput) HasId() bool`

HasId returns a boolean if a field has been set.

### GetRunnerImagePullSecretId

`func (o *CommonsUpdateRunnerImageInput) GetRunnerImagePullSecretId() string`

GetRunnerImagePullSecretId returns the RunnerImagePullSecretId field if non-nil, zero value otherwise.

### GetRunnerImagePullSecretIdOk

`func (o *CommonsUpdateRunnerImageInput) GetRunnerImagePullSecretIdOk() (*string, bool)`

GetRunnerImagePullSecretIdOk returns a tuple with the RunnerImagePullSecretId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerImagePullSecretId

`func (o *CommonsUpdateRunnerImageInput) SetRunnerImagePullSecretId(v string)`

SetRunnerImagePullSecretId sets RunnerImagePullSecretId field to given value.

### HasRunnerImagePullSecretId

`func (o *CommonsUpdateRunnerImageInput) HasRunnerImagePullSecretId() bool`

HasRunnerImagePullSecretId returns a boolean if a field has been set.

### GetSettings

`func (o *CommonsUpdateRunnerImageInput) GetSettings() CommonsRunnerImageSettings`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *CommonsUpdateRunnerImageInput) GetSettingsOk() (*CommonsRunnerImageSettings, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *CommonsUpdateRunnerImageInput) SetSettings(v CommonsRunnerImageSettings)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *CommonsUpdateRunnerImageInput) HasSettings() bool`

HasSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


