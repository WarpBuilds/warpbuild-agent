# WarpWorkflowsRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authorization** | Pointer to **map[string]string** | for github oauth flow | [optional] 
**RunnerId** | **string** |  | 
**WorkflowIds** | **[]string** |  | 

## Methods

### NewWarpWorkflowsRequest

`func NewWarpWorkflowsRequest(runnerId string, workflowIds []string, ) *WarpWorkflowsRequest`

NewWarpWorkflowsRequest instantiates a new WarpWorkflowsRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWarpWorkflowsRequestWithDefaults

`func NewWarpWorkflowsRequestWithDefaults() *WarpWorkflowsRequest`

NewWarpWorkflowsRequestWithDefaults instantiates a new WarpWorkflowsRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthorization

`func (o *WarpWorkflowsRequest) GetAuthorization() map[string]string`

GetAuthorization returns the Authorization field if non-nil, zero value otherwise.

### GetAuthorizationOk

`func (o *WarpWorkflowsRequest) GetAuthorizationOk() (*map[string]string, bool)`

GetAuthorizationOk returns a tuple with the Authorization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthorization

`func (o *WarpWorkflowsRequest) SetAuthorization(v map[string]string)`

SetAuthorization sets Authorization field to given value.

### HasAuthorization

`func (o *WarpWorkflowsRequest) HasAuthorization() bool`

HasAuthorization returns a boolean if a field has been set.

### GetRunnerId

`func (o *WarpWorkflowsRequest) GetRunnerId() string`

GetRunnerId returns the RunnerId field if non-nil, zero value otherwise.

### GetRunnerIdOk

`func (o *WarpWorkflowsRequest) GetRunnerIdOk() (*string, bool)`

GetRunnerIdOk returns a tuple with the RunnerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRunnerId

`func (o *WarpWorkflowsRequest) SetRunnerId(v string)`

SetRunnerId sets RunnerId field to given value.


### GetWorkflowIds

`func (o *WarpWorkflowsRequest) GetWorkflowIds() []string`

GetWorkflowIds returns the WorkflowIds field if non-nil, zero value otherwise.

### GetWorkflowIdsOk

`func (o *WarpWorkflowsRequest) GetWorkflowIdsOk() (*[]string, bool)`

GetWorkflowIdsOk returns a tuple with the WorkflowIds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowIds

`func (o *WarpWorkflowsRequest) SetWorkflowIds(v []string)`

SetWorkflowIds sets WorkflowIds field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


