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

// checks if the CommonsContainerRunnerImage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsContainerRunnerImage{}

// CommonsContainerRunnerImage struct for CommonsContainerRunnerImage
type CommonsContainerRunnerImage struct {
	Args []string `json:"args,omitempty"`
	Command *string `json:"command,omitempty"`
	Entrypoint *string `json:"entrypoint,omitempty"`
	EnvironmentVariables *map[string]string `json:"environment_variables,omitempty"`
	ImageRepository *string `json:"image_repository,omitempty"`
	ImageTag *string `json:"image_tag,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsContainerRunnerImage CommonsContainerRunnerImage

// NewCommonsContainerRunnerImage instantiates a new CommonsContainerRunnerImage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsContainerRunnerImage() *CommonsContainerRunnerImage {
	this := CommonsContainerRunnerImage{}
	return &this
}

// NewCommonsContainerRunnerImageWithDefaults instantiates a new CommonsContainerRunnerImage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsContainerRunnerImageWithDefaults() *CommonsContainerRunnerImage {
	this := CommonsContainerRunnerImage{}
	return &this
}

// GetArgs returns the Args field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetArgs() []string {
	if o == nil || IsNil(o.Args) {
		var ret []string
		return ret
	}
	return o.Args
}

// GetArgsOk returns a tuple with the Args field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetArgsOk() ([]string, bool) {
	if o == nil || IsNil(o.Args) {
		return nil, false
	}
	return o.Args, true
}

// HasArgs returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasArgs() bool {
	if o != nil && !IsNil(o.Args) {
		return true
	}

	return false
}

// SetArgs gets a reference to the given []string and assigns it to the Args field.
func (o *CommonsContainerRunnerImage) SetArgs(v []string) {
	o.Args = v
}

// GetCommand returns the Command field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetCommand() string {
	if o == nil || IsNil(o.Command) {
		var ret string
		return ret
	}
	return *o.Command
}

// GetCommandOk returns a tuple with the Command field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetCommandOk() (*string, bool) {
	if o == nil || IsNil(o.Command) {
		return nil, false
	}
	return o.Command, true
}

// HasCommand returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasCommand() bool {
	if o != nil && !IsNil(o.Command) {
		return true
	}

	return false
}

// SetCommand gets a reference to the given string and assigns it to the Command field.
func (o *CommonsContainerRunnerImage) SetCommand(v string) {
	o.Command = &v
}

// GetEntrypoint returns the Entrypoint field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetEntrypoint() string {
	if o == nil || IsNil(o.Entrypoint) {
		var ret string
		return ret
	}
	return *o.Entrypoint
}

// GetEntrypointOk returns a tuple with the Entrypoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetEntrypointOk() (*string, bool) {
	if o == nil || IsNil(o.Entrypoint) {
		return nil, false
	}
	return o.Entrypoint, true
}

// HasEntrypoint returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasEntrypoint() bool {
	if o != nil && !IsNil(o.Entrypoint) {
		return true
	}

	return false
}

// SetEntrypoint gets a reference to the given string and assigns it to the Entrypoint field.
func (o *CommonsContainerRunnerImage) SetEntrypoint(v string) {
	o.Entrypoint = &v
}

// GetEnvironmentVariables returns the EnvironmentVariables field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetEnvironmentVariables() map[string]string {
	if o == nil || IsNil(o.EnvironmentVariables) {
		var ret map[string]string
		return ret
	}
	return *o.EnvironmentVariables
}

// GetEnvironmentVariablesOk returns a tuple with the EnvironmentVariables field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetEnvironmentVariablesOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.EnvironmentVariables) {
		return nil, false
	}
	return o.EnvironmentVariables, true
}

// HasEnvironmentVariables returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasEnvironmentVariables() bool {
	if o != nil && !IsNil(o.EnvironmentVariables) {
		return true
	}

	return false
}

// SetEnvironmentVariables gets a reference to the given map[string]string and assigns it to the EnvironmentVariables field.
func (o *CommonsContainerRunnerImage) SetEnvironmentVariables(v map[string]string) {
	o.EnvironmentVariables = &v
}

// GetImageRepository returns the ImageRepository field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetImageRepository() string {
	if o == nil || IsNil(o.ImageRepository) {
		var ret string
		return ret
	}
	return *o.ImageRepository
}

// GetImageRepositoryOk returns a tuple with the ImageRepository field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetImageRepositoryOk() (*string, bool) {
	if o == nil || IsNil(o.ImageRepository) {
		return nil, false
	}
	return o.ImageRepository, true
}

// HasImageRepository returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasImageRepository() bool {
	if o != nil && !IsNil(o.ImageRepository) {
		return true
	}

	return false
}

// SetImageRepository gets a reference to the given string and assigns it to the ImageRepository field.
func (o *CommonsContainerRunnerImage) SetImageRepository(v string) {
	o.ImageRepository = &v
}

// GetImageTag returns the ImageTag field value if set, zero value otherwise.
func (o *CommonsContainerRunnerImage) GetImageTag() string {
	if o == nil || IsNil(o.ImageTag) {
		var ret string
		return ret
	}
	return *o.ImageTag
}

// GetImageTagOk returns a tuple with the ImageTag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsContainerRunnerImage) GetImageTagOk() (*string, bool) {
	if o == nil || IsNil(o.ImageTag) {
		return nil, false
	}
	return o.ImageTag, true
}

// HasImageTag returns a boolean if a field has been set.
func (o *CommonsContainerRunnerImage) HasImageTag() bool {
	if o != nil && !IsNil(o.ImageTag) {
		return true
	}

	return false
}

// SetImageTag gets a reference to the given string and assigns it to the ImageTag field.
func (o *CommonsContainerRunnerImage) SetImageTag(v string) {
	o.ImageTag = &v
}

func (o CommonsContainerRunnerImage) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsContainerRunnerImage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Args) {
		toSerialize["args"] = o.Args
	}
	if !IsNil(o.Command) {
		toSerialize["command"] = o.Command
	}
	if !IsNil(o.Entrypoint) {
		toSerialize["entrypoint"] = o.Entrypoint
	}
	if !IsNil(o.EnvironmentVariables) {
		toSerialize["environment_variables"] = o.EnvironmentVariables
	}
	if !IsNil(o.ImageRepository) {
		toSerialize["image_repository"] = o.ImageRepository
	}
	if !IsNil(o.ImageTag) {
		toSerialize["image_tag"] = o.ImageTag
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsContainerRunnerImage) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsContainerRunnerImage := _CommonsContainerRunnerImage{}

	err = json.Unmarshal(bytes, &varCommonsContainerRunnerImage)

	if err != nil {
		return err
	}

	*o = CommonsContainerRunnerImage(varCommonsContainerRunnerImage)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "args")
		delete(additionalProperties, "command")
		delete(additionalProperties, "entrypoint")
		delete(additionalProperties, "environment_variables")
		delete(additionalProperties, "image_repository")
		delete(additionalProperties, "image_tag")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsContainerRunnerImage struct {
	value *CommonsContainerRunnerImage
	isSet bool
}

func (v NullableCommonsContainerRunnerImage) Get() *CommonsContainerRunnerImage {
	return v.value
}

func (v *NullableCommonsContainerRunnerImage) Set(val *CommonsContainerRunnerImage) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsContainerRunnerImage) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsContainerRunnerImage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsContainerRunnerImage(val *CommonsContainerRunnerImage) *NullableCommonsContainerRunnerImage {
	return &NullableCommonsContainerRunnerImage{value: val, isSet: true}
}

func (v NullableCommonsContainerRunnerImage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsContainerRunnerImage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


