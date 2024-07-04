# CommonsUpdateRunnerInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** |  | [optional] 
**Configuration** | Pointer to [**CommonsRunnerSetConfiguration**](CommonsRunnerSetConfiguration.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsUpdateRunnerInput

`func NewCommonsUpdateRunnerInput() *CommonsUpdateRunnerInput`

NewCommonsUpdateRunnerInput instantiates a new CommonsUpdateRunnerInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsUpdateRunnerInputWithDefaults

`func NewCommonsUpdateRunnerInputWithDefaults() *CommonsUpdateRunnerInput`

NewCommonsUpdateRunnerInputWithDefaults instantiates a new CommonsUpdateRunnerInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *CommonsUpdateRunnerInput) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *CommonsUpdateRunnerInput) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *CommonsUpdateRunnerInput) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *CommonsUpdateRunnerInput) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetConfiguration

`func (o *CommonsUpdateRunnerInput) GetConfiguration() CommonsRunnerSetConfiguration`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *CommonsUpdateRunnerInput) GetConfigurationOk() (*CommonsRunnerSetConfiguration, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *CommonsUpdateRunnerInput) SetConfiguration(v CommonsRunnerSetConfiguration)`

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *CommonsUpdateRunnerInput) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.

### GetLabels

`func (o *CommonsUpdateRunnerInput) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CommonsUpdateRunnerInput) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CommonsUpdateRunnerInput) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CommonsUpdateRunnerInput) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetName

`func (o *CommonsUpdateRunnerInput) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommonsUpdateRunnerInput) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommonsUpdateRunnerInput) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommonsUpdateRunnerInput) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


