# CommonsRunnerInstance

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AllocatedAt** | Pointer to **string** |  | [optional] 
**AllocationFor** | Pointer to **string** |  | [optional] 
**AllocationMutexTimeout** | Pointer to **string** |  | [optional] 
**AllocationRequestedAt** | Pointer to **string** |  | [optional] 
**AllocationRequestedEventAt** | Pointer to **string** |  | [optional] 
**Cluster** | Pointer to **string** |  | [optional] 
**Configuration** | Pointer to [**CommonsRunnerInstanceConfiguration**](CommonsRunnerInstanceConfiguration.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**CreatedBy** | Pointer to **string** |  | [optional] 
**DynamicOverrides** | Pointer to [**CommonsDynamicOverrides**](CommonsDynamicOverrides.md) |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**FirstPolledAt** | Pointer to **string** |  | [optional] 
**Host** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LabelAttributes** | Pointer to [**CommonsLabelAttributes**](CommonsLabelAttributes.md) |  | [optional] 
**Labels** | Pointer to **[]string** |  | [optional] 
**LastJobProcessedId** | Pointer to **string** |  | [optional] 
**LastJobProcessedMeta** | Pointer to [**CommonsLastJobProcessedMeta**](CommonsLastJobProcessedMeta.md) |  | [optional] 
**LastPolledAt** | Pointer to **string** |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**PostResumeFirstPolledAt** | Pointer to **string** |  | [optional] 
**ProviderKindId** | Pointer to **string** |  | [optional] 
**PurgedAt** | Pointer to **string** |  | [optional] 
**PurgedReason** | Pointer to **string** |  | [optional] 
**RunnerFor** | Pointer to **string** |  | [optional] 
**RunnerSetId** | Pointer to **string** |  | [optional] 
**RunningStartedAt** | Pointer to **string** |  | [optional] 
**StackKind** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**SuspendedAt** | Pointer to **string** |  | [optional] 
**SuspendingStartedAt** | Pointer to **string** |  | [optional] 
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

### GetAllocatedAt

`func (o *CommonsRunnerInstance) GetAllocatedAt() string`

GetAllocatedAt returns the AllocatedAt field if non-nil, zero value otherwise.

### GetAllocatedAtOk

`func (o *CommonsRunnerInstance) GetAllocatedAtOk() (*string, bool)`

GetAllocatedAtOk returns a tuple with the AllocatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedAt

`func (o *CommonsRunnerInstance) SetAllocatedAt(v string)`

SetAllocatedAt sets AllocatedAt field to given value.

### HasAllocatedAt

`func (o *CommonsRunnerInstance) HasAllocatedAt() bool`

HasAllocatedAt returns a boolean if a field has been set.

### GetAllocationFor

`func (o *CommonsRunnerInstance) GetAllocationFor() string`

GetAllocationFor returns the AllocationFor field if non-nil, zero value otherwise.

### GetAllocationForOk

`func (o *CommonsRunnerInstance) GetAllocationForOk() (*string, bool)`

GetAllocationForOk returns a tuple with the AllocationFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationFor

`func (o *CommonsRunnerInstance) SetAllocationFor(v string)`

SetAllocationFor sets AllocationFor field to given value.

### HasAllocationFor

`func (o *CommonsRunnerInstance) HasAllocationFor() bool`

HasAllocationFor returns a boolean if a field has been set.

### GetAllocationMutexTimeout

`func (o *CommonsRunnerInstance) GetAllocationMutexTimeout() string`

GetAllocationMutexTimeout returns the AllocationMutexTimeout field if non-nil, zero value otherwise.

### GetAllocationMutexTimeoutOk

`func (o *CommonsRunnerInstance) GetAllocationMutexTimeoutOk() (*string, bool)`

GetAllocationMutexTimeoutOk returns a tuple with the AllocationMutexTimeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationMutexTimeout

`func (o *CommonsRunnerInstance) SetAllocationMutexTimeout(v string)`

SetAllocationMutexTimeout sets AllocationMutexTimeout field to given value.

### HasAllocationMutexTimeout

`func (o *CommonsRunnerInstance) HasAllocationMutexTimeout() bool`

HasAllocationMutexTimeout returns a boolean if a field has been set.

### GetAllocationRequestedAt

`func (o *CommonsRunnerInstance) GetAllocationRequestedAt() string`

GetAllocationRequestedAt returns the AllocationRequestedAt field if non-nil, zero value otherwise.

### GetAllocationRequestedAtOk

`func (o *CommonsRunnerInstance) GetAllocationRequestedAtOk() (*string, bool)`

GetAllocationRequestedAtOk returns a tuple with the AllocationRequestedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationRequestedAt

`func (o *CommonsRunnerInstance) SetAllocationRequestedAt(v string)`

SetAllocationRequestedAt sets AllocationRequestedAt field to given value.

### HasAllocationRequestedAt

`func (o *CommonsRunnerInstance) HasAllocationRequestedAt() bool`

HasAllocationRequestedAt returns a boolean if a field has been set.

### GetAllocationRequestedEventAt

`func (o *CommonsRunnerInstance) GetAllocationRequestedEventAt() string`

GetAllocationRequestedEventAt returns the AllocationRequestedEventAt field if non-nil, zero value otherwise.

### GetAllocationRequestedEventAtOk

`func (o *CommonsRunnerInstance) GetAllocationRequestedEventAtOk() (*string, bool)`

GetAllocationRequestedEventAtOk returns a tuple with the AllocationRequestedEventAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocationRequestedEventAt

`func (o *CommonsRunnerInstance) SetAllocationRequestedEventAt(v string)`

SetAllocationRequestedEventAt sets AllocationRequestedEventAt field to given value.

### HasAllocationRequestedEventAt

`func (o *CommonsRunnerInstance) HasAllocationRequestedEventAt() bool`

HasAllocationRequestedEventAt returns a boolean if a field has been set.

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

`func (o *CommonsRunnerInstance) GetConfiguration() CommonsRunnerInstanceConfiguration`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *CommonsRunnerInstance) GetConfigurationOk() (*CommonsRunnerInstanceConfiguration, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *CommonsRunnerInstance) SetConfiguration(v CommonsRunnerInstanceConfiguration)`

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

### GetDynamicOverrides

`func (o *CommonsRunnerInstance) GetDynamicOverrides() CommonsDynamicOverrides`

GetDynamicOverrides returns the DynamicOverrides field if non-nil, zero value otherwise.

### GetDynamicOverridesOk

`func (o *CommonsRunnerInstance) GetDynamicOverridesOk() (*CommonsDynamicOverrides, bool)`

GetDynamicOverridesOk returns a tuple with the DynamicOverrides field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicOverrides

`func (o *CommonsRunnerInstance) SetDynamicOverrides(v CommonsDynamicOverrides)`

SetDynamicOverrides sets DynamicOverrides field to given value.

### HasDynamicOverrides

`func (o *CommonsRunnerInstance) HasDynamicOverrides() bool`

HasDynamicOverrides returns a boolean if a field has been set.

### GetExternalId

`func (o *CommonsRunnerInstance) GetExternalId() string`

GetExternalId returns the ExternalId field if non-nil, zero value otherwise.

### GetExternalIdOk

`func (o *CommonsRunnerInstance) GetExternalIdOk() (*string, bool)`

GetExternalIdOk returns a tuple with the ExternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalId

`func (o *CommonsRunnerInstance) SetExternalId(v string)`

SetExternalId sets ExternalId field to given value.

### HasExternalId

`func (o *CommonsRunnerInstance) HasExternalId() bool`

HasExternalId returns a boolean if a field has been set.

### GetFirstPolledAt

`func (o *CommonsRunnerInstance) GetFirstPolledAt() string`

GetFirstPolledAt returns the FirstPolledAt field if non-nil, zero value otherwise.

### GetFirstPolledAtOk

`func (o *CommonsRunnerInstance) GetFirstPolledAtOk() (*string, bool)`

GetFirstPolledAtOk returns a tuple with the FirstPolledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstPolledAt

`func (o *CommonsRunnerInstance) SetFirstPolledAt(v string)`

SetFirstPolledAt sets FirstPolledAt field to given value.

### HasFirstPolledAt

`func (o *CommonsRunnerInstance) HasFirstPolledAt() bool`

HasFirstPolledAt returns a boolean if a field has been set.

### GetHost

`func (o *CommonsRunnerInstance) GetHost() string`

GetHost returns the Host field if non-nil, zero value otherwise.

### GetHostOk

`func (o *CommonsRunnerInstance) GetHostOk() (*string, bool)`

GetHostOk returns a tuple with the Host field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHost

`func (o *CommonsRunnerInstance) SetHost(v string)`

SetHost sets Host field to given value.

### HasHost

`func (o *CommonsRunnerInstance) HasHost() bool`

HasHost returns a boolean if a field has been set.

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

### GetLabelAttributes

`func (o *CommonsRunnerInstance) GetLabelAttributes() CommonsLabelAttributes`

GetLabelAttributes returns the LabelAttributes field if non-nil, zero value otherwise.

### GetLabelAttributesOk

`func (o *CommonsRunnerInstance) GetLabelAttributesOk() (*CommonsLabelAttributes, bool)`

GetLabelAttributesOk returns a tuple with the LabelAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabelAttributes

`func (o *CommonsRunnerInstance) SetLabelAttributes(v CommonsLabelAttributes)`

SetLabelAttributes sets LabelAttributes field to given value.

### HasLabelAttributes

`func (o *CommonsRunnerInstance) HasLabelAttributes() bool`

HasLabelAttributes returns a boolean if a field has been set.

### GetLabels

`func (o *CommonsRunnerInstance) GetLabels() []string`

GetLabels returns the Labels field if non-nil, zero value otherwise.

### GetLabelsOk

`func (o *CommonsRunnerInstance) GetLabelsOk() (*[]string, bool)`

GetLabelsOk returns a tuple with the Labels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabels

`func (o *CommonsRunnerInstance) SetLabels(v []string)`

SetLabels sets Labels field to given value.

### HasLabels

`func (o *CommonsRunnerInstance) HasLabels() bool`

HasLabels returns a boolean if a field has been set.

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

### GetLastJobProcessedMeta

`func (o *CommonsRunnerInstance) GetLastJobProcessedMeta() CommonsLastJobProcessedMeta`

GetLastJobProcessedMeta returns the LastJobProcessedMeta field if non-nil, zero value otherwise.

### GetLastJobProcessedMetaOk

`func (o *CommonsRunnerInstance) GetLastJobProcessedMetaOk() (*CommonsLastJobProcessedMeta, bool)`

GetLastJobProcessedMetaOk returns a tuple with the LastJobProcessedMeta field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastJobProcessedMeta

`func (o *CommonsRunnerInstance) SetLastJobProcessedMeta(v CommonsLastJobProcessedMeta)`

SetLastJobProcessedMeta sets LastJobProcessedMeta field to given value.

### HasLastJobProcessedMeta

`func (o *CommonsRunnerInstance) HasLastJobProcessedMeta() bool`

HasLastJobProcessedMeta returns a boolean if a field has been set.

### GetLastPolledAt

`func (o *CommonsRunnerInstance) GetLastPolledAt() string`

GetLastPolledAt returns the LastPolledAt field if non-nil, zero value otherwise.

### GetLastPolledAtOk

`func (o *CommonsRunnerInstance) GetLastPolledAtOk() (*string, bool)`

GetLastPolledAtOk returns a tuple with the LastPolledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPolledAt

`func (o *CommonsRunnerInstance) SetLastPolledAt(v string)`

SetLastPolledAt sets LastPolledAt field to given value.

### HasLastPolledAt

`func (o *CommonsRunnerInstance) HasLastPolledAt() bool`

HasLastPolledAt returns a boolean if a field has been set.

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

### GetPostResumeFirstPolledAt

`func (o *CommonsRunnerInstance) GetPostResumeFirstPolledAt() string`

GetPostResumeFirstPolledAt returns the PostResumeFirstPolledAt field if non-nil, zero value otherwise.

### GetPostResumeFirstPolledAtOk

`func (o *CommonsRunnerInstance) GetPostResumeFirstPolledAtOk() (*string, bool)`

GetPostResumeFirstPolledAtOk returns a tuple with the PostResumeFirstPolledAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostResumeFirstPolledAt

`func (o *CommonsRunnerInstance) SetPostResumeFirstPolledAt(v string)`

SetPostResumeFirstPolledAt sets PostResumeFirstPolledAt field to given value.

### HasPostResumeFirstPolledAt

`func (o *CommonsRunnerInstance) HasPostResumeFirstPolledAt() bool`

HasPostResumeFirstPolledAt returns a boolean if a field has been set.

### GetProviderKindId

`func (o *CommonsRunnerInstance) GetProviderKindId() string`

GetProviderKindId returns the ProviderKindId field if non-nil, zero value otherwise.

### GetProviderKindIdOk

`func (o *CommonsRunnerInstance) GetProviderKindIdOk() (*string, bool)`

GetProviderKindIdOk returns a tuple with the ProviderKindId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderKindId

`func (o *CommonsRunnerInstance) SetProviderKindId(v string)`

SetProviderKindId sets ProviderKindId field to given value.

### HasProviderKindId

`func (o *CommonsRunnerInstance) HasProviderKindId() bool`

HasProviderKindId returns a boolean if a field has been set.

### GetPurgedAt

`func (o *CommonsRunnerInstance) GetPurgedAt() string`

GetPurgedAt returns the PurgedAt field if non-nil, zero value otherwise.

### GetPurgedAtOk

`func (o *CommonsRunnerInstance) GetPurgedAtOk() (*string, bool)`

GetPurgedAtOk returns a tuple with the PurgedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurgedAt

`func (o *CommonsRunnerInstance) SetPurgedAt(v string)`

SetPurgedAt sets PurgedAt field to given value.

### HasPurgedAt

`func (o *CommonsRunnerInstance) HasPurgedAt() bool`

HasPurgedAt returns a boolean if a field has been set.

### GetPurgedReason

`func (o *CommonsRunnerInstance) GetPurgedReason() string`

GetPurgedReason returns the PurgedReason field if non-nil, zero value otherwise.

### GetPurgedReasonOk

`func (o *CommonsRunnerInstance) GetPurgedReasonOk() (*string, bool)`

GetPurgedReasonOk returns a tuple with the PurgedReason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPurgedReason

`func (o *CommonsRunnerInstance) SetPurgedReason(v string)`

SetPurgedReason sets PurgedReason field to given value.

### HasPurgedReason

`func (o *CommonsRunnerInstance) HasPurgedReason() bool`

HasPurgedReason returns a boolean if a field has been set.

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

### GetRunningStartedAt

`func (o *CommonsRunnerInstance) GetRunningStartedAt() string`

GetRunningStartedAt returns the RunningStartedAt field if non-nil, zero value otherwise.

### GetRunningStartedAtOk

`func (o *CommonsRunnerInstance) GetRunningStartedAtOk() (*string, bool)`

GetRunningStartedAtOk returns a tuple with the RunningStartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunningStartedAt

`func (o *CommonsRunnerInstance) SetRunningStartedAt(v string)`

SetRunningStartedAt sets RunningStartedAt field to given value.

### HasRunningStartedAt

`func (o *CommonsRunnerInstance) HasRunningStartedAt() bool`

HasRunningStartedAt returns a boolean if a field has been set.

### GetStackKind

`func (o *CommonsRunnerInstance) GetStackKind() string`

GetStackKind returns the StackKind field if non-nil, zero value otherwise.

### GetStackKindOk

`func (o *CommonsRunnerInstance) GetStackKindOk() (*string, bool)`

GetStackKindOk returns a tuple with the StackKind field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStackKind

`func (o *CommonsRunnerInstance) SetStackKind(v string)`

SetStackKind sets StackKind field to given value.

### HasStackKind

`func (o *CommonsRunnerInstance) HasStackKind() bool`

HasStackKind returns a boolean if a field has been set.

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

### GetSuspendedAt

`func (o *CommonsRunnerInstance) GetSuspendedAt() string`

GetSuspendedAt returns the SuspendedAt field if non-nil, zero value otherwise.

### GetSuspendedAtOk

`func (o *CommonsRunnerInstance) GetSuspendedAtOk() (*string, bool)`

GetSuspendedAtOk returns a tuple with the SuspendedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuspendedAt

`func (o *CommonsRunnerInstance) SetSuspendedAt(v string)`

SetSuspendedAt sets SuspendedAt field to given value.

### HasSuspendedAt

`func (o *CommonsRunnerInstance) HasSuspendedAt() bool`

HasSuspendedAt returns a boolean if a field has been set.

### GetSuspendingStartedAt

`func (o *CommonsRunnerInstance) GetSuspendingStartedAt() string`

GetSuspendingStartedAt returns the SuspendingStartedAt field if non-nil, zero value otherwise.

### GetSuspendingStartedAtOk

`func (o *CommonsRunnerInstance) GetSuspendingStartedAtOk() (*string, bool)`

GetSuspendingStartedAtOk returns a tuple with the SuspendingStartedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuspendingStartedAt

`func (o *CommonsRunnerInstance) SetSuspendingStartedAt(v string)`

SetSuspendingStartedAt sets SuspendingStartedAt field to given value.

### HasSuspendingStartedAt

`func (o *CommonsRunnerInstance) HasSuspendingStartedAt() bool`

HasSuspendingStartedAt returns a boolean if a field has been set.

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


