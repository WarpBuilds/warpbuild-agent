# CommonsListVCSRunnerGroupsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LastSyncedAt** | Pointer to **string** |  | [optional] 
**RunnerGroups** | Pointer to [**[]CommonsRunnerGroup**](CommonsRunnerGroup.md) |  | [optional] 

## Methods

### NewCommonsListVCSRunnerGroupsResponse

`func NewCommonsListVCSRunnerGroupsResponse() *CommonsListVCSRunnerGroupsResponse`

NewCommonsListVCSRunnerGroupsResponse instantiates a new CommonsListVCSRunnerGroupsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsListVCSRunnerGroupsResponseWithDefaults

`func NewCommonsListVCSRunnerGroupsResponseWithDefaults() *CommonsListVCSRunnerGroupsResponse`

NewCommonsListVCSRunnerGroupsResponseWithDefaults instantiates a new CommonsListVCSRunnerGroupsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLastSyncedAt

`func (o *CommonsListVCSRunnerGroupsResponse) GetLastSyncedAt() string`

GetLastSyncedAt returns the LastSyncedAt field if non-nil, zero value otherwise.

### GetLastSyncedAtOk

`func (o *CommonsListVCSRunnerGroupsResponse) GetLastSyncedAtOk() (*string, bool)`

GetLastSyncedAtOk returns a tuple with the LastSyncedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSyncedAt

`func (o *CommonsListVCSRunnerGroupsResponse) SetLastSyncedAt(v string)`

SetLastSyncedAt sets LastSyncedAt field to given value.

### HasLastSyncedAt

`func (o *CommonsListVCSRunnerGroupsResponse) HasLastSyncedAt() bool`

HasLastSyncedAt returns a boolean if a field has been set.

### GetRunnerGroups

`func (o *CommonsListVCSRunnerGroupsResponse) GetRunnerGroups() []CommonsRunnerGroup`

GetRunnerGroups returns the RunnerGroups field if non-nil, zero value otherwise.

### GetRunnerGroupsOk

`func (o *CommonsListVCSRunnerGroupsResponse) GetRunnerGroupsOk() (*[]CommonsRunnerGroup, bool)`

GetRunnerGroupsOk returns a tuple with the RunnerGroups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerGroups

`func (o *CommonsListVCSRunnerGroupsResponse) SetRunnerGroups(v []CommonsRunnerGroup)`

SetRunnerGroups sets RunnerGroups field to given value.

### HasRunnerGroups

`func (o *CommonsListVCSRunnerGroupsResponse) HasRunnerGroups() bool`

HasRunnerGroups returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


