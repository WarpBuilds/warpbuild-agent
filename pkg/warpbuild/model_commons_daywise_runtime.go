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

// checks if the CommonsDaywiseRuntime type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsDaywiseRuntime{}

// CommonsDaywiseRuntime struct for CommonsDaywiseRuntime
type CommonsDaywiseRuntime struct {
	Date *string `json:"date,omitempty"`
	RuntimeSeconds *int32 `json:"runtime_seconds,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsDaywiseRuntime CommonsDaywiseRuntime

// NewCommonsDaywiseRuntime instantiates a new CommonsDaywiseRuntime object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsDaywiseRuntime() *CommonsDaywiseRuntime {
	this := CommonsDaywiseRuntime{}
	return &this
}

// NewCommonsDaywiseRuntimeWithDefaults instantiates a new CommonsDaywiseRuntime object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsDaywiseRuntimeWithDefaults() *CommonsDaywiseRuntime {
	this := CommonsDaywiseRuntime{}
	return &this
}

// GetDate returns the Date field value if set, zero value otherwise.
func (o *CommonsDaywiseRuntime) GetDate() string {
	if o == nil || IsNil(o.Date) {
		var ret string
		return ret
	}
	return *o.Date
}

// GetDateOk returns a tuple with the Date field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsDaywiseRuntime) GetDateOk() (*string, bool) {
	if o == nil || IsNil(o.Date) {
		return nil, false
	}
	return o.Date, true
}

// HasDate returns a boolean if a field has been set.
func (o *CommonsDaywiseRuntime) HasDate() bool {
	if o != nil && !IsNil(o.Date) {
		return true
	}

	return false
}

// SetDate gets a reference to the given string and assigns it to the Date field.
func (o *CommonsDaywiseRuntime) SetDate(v string) {
	o.Date = &v
}

// GetRuntimeSeconds returns the RuntimeSeconds field value if set, zero value otherwise.
func (o *CommonsDaywiseRuntime) GetRuntimeSeconds() int32 {
	if o == nil || IsNil(o.RuntimeSeconds) {
		var ret int32
		return ret
	}
	return *o.RuntimeSeconds
}

// GetRuntimeSecondsOk returns a tuple with the RuntimeSeconds field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsDaywiseRuntime) GetRuntimeSecondsOk() (*int32, bool) {
	if o == nil || IsNil(o.RuntimeSeconds) {
		return nil, false
	}
	return o.RuntimeSeconds, true
}

// HasRuntimeSeconds returns a boolean if a field has been set.
func (o *CommonsDaywiseRuntime) HasRuntimeSeconds() bool {
	if o != nil && !IsNil(o.RuntimeSeconds) {
		return true
	}

	return false
}

// SetRuntimeSeconds gets a reference to the given int32 and assigns it to the RuntimeSeconds field.
func (o *CommonsDaywiseRuntime) SetRuntimeSeconds(v int32) {
	o.RuntimeSeconds = &v
}

func (o CommonsDaywiseRuntime) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsDaywiseRuntime) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Date) {
		toSerialize["date"] = o.Date
	}
	if !IsNil(o.RuntimeSeconds) {
		toSerialize["runtime_seconds"] = o.RuntimeSeconds
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsDaywiseRuntime) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsDaywiseRuntime := _CommonsDaywiseRuntime{}

	if err = json.Unmarshal(bytes, &varCommonsDaywiseRuntime); err == nil {
		*o = CommonsDaywiseRuntime(varCommonsDaywiseRuntime)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "date")
		delete(additionalProperties, "runtime_seconds")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsDaywiseRuntime struct {
	value *CommonsDaywiseRuntime
	isSet bool
}

func (v NullableCommonsDaywiseRuntime) Get() *CommonsDaywiseRuntime {
	return v.value
}

func (v *NullableCommonsDaywiseRuntime) Set(val *CommonsDaywiseRuntime) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsDaywiseRuntime) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsDaywiseRuntime) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsDaywiseRuntime(val *CommonsDaywiseRuntime) *NullableCommonsDaywiseRuntime {
	return &NullableCommonsDaywiseRuntime{value: val, isSet: true}
}

func (v NullableCommonsDaywiseRuntime) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsDaywiseRuntime) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

