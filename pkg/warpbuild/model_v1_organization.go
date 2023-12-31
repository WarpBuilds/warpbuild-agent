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

// checks if the V1Organization type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1Organization{}

// V1Organization struct for V1Organization
type V1Organization struct {
	AvatarUrl *string `json:"avatar_url,omitempty"`
	CreatedBy int32 `json:"created_by"`
	Id string `json:"id"`
	Name string `json:"name"`
	AdditionalProperties map[string]interface{}
}

type _V1Organization V1Organization

// NewV1Organization instantiates a new V1Organization object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1Organization(createdBy int32, id string, name string) *V1Organization {
	this := V1Organization{}
	this.CreatedBy = createdBy
	this.Id = id
	this.Name = name
	return &this
}

// NewV1OrganizationWithDefaults instantiates a new V1Organization object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1OrganizationWithDefaults() *V1Organization {
	this := V1Organization{}
	return &this
}

// GetAvatarUrl returns the AvatarUrl field value if set, zero value otherwise.
func (o *V1Organization) GetAvatarUrl() string {
	if o == nil || IsNil(o.AvatarUrl) {
		var ret string
		return ret
	}
	return *o.AvatarUrl
}

// GetAvatarUrlOk returns a tuple with the AvatarUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1Organization) GetAvatarUrlOk() (*string, bool) {
	if o == nil || IsNil(o.AvatarUrl) {
		return nil, false
	}
	return o.AvatarUrl, true
}

// HasAvatarUrl returns a boolean if a field has been set.
func (o *V1Organization) HasAvatarUrl() bool {
	if o != nil && !IsNil(o.AvatarUrl) {
		return true
	}

	return false
}

// SetAvatarUrl gets a reference to the given string and assigns it to the AvatarUrl field.
func (o *V1Organization) SetAvatarUrl(v string) {
	o.AvatarUrl = &v
}

// GetCreatedBy returns the CreatedBy field value
func (o *V1Organization) GetCreatedBy() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
func (o *V1Organization) GetCreatedByOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *V1Organization) SetCreatedBy(v int32) {
	o.CreatedBy = v
}

// GetId returns the Id field value
func (o *V1Organization) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *V1Organization) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *V1Organization) SetId(v string) {
	o.Id = v
}

// GetName returns the Name field value
func (o *V1Organization) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *V1Organization) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *V1Organization) SetName(v string) {
	o.Name = v
}

func (o V1Organization) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1Organization) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AvatarUrl) {
		toSerialize["avatar_url"] = o.AvatarUrl
	}
	toSerialize["created_by"] = o.CreatedBy
	toSerialize["id"] = o.Id
	toSerialize["name"] = o.Name

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1Organization) UnmarshalJSON(bytes []byte) (err error) {
	varV1Organization := _V1Organization{}

	err = json.Unmarshal(bytes, &varV1Organization)

	if err != nil {
		return err
	}

	*o = V1Organization(varV1Organization)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "avatar_url")
		delete(additionalProperties, "created_by")
		delete(additionalProperties, "id")
		delete(additionalProperties, "name")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1Organization struct {
	value *V1Organization
	isSet bool
}

func (v NullableV1Organization) Get() *V1Organization {
	return v.value
}

func (v *NullableV1Organization) Set(val *V1Organization) {
	v.value = val
	v.isSet = true
}

func (v NullableV1Organization) IsSet() bool {
	return v.isSet
}

func (v *NullableV1Organization) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1Organization(val *V1Organization) *NullableV1Organization {
	return &NullableV1Organization{value: val, isSet: true}
}

func (v NullableV1Organization) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1Organization) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


