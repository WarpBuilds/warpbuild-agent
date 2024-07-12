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

// checks if the CommonsRunnerImagePullSecret type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerImagePullSecret{}

// CommonsRunnerImagePullSecret struct for CommonsRunnerImagePullSecret
type CommonsRunnerImagePullSecret struct {
	Alias string `json:"alias"`
	Aws *CommonsRunnerImagePullSecretAWS `json:"aws,omitempty"`
	CreatedAt string `json:"created_at"`
	DockerCredentials *CommonsRunnerImagePullSecretDockerCredentials `json:"docker_credentials,omitempty"`
	Id string `json:"id"`
	OrganizationId string `json:"organization_id"`
	Type string `json:"type"`
	UpdatedAt string `json:"updated_at"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerImagePullSecret CommonsRunnerImagePullSecret

// NewCommonsRunnerImagePullSecret instantiates a new CommonsRunnerImagePullSecret object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerImagePullSecret(alias string, createdAt string, id string, organizationId string, type_ string, updatedAt string) *CommonsRunnerImagePullSecret {
	this := CommonsRunnerImagePullSecret{}
	this.Alias = alias
	this.CreatedAt = createdAt
	this.Id = id
	this.OrganizationId = organizationId
	this.Type = type_
	this.UpdatedAt = updatedAt
	return &this
}

// NewCommonsRunnerImagePullSecretWithDefaults instantiates a new CommonsRunnerImagePullSecret object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerImagePullSecretWithDefaults() *CommonsRunnerImagePullSecret {
	this := CommonsRunnerImagePullSecret{}
	return &this
}

// GetAlias returns the Alias field value
func (o *CommonsRunnerImagePullSecret) GetAlias() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Alias
}

// GetAliasOk returns a tuple with the Alias field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetAliasOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Alias, true
}

// SetAlias sets field value
func (o *CommonsRunnerImagePullSecret) SetAlias(v string) {
	o.Alias = v
}

// GetAws returns the Aws field value if set, zero value otherwise.
func (o *CommonsRunnerImagePullSecret) GetAws() CommonsRunnerImagePullSecretAWS {
	if o == nil || IsNil(o.Aws) {
		var ret CommonsRunnerImagePullSecretAWS
		return ret
	}
	return *o.Aws
}

// GetAwsOk returns a tuple with the Aws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetAwsOk() (*CommonsRunnerImagePullSecretAWS, bool) {
	if o == nil || IsNil(o.Aws) {
		return nil, false
	}
	return o.Aws, true
}

// HasAws returns a boolean if a field has been set.
func (o *CommonsRunnerImagePullSecret) HasAws() bool {
	if o != nil && !IsNil(o.Aws) {
		return true
	}

	return false
}

// SetAws gets a reference to the given CommonsRunnerImagePullSecretAWS and assigns it to the Aws field.
func (o *CommonsRunnerImagePullSecret) SetAws(v CommonsRunnerImagePullSecretAWS) {
	o.Aws = &v
}

// GetCreatedAt returns the CreatedAt field value
func (o *CommonsRunnerImagePullSecret) GetCreatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetCreatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CreatedAt, true
}

// SetCreatedAt sets field value
func (o *CommonsRunnerImagePullSecret) SetCreatedAt(v string) {
	o.CreatedAt = v
}

// GetDockerCredentials returns the DockerCredentials field value if set, zero value otherwise.
func (o *CommonsRunnerImagePullSecret) GetDockerCredentials() CommonsRunnerImagePullSecretDockerCredentials {
	if o == nil || IsNil(o.DockerCredentials) {
		var ret CommonsRunnerImagePullSecretDockerCredentials
		return ret
	}
	return *o.DockerCredentials
}

// GetDockerCredentialsOk returns a tuple with the DockerCredentials field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetDockerCredentialsOk() (*CommonsRunnerImagePullSecretDockerCredentials, bool) {
	if o == nil || IsNil(o.DockerCredentials) {
		return nil, false
	}
	return o.DockerCredentials, true
}

// HasDockerCredentials returns a boolean if a field has been set.
func (o *CommonsRunnerImagePullSecret) HasDockerCredentials() bool {
	if o != nil && !IsNil(o.DockerCredentials) {
		return true
	}

	return false
}

// SetDockerCredentials gets a reference to the given CommonsRunnerImagePullSecretDockerCredentials and assigns it to the DockerCredentials field.
func (o *CommonsRunnerImagePullSecret) SetDockerCredentials(v CommonsRunnerImagePullSecretDockerCredentials) {
	o.DockerCredentials = &v
}

// GetId returns the Id field value
func (o *CommonsRunnerImagePullSecret) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *CommonsRunnerImagePullSecret) SetId(v string) {
	o.Id = v
}

// GetOrganizationId returns the OrganizationId field value
func (o *CommonsRunnerImagePullSecret) GetOrganizationId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetOrganizationIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.OrganizationId, true
}

// SetOrganizationId sets field value
func (o *CommonsRunnerImagePullSecret) SetOrganizationId(v string) {
	o.OrganizationId = v
}

// GetType returns the Type field value
func (o *CommonsRunnerImagePullSecret) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *CommonsRunnerImagePullSecret) SetType(v string) {
	o.Type = v
}

// GetUpdatedAt returns the UpdatedAt field value
func (o *CommonsRunnerImagePullSecret) GetUpdatedAt() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value
// and a boolean to check if the value has been set.
func (o *CommonsRunnerImagePullSecret) GetUpdatedAtOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.UpdatedAt, true
}

// SetUpdatedAt sets field value
func (o *CommonsRunnerImagePullSecret) SetUpdatedAt(v string) {
	o.UpdatedAt = v
}

func (o CommonsRunnerImagePullSecret) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerImagePullSecret) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["alias"] = o.Alias
	if !IsNil(o.Aws) {
		toSerialize["aws"] = o.Aws
	}
	toSerialize["created_at"] = o.CreatedAt
	if !IsNil(o.DockerCredentials) {
		toSerialize["docker_credentials"] = o.DockerCredentials
	}
	toSerialize["id"] = o.Id
	toSerialize["organization_id"] = o.OrganizationId
	toSerialize["type"] = o.Type
	toSerialize["updated_at"] = o.UpdatedAt

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerImagePullSecret) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerImagePullSecret := _CommonsRunnerImagePullSecret{}

	err = json.Unmarshal(bytes, &varCommonsRunnerImagePullSecret)

	if err != nil {
		return err
	}

	*o = CommonsRunnerImagePullSecret(varCommonsRunnerImagePullSecret)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "alias")
		delete(additionalProperties, "aws")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "docker_credentials")
		delete(additionalProperties, "id")
		delete(additionalProperties, "organization_id")
		delete(additionalProperties, "type")
		delete(additionalProperties, "updated_at")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerImagePullSecret struct {
	value *CommonsRunnerImagePullSecret
	isSet bool
}

func (v NullableCommonsRunnerImagePullSecret) Get() *CommonsRunnerImagePullSecret {
	return v.value
}

func (v *NullableCommonsRunnerImagePullSecret) Set(val *CommonsRunnerImagePullSecret) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerImagePullSecret) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerImagePullSecret) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerImagePullSecret(val *CommonsRunnerImagePullSecret) *NullableCommonsRunnerImagePullSecret {
	return &NullableCommonsRunnerImagePullSecret{value: val, isSet: true}
}

func (v NullableCommonsRunnerImagePullSecret) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerImagePullSecret) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

