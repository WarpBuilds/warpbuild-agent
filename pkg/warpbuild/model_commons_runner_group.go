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

// checks if the CommonsRunnerGroup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerGroup{}

// CommonsRunnerGroup struct for CommonsRunnerGroup
type CommonsRunnerGroup struct {
	Id *int32 `json:"id,omitempty"`
	Inherited *bool `json:"inherited,omitempty"`
	IsDefault *bool `json:"is_default,omitempty"`
	Name *string `json:"name,omitempty"`
	RunnersUrl *string `json:"runners_url,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerGroup CommonsRunnerGroup

// NewCommonsRunnerGroup instantiates a new CommonsRunnerGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerGroup() *CommonsRunnerGroup {
	this := CommonsRunnerGroup{}
	return &this
}

// NewCommonsRunnerGroupWithDefaults instantiates a new CommonsRunnerGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerGroupWithDefaults() *CommonsRunnerGroup {
	this := CommonsRunnerGroup{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *CommonsRunnerGroup) SetId(v int32) {
	o.Id = &v
}

// GetInherited returns the Inherited field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetInherited() bool {
	if o == nil || IsNil(o.Inherited) {
		var ret bool
		return ret
	}
	return *o.Inherited
}

// GetInheritedOk returns a tuple with the Inherited field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetInheritedOk() (*bool, bool) {
	if o == nil || IsNil(o.Inherited) {
		return nil, false
	}
	return o.Inherited, true
}

// HasInherited returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasInherited() bool {
	if o != nil && !IsNil(o.Inherited) {
		return true
	}

	return false
}

// SetInherited gets a reference to the given bool and assigns it to the Inherited field.
func (o *CommonsRunnerGroup) SetInherited(v bool) {
	o.Inherited = &v
}

// GetIsDefault returns the IsDefault field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetIsDefault() bool {
	if o == nil || IsNil(o.IsDefault) {
		var ret bool
		return ret
	}
	return *o.IsDefault
}

// GetIsDefaultOk returns a tuple with the IsDefault field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetIsDefaultOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDefault) {
		return nil, false
	}
	return o.IsDefault, true
}

// HasIsDefault returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasIsDefault() bool {
	if o != nil && !IsNil(o.IsDefault) {
		return true
	}

	return false
}

// SetIsDefault gets a reference to the given bool and assigns it to the IsDefault field.
func (o *CommonsRunnerGroup) SetIsDefault(v bool) {
	o.IsDefault = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *CommonsRunnerGroup) SetName(v string) {
	o.Name = &v
}

// GetRunnersUrl returns the RunnersUrl field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetRunnersUrl() string {
	if o == nil || IsNil(o.RunnersUrl) {
		var ret string
		return ret
	}
	return *o.RunnersUrl
}

// GetRunnersUrlOk returns a tuple with the RunnersUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetRunnersUrlOk() (*string, bool) {
	if o == nil || IsNil(o.RunnersUrl) {
		return nil, false
	}
	return o.RunnersUrl, true
}

// HasRunnersUrl returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasRunnersUrl() bool {
	if o != nil && !IsNil(o.RunnersUrl) {
		return true
	}

	return false
}

// SetRunnersUrl gets a reference to the given string and assigns it to the RunnersUrl field.
func (o *CommonsRunnerGroup) SetRunnersUrl(v string) {
	o.RunnersUrl = &v
}

// GetVisibility returns the Visibility field value if set, zero value otherwise.
func (o *CommonsRunnerGroup) GetVisibility() string {
	if o == nil || IsNil(o.Visibility) {
		var ret string
		return ret
	}
	return *o.Visibility
}

// GetVisibilityOk returns a tuple with the Visibility field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerGroup) GetVisibilityOk() (*string, bool) {
	if o == nil || IsNil(o.Visibility) {
		return nil, false
	}
	return o.Visibility, true
}

// HasVisibility returns a boolean if a field has been set.
func (o *CommonsRunnerGroup) HasVisibility() bool {
	if o != nil && !IsNil(o.Visibility) {
		return true
	}

	return false
}

// SetVisibility gets a reference to the given string and assigns it to the Visibility field.
func (o *CommonsRunnerGroup) SetVisibility(v string) {
	o.Visibility = &v
}

func (o CommonsRunnerGroup) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerGroup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Inherited) {
		toSerialize["inherited"] = o.Inherited
	}
	if !IsNil(o.IsDefault) {
		toSerialize["is_default"] = o.IsDefault
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.RunnersUrl) {
		toSerialize["runners_url"] = o.RunnersUrl
	}
	if !IsNil(o.Visibility) {
		toSerialize["visibility"] = o.Visibility
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerGroup) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerGroup := _CommonsRunnerGroup{}

	err = json.Unmarshal(bytes, &varCommonsRunnerGroup)

	if err != nil {
		return err
	}

	*o = CommonsRunnerGroup(varCommonsRunnerGroup)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "inherited")
		delete(additionalProperties, "is_default")
		delete(additionalProperties, "name")
		delete(additionalProperties, "runners_url")
		delete(additionalProperties, "visibility")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerGroup struct {
	value *CommonsRunnerGroup
	isSet bool
}

func (v NullableCommonsRunnerGroup) Get() *CommonsRunnerGroup {
	return v.value
}

func (v *NullableCommonsRunnerGroup) Set(val *CommonsRunnerGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerGroup(val *CommonsRunnerGroup) *NullableCommonsRunnerGroup {
	return &NullableCommonsRunnerGroup{value: val, isSet: true}
}

func (v NullableCommonsRunnerGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


