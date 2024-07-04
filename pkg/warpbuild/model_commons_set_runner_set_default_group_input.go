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

// checks if the CommonsSetRunnerSetDefaultGroupInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsSetRunnerSetDefaultGroupInput{}

// CommonsSetRunnerSetDefaultGroupInput struct for CommonsSetRunnerSetDefaultGroupInput
type CommonsSetRunnerSetDefaultGroupInput struct {
	GroupId *int32 `json:"group_id,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsSetRunnerSetDefaultGroupInput CommonsSetRunnerSetDefaultGroupInput

// NewCommonsSetRunnerSetDefaultGroupInput instantiates a new CommonsSetRunnerSetDefaultGroupInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsSetRunnerSetDefaultGroupInput() *CommonsSetRunnerSetDefaultGroupInput {
	this := CommonsSetRunnerSetDefaultGroupInput{}
	return &this
}

// NewCommonsSetRunnerSetDefaultGroupInputWithDefaults instantiates a new CommonsSetRunnerSetDefaultGroupInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsSetRunnerSetDefaultGroupInputWithDefaults() *CommonsSetRunnerSetDefaultGroupInput {
	this := CommonsSetRunnerSetDefaultGroupInput{}
	return &this
}

// GetGroupId returns the GroupId field value if set, zero value otherwise.
func (o *CommonsSetRunnerSetDefaultGroupInput) GetGroupId() int32 {
	if o == nil || IsNil(o.GroupId) {
		var ret int32
		return ret
	}
	return *o.GroupId
}

// GetGroupIdOk returns a tuple with the GroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsSetRunnerSetDefaultGroupInput) GetGroupIdOk() (*int32, bool) {
	if o == nil || IsNil(o.GroupId) {
		return nil, false
	}
	return o.GroupId, true
}

// HasGroupId returns a boolean if a field has been set.
func (o *CommonsSetRunnerSetDefaultGroupInput) HasGroupId() bool {
	if o != nil && !IsNil(o.GroupId) {
		return true
	}

	return false
}

// SetGroupId gets a reference to the given int32 and assigns it to the GroupId field.
func (o *CommonsSetRunnerSetDefaultGroupInput) SetGroupId(v int32) {
	o.GroupId = &v
}

func (o CommonsSetRunnerSetDefaultGroupInput) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsSetRunnerSetDefaultGroupInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.GroupId) {
		toSerialize["group_id"] = o.GroupId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsSetRunnerSetDefaultGroupInput) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsSetRunnerSetDefaultGroupInput := _CommonsSetRunnerSetDefaultGroupInput{}

	if err = json.Unmarshal(bytes, &varCommonsSetRunnerSetDefaultGroupInput); err == nil {
		*o = CommonsSetRunnerSetDefaultGroupInput(varCommonsSetRunnerSetDefaultGroupInput)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "group_id")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsSetRunnerSetDefaultGroupInput struct {
	value *CommonsSetRunnerSetDefaultGroupInput
	isSet bool
}

func (v NullableCommonsSetRunnerSetDefaultGroupInput) Get() *CommonsSetRunnerSetDefaultGroupInput {
	return v.value
}

func (v *NullableCommonsSetRunnerSetDefaultGroupInput) Set(val *CommonsSetRunnerSetDefaultGroupInput) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsSetRunnerSetDefaultGroupInput) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsSetRunnerSetDefaultGroupInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsSetRunnerSetDefaultGroupInput(val *CommonsSetRunnerSetDefaultGroupInput) *NullableCommonsSetRunnerSetDefaultGroupInput {
	return &NullableCommonsSetRunnerSetDefaultGroupInput{value: val, isSet: true}
}

func (v NullableCommonsSetRunnerSetDefaultGroupInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsSetRunnerSetDefaultGroupInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

