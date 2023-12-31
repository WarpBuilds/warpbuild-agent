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

// checks if the CommonsUpcomingBill type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsUpcomingBill{}

// CommonsUpcomingBill struct for CommonsUpcomingBill
type CommonsUpcomingBill struct {
	Amount *int32 `json:"amount,omitempty"`
	Currency *string `json:"currency,omitempty"`
	Date *string `json:"date,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsUpcomingBill CommonsUpcomingBill

// NewCommonsUpcomingBill instantiates a new CommonsUpcomingBill object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsUpcomingBill() *CommonsUpcomingBill {
	this := CommonsUpcomingBill{}
	return &this
}

// NewCommonsUpcomingBillWithDefaults instantiates a new CommonsUpcomingBill object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsUpcomingBillWithDefaults() *CommonsUpcomingBill {
	this := CommonsUpcomingBill{}
	return &this
}

// GetAmount returns the Amount field value if set, zero value otherwise.
func (o *CommonsUpcomingBill) GetAmount() int32 {
	if o == nil || IsNil(o.Amount) {
		var ret int32
		return ret
	}
	return *o.Amount
}

// GetAmountOk returns a tuple with the Amount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpcomingBill) GetAmountOk() (*int32, bool) {
	if o == nil || IsNil(o.Amount) {
		return nil, false
	}
	return o.Amount, true
}

// HasAmount returns a boolean if a field has been set.
func (o *CommonsUpcomingBill) HasAmount() bool {
	if o != nil && !IsNil(o.Amount) {
		return true
	}

	return false
}

// SetAmount gets a reference to the given int32 and assigns it to the Amount field.
func (o *CommonsUpcomingBill) SetAmount(v int32) {
	o.Amount = &v
}

// GetCurrency returns the Currency field value if set, zero value otherwise.
func (o *CommonsUpcomingBill) GetCurrency() string {
	if o == nil || IsNil(o.Currency) {
		var ret string
		return ret
	}
	return *o.Currency
}

// GetCurrencyOk returns a tuple with the Currency field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpcomingBill) GetCurrencyOk() (*string, bool) {
	if o == nil || IsNil(o.Currency) {
		return nil, false
	}
	return o.Currency, true
}

// HasCurrency returns a boolean if a field has been set.
func (o *CommonsUpcomingBill) HasCurrency() bool {
	if o != nil && !IsNil(o.Currency) {
		return true
	}

	return false
}

// SetCurrency gets a reference to the given string and assigns it to the Currency field.
func (o *CommonsUpcomingBill) SetCurrency(v string) {
	o.Currency = &v
}

// GetDate returns the Date field value if set, zero value otherwise.
func (o *CommonsUpcomingBill) GetDate() string {
	if o == nil || IsNil(o.Date) {
		var ret string
		return ret
	}
	return *o.Date
}

// GetDateOk returns a tuple with the Date field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsUpcomingBill) GetDateOk() (*string, bool) {
	if o == nil || IsNil(o.Date) {
		return nil, false
	}
	return o.Date, true
}

// HasDate returns a boolean if a field has been set.
func (o *CommonsUpcomingBill) HasDate() bool {
	if o != nil && !IsNil(o.Date) {
		return true
	}

	return false
}

// SetDate gets a reference to the given string and assigns it to the Date field.
func (o *CommonsUpcomingBill) SetDate(v string) {
	o.Date = &v
}

func (o CommonsUpcomingBill) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsUpcomingBill) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Amount) {
		toSerialize["amount"] = o.Amount
	}
	if !IsNil(o.Currency) {
		toSerialize["currency"] = o.Currency
	}
	if !IsNil(o.Date) {
		toSerialize["date"] = o.Date
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsUpcomingBill) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsUpcomingBill := _CommonsUpcomingBill{}

	err = json.Unmarshal(bytes, &varCommonsUpcomingBill)

	if err != nil {
		return err
	}

	*o = CommonsUpcomingBill(varCommonsUpcomingBill)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "amount")
		delete(additionalProperties, "currency")
		delete(additionalProperties, "date")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsUpcomingBill struct {
	value *CommonsUpcomingBill
	isSet bool
}

func (v NullableCommonsUpcomingBill) Get() *CommonsUpcomingBill {
	return v.value
}

func (v *NullableCommonsUpcomingBill) Set(val *CommonsUpcomingBill) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsUpcomingBill) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsUpcomingBill) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsUpcomingBill(val *CommonsUpcomingBill) *NullableCommonsUpcomingBill {
	return &NullableCommonsUpcomingBill{value: val, isSet: true}
}

func (v NullableCommonsUpcomingBill) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsUpcomingBill) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


