# CommonsRunnerImageHook

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**File** | Pointer to **string** | File is the base64 encoded file. This must be a shell script which is encoded in base64. | [optional] 
**Type** | Pointer to **string** | Type is the type of hook. | [optional] 

## Methods

### NewCommonsRunnerImageHook

`func NewCommonsRunnerImageHook() *CommonsRunnerImageHook`

NewCommonsRunnerImageHook instantiates a new CommonsRunnerImageHook object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImageHookWithDefaults

`func NewCommonsRunnerImageHookWithDefaults() *CommonsRunnerImageHook`

NewCommonsRunnerImageHookWithDefaults instantiates a new CommonsRunnerImageHook object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFile

`func (o *CommonsRunnerImageHook) GetFile() string`

GetFile returns the File field if non-nil, zero value otherwise.

### GetFileOk

`func (o *CommonsRunnerImageHook) GetFileOk() (*string, bool)`

GetFileOk returns a tuple with the File field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFile

`func (o *CommonsRunnerImageHook) SetFile(v string)`

SetFile sets File field to given value.

### HasFile

`func (o *CommonsRunnerImageHook) HasFile() bool`

HasFile returns a boolean if a field has been set.

### GetType

`func (o *CommonsRunnerImageHook) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsRunnerImageHook) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsRunnerImageHook) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *CommonsRunnerImageHook) HasType() bool`

HasType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


