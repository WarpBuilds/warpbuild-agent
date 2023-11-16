# CommonsRunnerInstanceAllocationDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RunnerApplicationDetails** | Pointer to **map[string]interface{}** |  | [optional] 
**RunnerInstance** | Pointer to [**CommonsRunnerInstance**](CommonsRunnerInstance.md) |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

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

### GetRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplicationDetails() map[string]interface{}`

GetRunnerApplicationDetails returns the RunnerApplicationDetails field if non-nil, zero value otherwise.

### GetRunnerApplicationDetailsOk

`func (o *CommonsRunnerInstanceAllocationDetails) GetRunnerApplicationDetailsOk() (*map[string]interface{}, bool)`

GetRunnerApplicationDetailsOk returns a tuple with the RunnerApplicationDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) SetRunnerApplicationDetails(v map[string]interface{})`

SetRunnerApplicationDetails sets RunnerApplicationDetails field to given value.

### HasRunnerApplicationDetails

`func (o *CommonsRunnerInstanceAllocationDetails) HasRunnerApplicationDetails() bool`

HasRunnerApplicationDetails returns a boolean if a field has been set.

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


