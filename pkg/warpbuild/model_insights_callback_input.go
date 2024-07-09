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

// checks if the InsightsCallbackInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &InsightsCallbackInput{}

// InsightsCallbackInput struct for InsightsCallbackInput
type InsightsCallbackInput struct {
	Code string `json:"code"`
	InstallationId *string `json:"installation_id,omitempty"`
	SetupAction string `json:"setup_action"`
	AdditionalProperties map[string]interface{}
}

type _InsightsCallbackInput InsightsCallbackInput

// NewInsightsCallbackInput instantiates a new InsightsCallbackInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInsightsCallbackInput(code string, setupAction string) *InsightsCallbackInput {
	this := InsightsCallbackInput{}
	this.Code = code
	this.SetupAction = setupAction
	return &this
}

// NewInsightsCallbackInputWithDefaults instantiates a new InsightsCallbackInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInsightsCallbackInputWithDefaults() *InsightsCallbackInput {
	this := InsightsCallbackInput{}
	return &this
}

// GetCode returns the Code field value
func (o *InsightsCallbackInput) GetCode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Code
}

// GetCodeOk returns a tuple with the Code field value
// and a boolean to check if the value has been set.
func (o *InsightsCallbackInput) GetCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Code, true
}

// SetCode sets field value
func (o *InsightsCallbackInput) SetCode(v string) {
	o.Code = v
}

// GetInstallationId returns the InstallationId field value if set, zero value otherwise.
func (o *InsightsCallbackInput) GetInstallationId() string {
	if o == nil || IsNil(o.InstallationId) {
		var ret string
		return ret
	}
	return *o.InstallationId
}

// GetInstallationIdOk returns a tuple with the InstallationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InsightsCallbackInput) GetInstallationIdOk() (*string, bool) {
	if o == nil || IsNil(o.InstallationId) {
		return nil, false
	}
	return o.InstallationId, true
}

// HasInstallationId returns a boolean if a field has been set.
func (o *InsightsCallbackInput) HasInstallationId() bool {
	if o != nil && !IsNil(o.InstallationId) {
		return true
	}

	return false
}

// SetInstallationId gets a reference to the given string and assigns it to the InstallationId field.
func (o *InsightsCallbackInput) SetInstallationId(v string) {
	o.InstallationId = &v
}

// GetSetupAction returns the SetupAction field value
func (o *InsightsCallbackInput) GetSetupAction() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SetupAction
}

// GetSetupActionOk returns a tuple with the SetupAction field value
// and a boolean to check if the value has been set.
func (o *InsightsCallbackInput) GetSetupActionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SetupAction, true
}

// SetSetupAction sets field value
func (o *InsightsCallbackInput) SetSetupAction(v string) {
	o.SetupAction = v
}

func (o InsightsCallbackInput) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o InsightsCallbackInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["code"] = o.Code
	if !IsNil(o.InstallationId) {
		toSerialize["installation_id"] = o.InstallationId
	}
	toSerialize["setup_action"] = o.SetupAction

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *InsightsCallbackInput) UnmarshalJSON(bytes []byte) (err error) {
	varInsightsCallbackInput := _InsightsCallbackInput{}

	if err = json.Unmarshal(bytes, &varInsightsCallbackInput); err == nil {
		*o = InsightsCallbackInput(varInsightsCallbackInput)
	}

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "code")
		delete(additionalProperties, "installation_id")
		delete(additionalProperties, "setup_action")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableInsightsCallbackInput struct {
	value *InsightsCallbackInput
	isSet bool
}

func (v NullableInsightsCallbackInput) Get() *InsightsCallbackInput {
	return v.value
}

func (v *NullableInsightsCallbackInput) Set(val *InsightsCallbackInput) {
	v.value = val
	v.isSet = true
}

func (v NullableInsightsCallbackInput) IsSet() bool {
	return v.isSet
}

func (v *NullableInsightsCallbackInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInsightsCallbackInput(val *InsightsCallbackInput) *NullableInsightsCallbackInput {
	return &NullableInsightsCallbackInput{value: val, isSet: true}
}

func (v NullableInsightsCallbackInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInsightsCallbackInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


