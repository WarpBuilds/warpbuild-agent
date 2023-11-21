# CommonsRunnerInstance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Cluster** | Pointer to **string** |  | [optional] 
**Configuration** | Pointer to [**CommonsRunnerConfiguration**](CommonsRunnerConfiguration.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LastJobProcessedId** | Pointer to **string** |  | [optional] 
**LastPolled** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**RunnerFor** | Pointer to **string** |  | [optional] 
**RunnerSetId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 
**VcsIntegrationId** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsRunnerInstance

`func NewCommonsRunnerInstance() *CommonsRunnerInstance`

NewCommonsRunnerInstance instantiates a new CommonsRunnerInstance object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerInstanceWithDefaults

`func NewCommonsRunnerInstanceWithDefaults() *CommonsRunnerInstance`

NewCommonsRunnerInstanceWithDefaults instantiates a new CommonsRunnerInstance object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCluster

`func (o *CommonsRunnerInstance) GetCluster() string`

GetCluster returns the Cluster field if non-nil, zero value otherwise.

### GetClusterOk

`func (o *CommonsRunnerInstance) GetClusterOk() (*string, bool)`

GetClusterOk returns a tuple with the Cluster field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCluster

`func (o *CommonsRunnerInstance) SetCluster(v string)`

SetCluster sets Cluster field to given value.

### HasCluster

`func (o *CommonsRunnerInstance) HasCluster() bool`

HasCluster returns a boolean if a field has been set.

### GetConfiguration

`func (o *CommonsRunnerInstance) GetConfiguration() CommonsRunnerConfiguration`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *CommonsRunnerInstance) GetConfigurationOk() (*CommonsRunnerConfiguration, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *CommonsRunnerInstance) SetConfiguration(v CommonsRunnerConfiguration)`

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *CommonsRunnerInstance) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsRunnerInstance) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsRunnerInstance) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsRunnerInstance) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CommonsRunnerInstance) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetCreatedBy

`func (o *CommonsRunnerInstance) GetCreatedBy() string`

GetCreatedBy returns the CreatedBy field if non-nil, zero value otherwise.

### GetCreatedByOk

`func (o *CommonsRunnerInstance) GetCreatedByOk() (*string, bool)`

GetCreatedByOk returns a tuple with the CreatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedBy

`func (o *CommonsRunnerInstance) SetCreatedBy(v string)`

SetCreatedBy sets CreatedBy field to given value.

### HasCreatedBy

`func (o *CommonsRunnerInstance) HasCreatedBy() bool`

HasCreatedBy returns a boolean if a field has been set.

### GetId

`func (o *CommonsRunnerInstance) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsRunnerInstance) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsRunnerInstance) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsRunnerInstance) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastJobProcessedId

`func (o *CommonsRunnerInstance) GetLastJobProcessedId() string`

GetLastJobProcessedId returns the LastJobProcessedId field if non-nil, zero value otherwise.

### GetLastJobProcessedIdOk

`func (o *CommonsRunnerInstance) GetLastJobProcessedIdOk() (*string, bool)`

GetLastJobProcessedIdOk returns a tuple with the LastJobProcessedId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastJobProcessedId

`func (o *CommonsRunnerInstance) SetLastJobProcessedId(v string)`

SetLastJobProcessedId sets LastJobProcessedId field to given value.

### HasLastJobProcessedId

`func (o *CommonsRunnerInstance) HasLastJobProcessedId() bool`

HasLastJobProcessedId returns a boolean if a field has been set.

### GetLastPolled

`func (o *CommonsRunnerInstance) GetLastPolled() string`

GetLastPolled returns the LastPolled field if non-nil, zero value otherwise.

### GetLastPolledOk

`func (o *CommonsRunnerInstance) GetLastPolledOk() (*string, bool)`

GetLastPolledOk returns a tuple with the LastPolled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPolled

`func (o *CommonsRunnerInstance) SetLastPolled(v string)`

SetLastPolled sets LastPolled field to given value.

### HasLastPolled

`func (o *CommonsRunnerInstance) HasLastPolled() bool`

HasLastPolled returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsRunnerInstance) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsRunnerInstance) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsRunnerInstance) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsRunnerInstance) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetRunnerFor

`func (o *CommonsRunnerInstance) GetRunnerFor() string`

GetRunnerFor returns the RunnerFor field if non-nil, zero value otherwise.

### GetRunnerForOk

`func (o *CommonsRunnerInstance) GetRunnerForOk() (*string, bool)`

GetRunnerForOk returns a tuple with the RunnerFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerFor

`func (o *CommonsRunnerInstance) SetRunnerFor(v string)`

SetRunnerFor sets RunnerFor field to given value.

### HasRunnerFor

`func (o *CommonsRunnerInstance) HasRunnerFor() bool`

HasRunnerFor returns a boolean if a field has been set.

### GetRunnerSetId

`func (o *CommonsRunnerInstance) GetRunnerSetId() string`

GetRunnerSetId returns the RunnerSetId field if non-nil, zero value otherwise.

### GetRunnerSetIdOk

`func (o *CommonsRunnerInstance) GetRunnerSetIdOk() (*string, bool)`

GetRunnerSetIdOk returns a tuple with the RunnerSetId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerSetId

`func (o *CommonsRunnerInstance) SetRunnerSetId(v string)`

SetRunnerSetId sets RunnerSetId field to given value.

### HasRunnerSetId

`func (o *CommonsRunnerInstance) HasRunnerSetId() bool`

HasRunnerSetId returns a boolean if a field has been set.

### GetStatus

`func (o *CommonsRunnerInstance) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CommonsRunnerInstance) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CommonsRunnerInstance) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CommonsRunnerInstance) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *CommonsRunnerInstance) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsRunnerInstance) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsRunnerInstance) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *CommonsRunnerInstance) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.

### GetVcsIntegrationId

`func (o *CommonsRunnerInstance) GetVcsIntegrationId() string`

GetVcsIntegrationId returns the VcsIntegrationId field if non-nil, zero value otherwise.

### GetVcsIntegrationIdOk

`func (o *CommonsRunnerInstance) GetVcsIntegrationIdOk() (*string, bool)`

GetVcsIntegrationIdOk returns a tuple with the VcsIntegrationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsIntegrationId

`func (o *CommonsRunnerInstance) SetVcsIntegrationId(v string)`

SetVcsIntegrationId sets VcsIntegrationId field to given value.

### HasVcsIntegrationId

`func (o *CommonsRunnerInstance) HasVcsIntegrationId() bool`

HasVcsIntegrationId returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


