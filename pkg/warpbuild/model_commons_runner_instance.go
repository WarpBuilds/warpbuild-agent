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

// checks if the CommonsRunnerInstance type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CommonsRunnerInstance{}

// CommonsRunnerInstance struct for CommonsRunnerInstance
type CommonsRunnerInstance struct {
	AllocatedAt *string `json:"allocated_at,omitempty"`
	AllocationFor *string `json:"allocation_for,omitempty"`
	AllocationRequestedAt *string `json:"allocation_requested_at,omitempty"`
	AllocationRequestedEventAt *string `json:"allocation_requested_event_at,omitempty"`
	Cluster *string `json:"cluster,omitempty"`
	Configuration *CommonsRunnerInstanceConfiguration `json:"configuration,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	CreatedBy *string `json:"created_by,omitempty"`
	ExternalId *string `json:"external_id,omitempty"`
	FirstPolledAt *string `json:"first_polled_at,omitempty"`
	Host *string `json:"host,omitempty"`
	Id *string `json:"id,omitempty"`
	Labels []string `json:"labels,omitempty"`
	LastJobProcessedId *string `json:"last_job_processed_id,omitempty"`
	LastJobProcessedMeta *CommonsLastJobProcessedMeta `json:"last_job_processed_meta,omitempty"`
	LastPolledAt *string `json:"last_polled_at,omitempty"`
	OrganizationId *string `json:"organization_id,omitempty"`
	ProviderKind *string `json:"provider_kind,omitempty"`
	ProviderKindId *string `json:"provider_kind_id,omitempty"`
	PurgedAt *string `json:"purged_at,omitempty"`
	PurgedReason *string `json:"purged_reason,omitempty"`
	RunnerFor *string `json:"runner_for,omitempty"`
	RunnerSetId *string `json:"runner_set_id,omitempty"`
	RunningStartedAt *string `json:"running_started_at,omitempty"`
	Status *string `json:"status,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
	VcsIntegrationId *string `json:"vcs_integration_id,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _CommonsRunnerInstance CommonsRunnerInstance

// NewCommonsRunnerInstance instantiates a new CommonsRunnerInstance object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonsRunnerInstance() *CommonsRunnerInstance {
	this := CommonsRunnerInstance{}
	return &this
}

// NewCommonsRunnerInstanceWithDefaults instantiates a new CommonsRunnerInstance object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonsRunnerInstanceWithDefaults() *CommonsRunnerInstance {
	this := CommonsRunnerInstance{}
	return &this
}

// GetAllocatedAt returns the AllocatedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetAllocatedAt() string {
	if o == nil || IsNil(o.AllocatedAt) {
		var ret string
		return ret
	}
	return *o.AllocatedAt
}

// GetAllocatedAtOk returns a tuple with the AllocatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetAllocatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.AllocatedAt) {
		return nil, false
	}
	return o.AllocatedAt, true
}

// HasAllocatedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasAllocatedAt() bool {
	if o != nil && !IsNil(o.AllocatedAt) {
		return true
	}

	return false
}

// SetAllocatedAt gets a reference to the given string and assigns it to the AllocatedAt field.
func (o *CommonsRunnerInstance) SetAllocatedAt(v string) {
	o.AllocatedAt = &v
}

// GetAllocationFor returns the AllocationFor field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetAllocationFor() string {
	if o == nil || IsNil(o.AllocationFor) {
		var ret string
		return ret
	}
	return *o.AllocationFor
}

// GetAllocationForOk returns a tuple with the AllocationFor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetAllocationForOk() (*string, bool) {
	if o == nil || IsNil(o.AllocationFor) {
		return nil, false
	}
	return o.AllocationFor, true
}

// HasAllocationFor returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasAllocationFor() bool {
	if o != nil && !IsNil(o.AllocationFor) {
		return true
	}

	return false
}

// SetAllocationFor gets a reference to the given string and assigns it to the AllocationFor field.
func (o *CommonsRunnerInstance) SetAllocationFor(v string) {
	o.AllocationFor = &v
}

// GetAllocationRequestedAt returns the AllocationRequestedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetAllocationRequestedAt() string {
	if o == nil || IsNil(o.AllocationRequestedAt) {
		var ret string
		return ret
	}
	return *o.AllocationRequestedAt
}

// GetAllocationRequestedAtOk returns a tuple with the AllocationRequestedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetAllocationRequestedAtOk() (*string, bool) {
	if o == nil || IsNil(o.AllocationRequestedAt) {
		return nil, false
	}
	return o.AllocationRequestedAt, true
}

// HasAllocationRequestedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasAllocationRequestedAt() bool {
	if o != nil && !IsNil(o.AllocationRequestedAt) {
		return true
	}

	return false
}

// SetAllocationRequestedAt gets a reference to the given string and assigns it to the AllocationRequestedAt field.
func (o *CommonsRunnerInstance) SetAllocationRequestedAt(v string) {
	o.AllocationRequestedAt = &v
}

// GetAllocationRequestedEventAt returns the AllocationRequestedEventAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetAllocationRequestedEventAt() string {
	if o == nil || IsNil(o.AllocationRequestedEventAt) {
		var ret string
		return ret
	}
	return *o.AllocationRequestedEventAt
}

// GetAllocationRequestedEventAtOk returns a tuple with the AllocationRequestedEventAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetAllocationRequestedEventAtOk() (*string, bool) {
	if o == nil || IsNil(o.AllocationRequestedEventAt) {
		return nil, false
	}
	return o.AllocationRequestedEventAt, true
}

// HasAllocationRequestedEventAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasAllocationRequestedEventAt() bool {
	if o != nil && !IsNil(o.AllocationRequestedEventAt) {
		return true
	}

	return false
}

// SetAllocationRequestedEventAt gets a reference to the given string and assigns it to the AllocationRequestedEventAt field.
func (o *CommonsRunnerInstance) SetAllocationRequestedEventAt(v string) {
	o.AllocationRequestedEventAt = &v
}

// GetCluster returns the Cluster field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetCluster() string {
	if o == nil || IsNil(o.Cluster) {
		var ret string
		return ret
	}
	return *o.Cluster
}

// GetClusterOk returns a tuple with the Cluster field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetClusterOk() (*string, bool) {
	if o == nil || IsNil(o.Cluster) {
		return nil, false
	}
	return o.Cluster, true
}

// HasCluster returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasCluster() bool {
	if o != nil && !IsNil(o.Cluster) {
		return true
	}

	return false
}

// SetCluster gets a reference to the given string and assigns it to the Cluster field.
func (o *CommonsRunnerInstance) SetCluster(v string) {
	o.Cluster = &v
}

// GetConfiguration returns the Configuration field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetConfiguration() CommonsRunnerInstanceConfiguration {
	if o == nil || IsNil(o.Configuration) {
		var ret CommonsRunnerInstanceConfiguration
		return ret
	}
	return *o.Configuration
}

// GetConfigurationOk returns a tuple with the Configuration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetConfigurationOk() (*CommonsRunnerInstanceConfiguration, bool) {
	if o == nil || IsNil(o.Configuration) {
		return nil, false
	}
	return o.Configuration, true
}

// HasConfiguration returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasConfiguration() bool {
	if o != nil && !IsNil(o.Configuration) {
		return true
	}

	return false
}

// SetConfiguration gets a reference to the given CommonsRunnerInstanceConfiguration and assigns it to the Configuration field.
func (o *CommonsRunnerInstance) SetConfiguration(v CommonsRunnerInstanceConfiguration) {
	o.Configuration = &v
}

// GetCreatedAt returns the CreatedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetCreatedAt() string {
	if o == nil || IsNil(o.CreatedAt) {
		var ret string
		return ret
	}
	return *o.CreatedAt
}

// GetCreatedAtOk returns a tuple with the CreatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetCreatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedAt) {
		return nil, false
	}
	return o.CreatedAt, true
}

// HasCreatedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasCreatedAt() bool {
	if o != nil && !IsNil(o.CreatedAt) {
		return true
	}

	return false
}

// SetCreatedAt gets a reference to the given string and assigns it to the CreatedAt field.
func (o *CommonsRunnerInstance) SetCreatedAt(v string) {
	o.CreatedAt = &v
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *CommonsRunnerInstance) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}
	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *CommonsRunnerInstance) SetExternalId(v string) {
	o.ExternalId = &v
}

// GetFirstPolledAt returns the FirstPolledAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetFirstPolledAt() string {
	if o == nil || IsNil(o.FirstPolledAt) {
		var ret string
		return ret
	}
	return *o.FirstPolledAt
}

// GetFirstPolledAtOk returns a tuple with the FirstPolledAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetFirstPolledAtOk() (*string, bool) {
	if o == nil || IsNil(o.FirstPolledAt) {
		return nil, false
	}
	return o.FirstPolledAt, true
}

// HasFirstPolledAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasFirstPolledAt() bool {
	if o != nil && !IsNil(o.FirstPolledAt) {
		return true
	}

	return false
}

// SetFirstPolledAt gets a reference to the given string and assigns it to the FirstPolledAt field.
func (o *CommonsRunnerInstance) SetFirstPolledAt(v string) {
	o.FirstPolledAt = &v
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *CommonsRunnerInstance) SetHost(v string) {
	o.Host = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *CommonsRunnerInstance) SetId(v string) {
	o.Id = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetLabels() []string {
	if o == nil || IsNil(o.Labels) {
		var ret []string
		return ret
	}
	return o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetLabelsOk() ([]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given []string and assigns it to the Labels field.
func (o *CommonsRunnerInstance) SetLabels(v []string) {
	o.Labels = v
}

// GetLastJobProcessedId returns the LastJobProcessedId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetLastJobProcessedId() string {
	if o == nil || IsNil(o.LastJobProcessedId) {
		var ret string
		return ret
	}
	return *o.LastJobProcessedId
}

// GetLastJobProcessedIdOk returns a tuple with the LastJobProcessedId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetLastJobProcessedIdOk() (*string, bool) {
	if o == nil || IsNil(o.LastJobProcessedId) {
		return nil, false
	}
	return o.LastJobProcessedId, true
}

// HasLastJobProcessedId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasLastJobProcessedId() bool {
	if o != nil && !IsNil(o.LastJobProcessedId) {
		return true
	}

	return false
}

// SetLastJobProcessedId gets a reference to the given string and assigns it to the LastJobProcessedId field.
func (o *CommonsRunnerInstance) SetLastJobProcessedId(v string) {
	o.LastJobProcessedId = &v
}

// GetLastJobProcessedMeta returns the LastJobProcessedMeta field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetLastJobProcessedMeta() CommonsLastJobProcessedMeta {
	if o == nil || IsNil(o.LastJobProcessedMeta) {
		var ret CommonsLastJobProcessedMeta
		return ret
	}
	return *o.LastJobProcessedMeta
}

// GetLastJobProcessedMetaOk returns a tuple with the LastJobProcessedMeta field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetLastJobProcessedMetaOk() (*CommonsLastJobProcessedMeta, bool) {
	if o == nil || IsNil(o.LastJobProcessedMeta) {
		return nil, false
	}
	return o.LastJobProcessedMeta, true
}

// HasLastJobProcessedMeta returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasLastJobProcessedMeta() bool {
	if o != nil && !IsNil(o.LastJobProcessedMeta) {
		return true
	}

	return false
}

// SetLastJobProcessedMeta gets a reference to the given CommonsLastJobProcessedMeta and assigns it to the LastJobProcessedMeta field.
func (o *CommonsRunnerInstance) SetLastJobProcessedMeta(v CommonsLastJobProcessedMeta) {
	o.LastJobProcessedMeta = &v
}

// GetLastPolledAt returns the LastPolledAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetLastPolledAt() string {
	if o == nil || IsNil(o.LastPolledAt) {
		var ret string
		return ret
	}
	return *o.LastPolledAt
}

// GetLastPolledAtOk returns a tuple with the LastPolledAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetLastPolledAtOk() (*string, bool) {
	if o == nil || IsNil(o.LastPolledAt) {
		return nil, false
	}
	return o.LastPolledAt, true
}

// HasLastPolledAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasLastPolledAt() bool {
	if o != nil && !IsNil(o.LastPolledAt) {
		return true
	}

	return false
}

// SetLastPolledAt gets a reference to the given string and assigns it to the LastPolledAt field.
func (o *CommonsRunnerInstance) SetLastPolledAt(v string) {
	o.LastPolledAt = &v
}

// GetOrganizationId returns the OrganizationId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetOrganizationId() string {
	if o == nil || IsNil(o.OrganizationId) {
		var ret string
		return ret
	}
	return *o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetOrganizationIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrganizationId) {
		return nil, false
	}
	return o.OrganizationId, true
}

// HasOrganizationId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasOrganizationId() bool {
	if o != nil && !IsNil(o.OrganizationId) {
		return true
	}

	return false
}

// SetOrganizationId gets a reference to the given string and assigns it to the OrganizationId field.
func (o *CommonsRunnerInstance) SetOrganizationId(v string) {
	o.OrganizationId = &v
}

// GetProviderKind returns the ProviderKind field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetProviderKind() string {
	if o == nil || IsNil(o.ProviderKind) {
		var ret string
		return ret
	}
	return *o.ProviderKind
}

// GetProviderKindOk returns a tuple with the ProviderKind field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetProviderKindOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderKind) {
		return nil, false
	}
	return o.ProviderKind, true
}

// HasProviderKind returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasProviderKind() bool {
	if o != nil && !IsNil(o.ProviderKind) {
		return true
	}

	return false
}

// SetProviderKind gets a reference to the given string and assigns it to the ProviderKind field.
func (o *CommonsRunnerInstance) SetProviderKind(v string) {
	o.ProviderKind = &v
}

// GetProviderKindId returns the ProviderKindId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetProviderKindId() string {
	if o == nil || IsNil(o.ProviderKindId) {
		var ret string
		return ret
	}
	return *o.ProviderKindId
}

// GetProviderKindIdOk returns a tuple with the ProviderKindId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetProviderKindIdOk() (*string, bool) {
	if o == nil || IsNil(o.ProviderKindId) {
		return nil, false
	}
	return o.ProviderKindId, true
}

// HasProviderKindId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasProviderKindId() bool {
	if o != nil && !IsNil(o.ProviderKindId) {
		return true
	}

	return false
}

// SetProviderKindId gets a reference to the given string and assigns it to the ProviderKindId field.
func (o *CommonsRunnerInstance) SetProviderKindId(v string) {
	o.ProviderKindId = &v
}

// GetPurgedAt returns the PurgedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetPurgedAt() string {
	if o == nil || IsNil(o.PurgedAt) {
		var ret string
		return ret
	}
	return *o.PurgedAt
}

// GetPurgedAtOk returns a tuple with the PurgedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetPurgedAtOk() (*string, bool) {
	if o == nil || IsNil(o.PurgedAt) {
		return nil, false
	}
	return o.PurgedAt, true
}

// HasPurgedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasPurgedAt() bool {
	if o != nil && !IsNil(o.PurgedAt) {
		return true
	}

	return false
}

// SetPurgedAt gets a reference to the given string and assigns it to the PurgedAt field.
func (o *CommonsRunnerInstance) SetPurgedAt(v string) {
	o.PurgedAt = &v
}

// GetPurgedReason returns the PurgedReason field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetPurgedReason() string {
	if o == nil || IsNil(o.PurgedReason) {
		var ret string
		return ret
	}
	return *o.PurgedReason
}

// GetPurgedReasonOk returns a tuple with the PurgedReason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetPurgedReasonOk() (*string, bool) {
	if o == nil || IsNil(o.PurgedReason) {
		return nil, false
	}
	return o.PurgedReason, true
}

// HasPurgedReason returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasPurgedReason() bool {
	if o != nil && !IsNil(o.PurgedReason) {
		return true
	}

	return false
}

// SetPurgedReason gets a reference to the given string and assigns it to the PurgedReason field.
func (o *CommonsRunnerInstance) SetPurgedReason(v string) {
	o.PurgedReason = &v
}

// GetRunnerFor returns the RunnerFor field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetRunnerFor() string {
	if o == nil || IsNil(o.RunnerFor) {
		var ret string
		return ret
	}
	return *o.RunnerFor
}

// GetRunnerForOk returns a tuple with the RunnerFor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetRunnerForOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerFor) {
		return nil, false
	}
	return o.RunnerFor, true
}

// HasRunnerFor returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasRunnerFor() bool {
	if o != nil && !IsNil(o.RunnerFor) {
		return true
	}

	return false
}

// SetRunnerFor gets a reference to the given string and assigns it to the RunnerFor field.
func (o *CommonsRunnerInstance) SetRunnerFor(v string) {
	o.RunnerFor = &v
}

// GetRunnerSetId returns the RunnerSetId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetRunnerSetId() string {
	if o == nil || IsNil(o.RunnerSetId) {
		var ret string
		return ret
	}
	return *o.RunnerSetId
}

// GetRunnerSetIdOk returns a tuple with the RunnerSetId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetRunnerSetIdOk() (*string, bool) {
	if o == nil || IsNil(o.RunnerSetId) {
		return nil, false
	}
	return o.RunnerSetId, true
}

// HasRunnerSetId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasRunnerSetId() bool {
	if o != nil && !IsNil(o.RunnerSetId) {
		return true
	}

	return false
}

// SetRunnerSetId gets a reference to the given string and assigns it to the RunnerSetId field.
func (o *CommonsRunnerInstance) SetRunnerSetId(v string) {
	o.RunnerSetId = &v
}

// GetRunningStartedAt returns the RunningStartedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetRunningStartedAt() string {
	if o == nil || IsNil(o.RunningStartedAt) {
		var ret string
		return ret
	}
	return *o.RunningStartedAt
}

// GetRunningStartedAtOk returns a tuple with the RunningStartedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetRunningStartedAtOk() (*string, bool) {
	if o == nil || IsNil(o.RunningStartedAt) {
		return nil, false
	}
	return o.RunningStartedAt, true
}

// HasRunningStartedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasRunningStartedAt() bool {
	if o != nil && !IsNil(o.RunningStartedAt) {
		return true
	}

	return false
}

// SetRunningStartedAt gets a reference to the given string and assigns it to the RunningStartedAt field.
func (o *CommonsRunnerInstance) SetRunningStartedAt(v string) {
	o.RunningStartedAt = &v
}

// GetStatus returns the Status field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetStatus() string {
	if o == nil || IsNil(o.Status) {
		var ret string
		return ret
	}
	return *o.Status
}

// GetStatusOk returns a tuple with the Status field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetStatusOk() (*string, bool) {
	if o == nil || IsNil(o.Status) {
		return nil, false
	}
	return o.Status, true
}

// HasStatus returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasStatus() bool {
	if o != nil && !IsNil(o.Status) {
		return true
	}

	return false
}

// SetStatus gets a reference to the given string and assigns it to the Status field.
func (o *CommonsRunnerInstance) SetStatus(v string) {
	o.Status = &v
}

// GetUpdatedAt returns the UpdatedAt field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetUpdatedAt() string {
	if o == nil || IsNil(o.UpdatedAt) {
		var ret string
		return ret
	}
	return *o.UpdatedAt
}

// GetUpdatedAtOk returns a tuple with the UpdatedAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetUpdatedAtOk() (*string, bool) {
	if o == nil || IsNil(o.UpdatedAt) {
		return nil, false
	}
	return o.UpdatedAt, true
}

// HasUpdatedAt returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasUpdatedAt() bool {
	if o != nil && !IsNil(o.UpdatedAt) {
		return true
	}

	return false
}

// SetUpdatedAt gets a reference to the given string and assigns it to the UpdatedAt field.
func (o *CommonsRunnerInstance) SetUpdatedAt(v string) {
	o.UpdatedAt = &v
}

// GetVcsIntegrationId returns the VcsIntegrationId field value if set, zero value otherwise.
func (o *CommonsRunnerInstance) GetVcsIntegrationId() string {
	if o == nil || IsNil(o.VcsIntegrationId) {
		var ret string
		return ret
	}
	return *o.VcsIntegrationId
}

// GetVcsIntegrationIdOk returns a tuple with the VcsIntegrationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommonsRunnerInstance) GetVcsIntegrationIdOk() (*string, bool) {
	if o == nil || IsNil(o.VcsIntegrationId) {
		return nil, false
	}
	return o.VcsIntegrationId, true
}

// HasVcsIntegrationId returns a boolean if a field has been set.
func (o *CommonsRunnerInstance) HasVcsIntegrationId() bool {
	if o != nil && !IsNil(o.VcsIntegrationId) {
		return true
	}

	return false
}

// SetVcsIntegrationId gets a reference to the given string and assigns it to the VcsIntegrationId field.
func (o *CommonsRunnerInstance) SetVcsIntegrationId(v string) {
	o.VcsIntegrationId = &v
}

func (o CommonsRunnerInstance) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CommonsRunnerInstance) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AllocatedAt) {
		toSerialize["allocated_at"] = o.AllocatedAt
	}
	if !IsNil(o.AllocationFor) {
		toSerialize["allocation_for"] = o.AllocationFor
	}
	if !IsNil(o.AllocationRequestedAt) {
		toSerialize["allocation_requested_at"] = o.AllocationRequestedAt
	}
	if !IsNil(o.AllocationRequestedEventAt) {
		toSerialize["allocation_requested_event_at"] = o.AllocationRequestedEventAt
	}
	if !IsNil(o.Cluster) {
		toSerialize["cluster"] = o.Cluster
	}
	if !IsNil(o.Configuration) {
		toSerialize["configuration"] = o.Configuration
	}
	if !IsNil(o.CreatedAt) {
		toSerialize["created_at"] = o.CreatedAt
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["created_by"] = o.CreatedBy
	}
	if !IsNil(o.ExternalId) {
		toSerialize["external_id"] = o.ExternalId
	}
	if !IsNil(o.FirstPolledAt) {
		toSerialize["first_polled_at"] = o.FirstPolledAt
	}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.LastJobProcessedId) {
		toSerialize["last_job_processed_id"] = o.LastJobProcessedId
	}
	if !IsNil(o.LastJobProcessedMeta) {
		toSerialize["last_job_processed_meta"] = o.LastJobProcessedMeta
	}
	if !IsNil(o.LastPolledAt) {
		toSerialize["last_polled_at"] = o.LastPolledAt
	}
	if !IsNil(o.OrganizationId) {
		toSerialize["organization_id"] = o.OrganizationId
	}
	if !IsNil(o.ProviderKind) {
		toSerialize["provider_kind"] = o.ProviderKind
	}
	if !IsNil(o.ProviderKindId) {
		toSerialize["provider_kind_id"] = o.ProviderKindId
	}
	if !IsNil(o.PurgedAt) {
		toSerialize["purged_at"] = o.PurgedAt
	}
	if !IsNil(o.PurgedReason) {
		toSerialize["purged_reason"] = o.PurgedReason
	}
	if !IsNil(o.RunnerFor) {
		toSerialize["runner_for"] = o.RunnerFor
	}
	if !IsNil(o.RunnerSetId) {
		toSerialize["runner_set_id"] = o.RunnerSetId
	}
	if !IsNil(o.RunningStartedAt) {
		toSerialize["running_started_at"] = o.RunningStartedAt
	}
	if !IsNil(o.Status) {
		toSerialize["status"] = o.Status
	}
	if !IsNil(o.UpdatedAt) {
		toSerialize["updated_at"] = o.UpdatedAt
	}
	if !IsNil(o.VcsIntegrationId) {
		toSerialize["vcs_integration_id"] = o.VcsIntegrationId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *CommonsRunnerInstance) UnmarshalJSON(bytes []byte) (err error) {
	varCommonsRunnerInstance := _CommonsRunnerInstance{}

	err = json.Unmarshal(bytes, &varCommonsRunnerInstance)

	if err != nil {
		return err
	}

	*o = CommonsRunnerInstance(varCommonsRunnerInstance)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(bytes, &additionalProperties); err == nil {
		delete(additionalProperties, "allocated_at")
		delete(additionalProperties, "allocation_for")
		delete(additionalProperties, "allocation_requested_at")
		delete(additionalProperties, "allocation_requested_event_at")
		delete(additionalProperties, "cluster")
		delete(additionalProperties, "configuration")
		delete(additionalProperties, "created_at")
		delete(additionalProperties, "created_by")
		delete(additionalProperties, "external_id")
		delete(additionalProperties, "first_polled_at")
		delete(additionalProperties, "host")
		delete(additionalProperties, "id")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "last_job_processed_id")
		delete(additionalProperties, "last_job_processed_meta")
		delete(additionalProperties, "last_polled_at")
		delete(additionalProperties, "organization_id")
		delete(additionalProperties, "provider_kind")
		delete(additionalProperties, "provider_kind_id")
		delete(additionalProperties, "purged_at")
		delete(additionalProperties, "purged_reason")
		delete(additionalProperties, "runner_for")
		delete(additionalProperties, "runner_set_id")
		delete(additionalProperties, "running_started_at")
		delete(additionalProperties, "status")
		delete(additionalProperties, "updated_at")
		delete(additionalProperties, "vcs_integration_id")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonsRunnerInstance struct {
	value *CommonsRunnerInstance
	isSet bool
}

func (v NullableCommonsRunnerInstance) Get() *CommonsRunnerInstance {
	return v.value
}

func (v *NullableCommonsRunnerInstance) Set(val *CommonsRunnerInstance) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonsRunnerInstance) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonsRunnerInstance) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonsRunnerInstance(val *CommonsRunnerInstance) *NullableCommonsRunnerInstance {
	return &NullableCommonsRunnerInstance{value: val, isSet: true}
}

func (v NullableCommonsRunnerInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonsRunnerInstance) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


