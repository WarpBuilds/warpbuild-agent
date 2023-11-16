# CommonsListUsersResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Next** | Pointer to **int32** |  | [optional] 
**Page** | Pointer to **int32** |  | [optional] 
**PerPage** | Pointer to **int32** |  | [optional] 
**TotalPages** | Pointer to **int32** |  | [optional] 
**TotalRows** | Pointer to **int32** |  | [optional] 
**Users** | Pointer to [**[]V1User**](V1User.md) |  | [optional] 

## Methods

### NewCommonsListUsersResponse

`func NewCommonsListUsersResponse() *CommonsListUsersResponse`

NewCommonsListUsersResponse instantiates a new CommonsListUsersResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsListUsersResponseWithDefaults

`func NewCommonsListUsersResponseWithDefaults() *CommonsListUsersResponse`

NewCommonsListUsersResponseWithDefaults instantiates a new CommonsListUsersResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNext

`func (o *CommonsListUsersResponse) GetNext() int32`

GetNext returns the Next field if non-nil, zero value otherwise.

### GetNextOk

`func (o *CommonsListUsersResponse) GetNextOk() (*int32, bool)`

GetNextOk returns a tuple with the Next field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNext

`func (o *CommonsListUsersResponse) SetNext(v int32)`

SetNext sets Next field to given value.

### HasNext

`func (o *CommonsListUsersResponse) HasNext() bool`

HasNext returns a boolean if a field has been set.

### GetPage

`func (o *CommonsListUsersResponse) GetPage() int32`

GetPage returns the Page field if non-nil, zero value otherwise.

### GetPageOk

`func (o *CommonsListUsersResponse) GetPageOk() (*int32, bool)`

GetPageOk returns a tuple with the Page field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPage

`func (o *CommonsListUsersResponse) SetPage(v int32)`

SetPage sets Page field to given value.

### HasPage

`func (o *CommonsListUsersResponse) HasPage() bool`

HasPage returns a boolean if a field has been set.

### GetPerPage

`func (o *CommonsListUsersResponse) GetPerPage() int32`

GetPerPage returns the PerPage field if non-nil, zero value otherwise.

### GetPerPageOk

`func (o *CommonsListUsersResponse) GetPerPageOk() (*int32, bool)`

GetPerPageOk returns a tuple with the PerPage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPerPage

`func (o *CommonsListUsersResponse) SetPerPage(v int32)`

SetPerPage sets PerPage field to given value.

### HasPerPage

`func (o *CommonsListUsersResponse) HasPerPage() bool`

HasPerPage returns a boolean if a field has been set.

### GetTotalPages

`func (o *CommonsListUsersResponse) GetTotalPages() int32`

GetTotalPages returns the TotalPages field if non-nil, zero value otherwise.

### GetTotalPagesOk

`func (o *CommonsListUsersResponse) GetTotalPagesOk() (*int32, bool)`

GetTotalPagesOk returns a tuple with the TotalPages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalPages

`func (o *CommonsListUsersResponse) SetTotalPages(v int32)`

SetTotalPages sets TotalPages field to given value.

### HasTotalPages

`func (o *CommonsListUsersResponse) HasTotalPages() bool`

HasTotalPages returns a boolean if a field has been set.

### GetTotalRows

`func (o *CommonsListUsersResponse) GetTotalRows() int32`

GetTotalRows returns the TotalRows field if non-nil, zero value otherwise.

### GetTotalRowsOk

`func (o *CommonsListUsersResponse) GetTotalRowsOk() (*int32, bool)`

GetTotalRowsOk returns a tuple with the TotalRows field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalRows

`func (o *CommonsListUsersResponse) SetTotalRows(v int32)`

SetTotalRows sets TotalRows field to given value.

### HasTotalRows

`func (o *CommonsListUsersResponse) HasTotalRows() bool`

HasTotalRows returns a boolean if a field has been set.

### GetUsers

`func (o *CommonsListUsersResponse) GetUsers() []V1User`

GetUsers returns the Users field if non-nil, zero value otherwise.

### GetUsersOk

`func (o *CommonsListUsersResponse) GetUsersOk() (*[]V1User, bool)`

GetUsersOk returns a tuple with the Users field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsers

`func (o *CommonsListUsersResponse) SetUsers(v []V1User)`

SetUsers sets Users field to given value.

### HasUsers

`func (o *CommonsListUsersResponse) HasUsers() bool`

HasUsers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


