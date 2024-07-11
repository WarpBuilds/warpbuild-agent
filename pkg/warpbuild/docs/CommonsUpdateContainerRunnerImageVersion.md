# CommonsUpdateContainerRunnerImageVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DataDirSize** | Pointer to **string** |  | [optional] 
**ImageDigest** | Pointer to **string** |  | [optional] 
**ImageSize** | Pointer to **string** | ImageSize is the size of the image. Expected format: \&quot;4G\&quot; | [optional] 
**SnapshotId** | Pointer to **string** |  | [optional] 
**SnapshotSize** | Pointer to **string** | SnapshotSize is the size of the snapshot in human readable format. | [optional] 
**VolumeId** | Pointer to **string** |  | [optional] 
**VolumeSizeGb** | Pointer to **int32** | VolumeSizeGB is the size of the volume in GB. This is used when creating a snapshot out of the volume.  This is rounded up to the nearest GiB. | [optional] 

## Methods

### NewCommonsUpdateContainerRunnerImageVersion

`func NewCommonsUpdateContainerRunnerImageVersion() *CommonsUpdateContainerRunnerImageVersion`

NewCommonsUpdateContainerRunnerImageVersion instantiates a new CommonsUpdateContainerRunnerImageVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsUpdateContainerRunnerImageVersionWithDefaults

`func NewCommonsUpdateContainerRunnerImageVersionWithDefaults() *CommonsUpdateContainerRunnerImageVersion`

NewCommonsUpdateContainerRunnerImageVersionWithDefaults instantiates a new CommonsUpdateContainerRunnerImageVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDataDirSize

`func (o *CommonsUpdateContainerRunnerImageVersion) GetDataDirSize() string`

GetDataDirSize returns the DataDirSize field if non-nil, zero value otherwise.

### GetDataDirSizeOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetDataDirSizeOk() (*string, bool)`

GetDataDirSizeOk returns a tuple with the DataDirSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataDirSize

`func (o *CommonsUpdateContainerRunnerImageVersion) SetDataDirSize(v string)`

SetDataDirSize sets DataDirSize field to given value.

### HasDataDirSize

`func (o *CommonsUpdateContainerRunnerImageVersion) HasDataDirSize() bool`

HasDataDirSize returns a boolean if a field has been set.

### GetImageDigest

`func (o *CommonsUpdateContainerRunnerImageVersion) GetImageDigest() string`

GetImageDigest returns the ImageDigest field if non-nil, zero value otherwise.

### GetImageDigestOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetImageDigestOk() (*string, bool)`

GetImageDigestOk returns a tuple with the ImageDigest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageDigest

`func (o *CommonsUpdateContainerRunnerImageVersion) SetImageDigest(v string)`

SetImageDigest sets ImageDigest field to given value.

### HasImageDigest

`func (o *CommonsUpdateContainerRunnerImageVersion) HasImageDigest() bool`

HasImageDigest returns a boolean if a field has been set.

### GetImageSize

`func (o *CommonsUpdateContainerRunnerImageVersion) GetImageSize() string`

GetImageSize returns the ImageSize field if non-nil, zero value otherwise.

### GetImageSizeOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetImageSizeOk() (*string, bool)`

GetImageSizeOk returns a tuple with the ImageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSize

`func (o *CommonsUpdateContainerRunnerImageVersion) SetImageSize(v string)`

SetImageSize sets ImageSize field to given value.

### HasImageSize

`func (o *CommonsUpdateContainerRunnerImageVersion) HasImageSize() bool`

HasImageSize returns a boolean if a field has been set.

### GetSnapshotId

`func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotId() string`

GetSnapshotId returns the SnapshotId field if non-nil, zero value otherwise.

### GetSnapshotIdOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotIdOk() (*string, bool)`

GetSnapshotIdOk returns a tuple with the SnapshotId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnapshotId

`func (o *CommonsUpdateContainerRunnerImageVersion) SetSnapshotId(v string)`

SetSnapshotId sets SnapshotId field to given value.

### HasSnapshotId

`func (o *CommonsUpdateContainerRunnerImageVersion) HasSnapshotId() bool`

HasSnapshotId returns a boolean if a field has been set.

### GetSnapshotSize

`func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotSize() string`

GetSnapshotSize returns the SnapshotSize field if non-nil, zero value otherwise.

### GetSnapshotSizeOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotSizeOk() (*string, bool)`

GetSnapshotSizeOk returns a tuple with the SnapshotSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnapshotSize

`func (o *CommonsUpdateContainerRunnerImageVersion) SetSnapshotSize(v string)`

SetSnapshotSize sets SnapshotSize field to given value.

### HasSnapshotSize

`func (o *CommonsUpdateContainerRunnerImageVersion) HasSnapshotSize() bool`

HasSnapshotSize returns a boolean if a field has been set.

### GetVolumeId

`func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeId() string`

GetVolumeId returns the VolumeId field if non-nil, zero value otherwise.

### GetVolumeIdOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeIdOk() (*string, bool)`

GetVolumeIdOk returns a tuple with the VolumeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeId

`func (o *CommonsUpdateContainerRunnerImageVersion) SetVolumeId(v string)`

SetVolumeId sets VolumeId field to given value.

### HasVolumeId

`func (o *CommonsUpdateContainerRunnerImageVersion) HasVolumeId() bool`

HasVolumeId returns a boolean if a field has been set.

### GetVolumeSizeGb

`func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeSizeGb() int32`

GetVolumeSizeGb returns the VolumeSizeGb field if non-nil, zero value otherwise.

### GetVolumeSizeGbOk

`func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeSizeGbOk() (*int32, bool)`

GetVolumeSizeGbOk returns a tuple with the VolumeSizeGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeSizeGb

`func (o *CommonsUpdateContainerRunnerImageVersion) SetVolumeSizeGb(v int32)`

SetVolumeSizeGb sets VolumeSizeGb field to given value.

### HasVolumeSizeGb

`func (o *CommonsUpdateContainerRunnerImageVersion) HasVolumeSizeGb() bool`

HasVolumeSizeGb returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


