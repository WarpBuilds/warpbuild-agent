/*
Warp Builds API Docs

This is the docs for warp builds api for argonaut

API version: 0.4.0
Contact: support@swagger.io
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package warpbuild

import (
	"encoding/json"
)

// checks if the CommonsUpdateContainerRunnerImageVersion type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsUpdateContainerRunnerImageVersion{}

// CommonsUpdateContainerRunnerImageVersion struct for CommonsUpdateContainerRunnerImageVersion
type CommonsUpdateContainerRunnerImageVersion struct {
	DataDirSize *string `json:"data_dir_size,omitempty"`
	ImageDigest *string `json:"image_digest,omitempty"`
	// ImageSize is the size of the image. Expected format: \"4G\"
	ImageSize *string `json:"image_size,omitempty"`
	SnapshotId *string `json:"snapshot_id,omitempty"`
	// SnapshotSize is the size of the snapshot in human readable format.
	SnapshotSize *string `json:"snapshot_size,omitempty"`
	VolumeId *string `json:"volume_id,omitempty"`
	// VolumeSizeGB is the size of the volume in GB. This is used when creating a snapshot out of the volume.  This is rounded up to the nearest GiB.
	VolumeSizeGb *int32 `json:"volume_size_gb,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsUpdateContainerRunnerImageVersion CommonsUpdateContainerRunnerImageVersion

// NewCommonsUpdateContainerRunnerImageVersion instantiates a new CommonsUpdateContainerRunnerImageVersion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsUpdateContainerRunnerImageVersion() *CommonsUpdateContainerRunnerImageVersion {
	this := CommonsUpdateContainerRunnerImageVersion{}
	return &this
}

// NewCommonsUpdateContainerRunnerImageVersionWithDefaults instantiates a new CommonsUpdateContainerRunnerImageVersion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsUpdateContainerRunnerImageVersionWithDefaults() *CommonsUpdateContainerRunnerImageVersion {
	this := CommonsUpdateContainerRunnerImageVersion{}
	return &this
}

// GetDataDirSize returns the DataDirSize field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetDataDirSize() string {
	if o == nil || IsNil(o.DataDirSize) {
		var ret string
		return ret
	}
	return *o.DataDirSize
}

// GetDataDirSizeOk returns a tuple with the DataDirSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetDataDirSizeOk() (*string, bool) {
	if o == nil || IsNil(o.DataDirSize) {
		return nil, false
	}
	return o.DataDirSize, true
}

// HasDataDirSize returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasDataDirSize() bool {
	if o != nil && !IsNil(o.DataDirSize) {
		return true
	}

	return false
}

// SetDataDirSize gets a reference to the given string and assigns it to the DataDirSize field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetDataDirSize(v string) {
	o.DataDirSize = &v
}

// GetImageDigest returns the ImageDigest field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetImageDigest() string {
	if o == nil || IsNil(o.ImageDigest) {
		var ret string
		return ret
	}
	return *o.ImageDigest
}

// GetImageDigestOk returns a tuple with the ImageDigest field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetImageDigestOk() (*string, bool) {
	if o == nil || IsNil(o.ImageDigest) {
		return nil, false
	}
	return o.ImageDigest, true
}

// HasImageDigest returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasImageDigest() bool {
	if o != nil && !IsNil(o.ImageDigest) {
		return true
	}

	return false
}

// SetImageDigest gets a reference to the given string and assigns it to the ImageDigest field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetImageDigest(v string) {
	o.ImageDigest = &v
}

// GetImageSize returns the ImageSize field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetImageSize() string {
	if o == nil || IsNil(o.ImageSize) {
		var ret string
		return ret
	}
	return *o.ImageSize
}

// GetImageSizeOk returns a tuple with the ImageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetImageSizeOk() (*string, bool) {
	if o == nil || IsNil(o.ImageSize) {
		return nil, false
	}
	return o.ImageSize, true
}

// HasImageSize returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasImageSize() bool {
	if o != nil && !IsNil(o.ImageSize) {
		return true
	}

	return false
}

// SetImageSize gets a reference to the given string and assigns it to the ImageSize field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetImageSize(v string) {
	o.ImageSize = &v
}

// GetSnapshotId returns the SnapshotId field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotId() string {
	if o == nil || IsNil(o.SnapshotId) {
		var ret string
		return ret
	}
	return *o.SnapshotId
}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotIdOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotId) {
		return nil, false
	}
	return o.SnapshotId, true
}

// HasSnapshotId returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasSnapshotId() bool {
	if o != nil && !IsNil(o.SnapshotId) {
		return true
	}

	return false
}

// SetSnapshotId gets a reference to the given string and assigns it to the SnapshotId field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetSnapshotId(v string) {
	o.SnapshotId = &v
}

// GetSnapshotSize returns the SnapshotSize field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotSize() string {
	if o == nil || IsNil(o.SnapshotSize) {
		var ret string
		return ret
	}
	return *o.SnapshotSize
}

// GetSnapshotSizeOk returns a tuple with the SnapshotSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetSnapshotSizeOk() (*string, bool) {
	if o == nil || IsNil(o.SnapshotSize) {
		return nil, false
	}
	return o.SnapshotSize, true
}

// HasSnapshotSize returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasSnapshotSize() bool {
	if o != nil && !IsNil(o.SnapshotSize) {
		return true
	}

	return false
}

// SetSnapshotSize gets a reference to the given string and assigns it to the SnapshotSize field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetSnapshotSize(v string) {
	o.SnapshotSize = &v
}

// GetVolumeId returns the VolumeId field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeId() string {
	if o == nil || IsNil(o.VolumeId) {
		var ret string
		return ret
	}
	return *o.VolumeId
}

// GetVolumeIdOk returns a tuple with the VolumeId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeIdOk() (*string, bool) {
	if o == nil || IsNil(o.VolumeId) {
		return nil, false
	}
	return o.VolumeId, true
}

// HasVolumeId returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasVolumeId() bool {
	if o != nil && !IsNil(o.VolumeId) {
		return true
	}

	return false
}

// SetVolumeId gets a reference to the given string and assigns it to the VolumeId field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetVolumeId(v string) {
	o.VolumeId = &v
}

// GetVolumeSizeGb returns the VolumeSizeGb field value if set, zero value otherwise.
func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeSizeGb() int32 {
	if o == nil || IsNil(o.VolumeSizeGb) {
		var ret int32
		return ret
	}
	return *o.VolumeSizeGb
}

// GetVolumeSizeGbOk returns a tuple with the VolumeSizeGb field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) GetVolumeSizeGbOk() (*int32, bool) {
	if o == nil || IsNil(o.VolumeSizeGb) {
		return nil, false
	}
	return o.VolumeSizeGb, true
}

// HasVolumeSizeGb returns a boolean if a field has been set.
func (o *CommonsUpdateContainerRunnerImageVersion) HasVolumeSizeGb() bool {
	if o != nil && !IsNil(o.VolumeSizeGb) {
		return true
	}

	return false
}

// SetVolumeSizeGb gets a reference to the given int32 and assigns it to the VolumeSizeGb field.
func (o *CommonsUpdateContainerRunnerImageVersion) SetVolumeSizeGb(v int32) {
	o.VolumeSizeGb = &v
}

func (o CommonsUpdateContainerRunnerImageVersion) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsUpdateContainerRunnerImageVersion) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.DataDirSize) {
		toSerialize["data_dir_size"] = o.DataDirSize
	}
	if !IsNil(o.ImageDigest) {
		toSerialize["image_digest"] = o.ImageDigest
	}
	if !IsNil(o.ImageSize) {
		toSerialize["image_size"] = o.ImageSize
	}
	if !IsNil(o.SnapshotId) {
		toSerialize["snapshot_id"] = o.SnapshotId
	}
	if !IsNil(o.SnapshotSize) {
		toSerialize["snapshot_size"] = o.SnapshotSize
	}
	if !IsNil(o.VolumeId) {
		toSerialize["volume_id"] = o.VolumeId
	}
	if !IsNil(o.VolumeSizeGb) {
		toSerialize["volume_size_gb"] = o.VolumeSizeGb
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsUpdateContainerRunnerImageVersion) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsUpdateContainerRunnerImageVersion := _CommonsUpdateContainerRunnerImageVersion{}

	err = json.Unmarshal(bytes, &varCommonsUpdateContainerRunnerImageVersion)

	if err != nil {
		return err
	}

	*o = CommonsUpdateContainerRunnerImageVersion(varCommonsUpdateContainerRunnerImageVersion)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "data_dir_size")
		delete(additionalProperties, "image_digest")
		delete(additionalProperties, "image_size")
		delete(additionalProperties, "snapshot_id")
		delete(additionalProperties, "snapshot_size")
		delete(additionalProperties, "volume_id")
		delete(additionalProperties, "volume_size_gb")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsUpdateContainerRunnerImageVersion struct {
	value *CommonsUpdateContainerRunnerImageVersion
	isSet bool
}

func (v NullableCommonsUpdateContainerRunnerImageVersion) Get() *CommonsUpdateContainerRunnerImageVersion {
	return v.value
}

func (v *NullableCommonsUpdateContainerRunnerImageVersion) Set(val *CommonsUpdateContainerRunnerImageVersion) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsUpdateContainerRunnerImageVersion) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsUpdateContainerRunnerImageVersion) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsUpdateContainerRunnerImageVersion(val *CommonsUpdateContainerRunnerImageVersion) *NullableCommonsUpdateContainerRunnerImageVersion {
	return &NullableCommonsUpdateContainerRunnerImageVersion{value: val, isSet: true}
}

func (v NullableCommonsUpdateContainerRunnerImageVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsUpdateContainerRunnerImageVersion) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


