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

// checks if the CommonsUpdateBillingInfoInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsUpdateBillingInfoInput{}

// CommonsUpdateBillingInfoInput struct for CommonsUpdateBillingInfoInput
type CommonsUpdateBillingInfoInput struct {
	Email *string `json:"email,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsUpdateBillingInfoInput CommonsUpdateBillingInfoInput

// NewCommonsUpdateBillingInfoInput instantiates a new CommonsUpdateBillingInfoInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsUpdateBillingInfoInput() *CommonsUpdateBillingInfoInput {
	this := CommonsUpdateBillingInfoInput{}
	return &this
}

// NewCommonsUpdateBillingInfoInputWithDefaults instantiates a new CommonsUpdateBillingInfoInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsUpdateBillingInfoInputWithDefaults() *CommonsUpdateBillingInfoInput {
	this := CommonsUpdateBillingInfoInput{}
	return &this
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *CommonsUpdateBillingInfoInput) GetEmail() string {
	if o == nil || IsNil(o.Email) {
		var ret string
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpdateBillingInfoInput) GetEmailOk() (*string, bool) {
	if o == nil || IsNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *CommonsUpdateBillingInfoInput) HasEmail() bool {
	if o != nil && !IsNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given string and assigns it to the Email field.
func (o *CommonsUpdateBillingInfoInput) SetEmail(v string) {
	o.Email = &v
}

func (o CommonsUpdateBillingInfoInput) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsUpdateBillingInfoInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Email) {
		toSerialize["email"] = o.Email
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsUpdateBillingInfoInput) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsUpdateBillingInfoInput := _CommonsUpdateBillingInfoInput{}

	if err = json.Unmarshal(bytes, &varCommonsUpdateBillingInfoInput); err == nil {
		*o = CommonsUpdateBillingInfoInput(varCommonsUpdateBillingInfoInput)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "email")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsUpdateBillingInfoInput struct {
	value *CommonsUpdateBillingInfoInput
	isSet bool
}

func (v NullableCommonsUpdateBillingInfoInput) Get() *CommonsUpdateBillingInfoInput {
	return v.value
}

func (v *NullableCommonsUpdateBillingInfoInput) Set(val *CommonsUpdateBillingInfoInput) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsUpdateBillingInfoInput) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsUpdateBillingInfoInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsUpdateBillingInfoInput(val *CommonsUpdateBillingInfoInput) *NullableCommonsUpdateBillingInfoInput {
	return &NullableCommonsUpdateBillingInfoInput{value: val, isSet: true}
}

func (v NullableCommonsUpdateBillingInfoInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsUpdateBillingInfoInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

