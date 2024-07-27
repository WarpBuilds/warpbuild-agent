# CommonsContainerRunnerImageVersion

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DataDirSize** | Pointer to **string** |  | [optional] 
**ImageDigest** | Pointer to **string** |  | [optional] 
**ImageRepository** | Pointer to **string** |  | [optional] 
**ImageSize** | Pointer to **string** |  | [optional] 
**ImageTag** | Pointer to **string** |  | [optional] 
**ImageUri** | Pointer to **string** | ImageURI is the full uri of the image. This can be used to download the image  +OutputOnly | [optional] 
**SnapshotId** | Pointer to **string** |  | [optional] 
**SnapshotSize** | Pointer to **string** |  | [optional] 
**VolumeId** | Pointer to **string** |  | [optional] 
**VolumeSizeGb** | Pointer to **int32** | VolumeSizeGB is the size of the volume in GB. This is used when creating a snapshot out of the volume.  This is rounded up to the nearest GiB. | [optional] 

## Methods

### NewCommonsContainerRunnerImageVersion

`func NewCommonsContainerRunnerImageVersion() *CommonsContainerRunnerImageVersion`

NewCommonsContainerRunnerImageVersion instantiates a new CommonsContainerRunnerImageVersion object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsContainerRunnerImageVersionWithDefaults

`func NewCommonsContainerRunnerImageVersionWithDefaults() *CommonsContainerRunnerImageVersion`

NewCommonsContainerRunnerImageVersionWithDefaults instantiates a new CommonsContainerRunnerImageVersion object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDataDirSize

`func (o *CommonsContainerRunnerImageVersion) GetDataDirSize() string`

GetDataDirSize returns the DataDirSize field if non-nil, zero value otherwise.

### GetDataDirSizeOk

`func (o *CommonsContainerRunnerImageVersion) GetDataDirSizeOk() (*string, bool)`

GetDataDirSizeOk returns a tuple with the DataDirSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataDirSize

`func (o *CommonsContainerRunnerImageVersion) SetDataDirSize(v string)`

SetDataDirSize sets DataDirSize field to given value.

### HasDataDirSize

`func (o *CommonsContainerRunnerImageVersion) HasDataDirSize() bool`

HasDataDirSize returns a boolean if a field has been set.

### GetImageDigest

`func (o *CommonsContainerRunnerImageVersion) GetImageDigest() string`

GetImageDigest returns the ImageDigest field if non-nil, zero value otherwise.

### GetImageDigestOk

`func (o *CommonsContainerRunnerImageVersion) GetImageDigestOk() (*string, bool)`

GetImageDigestOk returns a tuple with the ImageDigest field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageDigest

`func (o *CommonsContainerRunnerImageVersion) SetImageDigest(v string)`

SetImageDigest sets ImageDigest field to given value.

### HasImageDigest

`func (o *CommonsContainerRunnerImageVersion) HasImageDigest() bool`

HasImageDigest returns a boolean if a field has been set.

### GetImageRepository

`func (o *CommonsContainerRunnerImageVersion) GetImageRepository() string`

GetImageRepository returns the ImageRepository field if non-nil, zero value otherwise.

### GetImageRepositoryOk

`func (o *CommonsContainerRunnerImageVersion) GetImageRepositoryOk() (*string, bool)`

GetImageRepositoryOk returns a tuple with the ImageRepository field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageRepository

`func (o *CommonsContainerRunnerImageVersion) SetImageRepository(v string)`

SetImageRepository sets ImageRepository field to given value.

### HasImageRepository

`func (o *CommonsContainerRunnerImageVersion) HasImageRepository() bool`

HasImageRepository returns a boolean if a field has been set.

### GetImageSize

`func (o *CommonsContainerRunnerImageVersion) GetImageSize() string`

GetImageSize returns the ImageSize field if non-nil, zero value otherwise.

### GetImageSizeOk

`func (o *CommonsContainerRunnerImageVersion) GetImageSizeOk() (*string, bool)`

GetImageSizeOk returns a tuple with the ImageSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageSize

`func (o *CommonsContainerRunnerImageVersion) SetImageSize(v string)`

SetImageSize sets ImageSize field to given value.

### HasImageSize

`func (o *CommonsContainerRunnerImageVersion) HasImageSize() bool`

HasImageSize returns a boolean if a field has been set.

### GetImageTag

`func (o *CommonsContainerRunnerImageVersion) GetImageTag() string`

GetImageTag returns the ImageTag field if non-nil, zero value otherwise.

### GetImageTagOk

`func (o *CommonsContainerRunnerImageVersion) GetImageTagOk() (*string, bool)`

GetImageTagOk returns a tuple with the ImageTag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageTag

`func (o *CommonsContainerRunnerImageVersion) SetImageTag(v string)`

SetImageTag sets ImageTag field to given value.

### HasImageTag

`func (o *CommonsContainerRunnerImageVersion) HasImageTag() bool`

HasImageTag returns a boolean if a field has been set.

### GetImageUri

`func (o *CommonsContainerRunnerImageVersion) GetImageUri() string`

GetImageUri returns the ImageUri field if non-nil, zero value otherwise.

### GetImageUriOk

`func (o *CommonsContainerRunnerImageVersion) GetImageUriOk() (*string, bool)`

GetImageUriOk returns a tuple with the ImageUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageUri

`func (o *CommonsContainerRunnerImageVersion) SetImageUri(v string)`

SetImageUri sets ImageUri field to given value.

### HasImageUri

`func (o *CommonsContainerRunnerImageVersion) HasImageUri() bool`

HasImageUri returns a boolean if a field has been set.

### GetSnapshotId

`func (o *CommonsContainerRunnerImageVersion) GetSnapshotId() string`

GetSnapshotId returns the SnapshotId field if non-nil, zero value otherwise.

### GetSnapshotIdOk

`func (o *CommonsContainerRunnerImageVersion) GetSnapshotIdOk() (*string, bool)`

GetSnapshotIdOk returns a tuple with the SnapshotId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnapshotId

`func (o *CommonsContainerRunnerImageVersion) SetSnapshotId(v string)`

SetSnapshotId sets SnapshotId field to given value.

### HasSnapshotId

`func (o *CommonsContainerRunnerImageVersion) HasSnapshotId() bool`

HasSnapshotId returns a boolean if a field has been set.

### GetSnapshotSize

`func (o *CommonsContainerRunnerImageVersion) GetSnapshotSize() string`

GetSnapshotSize returns the SnapshotSize field if non-nil, zero value otherwise.

### GetSnapshotSizeOk

`func (o *CommonsContainerRunnerImageVersion) GetSnapshotSizeOk() (*string, bool)`

GetSnapshotSizeOk returns a tuple with the SnapshotSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnapshotSize

`func (o *CommonsContainerRunnerImageVersion) SetSnapshotSize(v string)`

SetSnapshotSize sets SnapshotSize field to given value.

### HasSnapshotSize

`func (o *CommonsContainerRunnerImageVersion) HasSnapshotSize() bool`

HasSnapshotSize returns a boolean if a field has been set.

### GetVolumeId

`func (o *CommonsContainerRunnerImageVersion) GetVolumeId() string`

GetVolumeId returns the VolumeId field if non-nil, zero value otherwise.

### GetVolumeIdOk

`func (o *CommonsContainerRunnerImageVersion) GetVolumeIdOk() (*string, bool)`

GetVolumeIdOk returns a tuple with the VolumeId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeId

`func (o *CommonsContainerRunnerImageVersion) SetVolumeId(v string)`

SetVolumeId sets VolumeId field to given value.

### HasVolumeId

`func (o *CommonsContainerRunnerImageVersion) HasVolumeId() bool`

HasVolumeId returns a boolean if a field has been set.

### GetVolumeSizeGb

`func (o *CommonsContainerRunnerImageVersion) GetVolumeSizeGb() int32`

GetVolumeSizeGb returns the VolumeSizeGb field if non-nil, zero value otherwise.

### GetVolumeSizeGbOk

`func (o *CommonsContainerRunnerImageVersion) GetVolumeSizeGbOk() (*int32, bool)`

GetVolumeSizeGbOk returns a tuple with the VolumeSizeGb field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVolumeSizeGb

`func (o *CommonsContainerRunnerImageVersion) SetVolumeSizeGb(v int32)`

SetVolumeSizeGb sets VolumeSizeGb field to given value.

### HasVolumeSizeGb

`func (o *CommonsContainerRunnerImageVersion) HasVolumeSizeGb() bool`

HasVolumeSizeGb returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


