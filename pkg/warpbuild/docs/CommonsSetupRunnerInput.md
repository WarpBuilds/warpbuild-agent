# CommonsSetupRunnerInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** |  | [optional] 
<<<<<<< HEAD
**Configuration** | Pointer to [**CommonsRunnerConfiguration**](CommonsRunnerConfiguration.md) |  | [optional] 
**Labels** | Pointer to **map[string]interface{}** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
=======
**Configuration** | Pointer to [**CommonsRunnerSetConfiguration**](CommonsRunnerSetConfiguration.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**StockRunnerId** | Pointer to **string** |  | [optional] 
>>>>>>> prajjwal-warp-323
**VcsIntegrationId** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsSetupRunnerInput

`func NewCommonsSetupRunnerInput() *CommonsSetupRunnerInput`

NewCommonsSetupRunnerInput instantiates a new CommonsSetupRunnerInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsSetupRunnerInputWithDefaults

`func NewCommonsSetupRunnerInputWithDefaults() *CommonsSetupRunnerInput`

NewCommonsSetupRunnerInputWithDefaults instantiates a new CommonsSetupRunnerInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *CommonsSetupRunnerInput) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *CommonsSetupRunnerInput) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *CommonsSetupRunnerInput) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *CommonsSetupRunnerInput) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetConfiguration

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) GetConfiguration() CommonsRunnerConfiguration`
=======
`func (o *CommonsSetupRunnerInput) GetConfiguration() CommonsRunnerSetConfiguration`
>>>>>>> prajjwal-warp-323

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) GetConfigurationOk() (*CommonsRunnerConfiguration, bool)`
=======
`func (o *CommonsSetupRunnerInput) GetConfigurationOk() (*CommonsRunnerSetConfiguration, bool)`
>>>>>>> prajjwal-warp-323

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) SetConfiguration(v CommonsRunnerConfiguration)`
=======
`func (o *CommonsSetupRunnerInput) SetConfiguration(v CommonsRunnerSetConfiguration)`
>>>>>>> prajjwal-warp-323

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *CommonsSetupRunnerInput) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.

### GetLabels

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) GetLabels() map[string]interface{}`
=======
`func (o *CommonsSetupRunnerInput) GetLabels() []string`
>>>>>>> prajjwal-warp-323

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) GetLabelsOk() (*map[string]interface{}, bool)`
=======
`func (o *CommonsSetupRunnerInput) GetLabelsOk() (*[]string, bool)`
>>>>>>> prajjwal-warp-323

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

<<<<<<< HEAD
`func (o *CommonsSetupRunnerInput) SetLabels(v map[string]interface{})`
=======
`func (o *CommonsSetupRunnerInput) SetLabels(v []string)`
>>>>>>> prajjwal-warp-323

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CommonsSetupRunnerInput) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetName

`func (o *CommonsSetupRunnerInput) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommonsSetupRunnerInput) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommonsSetupRunnerInput) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommonsSetupRunnerInput) HasName() bool`

HasName returns a boolean if a field has been set.

<<<<<<< HEAD
=======
### GetStockRunnerId

`func (o *CommonsSetupRunnerInput) GetStockRunnerId() string`

GetStockRunnerId returns the StockRunnerId field if non-nil, zero value otherwise.

### GetStockRunnerIdOk

`func (o *CommonsSetupRunnerInput) GetStockRunnerIdOk() (*string, bool)`

GetStockRunnerIdOk returns a tuple with the StockRunnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStockRunnerId

`func (o *CommonsSetupRunnerInput) SetStockRunnerId(v string)`

SetStockRunnerId sets StockRunnerId field to given value.

### HasStockRunnerId

`func (o *CommonsSetupRunnerInput) HasStockRunnerId() bool`

HasStockRunnerId returns a boolean if a field has been set.

>>>>>>> prajjwal-warp-323
### GetVcsIntegrationId

`func (o *CommonsSetupRunnerInput) GetVcsIntegrationId() string`

GetVcsIntegrationId returns the VcsIntegrationId field if non-nil, zero value otherwise.

### GetVcsIntegrationIdOk

`func (o *CommonsSetupRunnerInput) GetVcsIntegrationIdOk() (*string, bool)`

GetVcsIntegrationIdOk returns a tuple with the VcsIntegrationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsIntegrationId

`func (o *CommonsSetupRunnerInput) SetVcsIntegrationId(v string)`

SetVcsIntegrationId sets VcsIntegrationId field to given value.

### HasVcsIntegrationId

`func (o *CommonsSetupRunnerInput) HasVcsIntegrationId() bool`

HasVcsIntegrationId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


