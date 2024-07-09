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

// checks if the CommonsRunnerwiseRuntime type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerwiseRuntime{}

// CommonsRunnerwiseRuntime struct for CommonsRunnerwiseRuntime
type CommonsRunnerwiseRuntime struct {
	RunnerId *string `json:"runner_id,omitempty"`
	RuntimeSeconds *int32 `json:"runtime_seconds,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerwiseRuntime CommonsRunnerwiseRuntime

// NewCommonsRunnerwiseRuntime instantiates a new CommonsRunnerwiseRuntime object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerwiseRuntime() *CommonsRunnerwiseRuntime {
	this := CommonsRunnerwiseRuntime{}
	return &this
}

// NewCommonsRunnerwiseRuntimeWithDefaults instantiates a new CommonsRunnerwiseRuntime object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerwiseRuntimeWithDefaults() *CommonsRunnerwiseRuntime {
	this := CommonsRunnerwiseRuntime{}
	return &this
}

// GetRunnerId returns the RunnerId field value if set, zero value otherwise.
func (o *CommonsRunnerwiseRuntime) GetRunnerId() string {
	if o == nil || IsNil(o.RunnerId) {
		var ret string
		return ret
	}
	return *o.RunnerId
}

// GetRunnerIdOk returns a tuple with the RunnerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerwiseRuntime) GetRunnerIdOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerId) {
		return nil, false
	}
	return o.RunnerId, true
}

// HasRunnerId returns a boolean if a field has been set.
func (o *CommonsRunnerwiseRuntime) HasRunnerId() bool {
	if o != nil && !IsNil(o.RunnerId) {
		return true
	}

	return false
}

// SetRunnerId gets a reference to the given string and assigns it to the RunnerId field.
func (o *CommonsRunnerwiseRuntime) SetRunnerId(v string) {
	o.RunnerId = &v
}

// GetRuntimeSeconds returns the RuntimeSeconds field value if set, zero value otherwise.
func (o *CommonsRunnerwiseRuntime) GetRuntimeSeconds() int32 {
	if o == nil || IsNil(o.RuntimeSeconds) {
		var ret int32
		return ret
	}
	return *o.RuntimeSeconds
}

// GetRuntimeSecondsOk returns a tuple with the RuntimeSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerwiseRuntime) GetRuntimeSecondsOk() (*int32, bool) {
	if o == nil || IsNil(o.RuntimeSeconds) {
		return nil, false
	}
	return o.RuntimeSeconds, true
}

// HasRuntimeSeconds returns a boolean if a field has been set.
func (o *CommonsRunnerwiseRuntime) HasRuntimeSeconds() bool {
	if o != nil && !IsNil(o.RuntimeSeconds) {
		return true
	}

	return false
}

// SetRuntimeSeconds gets a reference to the given int32 and assigns it to the RuntimeSeconds field.
func (o *CommonsRunnerwiseRuntime) SetRuntimeSeconds(v int32) {
	o.RuntimeSeconds = &v
}

func (o CommonsRunnerwiseRuntime) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerwiseRuntime) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.RunnerId) {
		toSerialize["runner_id"] = o.RunnerId
	}
	if !IsNil(o.RuntimeSeconds) {
		toSerialize["runtime_seconds"] = o.RuntimeSeconds
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerwiseRuntime) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerwiseRuntime := _CommonsRunnerwiseRuntime{}

	if err = json.Unmarshal(bytes, &varCommonsRunnerwiseRuntime); err == nil {
		*o = CommonsRunnerwiseRuntime(varCommonsRunnerwiseRuntime)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "runner_id")
		delete(additionalProperties, "runtime_seconds")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerwiseRuntime struct {
	value *CommonsRunnerwiseRuntime
	isSet bool
}

func (v NullableCommonsRunnerwiseRuntime) Get() *CommonsRunnerwiseRuntime {
	return v.value
}

func (v *NullableCommonsRunnerwiseRuntime) Set(val *CommonsRunnerwiseRuntime) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerwiseRuntime) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerwiseRuntime) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerwiseRuntime(val *CommonsRunnerwiseRuntime) *NullableCommonsRunnerwiseRuntime {
	return &NullableCommonsRunnerwiseRuntime{value: val, isSet: true}
}

func (v NullableCommonsRunnerwiseRuntime) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerwiseRuntime) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


