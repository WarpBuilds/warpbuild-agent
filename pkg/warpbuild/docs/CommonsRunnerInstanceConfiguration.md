# CommonsRunnerInstanceConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CapacityType** | Pointer to **string** |  | [optional] 
**Image** | Pointer to **string** | Refer | [optional] 
**ProviderSkuMapping** | Pointer to [**[]CommonsProviderInstanceSkuMapping**](CommonsProviderInstanceSkuMapping.md) |  | [optional] 
**Sku** | Pointer to [**CommonsInstanceSku**](CommonsInstanceSku.md) |  | [optional] 
**StockRunnerSetId** | Pointer to **string** |  | [optional] 
**Storage** | Pointer to [**CommonsStorage**](CommonsStorage.md) |  | [optional] 

## Methods

### NewCommonsRunnerInstanceConfiguration

`func NewCommonsRunnerInstanceConfiguration() *CommonsRunnerInstanceConfiguration`

NewCommonsRunnerInstanceConfiguration instantiates a new CommonsRunnerInstanceConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerInstanceConfigurationWithDefaults

`func NewCommonsRunnerInstanceConfigurationWithDefaults() *CommonsRunnerInstanceConfiguration`

NewCommonsRunnerInstanceConfigurationWithDefaults instantiates a new CommonsRunnerInstanceConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCapacityType

`func (o *CommonsRunnerInstanceConfiguration) GetCapacityType() string`

GetCapacityType returns the CapacityType field if non-nil, zero value otherwise.

### GetCapacityTypeOk

`func (o *CommonsRunnerInstanceConfiguration) GetCapacityTypeOk() (*string, bool)`

GetCapacityTypeOk returns a tuple with the CapacityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCapacityType

`func (o *CommonsRunnerInstanceConfiguration) SetCapacityType(v string)`

SetCapacityType sets CapacityType field to given value.

### HasCapacityType

`func (o *CommonsRunnerInstanceConfiguration) HasCapacityType() bool`

HasCapacityType returns a boolean if a field has been set.

### GetImage

`func (o *CommonsRunnerInstanceConfiguration) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *CommonsRunnerInstanceConfiguration) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *CommonsRunnerInstanceConfiguration) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *CommonsRunnerInstanceConfiguration) HasImage() bool`

HasImage returns a boolean if a field has been set.

### GetProviderSkuMapping

`func (o *CommonsRunnerInstanceConfiguration) GetProviderSkuMapping() []CommonsProviderInstanceSkuMapping`

GetProviderSkuMapping returns the ProviderSkuMapping field if non-nil, zero value otherwise.

### GetProviderSkuMappingOk

`func (o *CommonsRunnerInstanceConfiguration) GetProviderSkuMappingOk() (*[]CommonsProviderInstanceSkuMapping, bool)`

GetProviderSkuMappingOk returns a tuple with the ProviderSkuMapping field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProviderSkuMapping

`func (o *CommonsRunnerInstanceConfiguration) SetProviderSkuMapping(v []CommonsProviderInstanceSkuMapping)`

SetProviderSkuMapping sets ProviderSkuMapping field to given value.

### HasProviderSkuMapping

`func (o *CommonsRunnerInstanceConfiguration) HasProviderSkuMapping() bool`

HasProviderSkuMapping returns a boolean if a field has been set.

### GetSku

`func (o *CommonsRunnerInstanceConfiguration) GetSku() CommonsInstanceSku`

GetSku returns the Sku field if non-nil, zero value otherwise.

### GetSkuOk

`func (o *CommonsRunnerInstanceConfiguration) GetSkuOk() (*CommonsInstanceSku, bool)`

GetSkuOk returns a tuple with the Sku field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSku

`func (o *CommonsRunnerInstanceConfiguration) SetSku(v CommonsInstanceSku)`

SetSku sets Sku field to given value.

### HasSku

`func (o *CommonsRunnerInstanceConfiguration) HasSku() bool`

HasSku returns a boolean if a field has been set.

### GetStockRunnerSetId

`func (o *CommonsRunnerInstanceConfiguration) GetStockRunnerSetId() string`

GetStockRunnerSetId returns the StockRunnerSetId field if non-nil, zero value otherwise.

### GetStockRunnerSetIdOk

`func (o *CommonsRunnerInstanceConfiguration) GetStockRunnerSetIdOk() (*string, bool)`

GetStockRunnerSetIdOk returns a tuple with the StockRunnerSetId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStockRunnerSetId

`func (o *CommonsRunnerInstanceConfiguration) SetStockRunnerSetId(v string)`

SetStockRunnerSetId sets StockRunnerSetId field to given value.

### HasStockRunnerSetId

`func (o *CommonsRunnerInstanceConfiguration) HasStockRunnerSetId() bool`

HasStockRunnerSetId returns a boolean if a field has been set.

### GetStorage

`func (o *CommonsRunnerInstanceConfiguration) GetStorage() CommonsStorage`

GetStorage returns the Storage field if non-nil, zero value otherwise.

### GetStorageOk

`func (o *CommonsRunnerInstanceConfiguration) GetStorageOk() (*CommonsStorage, bool)`

GetStorageOk returns a tuple with the Storage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStorage

`func (o *CommonsRunnerInstanceConfiguration) SetStorage(v CommonsStorage)`

SetStorage sets Storage field to given value.

### HasStorage

`func (o *CommonsRunnerInstanceConfiguration) HasStorage() bool`

HasStorage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


