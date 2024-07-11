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

// checks if the TokenRefreshRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TokenRefreshRequest{}

// TokenRefreshRequest struct for TokenRefreshRequest
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	AdditionalProperties map[string]interface{}
}

type _TokenRefreshRequest TokenRefreshRequest

// NewTokenRefreshRequest instantiates a new TokenRefreshRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTokenRefreshRequest(refreshToken string) *TokenRefreshRequest {
	this := TokenRefreshRequest{}
	this.RefreshToken = refreshToken
	return &this
}

// NewTokenRefreshRequestWithDefaults instantiates a new TokenRefreshRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTokenRefreshRequestWithDefaults() *TokenRefreshRequest {
	this := TokenRefreshRequest{}
	return &this
}

// GetRefreshToken returns the RefreshToken field value
func (o *TokenRefreshRequest) GetRefreshToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RefreshToken
}

// GetRefreshTokenOk returns a tuple with the RefreshToken field value
// and a boolean to check if the value has been set.
func (o *TokenRefreshRequest) GetRefreshTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RefreshToken, true
}

// SetRefreshToken sets field value
func (o *TokenRefreshRequest) SetRefreshToken(v string) {
	o.RefreshToken = v
}

func (o TokenRefreshRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TokenRefreshRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["refresh_token"] = o.RefreshToken

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TokenRefreshRequest) UnmarshalJSON(bytes []byte) (err error) {
	varTokenRefreshRequest := _TokenRefreshRequest{}

	err = json.Unmarshal(bytes, &varTokenRefreshRequest)

	if err != nil {
		return err
	}

	*o = TokenRefreshRequest(varTokenRefreshRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "refresh_token")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTokenRefreshRequest struct {
	value *TokenRefreshRequest
	isSet bool
}

func (v NullableTokenRefreshRequest) Get() *TokenRefreshRequest {
	return v.value
}

func (v *NullableTokenRefreshRequest) Set(val *TokenRefreshRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableTokenRefreshRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableTokenRefreshRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTokenRefreshRequest(val *TokenRefreshRequest) *NullableTokenRefreshRequest {
	return &NullableTokenRefreshRequest{value: val, isSet: true}
}

func (v NullableTokenRefreshRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTokenRefreshRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


