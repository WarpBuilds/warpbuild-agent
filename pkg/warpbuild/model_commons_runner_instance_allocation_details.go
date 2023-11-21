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

// checks if the CommonsRunnerInstanceAllocationDetails type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerInstanceAllocationDetails{}

// CommonsRunnerInstanceAllocationDetails struct for CommonsRunnerInstanceAllocationDetails
type CommonsRunnerInstanceAllocationDetails struct {
	GhRunnerApplicationDetails *CommonsGithubRunnerApplicationDetails `json:"gh_runner_application_details,omitempty"`
	RunnerApplication *string `json:"runner_application,omitempty"`
	RunnerInstance *CommonsRunnerInstance `json:"runner_instance,omitempty"`
	Status *string `json:"status,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerInstanceAllocationDetails CommonsRunnerInstanceAllocationDetails

// NewCommonsRunnerInstanceAllocationDetails instantiates a new CommonsRunnerInstanceAllocationDetails object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerInstanceAllocationDetails() *CommonsRunnerInstanceAllocationDetails {
	this := CommonsRunnerInstanceAllocationDetails{}
	return &this
}

// NewCommonsRunnerInstanceAllocationDetailsWithDefaults instantiates a new CommonsRunnerInstanceAllocationDetails object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerInstanceAllocationDetailsWithDefaults() *CommonsRunnerInstanceAllocationDetails {
	this := CommonsRunnerInstanceAllocationDetails{}
	return &this
}

// GetGhRunnerApplicationDetails returns the GhRunnerApplicationDetails field value if set, zero value otherwise.
func (o *CommonsRunnerInstanceAllocationDetails) GetGhRunnerApplicationDetails() CommonsGithubRunnerApplicationDetails {
	if o == nil || IsNil(o.GhRunnerApplicationDetails) {
		var ret CommonsGithubRunnerApplicationDetails
		return ret
	}
	return *o.GhRunnerApplicationDetails
}

// GetGhRunnerApplicationDetailsOk returns a tuple with the GhRunnerApplicationDetails field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstanceAllocationDetails) GetGhRunnerApplicationDetailsOk() (*CommonsGithubRunnerApplicationDetails, bool) {
	if o == nil || IsNil(o.GhRunnerApplicationDetails) {
		return nil, false
	}
	return o.GhRunnerApplicationDetails, true
}

// HasGhRunnerApplicationDetails returns a boolean if a field has been set.
func (o *CommonsRunnerInstanceAllocationDetails) HasGhRunnerApplicationDetails() bool {
	if o != nil && !IsNil(o.GhRunnerApplicationDetails) {
		return true
	}

	return false
}

// SetGhRunnerApplicationDetails gets a reference to the given CommonsGithubRunnerApplicationDetails and assigns it to the GhRunnerApplicationDetails field.
func (o *CommonsRunnerInstanceAllocationDetails) SetGhRunnerApplicationDetails(v CommonsGithubRunnerApplicationDetails) {
	o.GhRunnerApplicationDetails = &v
}

// GetRunnerApplication returns the RunnerApplication field value if set, zero value otherwise.
func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplication() string {
	if o == nil || IsNil(o.RunnerApplication) {
		var ret string
		return ret
	}
	return *o.RunnerApplication
}

// GetRunnerApplicationOk returns a tuple with the RunnerApplication field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplicationOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerApplication) {
		return nil, false
	}
	return o.RunnerApplication, true
}

// HasRunnerApplication returns a boolean if a field has been set.
func (o *CommonsRunnerInstanceAllocationDetails) HasRunnerApplication() bool {
	if o != nil && !IsNil(o.RunnerApplication) {
		return true
	}

	return false
}

// SetRunnerApplication gets a reference to the given string and assigns it to the RunnerApplication field.
func (o *CommonsRunnerInstanceAllocationDetails) SetRunnerApplication(v string) {
	o.RunnerApplication = &v
}

// GetRunnerInstance returns the RunnerInstance field value if set, zero value otherwise.
func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerInstance() CommonsRunnerInstance {
	if o == nil || IsNil(o.RunnerInstance) {
		var ret CommonsRunnerInstance
		return ret
	}
	return *o.RunnerInstance
}

// GetRunnerInstanceOk returns a tuple with the RunnerInstance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerInstanceOk() (*CommonsRunnerInstance, bool) {
	if o == nil || IsNil(o.RunnerInstance) {
		return nil, false
	}
	return o.RunnerInstance, true
}

// HasRunnerInstance returns a boolean if a field has been set.
func (o *CommonsRunnerInstanceAllocationDetails) HasRunnerInstance() bool {
	if o != nil && !IsNil(o.RunnerInstance) {
		return true
	}

	return false
}

// SetRunnerInstance gets a reference to the given CommonsRunnerInstance and assigns it to the RunnerInstance field.
func (o *CommonsRunnerInstanceAllocationDetails) SetRunnerInstance(v CommonsRunnerInstance) {
	o.RunnerInstance = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *CommonsRunnerInstanceAllocationDetails) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstanceAllocationDetails) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *CommonsRunnerInstanceAllocationDetails) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *CommonsRunnerInstanceAllocationDetails) SetStatus(v string) {
	o.Status = &v
}

func (o CommonsRunnerInstanceAllocationDetails) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerInstanceAllocationDetails) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.GhRunnerApplicationDetails) {
		toSerialize["gh_runner_application_details"] = o.GhRunnerApplicationDetails
	}
	if !IsNil(o.RunnerApplication) {
		toSerialize["runner_application"] = o.RunnerApplication
	}
	if !IsNil(o.RunnerInstance) {
		toSerialize["runner_instance"] = o.RunnerInstance
	}
	if !IsNil(o.Status) {
		toSerialize["status"] = o.Status
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerInstanceAllocationDetails) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerInstanceAllocationDetails := _CommonsRunnerInstanceAllocationDetails{}

	err = json.Unmarshal(bytes, &varCommonsRunnerInstanceAllocationDetails)

	if err != nil {
		return err
	}

	*o = CommonsRunnerInstanceAllocationDetails(varCommonsRunnerInstanceAllocationDetails)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "gh_runner_application_details")
		delete(additionalProperties, "runner_application")
		delete(additionalProperties, "runner_instance")
		delete(additionalProperties, "status")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerInstanceAllocationDetails struct {
	value *CommonsRunnerInstanceAllocationDetails
	isSet bool
}

func (v NullableCommonsRunnerInstanceAllocationDetails) Get() *CommonsRunnerInstanceAllocationDetails {
	return v.value
}

func (v *NullableCommonsRunnerInstanceAllocationDetails) Set(val *CommonsRunnerInstanceAllocationDetails) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerInstanceAllocationDetails) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerInstanceAllocationDetails) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerInstanceAllocationDetails(val *CommonsRunnerInstanceAllocationDetails) *NullableCommonsRunnerInstanceAllocationDetails {
	return &NullableCommonsRunnerInstanceAllocationDetails{value: val, isSet: true}
}

func (v NullableCommonsRunnerInstanceAllocationDetails) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerInstanceAllocationDetails) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


