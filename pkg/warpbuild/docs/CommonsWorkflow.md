# CommonsWorkflow

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ExternalId** | Pointer to **string** |  | [optional] 
**ExternalRepoEntity** | Pointer to **string** |  | [optional] 
**ExternalRepoId** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**IsRepoPublic** | Pointer to **bool** |  | [optional] 
**JobsRunnerInfo** | Pointer to [**[]CommonsJobRunnerInfo**](CommonsJobRunnerInfo.md) |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**ProviderIntegrationId** | Pointer to **string** |  | [optional] 
**RunnerInfo** | Pointer to [**CommonsRunnerInfo**](CommonsRunnerInfo.md) |  | [optional] 
**Stats** | Pointer to [**CommonsWorkflowStats**](CommonsWorkflowStats.md) |  | [optional] 
**Url** | Pointer to **string** |  | [optional] 
**WarpPrId** | Pointer to **string** |  | [optional] 
**WarpPrLink** | Pointer to **string** |  | [optional] 
**WarpPrRunner** | Pointer to **string** |  | [optional] 
**WarpStatus** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsWorkflow

`func NewCommonsWorkflow() *CommonsWorkflow`

NewCommonsWorkflow instantiates a new CommonsWorkflow object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsWorkflowWithDefaults

`func NewCommonsWorkflowWithDefaults() *CommonsWorkflow`

NewCommonsWorkflowWithDefaults instantiates a new CommonsWorkflow object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetExternalId

`func (o *CommonsWorkflow) GetExternalId() string`

GetExternalId returns the ExternalId field if non-nil, zero value otherwise.

### GetExternalIdOk

`func (o *CommonsWorkflow) GetExternalIdOk() (*string, bool)`

GetExternalIdOk returns a tuple with the ExternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalId

`func (o *CommonsWorkflow) SetExternalId(v string)`

SetExternalId sets ExternalId field to given value.

### HasExternalId

`func (o *CommonsWorkflow) HasExternalId() bool`

HasExternalId returns a boolean if a field has been set.

### GetExternalRepoEntity

`func (o *CommonsWorkflow) GetExternalRepoEntity() string`

GetExternalRepoEntity returns the ExternalRepoEntity field if non-nil, zero value otherwise.

### GetExternalRepoEntityOk

`func (o *CommonsWorkflow) GetExternalRepoEntityOk() (*string, bool)`

GetExternalRepoEntityOk returns a tuple with the ExternalRepoEntity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalRepoEntity

`func (o *CommonsWorkflow) SetExternalRepoEntity(v string)`

SetExternalRepoEntity sets ExternalRepoEntity field to given value.

### HasExternalRepoEntity

`func (o *CommonsWorkflow) HasExternalRepoEntity() bool`

HasExternalRepoEntity returns a boolean if a field has been set.

### GetExternalRepoId

`func (o *CommonsWorkflow) GetExternalRepoId() string`

GetExternalRepoId returns the ExternalRepoId field if non-nil, zero value otherwise.

### GetExternalRepoIdOk

`func (o *CommonsWorkflow) GetExternalRepoIdOk() (*string, bool)`

GetExternalRepoIdOk returns a tuple with the ExternalRepoId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalRepoId

`func (o *CommonsWorkflow) SetExternalRepoId(v string)`

SetExternalRepoId sets ExternalRepoId field to given value.

### HasExternalRepoId

`func (o *CommonsWorkflow) HasExternalRepoId() bool`

HasExternalRepoId returns a boolean if a field has been set.

### GetId

`func (o *CommonsWorkflow) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsWorkflow) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsWorkflow) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsWorkflow) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIsRepoPublic

`func (o *CommonsWorkflow) GetIsRepoPublic() bool`

GetIsRepoPublic returns the IsRepoPublic field if non-nil, zero value otherwise.

### GetIsRepoPublicOk

`func (o *CommonsWorkflow) GetIsRepoPublicOk() (*bool, bool)`

GetIsRepoPublicOk returns a tuple with the IsRepoPublic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsRepoPublic

`func (o *CommonsWorkflow) SetIsRepoPublic(v bool)`

SetIsRepoPublic sets IsRepoPublic field to given value.

### HasIsRepoPublic

`func (o *CommonsWorkflow) HasIsRepoPublic() bool`

HasIsRepoPublic returns a boolean if a field has been set.

### GetJobsRunnerInfo

`func (o *CommonsWorkflow) GetJobsRunnerInfo() []CommonsJobRunnerInfo`

GetJobsRunnerInfo returns the JobsRunnerInfo field if non-nil, zero value otherwise.

### GetJobsRunnerInfoOk

`func (o *CommonsWorkflow) GetJobsRunnerInfoOk() (*[]CommonsJobRunnerInfo, bool)`

GetJobsRunnerInfoOk returns a tuple with the JobsRunnerInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJobsRunnerInfo

`func (o *CommonsWorkflow) SetJobsRunnerInfo(v []CommonsJobRunnerInfo)`

SetJobsRunnerInfo sets JobsRunnerInfo field to given value.

### HasJobsRunnerInfo

`func (o *CommonsWorkflow) HasJobsRunnerInfo() bool`

HasJobsRunnerInfo returns a boolean if a field has been set.

### GetName

`func (o *CommonsWorkflow) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommonsWorkflow) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommonsWorkflow) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommonsWorkflow) HasName() bool`

HasName returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsWorkflow) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsWorkflow) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsWorkflow) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsWorkflow) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetProvider

`func (o *CommonsWorkflow) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *CommonsWorkflow) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *CommonsWorkflow) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *CommonsWorkflow) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetProviderIntegrationId

`func (o *CommonsWorkflow) GetProviderIntegrationId() string`

GetProviderIntegrationId returns the ProviderIntegrationId field if non-nil, zero value otherwise.

### GetProviderIntegrationIdOk

`func (o *CommonsWorkflow) GetProviderIntegrationIdOk() (*string, bool)`

GetProviderIntegrationIdOk returns a tuple with the ProviderIntegrationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderIntegrationId

`func (o *CommonsWorkflow) SetProviderIntegrationId(v string)`

SetProviderIntegrationId sets ProviderIntegrationId field to given value.

### HasProviderIntegrationId

`func (o *CommonsWorkflow) HasProviderIntegrationId() bool`

HasProviderIntegrationId returns a boolean if a field has been set.

### GetRunnerInfo

`func (o *CommonsWorkflow) GetRunnerInfo() CommonsRunnerInfo`

GetRunnerInfo returns the RunnerInfo field if non-nil, zero value otherwise.

### GetRunnerInfoOk

`func (o *CommonsWorkflow) GetRunnerInfoOk() (*CommonsRunnerInfo, bool)`

GetRunnerInfoOk returns a tuple with the RunnerInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerInfo

`func (o *CommonsWorkflow) SetRunnerInfo(v CommonsRunnerInfo)`

SetRunnerInfo sets RunnerInfo field to given value.

### HasRunnerInfo

`func (o *CommonsWorkflow) HasRunnerInfo() bool`

HasRunnerInfo returns a boolean if a field has been set.

### GetStats

`func (o *CommonsWorkflow) GetStats() CommonsWorkflowStats`

GetStats returns the Stats field if non-nil, zero value otherwise.

### GetStatsOk

`func (o *CommonsWorkflow) GetStatsOk() (*CommonsWorkflowStats, bool)`

GetStatsOk returns a tuple with the Stats field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStats

`func (o *CommonsWorkflow) SetStats(v CommonsWorkflowStats)`

SetStats sets Stats field to given value.

### HasStats

`func (o *CommonsWorkflow) HasStats() bool`

HasStats returns a boolean if a field has been set.

### GetUrl

`func (o *CommonsWorkflow) GetUrl() string`

GetUrl returns the Url field if non-nil, zero value otherwise.

### GetUrlOk

`func (o *CommonsWorkflow) GetUrlOk() (*string, bool)`

GetUrlOk returns a tuple with the Url field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrl

`func (o *CommonsWorkflow) SetUrl(v string)`

SetUrl sets Url field to given value.

### HasUrl

`func (o *CommonsWorkflow) HasUrl() bool`

HasUrl returns a boolean if a field has been set.

### GetWarpPrId

`func (o *CommonsWorkflow) GetWarpPrId() string`

GetWarpPrId returns the WarpPrId field if non-nil, zero value otherwise.

### GetWarpPrIdOk

`func (o *CommonsWorkflow) GetWarpPrIdOk() (*string, bool)`

GetWarpPrIdOk returns a tuple with the WarpPrId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarpPrId

`func (o *CommonsWorkflow) SetWarpPrId(v string)`

SetWarpPrId sets WarpPrId field to given value.

### HasWarpPrId

`func (o *CommonsWorkflow) HasWarpPrId() bool`

HasWarpPrId returns a boolean if a field has been set.

### GetWarpPrLink

`func (o *CommonsWorkflow) GetWarpPrLink() string`

GetWarpPrLink returns the WarpPrLink field if non-nil, zero value otherwise.

### GetWarpPrLinkOk

`func (o *CommonsWorkflow) GetWarpPrLinkOk() (*string, bool)`

GetWarpPrLinkOk returns a tuple with the WarpPrLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarpPrLink

`func (o *CommonsWorkflow) SetWarpPrLink(v string)`

SetWarpPrLink sets WarpPrLink field to given value.

### HasWarpPrLink

`func (o *CommonsWorkflow) HasWarpPrLink() bool`

HasWarpPrLink returns a boolean if a field has been set.

### GetWarpPrRunner

`func (o *CommonsWorkflow) GetWarpPrRunner() string`

GetWarpPrRunner returns the WarpPrRunner field if non-nil, zero value otherwise.

### GetWarpPrRunnerOk

`func (o *CommonsWorkflow) GetWarpPrRunnerOk() (*string, bool)`

GetWarpPrRunnerOk returns a tuple with the WarpPrRunner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarpPrRunner

`func (o *CommonsWorkflow) SetWarpPrRunner(v string)`

SetWarpPrRunner sets WarpPrRunner field to given value.

### HasWarpPrRunner

`func (o *CommonsWorkflow) HasWarpPrRunner() bool`

HasWarpPrRunner returns a boolean if a field has been set.

### GetWarpStatus

`func (o *CommonsWorkflow) GetWarpStatus() string`

GetWarpStatus returns the WarpStatus field if non-nil, zero value otherwise.

### GetWarpStatusOk

`func (o *CommonsWorkflow) GetWarpStatusOk() (*string, bool)`

GetWarpStatusOk returns a tuple with the WarpStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarpStatus

`func (o *CommonsWorkflow) SetWarpStatus(v string)`

SetWarpStatus sets WarpStatus field to given value.

### HasWarpStatus

`func (o *CommonsWorkflow) HasWarpStatus() bool`

HasWarpStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


