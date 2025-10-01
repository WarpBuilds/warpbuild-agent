# CommonsByocSku

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Arch** | Pointer to **string** |  | [optional] 
**InstanceTypes** | Pointer to **[]string** |  | [optional] 
**IsPublic** | Pointer to **bool** |  | [optional] 
**NetworkTier** | Pointer to **string** | NetworkTier is the network tier associated with the instance. This is only applicable for GCP based runners. Passing this for other providers will have no effect.  Possible values:  - PREMIUM: A premium network tier with fast network performance.  - STANDARD: A standard network tier with basic features.  Refer: https://cloud.google.com/network-tiers  Default value is automatically picked if nothing is passed in case of GCP.  +Default: STANDARD | [optional] 
**RoleArn** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsByocSku

`func NewCommonsByocSku() *CommonsByocSku`

NewCommonsByocSku instantiates a new CommonsByocSku object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsByocSkuWithDefaults

`func NewCommonsByocSkuWithDefaults() *CommonsByocSku`

NewCommonsByocSkuWithDefaults instantiates a new CommonsByocSku object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetArch

`func (o *CommonsByocSku) GetArch() string`

GetArch returns the Arch field if non-nil, zero value otherwise.

### GetArchOk

`func (o *CommonsByocSku) GetArchOk() (*string, bool)`

GetArchOk returns a tuple with the Arch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetArch

`func (o *CommonsByocSku) SetArch(v string)`

SetArch sets Arch field to given value.

### HasArch

`func (o *CommonsByocSku) HasArch() bool`

HasArch returns a boolean if a field has been set.

### GetInstanceTypes

`func (o *CommonsByocSku) GetInstanceTypes() []string`

GetInstanceTypes returns the InstanceTypes field if non-nil, zero value otherwise.

### GetInstanceTypesOk

`func (o *CommonsByocSku) GetInstanceTypesOk() (*[]string, bool)`

GetInstanceTypesOk returns a tuple with the InstanceTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceTypes

`func (o *CommonsByocSku) SetInstanceTypes(v []string)`

SetInstanceTypes sets InstanceTypes field to given value.

### HasInstanceTypes

`func (o *CommonsByocSku) HasInstanceTypes() bool`

HasInstanceTypes returns a boolean if a field has been set.

### GetIsPublic

`func (o *CommonsByocSku) GetIsPublic() bool`

GetIsPublic returns the IsPublic field if non-nil, zero value otherwise.

### GetIsPublicOk

`func (o *CommonsByocSku) GetIsPublicOk() (*bool, bool)`

GetIsPublicOk returns a tuple with the IsPublic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPublic

`func (o *CommonsByocSku) SetIsPublic(v bool)`

SetIsPublic sets IsPublic field to given value.

### HasIsPublic

`func (o *CommonsByocSku) HasIsPublic() bool`

HasIsPublic returns a boolean if a field has been set.

### GetNetworkTier

`func (o *CommonsByocSku) GetNetworkTier() string`

GetNetworkTier returns the NetworkTier field if non-nil, zero value otherwise.

### GetNetworkTierOk

`func (o *CommonsByocSku) GetNetworkTierOk() (*string, bool)`

GetNetworkTierOk returns a tuple with the NetworkTier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNetworkTier

`func (o *CommonsByocSku) SetNetworkTier(v string)`

SetNetworkTier sets NetworkTier field to given value.

### HasNetworkTier

`func (o *CommonsByocSku) HasNetworkTier() bool`

HasNetworkTier returns a boolean if a field has been set.

### GetRoleArn

`func (o *CommonsByocSku) GetRoleArn() string`

GetRoleArn returns the RoleArn field if non-nil, zero value otherwise.

### GetRoleArnOk

`func (o *CommonsByocSku) GetRoleArnOk() (*string, bool)`

GetRoleArnOk returns a tuple with the RoleArn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoleArn

`func (o *CommonsByocSku) SetRoleArn(v string)`

SetRoleArn sets RoleArn field to given value.

### HasRoleArn

`func (o *CommonsByocSku) HasRoleArn() bool`

HasRoleArn returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


