# CommonsBalanceDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FreeBalance** | Pointer to [**CommonsBalance**](CommonsBalance.md) |  | [optional] 
**PaidBalance** | Pointer to [**CommonsBalance**](CommonsBalance.md) |  | [optional] 

## Methods

### NewCommonsBalanceDetails

`func NewCommonsBalanceDetails() *CommonsBalanceDetails`

NewCommonsBalanceDetails instantiates a new CommonsBalanceDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsBalanceDetailsWithDefaults

`func NewCommonsBalanceDetailsWithDefaults() *CommonsBalanceDetails`

NewCommonsBalanceDetailsWithDefaults instantiates a new CommonsBalanceDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFreeBalance

`func (o *CommonsBalanceDetails) GetFreeBalance() CommonsBalance`

GetFreeBalance returns the FreeBalance field if non-nil, zero value otherwise.

### GetFreeBalanceOk

`func (o *CommonsBalanceDetails) GetFreeBalanceOk() (*CommonsBalance, bool)`

GetFreeBalanceOk returns a tuple with the FreeBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFreeBalance

`func (o *CommonsBalanceDetails) SetFreeBalance(v CommonsBalance)`

SetFreeBalance sets FreeBalance field to given value.

### HasFreeBalance

`func (o *CommonsBalanceDetails) HasFreeBalance() bool`

HasFreeBalance returns a boolean if a field has been set.

### GetPaidBalance

`func (o *CommonsBalanceDetails) GetPaidBalance() CommonsBalance`

GetPaidBalance returns the PaidBalance field if non-nil, zero value otherwise.

### GetPaidBalanceOk

`func (o *CommonsBalanceDetails) GetPaidBalanceOk() (*CommonsBalance, bool)`

GetPaidBalanceOk returns a tuple with the PaidBalance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaidBalance

`func (o *CommonsBalanceDetails) SetPaidBalance(v CommonsBalance)`

SetPaidBalance sets PaidBalance field to given value.

### HasPaidBalance

`func (o *CommonsBalanceDetails) HasPaidBalance() bool`

HasPaidBalance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


