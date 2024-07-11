# CommonsRunnerImageSettings

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PurgeImageVersionsOffset** | Pointer to **int32** | PurgeImageVersionsOffset is the number of versions to keep. Each time a new version is created, the oldest version is purged. Allowed values range [1, inf).  Default value is 2. | [optional] 

## Methods

### NewCommonsRunnerImageSettings

`func NewCommonsRunnerImageSettings() *CommonsRunnerImageSettings`

NewCommonsRunnerImageSettings instantiates a new CommonsRunnerImageSettings object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImageSettingsWithDefaults

`func NewCommonsRunnerImageSettingsWithDefaults() *CommonsRunnerImageSettings`

NewCommonsRunnerImageSettingsWithDefaults instantiates a new CommonsRunnerImageSettings object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPurgeImageVersionsOffset

`func (o *CommonsRunnerImageSettings) GetPurgeImageVersionsOffset() int32`

GetPurgeImageVersionsOffset returns the PurgeImageVersionsOffset field if non-nil, zero value otherwise.

### GetPurgeImageVersionsOffsetOk

`func (o *CommonsRunnerImageSettings) GetPurgeImageVersionsOffsetOk() (*int32, bool)`

GetPurgeImageVersionsOffsetOk returns a tuple with the PurgeImageVersionsOffset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurgeImageVersionsOffset

`func (o *CommonsRunnerImageSettings) SetPurgeImageVersionsOffset(v int32)`

SetPurgeImageVersionsOffset sets PurgeImageVersionsOffset field to given value.

### HasPurgeImageVersionsOffset

`func (o *CommonsRunnerImageSettings) HasPurgeImageVersionsOffset() bool`

HasPurgeImageVersionsOffset returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


