# CommonsPaymentMethod

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CountryCode** | Pointer to **string** |  | [optional] 
**ExpMonth** | Pointer to **int32** |  | [optional] 
**ExpYear** | Pointer to **int32** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**Iin** | Pointer to **string** |  | [optional] 
**IsDefault** | Pointer to **bool** |  | [optional] 
**Issuer** | Pointer to **string** |  | [optional] 
**Last4** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Network** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsPaymentMethod

`func NewCommonsPaymentMethod() *CommonsPaymentMethod`

NewCommonsPaymentMethod instantiates a new CommonsPaymentMethod object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsPaymentMethodWithDefaults

`func NewCommonsPaymentMethodWithDefaults() *CommonsPaymentMethod`

NewCommonsPaymentMethodWithDefaults instantiates a new CommonsPaymentMethod object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCountryCode

`func (o *CommonsPaymentMethod) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *CommonsPaymentMethod) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *CommonsPaymentMethod) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.

### HasCountryCode

`func (o *CommonsPaymentMethod) HasCountryCode() bool`

HasCountryCode returns a boolean if a field has been set.

### GetExpMonth

`func (o *CommonsPaymentMethod) GetExpMonth() int32`

GetExpMonth returns the ExpMonth field if non-nil, zero value otherwise.

### GetExpMonthOk

`func (o *CommonsPaymentMethod) GetExpMonthOk() (*int32, bool)`

GetExpMonthOk returns a tuple with the ExpMonth field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpMonth

`func (o *CommonsPaymentMethod) SetExpMonth(v int32)`

SetExpMonth sets ExpMonth field to given value.

### HasExpMonth

`func (o *CommonsPaymentMethod) HasExpMonth() bool`

HasExpMonth returns a boolean if a field has been set.

### GetExpYear

`func (o *CommonsPaymentMethod) GetExpYear() int32`

GetExpYear returns the ExpYear field if non-nil, zero value otherwise.

### GetExpYearOk

`func (o *CommonsPaymentMethod) GetExpYearOk() (*int32, bool)`

GetExpYearOk returns a tuple with the ExpYear field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpYear

`func (o *CommonsPaymentMethod) SetExpYear(v int32)`

SetExpYear sets ExpYear field to given value.

### HasExpYear

`func (o *CommonsPaymentMethod) HasExpYear() bool`

HasExpYear returns a boolean if a field has been set.

### GetExternalId

`func (o *CommonsPaymentMethod) GetExternalId() string`

GetExternalId returns the ExternalId field if non-nil, zero value otherwise.

### GetExternalIdOk

`func (o *CommonsPaymentMethod) GetExternalIdOk() (*string, bool)`

GetExternalIdOk returns a tuple with the ExternalId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalId

`func (o *CommonsPaymentMethod) SetExternalId(v string)`

SetExternalId sets ExternalId field to given value.

### HasExternalId

`func (o *CommonsPaymentMethod) HasExternalId() bool`

HasExternalId returns a boolean if a field has been set.

### GetIin

`func (o *CommonsPaymentMethod) GetIin() string`

GetIin returns the Iin field if non-nil, zero value otherwise.

### GetIinOk

`func (o *CommonsPaymentMethod) GetIinOk() (*string, bool)`

GetIinOk returns a tuple with the Iin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIin

`func (o *CommonsPaymentMethod) SetIin(v string)`

SetIin sets Iin field to given value.

### HasIin

`func (o *CommonsPaymentMethod) HasIin() bool`

HasIin returns a boolean if a field has been set.

### GetIsDefault

`func (o *CommonsPaymentMethod) GetIsDefault() bool`

GetIsDefault returns the IsDefault field if non-nil, zero value otherwise.

### GetIsDefaultOk

`func (o *CommonsPaymentMethod) GetIsDefaultOk() (*bool, bool)`

GetIsDefaultOk returns a tuple with the IsDefault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDefault

`func (o *CommonsPaymentMethod) SetIsDefault(v bool)`

SetIsDefault sets IsDefault field to given value.

### HasIsDefault

`func (o *CommonsPaymentMethod) HasIsDefault() bool`

HasIsDefault returns a boolean if a field has been set.

### GetIssuer

`func (o *CommonsPaymentMethod) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *CommonsPaymentMethod) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *CommonsPaymentMethod) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### HasIssuer

`func (o *CommonsPaymentMethod) HasIssuer() bool`

HasIssuer returns a boolean if a field has been set.

### GetLast4

`func (o *CommonsPaymentMethod) GetLast4() string`

GetLast4 returns the Last4 field if non-nil, zero value otherwise.

### GetLast4Ok

`func (o *CommonsPaymentMethod) GetLast4Ok() (*string, bool)`

GetLast4Ok returns a tuple with the Last4 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLast4

`func (o *CommonsPaymentMethod) SetLast4(v string)`

SetLast4 sets Last4 field to given value.

### HasLast4

`func (o *CommonsPaymentMethod) HasLast4() bool`

HasLast4 returns a boolean if a field has been set.

### GetName

`func (o *CommonsPaymentMethod) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommonsPaymentMethod) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommonsPaymentMethod) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommonsPaymentMethod) HasName() bool`

HasName returns a boolean if a field has been set.

### GetNetwork

`func (o *CommonsPaymentMethod) GetNetwork() string`

GetNetwork returns the Network field if non-nil, zero value otherwise.

### GetNetworkOk

`func (o *CommonsPaymentMethod) GetNetworkOk() (*string, bool)`

GetNetworkOk returns a tuple with the Network field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetwork

`func (o *CommonsPaymentMethod) SetNetwork(v string)`

SetNetwork sets Network field to given value.

### HasNetwork

`func (o *CommonsPaymentMethod) HasNetwork() bool`

HasNetwork returns a boolean if a field has been set.

### GetType

`func (o *CommonsPaymentMethod) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsPaymentMethod) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsPaymentMethod) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *CommonsPaymentMethod) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


