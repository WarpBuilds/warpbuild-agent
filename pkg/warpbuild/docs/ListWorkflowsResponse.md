# ListWorkflowsResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IsPulling** | Pointer to **bool** |  | [optional] 
**LastSyncedAt** | Pointer to **string** |  | [optional] 
**Workflows** | Pointer to [**[]CommonsWorkflow**](CommonsWorkflow.md) |  | [optional] 

## Methods

### NewListWorkflowsResponse

`func NewListWorkflowsResponse() *ListWorkflowsResponse`

NewListWorkflowsResponse instantiates a new ListWorkflowsResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewListWorkflowsResponseWithDefaults

`func NewListWorkflowsResponseWithDefaults() *ListWorkflowsResponse`

NewListWorkflowsResponseWithDefaults instantiates a new ListWorkflowsResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsPulling

`func (o *ListWorkflowsResponse) GetIsPulling() bool`

GetIsPulling returns the IsPulling field if non-nil, zero value otherwise.

### GetIsPullingOk

`func (o *ListWorkflowsResponse) GetIsPullingOk() (*bool, bool)`

GetIsPullingOk returns a tuple with the IsPulling field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPulling

`func (o *ListWorkflowsResponse) SetIsPulling(v bool)`

SetIsPulling sets IsPulling field to given value.

### HasIsPulling

`func (o *ListWorkflowsResponse) HasIsPulling() bool`

HasIsPulling returns a boolean if a field has been set.

### GetLastSyncedAt

`func (o *ListWorkflowsResponse) GetLastSyncedAt() string`

GetLastSyncedAt returns the LastSyncedAt field if non-nil, zero value otherwise.

### GetLastSyncedAtOk

`func (o *ListWorkflowsResponse) GetLastSyncedAtOk() (*string, bool)`

GetLastSyncedAtOk returns a tuple with the LastSyncedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastSyncedAt

`func (o *ListWorkflowsResponse) SetLastSyncedAt(v string)`

SetLastSyncedAt sets LastSyncedAt field to given value.

### HasLastSyncedAt

`func (o *ListWorkflowsResponse) HasLastSyncedAt() bool`

HasLastSyncedAt returns a boolean if a field has been set.

### GetWorkflows

`func (o *ListWorkflowsResponse) GetWorkflows() []CommonsWorkflow`

GetWorkflows returns the Workflows field if non-nil, zero value otherwise.

### GetWorkflowsOk

`func (o *ListWorkflowsResponse) GetWorkflowsOk() (*[]CommonsWorkflow, bool)`

GetWorkflowsOk returns a tuple with the Workflows field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflows

`func (o *ListWorkflowsResponse) SetWorkflows(v []CommonsWorkflow)`

SetWorkflows sets Workflows field to given value.

### HasWorkflows

`func (o *ListWorkflowsResponse) HasWorkflows() bool`

HasWorkflows returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


