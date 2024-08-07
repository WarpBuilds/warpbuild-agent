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

// checks if the CommonsRunnerSetConfiguration type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerSetConfiguration{}

// CommonsRunnerSetConfiguration struct for CommonsRunnerSetConfiguration
type CommonsRunnerSetConfiguration struct {
	CapacityType *string `json:"capacity_type,omitempty"`
	Image *string `json:"image,omitempty"`
	Sku *string `json:"sku,omitempty"`
	Storage *CommonsStorage `json:"storage,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerSetConfiguration CommonsRunnerSetConfiguration

// NewCommonsRunnerSetConfiguration instantiates a new CommonsRunnerSetConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerSetConfiguration() *CommonsRunnerSetConfiguration {
	this := CommonsRunnerSetConfiguration{}
	return &this
}

// NewCommonsRunnerSetConfigurationWithDefaults instantiates a new CommonsRunnerSetConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerSetConfigurationWithDefaults() *CommonsRunnerSetConfiguration {
	this := CommonsRunnerSetConfiguration{}
	return &this
}

// GetCapacityType returns the CapacityType field value if set, zero value otherwise.
func (o *CommonsRunnerSetConfiguration) GetCapacityType() string {
	if o == nil || IsNil(o.CapacityType) {
		var ret string
		return ret
	}
	return *o.CapacityType
}

// GetCapacityTypeOk returns a tuple with the CapacityType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerSetConfiguration) GetCapacityTypeOk() (*string, bool) {
	if o == nil || IsNil(o.CapacityType) {
		return nil, false
	}
	return o.CapacityType, true
}

// HasCapacityType returns a boolean if a field has been set.
func (o *CommonsRunnerSetConfiguration) HasCapacityType() bool {
	if o != nil && !IsNil(o.CapacityType) {
		return true
	}

	return false
}

// SetCapacityType gets a reference to the given string and assigns it to the CapacityType field.
func (o *CommonsRunnerSetConfiguration) SetCapacityType(v string) {
	o.CapacityType = &v
}

// GetImage returns the Image field value if set, zero value otherwise.
func (o *CommonsRunnerSetConfiguration) GetImage() string {
	if o == nil || IsNil(o.Image) {
		var ret string
		return ret
	}
	return *o.Image
}

// GetImageOk returns a tuple with the Image field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerSetConfiguration) GetImageOk() (*string, bool) {
	if o == nil || IsNil(o.Image) {
		return nil, false
	}
	return o.Image, true
}

// HasImage returns a boolean if a field has been set.
func (o *CommonsRunnerSetConfiguration) HasImage() bool {
	if o != nil && !IsNil(o.Image) {
		return true
	}

	return false
}

// SetImage gets a reference to the given string and assigns it to the Image field.
func (o *CommonsRunnerSetConfiguration) SetImage(v string) {
	o.Image = &v
}

// GetSku returns the Sku field value if set, zero value otherwise.
func (o *CommonsRunnerSetConfiguration) GetSku() string {
	if o == nil || IsNil(o.Sku) {
		var ret string
		return ret
	}
	return *o.Sku
}

// GetSkuOk returns a tuple with the Sku field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerSetConfiguration) GetSkuOk() (*string, bool) {
	if o == nil || IsNil(o.Sku) {
		return nil, false
	}
	return o.Sku, true
}

// HasSku returns a boolean if a field has been set.
func (o *CommonsRunnerSetConfiguration) HasSku() bool {
	if o != nil && !IsNil(o.Sku) {
		return true
	}

	return false
}

// SetSku gets a reference to the given string and assigns it to the Sku field.
func (o *CommonsRunnerSetConfiguration) SetSku(v string) {
	o.Sku = &v
}

// GetStorage returns the Storage field value if set, zero value otherwise.
func (o *CommonsRunnerSetConfiguration) GetStorage() CommonsStorage {
	if o == nil || IsNil(o.Storage) {
		var ret CommonsStorage
		return ret
	}
	return *o.Storage
}

// GetStorageOk returns a tuple with the Storage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerSetConfiguration) GetStorageOk() (*CommonsStorage, bool) {
	if o == nil || IsNil(o.Storage) {
		return nil, false
	}
	return o.Storage, true
}

// HasStorage returns a boolean if a field has been set.
func (o *CommonsRunnerSetConfiguration) HasStorage() bool {
	if o != nil && !IsNil(o.Storage) {
		return true
	}

	return false
}

// SetStorage gets a reference to the given CommonsStorage and assigns it to the Storage field.
func (o *CommonsRunnerSetConfiguration) SetStorage(v CommonsStorage) {
	o.Storage = &v
}

func (o CommonsRunnerSetConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerSetConfiguration) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CapacityType) {
		toSerialize["capacity_type"] = o.CapacityType
	}
	if !IsNil(o.Image) {
		toSerialize["image"] = o.Image
	}
	if !IsNil(o.Sku) {
		toSerialize["sku"] = o.Sku
	}
	if !IsNil(o.Storage) {
		toSerialize["storage"] = o.Storage
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerSetConfiguration) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerSetConfiguration := _CommonsRunnerSetConfiguration{}

	err = json.Unmarshal(bytes, &varCommonsRunnerSetConfiguration)

	if err != nil {
		return err
	}

	*o = CommonsRunnerSetConfiguration(varCommonsRunnerSetConfiguration)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "capacity_type")
		delete(additionalProperties, "image")
		delete(additionalProperties, "sku")
		delete(additionalProperties, "storage")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerSetConfiguration struct {
	value *CommonsRunnerSetConfiguration
	isSet bool
}

func (v NullableCommonsRunnerSetConfiguration) Get() *CommonsRunnerSetConfiguration {
	return v.value
}

func (v *NullableCommonsRunnerSetConfiguration) Set(val *CommonsRunnerSetConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerSetConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerSetConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerSetConfiguration(val *CommonsRunnerSetConfiguration) *NullableCommonsRunnerSetConfiguration {
	return &NullableCommonsRunnerSetConfiguration{value: val, isSet: true}
}

func (v NullableCommonsRunnerSetConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerSetConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


