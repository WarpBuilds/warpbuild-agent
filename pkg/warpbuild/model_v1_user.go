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

// checks if the V1User type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1User{}

// V1User struct for V1User
type V1User struct {
	AvatarUrl *string `json:"avatar_url,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	Email *string `json:"email,omitempty"`
	Id *int32 `json:"id,omitempty"`
	Login *string `json:"login,omitempty"`
	Name *string `json:"name,omitempty"`
	SignupProvider *string `json:"signup_provider,omitempty"`
	SignupProviderUserId *string `json:"signup_provider_user_id,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1User V1User

// NewV1User instantiates a new V1User object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1User() *V1User {
	this := V1User{}
	return &this
}

// NewV1UserWithDefaults instantiates a new V1User object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1UserWithDefaults() *V1User {
	this := V1User{}
	return &this
}

// GetAvatarUrl returns the AvatarUrl field value if set, zero value otherwise.
func (o *V1User) GetAvatarUrl() string {
	if o == nil || IsNil(o.AvatarUrl) {
		var ret string
		return ret
	}
	return *o.AvatarUrl
}

// GetAvatarUrlOk returns a tuple with the AvatarUrl field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetAvatarUrlOk() (*string, bool) {
	if o == nil || IsNil(o.AvatarUrl) {
		return nil, false
	}
	return o.AvatarUrl, true
}

// HasAvatarUrl returns a boolean if a field has been set.
func (o *V1User) HasAvatarUrl() bool {
	if o != nil && !IsNil(o.AvatarUrl) {
		return true
	}

	return false
}

// SetAvatarUrl gets a reference to the given string and assigns it to the AvatarUrl field.
func (o *V1User) SetAvatarUrl(v string) {
	o.AvatarUrl = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *V1User) GetCreatedAt() string {
	if o == nil || IsNil(o.CreatedAt) {
		var ret string
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetCreatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *V1User) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given string and assigns it to the CreatedAt field.
func (o *V1User) SetCreatedAt(v string) {
	o.CreatedAt = &v
}

// GetEmail returns the Email field value if set, zero value otherwise.
func (o *V1User) GetEmail() string {
	if o == nil || IsNil(o.Email) {
		var ret string
		return ret
	}
	return *o.Email
}

// GetEmailOk returns a tuple with the Email field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetEmailOk() (*string, bool) {
	if o == nil || IsNil(o.Email) {
		return nil, false
	}
	return o.Email, true
}

// HasEmail returns a boolean if a field has been set.
func (o *V1User) HasEmail() bool {
	if o != nil && !IsNil(o.Email) {
		return true
	}

	return false
}

// SetEmail gets a reference to the given string and assigns it to the Email field.
func (o *V1User) SetEmail(v string) {
	o.Email = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *V1User) GetId() int32 {
	if o == nil || IsNil(o.Id) {
		var ret int32
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetIdOk() (*int32, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *V1User) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given int32 and assigns it to the Id field.
func (o *V1User) SetId(v int32) {
	o.Id = &v
}

// GetLogin returns the Login field value if set, zero value otherwise.
func (o *V1User) GetLogin() string {
	if o == nil || IsNil(o.Login) {
		var ret string
		return ret
	}
	return *o.Login
}

// GetLoginOk returns a tuple with the Login field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetLoginOk() (*string, bool) {
	if o == nil || IsNil(o.Login) {
		return nil, false
	}
	return o.Login, true
}

// HasLogin returns a boolean if a field has been set.
func (o *V1User) HasLogin() bool {
	if o != nil && !IsNil(o.Login) {
		return true
	}

	return false
}

// SetLogin gets a reference to the given string and assigns it to the Login field.
func (o *V1User) SetLogin(v string) {
	o.Login = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1User) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1User) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1User) SetName(v string) {
	o.Name = &v
}

// GetSignupProvider returns the SignupProvider field value if set, zero value otherwise.
func (o *V1User) GetSignupProvider() string {
	if o == nil || IsNil(o.SignupProvider) {
		var ret string
		return ret
	}
	return *o.SignupProvider
}

// GetSignupProviderOk returns a tuple with the SignupProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetSignupProviderOk() (*string, bool) {
	if o == nil || IsNil(o.SignupProvider) {
		return nil, false
	}
	return o.SignupProvider, true
}

// HasSignupProvider returns a boolean if a field has been set.
func (o *V1User) HasSignupProvider() bool {
	if o != nil && !IsNil(o.SignupProvider) {
		return true
	}

	return false
}

// SetSignupProvider gets a reference to the given string and assigns it to the SignupProvider field.
func (o *V1User) SetSignupProvider(v string) {
	o.SignupProvider = &v
}

// GetSignupProviderUserId returns the SignupProviderUserId field value if set, zero value otherwise.
func (o *V1User) GetSignupProviderUserId() string {
	if o == nil || IsNil(o.SignupProviderUserId) {
		var ret string
		return ret
	}
	return *o.SignupProviderUserId
}

// GetSignupProviderUserIdOk returns a tuple with the SignupProviderUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetSignupProviderUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.SignupProviderUserId) {
		return nil, false
	}
	return o.SignupProviderUserId, true
}

// HasSignupProviderUserId returns a boolean if a field has been set.
func (o *V1User) HasSignupProviderUserId() bool {
	if o != nil && !IsNil(o.SignupProviderUserId) {
		return true
	}

	return false
}

// SetSignupProviderUserId gets a reference to the given string and assigns it to the SignupProviderUserId field.
func (o *V1User) SetSignupProviderUserId(v string) {
	o.SignupProviderUserId = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *V1User) GetUpdatedAt() string {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret string
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1User) GetUpdatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *V1User) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given string and assigns it to the UpdatedAt field.
func (o *V1User) SetUpdatedAt(v string) {
	o.UpdatedAt = &v
}

func (o V1User) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1User) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AvatarUrl) {
		toSerialize["avatar_url"] = o.AvatarUrl
	}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.Email) {
		toSerialize["email"] = o.Email
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Login) {
		toSerialize["login"] = o.Login
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.SignupProvider) {
		toSerialize["signup_provider"] = o.SignupProvider
	}
	if !IsNil(o.SignupProviderUserId) {
		toSerialize["signup_provider_user_id"] = o.SignupProviderUserId
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1User) UnmarshalJSON(bytes []byte) (err error) {
	varV1User := _V1User{}

	err = json.Unmarshal(bytes, &varV1User)

	if err != nil {
		return err
	}

	*o = V1User(varV1User)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "avatar_url")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "email")
		delete(additionalProperties, "id")
		delete(additionalProperties, "login")
		delete(additionalProperties, "name")
		delete(additionalProperties, "signup_provider")
		delete(additionalProperties, "signup_provider_user_id")
		delete(additionalProperties, "updated_at")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1User struct {
	value *V1User
	isSet bool
}

func (v NullableV1User) Get() *V1User {
	return v.value
}

func (v *NullableV1User) Set(val *V1User) {
	v.value = val
	v.isSet = true
}

func (v NullableV1User) IsSet() bool {
	return v.isSet
}

func (v *NullableV1User) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1User(val *V1User) *NullableV1User {
	return &NullableV1User{value: val, isSet: true}
}

func (v NullableV1User) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1User) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


