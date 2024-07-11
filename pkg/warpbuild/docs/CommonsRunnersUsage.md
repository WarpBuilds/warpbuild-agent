# CommonsRunnersUsage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Daywise** | Pointer to [**[]CommonsDaywiseRuntime**](CommonsDaywiseRuntime.md) |  | [optional] 
**Runnerwise** | Pointer to [**[]CommonsRunnerwiseRuntime**](CommonsRunnerwiseRuntime.md) |  | [optional] 
**TotalJobCount** | Pointer to **int32** |  | [optional] 
**TotalRuntimeSeconds** | Pointer to **int32** |  | [optional] 

## Methods

### NewCommonsRunnersUsage

`func NewCommonsRunnersUsage() *CommonsRunnersUsage`

NewCommonsRunnersUsage instantiates a new CommonsRunnersUsage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnersUsageWithDefaults

`func NewCommonsRunnersUsageWithDefaults() *CommonsRunnersUsage`

NewCommonsRunnersUsageWithDefaults instantiates a new CommonsRunnersUsage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDaywise

`func (o *CommonsRunnersUsage) GetDaywise() []CommonsDaywiseRuntime`

GetDaywise returns the Daywise field if non-nil, zero value otherwise.

### GetDaywiseOk

`func (o *CommonsRunnersUsage) GetDaywiseOk() (*[]CommonsDaywiseRuntime, bool)`

GetDaywiseOk returns a tuple with the Daywise field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDaywise

`func (o *CommonsRunnersUsage) SetDaywise(v []CommonsDaywiseRuntime)`

SetDaywise sets Daywise field to given value.

### HasDaywise

`func (o *CommonsRunnersUsage) HasDaywise() bool`

HasDaywise returns a boolean if a field has been set.

### GetRunnerwise

`func (o *CommonsRunnersUsage) GetRunnerwise() []CommonsRunnerwiseRuntime`

GetRunnerwise returns the Runnerwise field if non-nil, zero value otherwise.

### GetRunnerwiseOk

`func (o *CommonsRunnersUsage) GetRunnerwiseOk() (*[]CommonsRunnerwiseRuntime, bool)`

GetRunnerwiseOk returns a tuple with the Runnerwise field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerwise

`func (o *CommonsRunnersUsage) SetRunnerwise(v []CommonsRunnerwiseRuntime)`

SetRunnerwise sets Runnerwise field to given value.

### HasRunnerwise

`func (o *CommonsRunnersUsage) HasRunnerwise() bool`

HasRunnerwise returns a boolean if a field has been set.

### GetTotalJobCount

`func (o *CommonsRunnersUsage) GetTotalJobCount() int32`

GetTotalJobCount returns the TotalJobCount field if non-nil, zero value otherwise.

### GetTotalJobCountOk

`func (o *CommonsRunnersUsage) GetTotalJobCountOk() (*int32, bool)`

GetTotalJobCountOk returns a tuple with the TotalJobCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalJobCount

`func (o *CommonsRunnersUsage) SetTotalJobCount(v int32)`

SetTotalJobCount sets TotalJobCount field to given value.

### HasTotalJobCount

`func (o *CommonsRunnersUsage) HasTotalJobCount() bool`

HasTotalJobCount returns a boolean if a field has been set.

### GetTotalRuntimeSeconds

`func (o *CommonsRunnersUsage) GetTotalRuntimeSeconds() int32`

GetTotalRuntimeSeconds returns the TotalRuntimeSeconds field if non-nil, zero value otherwise.

### GetTotalRuntimeSecondsOk

`func (o *CommonsRunnersUsage) GetTotalRuntimeSecondsOk() (*int32, bool)`

GetTotalRuntimeSecondsOk returns a tuple with the TotalRuntimeSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalRuntimeSeconds

`func (o *CommonsRunnersUsage) SetTotalRuntimeSeconds(v int32)`

SetTotalRuntimeSeconds sets TotalRuntimeSeconds field to given value.

### HasTotalRuntimeSeconds

`func (o *CommonsRunnersUsage) HasTotalRuntimeSeconds() bool`

HasTotalRuntimeSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


