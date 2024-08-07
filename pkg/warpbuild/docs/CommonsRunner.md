# CommonsRunner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Active** | Pointer to **bool** |  | [optional] 
**Configuration** | Pointer to [**CommonsRunnerSetConfiguration**](CommonsRunnerSetConfiguration.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**ProviderId** | Pointer to **string** |  | [optional] 
**StockRunnerId** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**VcsIntegrationId** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsRunner

`func NewCommonsRunner() *CommonsRunner`

NewCommonsRunner instantiates a new CommonsRunner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerWithDefaults

`func NewCommonsRunnerWithDefaults() *CommonsRunner`

NewCommonsRunnerWithDefaults instantiates a new CommonsRunner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetActive

`func (o *CommonsRunner) GetActive() bool`

GetActive returns the Active field if non-nil, zero value otherwise.

### GetActiveOk

`func (o *CommonsRunner) GetActiveOk() (*bool, bool)`

GetActiveOk returns a tuple with the Active field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActive

`func (o *CommonsRunner) SetActive(v bool)`

SetActive sets Active field to given value.

### HasActive

`func (o *CommonsRunner) HasActive() bool`

HasActive returns a boolean if a field has been set.

### GetConfiguration

`func (o *CommonsRunner) GetConfiguration() CommonsRunnerSetConfiguration`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *CommonsRunner) GetConfigurationOk() (*CommonsRunnerSetConfiguration, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *CommonsRunner) SetConfiguration(v CommonsRunnerSetConfiguration)`

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *CommonsRunner) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsRunner) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsRunner) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsRunner) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CommonsRunner) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetId

`func (o *CommonsRunner) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsRunner) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsRunner) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsRunner) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLabels

`func (o *CommonsRunner) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CommonsRunner) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CommonsRunner) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CommonsRunner) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

### GetName

`func (o *CommonsRunner) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommonsRunner) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommonsRunner) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommonsRunner) HasName() bool`

HasName returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsRunner) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsRunner) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsRunner) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsRunner) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetProviderId

`func (o *CommonsRunner) GetProviderId() string`

GetProviderId returns the ProviderId field if non-nil, zero value otherwise.

### GetProviderIdOk

`func (o *CommonsRunner) GetProviderIdOk() (*string, bool)`

GetProviderIdOk returns a tuple with the ProviderId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderId

`func (o *CommonsRunner) SetProviderId(v string)`

SetProviderId sets ProviderId field to given value.

### HasProviderId

`func (o *CommonsRunner) HasProviderId() bool`

HasProviderId returns a boolean if a field has been set.

### GetStockRunnerId

`func (o *CommonsRunner) GetStockRunnerId() string`

GetStockRunnerId returns the StockRunnerId field if non-nil, zero value otherwise.

### GetStockRunnerIdOk

`func (o *CommonsRunner) GetStockRunnerIdOk() (*string, bool)`

GetStockRunnerIdOk returns a tuple with the StockRunnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStockRunnerId

`func (o *CommonsRunner) SetStockRunnerId(v string)`

SetStockRunnerId sets StockRunnerId field to given value.

### HasStockRunnerId

`func (o *CommonsRunner) HasStockRunnerId() bool`

HasStockRunnerId returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *CommonsRunner) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsRunner) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsRunner) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *CommonsRunner) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetVcsIntegrationId

`func (o *CommonsRunner) GetVcsIntegrationId() string`

GetVcsIntegrationId returns the VcsIntegrationId field if non-nil, zero value otherwise.

### GetVcsIntegrationIdOk

`func (o *CommonsRunner) GetVcsIntegrationIdOk() (*string, bool)`

GetVcsIntegrationIdOk returns a tuple with the VcsIntegrationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsIntegrationId

`func (o *CommonsRunner) SetVcsIntegrationId(v string)`

SetVcsIntegrationId sets VcsIntegrationId field to given value.

### HasVcsIntegrationId

`func (o *CommonsRunner) HasVcsIntegrationId() bool`

HasVcsIntegrationId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


