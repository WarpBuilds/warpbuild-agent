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

// checks if the CommonsRunnersCostSummary type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnersCostSummary{}

// CommonsRunnersCostSummary struct for CommonsRunnersCostSummary
type CommonsRunnersCostSummary struct {
	Arm64BilledMinutes *int32 `json:"arm64_billed_minutes,omitempty"`
	TotalCost *float32 `json:"total_cost,omitempty"`
	X64BilledMinutes *int32 `json:"x64_billed_minutes,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnersCostSummary CommonsRunnersCostSummary

// NewCommonsRunnersCostSummary instantiates a new CommonsRunnersCostSummary object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnersCostSummary() *CommonsRunnersCostSummary {
	this := CommonsRunnersCostSummary{}
	return &this
}

// NewCommonsRunnersCostSummaryWithDefaults instantiates a new CommonsRunnersCostSummary object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnersCostSummaryWithDefaults() *CommonsRunnersCostSummary {
	this := CommonsRunnersCostSummary{}
	return &this
}

// GetArm64BilledMinutes returns the Arm64BilledMinutes field value if set, zero value otherwise.
func (o *CommonsRunnersCostSummary) GetArm64BilledMinutes() int32 {
	if o == nil || IsNil(o.Arm64BilledMinutes) {
		var ret int32
		return ret
	}
	return *o.Arm64BilledMinutes
}

// GetArm64BilledMinutesOk returns a tuple with the Arm64BilledMinutes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnersCostSummary) GetArm64BilledMinutesOk() (*int32, bool) {
	if o == nil || IsNil(o.Arm64BilledMinutes) {
		return nil, false
	}
	return o.Arm64BilledMinutes, true
}

// HasArm64BilledMinutes returns a boolean if a field has been set.
func (o *CommonsRunnersCostSummary) HasArm64BilledMinutes() bool {
	if o != nil && !IsNil(o.Arm64BilledMinutes) {
		return true
	}

	return false
}

// SetArm64BilledMinutes gets a reference to the given int32 and assigns it to the Arm64BilledMinutes field.
func (o *CommonsRunnersCostSummary) SetArm64BilledMinutes(v int32) {
	o.Arm64BilledMinutes = &v
}

// GetTotalCost returns the TotalCost field value if set, zero value otherwise.
func (o *CommonsRunnersCostSummary) GetTotalCost() float32 {
	if o == nil || IsNil(o.TotalCost) {
		var ret float32
		return ret
	}
	return *o.TotalCost
}

// GetTotalCostOk returns a tuple with the TotalCost field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnersCostSummary) GetTotalCostOk() (*float32, bool) {
	if o == nil || IsNil(o.TotalCost) {
		return nil, false
	}
	return o.TotalCost, true
}

// HasTotalCost returns a boolean if a field has been set.
func (o *CommonsRunnersCostSummary) HasTotalCost() bool {
	if o != nil && !IsNil(o.TotalCost) {
		return true
	}

	return false
}

// SetTotalCost gets a reference to the given float32 and assigns it to the TotalCost field.
func (o *CommonsRunnersCostSummary) SetTotalCost(v float32) {
	o.TotalCost = &v
}

// GetX64BilledMinutes returns the X64BilledMinutes field value if set, zero value otherwise.
func (o *CommonsRunnersCostSummary) GetX64BilledMinutes() int32 {
	if o == nil || IsNil(o.X64BilledMinutes) {
		var ret int32
		return ret
	}
	return *o.X64BilledMinutes
}

// GetX64BilledMinutesOk returns a tuple with the X64BilledMinutes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnersCostSummary) GetX64BilledMinutesOk() (*int32, bool) {
	if o == nil || IsNil(o.X64BilledMinutes) {
		return nil, false
	}
	return o.X64BilledMinutes, true
}

// HasX64BilledMinutes returns a boolean if a field has been set.
func (o *CommonsRunnersCostSummary) HasX64BilledMinutes() bool {
	if o != nil && !IsNil(o.X64BilledMinutes) {
		return true
	}

	return false
}

// SetX64BilledMinutes gets a reference to the given int32 and assigns it to the X64BilledMinutes field.
func (o *CommonsRunnersCostSummary) SetX64BilledMinutes(v int32) {
	o.X64BilledMinutes = &v
}

func (o CommonsRunnersCostSummary) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnersCostSummary) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Arm64BilledMinutes) {
		toSerialize["arm64_billed_minutes"] = o.Arm64BilledMinutes
	}
	if !IsNil(o.TotalCost) {
		toSerialize["total_cost"] = o.TotalCost
	}
	if !IsNil(o.X64BilledMinutes) {
		toSerialize["x64_billed_minutes"] = o.X64BilledMinutes
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnersCostSummary) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnersCostSummary := _CommonsRunnersCostSummary{}

	if err = json.Unmarshal(bytes, &varCommonsRunnersCostSummary); err == nil {
		*o = CommonsRunnersCostSummary(varCommonsRunnersCostSummary)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "arm64_billed_minutes")
		delete(additionalProperties, "total_cost")
		delete(additionalProperties, "x64_billed_minutes")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnersCostSummary struct {
	value *CommonsRunnersCostSummary
	isSet bool
}

func (v NullableCommonsRunnersCostSummary) Get() *CommonsRunnersCostSummary {
	return v.value
}

func (v *NullableCommonsRunnersCostSummary) Set(val *CommonsRunnersCostSummary) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnersCostSummary) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnersCostSummary) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnersCostSummary(val *CommonsRunnersCostSummary) *NullableCommonsRunnersCostSummary {
	return &NullableCommonsRunnersCostSummary{value: val, isSet: true}
}

func (v NullableCommonsRunnersCostSummary) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnersCostSummary) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

