# CommonsRunnerInstanceAllocationDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**GhRunnerApplicationDetails** | Pointer to [**CommonsGithubRunnerApplicationDetails**](CommonsGithubRunnerApplicationDetails.md) |  | [optional] 
**RunnerApplication** | Pointer to **string** |  | [optional] 
**RunnerInstance** | Pointer to [**CommonsRunnerInstance**](CommonsRunnerInstance.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**TelemetryEnabled** | Pointer to **bool** |  | [optional] 

## Methods

### NewCommonsRunnerInstanceAllocationDetails

`func NewCommonsRunnerInstanceAllocationDetails() *CommonsRunnerInstanceAllocationDetails`

NewCommonsRunnerInstanceAllocationDetails instantiates a new CommonsRunnerInstanceAllocationDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerInstanceAllocationDetailsWithDefaults

`func NewCommonsRunnerInstanceAllocationDetailsWithDefaults() *CommonsRunnerInstanceAllocationDetails`

NewCommonsRunnerInstanceAllocationDetailsWithDefaults instantiates a new CommonsRunnerInstanceAllocationDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGhRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) GetGhRunnerApplicationDetails() CommonsGithubRunnerApplicationDetails`

GetGhRunnerApplicationDetails returns the GhRunnerApplicationDetails field if non-nil, zero value otherwise.

### GetGhRunnerApplicationDetailsOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetGhRunnerApplicationDetailsOk() (*CommonsGithubRunnerApplicationDetails, bool)`

GetGhRunnerApplicationDetailsOk returns a tuple with the GhRunnerApplicationDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGhRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) SetGhRunnerApplicationDetails(v CommonsGithubRunnerApplicationDetails)`

SetGhRunnerApplicationDetails sets GhRunnerApplicationDetails field to given value.

### HasGhRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) HasGhRunnerApplicationDetails() bool`

HasGhRunnerApplicationDetails returns a boolean if a field has been set.

### GetRunnerApplication

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplication() string`

GetRunnerApplication returns the RunnerApplication field if non-nil, zero value otherwise.

### GetRunnerApplicationOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplicationOk() (*string, bool)`

GetRunnerApplicationOk returns a tuple with the RunnerApplication field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerApplication

`func (o *CommonsRunnerInstanceAllocationDetails) SetRunnerApplication(v string)`

SetRunnerApplication sets RunnerApplication field to given value.

### HasRunnerApplication

`func (o *CommonsRunnerInstanceAllocationDetails) HasRunnerApplication() bool`

HasRunnerApplication returns a boolean if a field has been set.

### GetRunnerInstance

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerInstance() CommonsRunnerInstance`

GetRunnerInstance returns the RunnerInstance field if non-nil, zero value otherwise.

### GetRunnerInstanceOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerInstanceOk() (*CommonsRunnerInstance, bool)`

GetRunnerInstanceOk returns a tuple with the RunnerInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerInstance

`func (o *CommonsRunnerInstanceAllocationDetails) SetRunnerInstance(v CommonsRunnerInstance)`

SetRunnerInstance sets RunnerInstance field to given value.

### HasRunnerInstance

`func (o *CommonsRunnerInstanceAllocationDetails) HasRunnerInstance() bool`

HasRunnerInstance returns a boolean if a field has been set.

### GetStatus

`func (o *CommonsRunnerInstanceAllocationDetails) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CommonsRunnerInstanceAllocationDetails) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CommonsRunnerInstanceAllocationDetails) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTelemetryEnabled

`func (o *CommonsRunnerInstanceAllocationDetails) GetTelemetryEnabled() bool`

GetTelemetryEnabled returns the TelemetryEnabled field if non-nil, zero value otherwise.

### GetTelemetryEnabledOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetTelemetryEnabledOk() (*bool, bool)`

GetTelemetryEnabledOk returns a tuple with the TelemetryEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTelemetryEnabled

`func (o *CommonsRunnerInstanceAllocationDetails) SetTelemetryEnabled(v bool)`

SetTelemetryEnabled sets TelemetryEnabled field to given value.

### HasTelemetryEnabled

`func (o *CommonsRunnerInstanceAllocationDetails) HasTelemetryEnabled() bool`

HasTelemetryEnabled returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


