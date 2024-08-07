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

// checks if the CommonsRunnerImageVersion type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerImageVersion{}

// CommonsRunnerImageVersion struct for CommonsRunnerImageVersion
type CommonsRunnerImageVersion struct {
	Arch *string `json:"arch,omitempty"`
	ContainerRunnerImage *CommonsContainerRunnerImageVersion `json:"container_runner_image,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	// ExternalID is the ID of the runner image version in the external system.
	ExternalId *string `json:"external_id,omitempty"`
	Id *string `json:"id,omitempty"`
	OrganizationId *string `json:"organization_id,omitempty"`
	Os *string `json:"os,omitempty"`
	RunnerImageId *string `json:"runner_image_id,omitempty"`
	RunnerImagePullSecretId *string `json:"runner_image_pull_secret_id,omitempty"`
	Status *string `json:"status,omitempty"`
	Type *string `json:"type,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	VersionTimeId *int32 `json:"version_time_id,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerImageVersion CommonsRunnerImageVersion

// NewCommonsRunnerImageVersion instantiates a new CommonsRunnerImageVersion object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerImageVersion() *CommonsRunnerImageVersion {
	this := CommonsRunnerImageVersion{}
	return &this
}

// NewCommonsRunnerImageVersionWithDefaults instantiates a new CommonsRunnerImageVersion object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerImageVersionWithDefaults() *CommonsRunnerImageVersion {
	this := CommonsRunnerImageVersion{}
	return &this
}

// GetArch returns the Arch field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetArch() string {
	if o == nil || IsNil(o.Arch) {
		var ret string
		return ret
	}
	return *o.Arch
}

// GetArchOk returns a tuple with the Arch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetArchOk() (*string, bool) {
	if o == nil || IsNil(o.Arch) {
		return nil, false
	}
	return o.Arch, true
}

// HasArch returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasArch() bool {
	if o != nil && !IsNil(o.Arch) {
		return true
	}

	return false
}

// SetArch gets a reference to the given string and assigns it to the Arch field.
func (o *CommonsRunnerImageVersion) SetArch(v string) {
	o.Arch = &v
}

// GetContainerRunnerImage returns the ContainerRunnerImage field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetContainerRunnerImage() CommonsContainerRunnerImageVersion {
	if o == nil || IsNil(o.ContainerRunnerImage) {
		var ret CommonsContainerRunnerImageVersion
		return ret
	}
	return *o.ContainerRunnerImage
}

// GetContainerRunnerImageOk returns a tuple with the ContainerRunnerImage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetContainerRunnerImageOk() (*CommonsContainerRunnerImageVersion, bool) {
	if o == nil || IsNil(o.ContainerRunnerImage) {
		return nil, false
	}
	return o.ContainerRunnerImage, true
}

// HasContainerRunnerImage returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasContainerRunnerImage() bool {
	if o != nil && !IsNil(o.ContainerRunnerImage) {
		return true
	}

	return false
}

// SetContainerRunnerImage gets a reference to the given CommonsContainerRunnerImageVersion and assigns it to the ContainerRunnerImage field.
func (o *CommonsRunnerImageVersion) SetContainerRunnerImage(v CommonsContainerRunnerImageVersion) {
	o.ContainerRunnerImage = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetCreatedAt() string {
	if o == nil || IsNil(o.CreatedAt) {
		var ret string
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetCreatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given string and assigns it to the CreatedAt field.
func (o *CommonsRunnerImageVersion) SetCreatedAt(v string) {
	o.CreatedAt = &v
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}
	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *CommonsRunnerImageVersion) SetExternalId(v string) {
	o.ExternalId = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CommonsRunnerImageVersion) SetId(v string) {
	o.Id = &v
}

// GetOrganizationId returns the OrganizationId field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetOrganizationId() string {
	if o == nil || IsNil(o.OrganizationId) {
		var ret string
		return ret
	}
	return *o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetOrganizationIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrganizationId) {
		return nil, false
	}
	return o.OrganizationId, true
}

// HasOrganizationId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasOrganizationId() bool {
	if o != nil && !IsNil(o.OrganizationId) {
		return true
	}

	return false
}

// SetOrganizationId gets a reference to the given string and assigns it to the OrganizationId field.
func (o *CommonsRunnerImageVersion) SetOrganizationId(v string) {
	o.OrganizationId = &v
}

// GetOs returns the Os field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetOs() string {
	if o == nil || IsNil(o.Os) {
		var ret string
		return ret
	}
	return *o.Os
}

// GetOsOk returns a tuple with the Os field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetOsOk() (*string, bool) {
	if o == nil || IsNil(o.Os) {
		return nil, false
	}
	return o.Os, true
}

// HasOs returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasOs() bool {
	if o != nil && !IsNil(o.Os) {
		return true
	}

	return false
}

// SetOs gets a reference to the given string and assigns it to the Os field.
func (o *CommonsRunnerImageVersion) SetOs(v string) {
	o.Os = &v
}

// GetRunnerImageId returns the RunnerImageId field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetRunnerImageId() string {
	if o == nil || IsNil(o.RunnerImageId) {
		var ret string
		return ret
	}
	return *o.RunnerImageId
}

// GetRunnerImageIdOk returns a tuple with the RunnerImageId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetRunnerImageIdOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerImageId) {
		return nil, false
	}
	return o.RunnerImageId, true
}

// HasRunnerImageId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasRunnerImageId() bool {
	if o != nil && !IsNil(o.RunnerImageId) {
		return true
	}

	return false
}

// SetRunnerImageId gets a reference to the given string and assigns it to the RunnerImageId field.
func (o *CommonsRunnerImageVersion) SetRunnerImageId(v string) {
	o.RunnerImageId = &v
}

// GetRunnerImagePullSecretId returns the RunnerImagePullSecretId field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetRunnerImagePullSecretId() string {
	if o == nil || IsNil(o.RunnerImagePullSecretId) {
		var ret string
		return ret
	}
	return *o.RunnerImagePullSecretId
}

// GetRunnerImagePullSecretIdOk returns a tuple with the RunnerImagePullSecretId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetRunnerImagePullSecretIdOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerImagePullSecretId) {
		return nil, false
	}
	return o.RunnerImagePullSecretId, true
}

// HasRunnerImagePullSecretId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasRunnerImagePullSecretId() bool {
	if o != nil && !IsNil(o.RunnerImagePullSecretId) {
		return true
	}

	return false
}

// SetRunnerImagePullSecretId gets a reference to the given string and assigns it to the RunnerImagePullSecretId field.
func (o *CommonsRunnerImageVersion) SetRunnerImagePullSecretId(v string) {
	o.RunnerImagePullSecretId = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *CommonsRunnerImageVersion) SetStatus(v string) {
	o.Status = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *CommonsRunnerImageVersion) SetType(v string) {
	o.Type = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetUpdatedAt() string {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret string
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetUpdatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given string and assigns it to the UpdatedAt field.
func (o *CommonsRunnerImageVersion) SetUpdatedAt(v string) {
	o.UpdatedAt = &v
}

// GetVersionTimeId returns the VersionTimeId field value if set, zero value otherwise.
func (o *CommonsRunnerImageVersion) GetVersionTimeId() int32 {
	if o == nil || IsNil(o.VersionTimeId) {
		var ret int32
		return ret
	}
	return *o.VersionTimeId
}

// GetVersionTimeIdOk returns a tuple with the VersionTimeId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImageVersion) GetVersionTimeIdOk() (*int32, bool) {
	if o == nil || IsNil(o.VersionTimeId) {
		return nil, false
	}
	return o.VersionTimeId, true
}

// HasVersionTimeId returns a boolean if a field has been set.
func (o *CommonsRunnerImageVersion) HasVersionTimeId() bool {
	if o != nil && !IsNil(o.VersionTimeId) {
		return true
	}

	return false
}

// SetVersionTimeId gets a reference to the given int32 and assigns it to the VersionTimeId field.
func (o *CommonsRunnerImageVersion) SetVersionTimeId(v int32) {
	o.VersionTimeId = &v
}

func (o CommonsRunnerImageVersion) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerImageVersion) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Arch) {
		toSerialize["arch"] = o.Arch
	}
	if !IsNil(o.ContainerRunnerImage) {
		toSerialize["container_runner_image"] = o.ContainerRunnerImage
	}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.ExternalId) {
		toSerialize["external_id"] = o.ExternalId
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.OrganizationId) {
		toSerialize["organization_id"] = o.OrganizationId
	}
	if !IsNil(o.Os) {
		toSerialize["os"] = o.Os
	}
	if !IsNil(o.RunnerImageId) {
		toSerialize["runner_image_id"] = o.RunnerImageId
	}
	if !IsNil(o.RunnerImagePullSecretId) {
		toSerialize["runner_image_pull_secret_id"] = o.RunnerImagePullSecretId
	}
	if !IsNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !IsNil(o.VersionTimeId) {
		toSerialize["version_time_id"] = o.VersionTimeId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerImageVersion) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerImageVersion := _CommonsRunnerImageVersion{}

	err = json.Unmarshal(bytes, &varCommonsRunnerImageVersion)

	if err != nil {
		return err
	}

	*o = CommonsRunnerImageVersion(varCommonsRunnerImageVersion)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "arch")
		delete(additionalProperties, "container_runner_image")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "external_id")
		delete(additionalProperties, "id")
		delete(additionalProperties, "organization_id")
		delete(additionalProperties, "os")
		delete(additionalProperties, "runner_image_id")
		delete(additionalProperties, "runner_image_pull_secret_id")
		delete(additionalProperties, "status")
		delete(additionalProperties, "type")
		delete(additionalProperties, "updated_at")
		delete(additionalProperties, "version_time_id")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerImageVersion struct {
	value *CommonsRunnerImageVersion
	isSet bool
}

func (v NullableCommonsRunnerImageVersion) Get() *CommonsRunnerImageVersion {
	return v.value
}

func (v *NullableCommonsRunnerImageVersion) Set(val *CommonsRunnerImageVersion) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerImageVersion) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerImageVersion) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerImageVersion(val *CommonsRunnerImageVersion) *NullableCommonsRunnerImageVersion {
	return &NullableCommonsRunnerImageVersion{value: val, isSet: true}
}

func (v NullableCommonsRunnerImageVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerImageVersion) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


