# CommonsStorage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DiskType** | Pointer to **string** | DiskType is the disk type associated with the instance.  This is only applicable for GCP based runners for now. Passing this for other providers will have no effect.  Possible values:  - pd-standard: Standard persistent disk.  - pd-ssd: SSD persistent disk.  - pd-balanced: Balanced persistent disk.  - pd-extreme: Extreme persistent disk. (to be added later)  - hyperdisk-balanced: Hyperdisk balanced persistent disk.  - hyperdisk-extreme: Hyperdisk extreme persistent disk.  Refer: https://cloud.google.com/compute/docs/disks#disk-types  Default value is automatically picked if nothing is passed in case of GCP.  +Default: pd-ssd | [optional] 
**Iops** | Pointer to **int32** | IOPS is the IOPS of the storage.  For GCP, This is not applicable. Any passed value will be ignored. | [optional] 
**PerformanceTier** | Pointer to **string** | PerformanceTier is the provider specific performance tier of the storage. This is applicable for Azure as of now, can be extended to other providers in the future. Passing this for other providers will have no effect.  Refer: https://learn.microsoft.com/en-us/azure/virtual-machines/disks-change-performance  Possible values:  - P15 - P20 - P30 - P40 - P50  Default value is automatically picked if nothing is passed in case of Azure.  +Default: P15 | [optional] 
**Size** | Pointer to **int32** |  | [optional] 
**Throughput** | Pointer to **int32** | Throughput is the throughput of the storage.  For GCP, This is not applicable. Any passed value will be ignored. | [optional] 
**Tier** | Pointer to **string** | Tier is the tier of the storage.  If GCP based, you must set this to custom. | [optional] 

## Methods

### NewCommonsStorage

`func NewCommonsStorage() *CommonsStorage`

NewCommonsStorage instantiates a new CommonsStorage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsStorageWithDefaults

`func NewCommonsStorageWithDefaults() *CommonsStorage`

NewCommonsStorageWithDefaults instantiates a new CommonsStorage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDiskType

`func (o *CommonsStorage) GetDiskType() string`

GetDiskType returns the DiskType field if non-nil, zero value otherwise.

### GetDiskTypeOk

`func (o *CommonsStorage) GetDiskTypeOk() (*string, bool)`

GetDiskTypeOk returns a tuple with the DiskType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiskType

`func (o *CommonsStorage) SetDiskType(v string)`

SetDiskType sets DiskType field to given value.

### HasDiskType

`func (o *CommonsStorage) HasDiskType() bool`

HasDiskType returns a boolean if a field has been set.

### GetIops

`func (o *CommonsStorage) GetIops() int32`

GetIops returns the Iops field if non-nil, zero value otherwise.

### GetIopsOk

`func (o *CommonsStorage) GetIopsOk() (*int32, bool)`

GetIopsOk returns a tuple with the Iops field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIops

`func (o *CommonsStorage) SetIops(v int32)`

SetIops sets Iops field to given value.

### HasIops

`func (o *CommonsStorage) HasIops() bool`

HasIops returns a boolean if a field has been set.

### GetPerformanceTier

`func (o *CommonsStorage) GetPerformanceTier() string`

GetPerformanceTier returns the PerformanceTier field if non-nil, zero value otherwise.

### GetPerformanceTierOk

`func (o *CommonsStorage) GetPerformanceTierOk() (*string, bool)`

GetPerformanceTierOk returns a tuple with the PerformanceTier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPerformanceTier

`func (o *CommonsStorage) SetPerformanceTier(v string)`

SetPerformanceTier sets PerformanceTier field to given value.

### HasPerformanceTier

`func (o *CommonsStorage) HasPerformanceTier() bool`

HasPerformanceTier returns a boolean if a field has been set.

### GetSize

`func (o *CommonsStorage) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *CommonsStorage) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *CommonsStorage) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *CommonsStorage) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetThroughput

`func (o *CommonsStorage) GetThroughput() int32`

GetThroughput returns the Throughput field if non-nil, zero value otherwise.

### GetThroughputOk

`func (o *CommonsStorage) GetThroughputOk() (*int32, bool)`

GetThroughputOk returns a tuple with the Throughput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThroughput

`func (o *CommonsStorage) SetThroughput(v int32)`

SetThroughput sets Throughput field to given value.

### HasThroughput

`func (o *CommonsStorage) HasThroughput() bool`

HasThroughput returns a boolean if a field has been set.

### GetTier

`func (o *CommonsStorage) GetTier() string`

GetTier returns the Tier field if non-nil, zero value otherwise.

### GetTierOk

`func (o *CommonsStorage) GetTierOk() (*string, bool)`

GetTierOk returns a tuple with the Tier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTier

`func (o *CommonsStorage) SetTier(v string)`

SetTier sets Tier field to given value.

### HasTier

`func (o *CommonsStorage) HasTier() bool`

HasTier returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


