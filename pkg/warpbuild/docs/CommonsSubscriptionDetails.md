# CommonsSubscriptionDetails

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BalanceDetails** | Pointer to [**CommonsBalanceDetails**](CommonsBalanceDetails.md) |  | [optional] 
**CreatedAt** | Pointer to **string** |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**LastPaymentAt** | Pointer to **string** |  | [optional] 
**LastPaymentDetails** | Pointer to [**CommonsPaymentDetails**](CommonsPaymentDetails.md) |  | [optional] 
**OrganizationId** | Pointer to **string** |  | [optional] 
**PaymentMethods** | Pointer to [**[]CommonsPaymentMethod**](CommonsPaymentMethod.md) |  | [optional] 
**PgCustomerId** | Pointer to **string** |  | [optional] 
**PgSubscriptionId** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**SubscribedAt** | Pointer to **string** |  | [optional] 
**UpcomingBill** | Pointer to [**CommonsUpcomingBill**](CommonsUpcomingBill.md) |  | [optional] 
**UpdatedAt** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsSubscriptionDetails

`func NewCommonsSubscriptionDetails() *CommonsSubscriptionDetails`

NewCommonsSubscriptionDetails instantiates a new CommonsSubscriptionDetails object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsSubscriptionDetailsWithDefaults

`func NewCommonsSubscriptionDetailsWithDefaults() *CommonsSubscriptionDetails`

NewCommonsSubscriptionDetailsWithDefaults instantiates a new CommonsSubscriptionDetails object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBalanceDetails

`func (o *CommonsSubscriptionDetails) GetBalanceDetails() CommonsBalanceDetails`

GetBalanceDetails returns the BalanceDetails field if non-nil, zero value otherwise.

### GetBalanceDetailsOk

`func (o *CommonsSubscriptionDetails) GetBalanceDetailsOk() (*CommonsBalanceDetails, bool)`

GetBalanceDetailsOk returns a tuple with the BalanceDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBalanceDetails

`func (o *CommonsSubscriptionDetails) SetBalanceDetails(v CommonsBalanceDetails)`

SetBalanceDetails sets BalanceDetails field to given value.

### HasBalanceDetails

`func (o *CommonsSubscriptionDetails) HasBalanceDetails() bool`

HasBalanceDetails returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsSubscriptionDetails) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsSubscriptionDetails) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsSubscriptionDetails) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *CommonsSubscriptionDetails) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetId

`func (o *CommonsSubscriptionDetails) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsSubscriptionDetails) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsSubscriptionDetails) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *CommonsSubscriptionDetails) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLastPaymentAt

`func (o *CommonsSubscriptionDetails) GetLastPaymentAt() string`

GetLastPaymentAt returns the LastPaymentAt field if non-nil, zero value otherwise.

### GetLastPaymentAtOk

`func (o *CommonsSubscriptionDetails) GetLastPaymentAtOk() (*string, bool)`

GetLastPaymentAtOk returns a tuple with the LastPaymentAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPaymentAt

`func (o *CommonsSubscriptionDetails) SetLastPaymentAt(v string)`

SetLastPaymentAt sets LastPaymentAt field to given value.

### HasLastPaymentAt

`func (o *CommonsSubscriptionDetails) HasLastPaymentAt() bool`

HasLastPaymentAt returns a boolean if a field has been set.

### GetLastPaymentDetails

`func (o *CommonsSubscriptionDetails) GetLastPaymentDetails() CommonsPaymentDetails`

GetLastPaymentDetails returns the LastPaymentDetails field if non-nil, zero value otherwise.

### GetLastPaymentDetailsOk

`func (o *CommonsSubscriptionDetails) GetLastPaymentDetailsOk() (*CommonsPaymentDetails, bool)`

GetLastPaymentDetailsOk returns a tuple with the LastPaymentDetails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastPaymentDetails

`func (o *CommonsSubscriptionDetails) SetLastPaymentDetails(v CommonsPaymentDetails)`

SetLastPaymentDetails sets LastPaymentDetails field to given value.

### HasLastPaymentDetails

`func (o *CommonsSubscriptionDetails) HasLastPaymentDetails() bool`

HasLastPaymentDetails returns a boolean if a field has been set.

### GetOrganizationId

`func (o *CommonsSubscriptionDetails) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsSubscriptionDetails) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsSubscriptionDetails) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *CommonsSubscriptionDetails) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetPaymentMethods

`func (o *CommonsSubscriptionDetails) GetPaymentMethods() []CommonsPaymentMethod`

GetPaymentMethods returns the PaymentMethods field if non-nil, zero value otherwise.

### GetPaymentMethodsOk

`func (o *CommonsSubscriptionDetails) GetPaymentMethodsOk() (*[]CommonsPaymentMethod, bool)`

GetPaymentMethodsOk returns a tuple with the PaymentMethods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPaymentMethods

`func (o *CommonsSubscriptionDetails) SetPaymentMethods(v []CommonsPaymentMethod)`

SetPaymentMethods sets PaymentMethods field to given value.

### HasPaymentMethods

`func (o *CommonsSubscriptionDetails) HasPaymentMethods() bool`

HasPaymentMethods returns a boolean if a field has been set.

### GetPgCustomerId

`func (o *CommonsSubscriptionDetails) GetPgCustomerId() string`

GetPgCustomerId returns the PgCustomerId field if non-nil, zero value otherwise.

### GetPgCustomerIdOk

`func (o *CommonsSubscriptionDetails) GetPgCustomerIdOk() (*string, bool)`

GetPgCustomerIdOk returns a tuple with the PgCustomerId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPgCustomerId

`func (o *CommonsSubscriptionDetails) SetPgCustomerId(v string)`

SetPgCustomerId sets PgCustomerId field to given value.

### HasPgCustomerId

`func (o *CommonsSubscriptionDetails) HasPgCustomerId() bool`

HasPgCustomerId returns a boolean if a field has been set.

### GetPgSubscriptionId

`func (o *CommonsSubscriptionDetails) GetPgSubscriptionId() string`

GetPgSubscriptionId returns the PgSubscriptionId field if non-nil, zero value otherwise.

### GetPgSubscriptionIdOk

`func (o *CommonsSubscriptionDetails) GetPgSubscriptionIdOk() (*string, bool)`

GetPgSubscriptionIdOk returns a tuple with the PgSubscriptionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPgSubscriptionId

`func (o *CommonsSubscriptionDetails) SetPgSubscriptionId(v string)`

SetPgSubscriptionId sets PgSubscriptionId field to given value.

### HasPgSubscriptionId

`func (o *CommonsSubscriptionDetails) HasPgSubscriptionId() bool`

HasPgSubscriptionId returns a boolean if a field has been set.

### GetStatus

`func (o *CommonsSubscriptionDetails) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *CommonsSubscriptionDetails) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *CommonsSubscriptionDetails) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *CommonsSubscriptionDetails) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetSubscribedAt

`func (o *CommonsSubscriptionDetails) GetSubscribedAt() string`

GetSubscribedAt returns the SubscribedAt field if non-nil, zero value otherwise.

### GetSubscribedAtOk

`func (o *CommonsSubscriptionDetails) GetSubscribedAtOk() (*string, bool)`

GetSubscribedAtOk returns a tuple with the SubscribedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubscribedAt

`func (o *CommonsSubscriptionDetails) SetSubscribedAt(v string)`

SetSubscribedAt sets SubscribedAt field to given value.

### HasSubscribedAt

`func (o *CommonsSubscriptionDetails) HasSubscribedAt() bool`

HasSubscribedAt returns a boolean if a field has been set.

### GetUpcomingBill

`func (o *CommonsSubscriptionDetails) GetUpcomingBill() CommonsUpcomingBill`

GetUpcomingBill returns the UpcomingBill field if non-nil, zero value otherwise.

### GetUpcomingBillOk

`func (o *CommonsSubscriptionDetails) GetUpcomingBillOk() (*CommonsUpcomingBill, bool)`

GetUpcomingBillOk returns a tuple with the UpcomingBill field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpcomingBill

`func (o *CommonsSubscriptionDetails) SetUpcomingBill(v CommonsUpcomingBill)`

SetUpcomingBill sets UpcomingBill field to given value.

### HasUpcomingBill

`func (o *CommonsSubscriptionDetails) HasUpcomingBill() bool`

HasUpcomingBill returns a boolean if a field has been set.

### GetUpdatedAt

`func (o *CommonsSubscriptionDetails) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsSubscriptionDetails) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsSubscriptionDetails) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.

### HasUpdatedAt

`func (o *CommonsSubscriptionDetails) HasUpdatedAt() bool`

HasUpdatedAt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


