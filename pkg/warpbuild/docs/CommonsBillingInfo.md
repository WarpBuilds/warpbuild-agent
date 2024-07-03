# CommonsBillingInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AppliedCoupon** | Pointer to [**CommonsCoupon**](CommonsCoupon.md) |  | [optional] 
**Email** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsBillingInfo

`func NewCommonsBillingInfo() *CommonsBillingInfo`

NewCommonsBillingInfo instantiates a new CommonsBillingInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsBillingInfoWithDefaults

`func NewCommonsBillingInfoWithDefaults() *CommonsBillingInfo`

NewCommonsBillingInfoWithDefaults instantiates a new CommonsBillingInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAppliedCoupon

`func (o *CommonsBillingInfo) GetAppliedCoupon() CommonsCoupon`

GetAppliedCoupon returns the AppliedCoupon field if non-nil, zero value otherwise.

### GetAppliedCouponOk

`func (o *CommonsBillingInfo) GetAppliedCouponOk() (*CommonsCoupon, bool)`

GetAppliedCouponOk returns a tuple with the AppliedCoupon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAppliedCoupon

`func (o *CommonsBillingInfo) SetAppliedCoupon(v CommonsCoupon)`

SetAppliedCoupon sets AppliedCoupon field to given value.

### HasAppliedCoupon

`func (o *CommonsBillingInfo) HasAppliedCoupon() bool`

HasAppliedCoupon returns a boolean if a field has been set.

### GetEmail

`func (o *CommonsBillingInfo) GetEmail() string`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *CommonsBillingInfo) GetEmailOk() (*string, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *CommonsBillingInfo) SetEmail(v string)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *CommonsBillingInfo) HasEmail() bool`

HasEmail returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


