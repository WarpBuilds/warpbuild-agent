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

// checks if the CommonsProviderInstanceSkuMapping type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsProviderInstanceSkuMapping{}

// CommonsProviderInstanceSkuMapping struct for CommonsProviderInstanceSkuMapping
type CommonsProviderInstanceSkuMapping struct {
	Id *string `json:"id,omitempty"`
	Priority *int32 `json:"priority,omitempty"`
	Provider *string `json:"provider,omitempty"`
	ProviderSku *string `json:"provider_sku,omitempty"`
	ProviderSkuMeta map[string]interface{} `json:"provider_sku_meta,omitempty"`
	SkuId *string `json:"sku_id,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsProviderInstanceSkuMapping CommonsProviderInstanceSkuMapping

// NewCommonsProviderInstanceSkuMapping instantiates a new CommonsProviderInstanceSkuMapping object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsProviderInstanceSkuMapping() *CommonsProviderInstanceSkuMapping {
	this := CommonsProviderInstanceSkuMapping{}
	return &this
}

// NewCommonsProviderInstanceSkuMappingWithDefaults instantiates a new CommonsProviderInstanceSkuMapping object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsProviderInstanceSkuMappingWithDefaults() *CommonsProviderInstanceSkuMapping {
	this := CommonsProviderInstanceSkuMapping{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CommonsProviderInstanceSkuMapping) SetId(v string) {
	o.Id = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetPriority() int32 {
	if o == nil || IsNil(o.Priority) {
		var ret int32
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetPriorityOk() (*int32, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}
	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given int32 and assigns it to the Priority field.
func (o *CommonsProviderInstanceSkuMapping) SetPriority(v int32) {
	o.Priority = &v
}

// GetProvider returns the Provider field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetProvider() string {
	if o == nil || IsNil(o.Provider) {
		var ret string
		return ret
	}
	return *o.Provider
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetProviderOk() (*string, bool) {
	if o == nil || IsNil(o.Provider) {
		return nil, false
	}
	return o.Provider, true
}

// HasProvider returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasProvider() bool {
	if o != nil && !IsNil(o.Provider) {
		return true
	}

	return false
}

// SetProvider gets a reference to the given string and assigns it to the Provider field.
func (o *CommonsProviderInstanceSkuMapping) SetProvider(v string) {
	o.Provider = &v
}

// GetProviderSku returns the ProviderSku field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetProviderSku() string {
	if o == nil || IsNil(o.ProviderSku) {
		var ret string
		return ret
	}
	return *o.ProviderSku
}

// GetProviderSkuOk returns a tuple with the ProviderSku field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetProviderSkuOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderSku) {
		return nil, false
	}
	return o.ProviderSku, true
}

// HasProviderSku returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasProviderSku() bool {
	if o != nil && !IsNil(o.ProviderSku) {
		return true
	}

	return false
}

// SetProviderSku gets a reference to the given string and assigns it to the ProviderSku field.
func (o *CommonsProviderInstanceSkuMapping) SetProviderSku(v string) {
	o.ProviderSku = &v
}

// GetProviderSkuMeta returns the ProviderSkuMeta field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetProviderSkuMeta() map[string]interface{} {
	if o == nil || IsNil(o.ProviderSkuMeta) {
		var ret map[string]interface{}
		return ret
	}
	return o.ProviderSkuMeta
}

// GetProviderSkuMetaOk returns a tuple with the ProviderSkuMeta field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetProviderSkuMetaOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.ProviderSkuMeta) {
		return map[string]interface{}{}, false
	}
	return o.ProviderSkuMeta, true
}

// HasProviderSkuMeta returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasProviderSkuMeta() bool {
	if o != nil && !IsNil(o.ProviderSkuMeta) {
		return true
	}

	return false
}

// SetProviderSkuMeta gets a reference to the given map[string]interface{} and assigns it to the ProviderSkuMeta field.
func (o *CommonsProviderInstanceSkuMapping) SetProviderSkuMeta(v map[string]interface{}) {
	o.ProviderSkuMeta = v
}

// GetSkuId returns the SkuId field value if set, zero value otherwise.
func (o *CommonsProviderInstanceSkuMapping) GetSkuId() string {
	if o == nil || IsNil(o.SkuId) {
		var ret string
		return ret
	}
	return *o.SkuId
}

// GetSkuIdOk returns a tuple with the SkuId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsProviderInstanceSkuMapping) GetSkuIdOk() (*string, bool) {
	if o == nil || IsNil(o.SkuId) {
		return nil, false
	}
	return o.SkuId, true
}

// HasSkuId returns a boolean if a field has been set.
func (o *CommonsProviderInstanceSkuMapping) HasSkuId() bool {
	if o != nil && !IsNil(o.SkuId) {
		return true
	}

	return false
}

// SetSkuId gets a reference to the given string and assigns it to the SkuId field.
func (o *CommonsProviderInstanceSkuMapping) SetSkuId(v string) {
	o.SkuId = &v
}

func (o CommonsProviderInstanceSkuMapping) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsProviderInstanceSkuMapping) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Priority) {
		toSerialize["priority"] = o.Priority
	}
	if !IsNil(o.Provider) {
		toSerialize["provider"] = o.Provider
	}
	if !IsNil(o.ProviderSku) {
		toSerialize["provider_sku"] = o.ProviderSku
	}
	if !IsNil(o.ProviderSkuMeta) {
		toSerialize["provider_sku_meta"] = o.ProviderSkuMeta
	}
	if !IsNil(o.SkuId) {
		toSerialize["sku_id"] = o.SkuId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsProviderInstanceSkuMapping) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsProviderInstanceSkuMapping := _CommonsProviderInstanceSkuMapping{}

	err = json.Unmarshal(bytes, &varCommonsProviderInstanceSkuMapping)

	if err != nil {
		return err
	}

	*o = CommonsProviderInstanceSkuMapping(varCommonsProviderInstanceSkuMapping)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "priority")
		delete(additionalProperties, "provider")
		delete(additionalProperties, "provider_sku")
		delete(additionalProperties, "provider_sku_meta")
		delete(additionalProperties, "sku_id")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsProviderInstanceSkuMapping struct {
	value *CommonsProviderInstanceSkuMapping
	isSet bool
}

func (v NullableCommonsProviderInstanceSkuMapping) Get() *CommonsProviderInstanceSkuMapping {
	return v.value
}

func (v *NullableCommonsProviderInstanceSkuMapping) Set(val *CommonsProviderInstanceSkuMapping) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsProviderInstanceSkuMapping) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsProviderInstanceSkuMapping) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsProviderInstanceSkuMapping(val *CommonsProviderInstanceSkuMapping) *NullableCommonsProviderInstanceSkuMapping {
	return &NullableCommonsProviderInstanceSkuMapping{value: val, isSet: true}
}

func (v NullableCommonsProviderInstanceSkuMapping) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsProviderInstanceSkuMapping) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


