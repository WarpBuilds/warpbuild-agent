# CommonsCreateRunnerImagePullSecretInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alias** | **string** |  | 
**Aws** | Pointer to [**CommonsRunnerImagePullSecretAWS**](CommonsRunnerImagePullSecretAWS.md) |  | [optional] 
**DockerCredentials** | Pointer to [**CommonsRunnerImagePullSecretDockerCredentials**](CommonsRunnerImagePullSecretDockerCredentials.md) |  | [optional] 
**Type** | **string** |  | 

## Methods

### NewCommonsCreateRunnerImagePullSecretInput

`func NewCommonsCreateRunnerImagePullSecretInput(alias string, type_ string, ) *CommonsCreateRunnerImagePullSecretInput`

NewCommonsCreateRunnerImagePullSecretInput instantiates a new CommonsCreateRunnerImagePullSecretInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsCreateRunnerImagePullSecretInputWithDefaults

`func NewCommonsCreateRunnerImagePullSecretInputWithDefaults() *CommonsCreateRunnerImagePullSecretInput`

NewCommonsCreateRunnerImagePullSecretInputWithDefaults instantiates a new CommonsCreateRunnerImagePullSecretInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlias

`func (o *CommonsCreateRunnerImagePullSecretInput) GetAlias() string`

GetAlias returns the Alias field if non-nil, zero value otherwise.

### GetAliasOk

`func (o *CommonsCreateRunnerImagePullSecretInput) GetAliasOk() (*string, bool)`

GetAliasOk returns a tuple with the Alias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlias

`func (o *CommonsCreateRunnerImagePullSecretInput) SetAlias(v string)`

SetAlias sets Alias field to given value.


### GetAws

`func (o *CommonsCreateRunnerImagePullSecretInput) GetAws() CommonsRunnerImagePullSecretAWS`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *CommonsCreateRunnerImagePullSecretInput) GetAwsOk() (*CommonsRunnerImagePullSecretAWS, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *CommonsCreateRunnerImagePullSecretInput) SetAws(v CommonsRunnerImagePullSecretAWS)`

SetAws sets Aws field to given value.

### HasAws

`func (o *CommonsCreateRunnerImagePullSecretInput) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetDockerCredentials

`func (o *CommonsCreateRunnerImagePullSecretInput) GetDockerCredentials() CommonsRunnerImagePullSecretDockerCredentials`

GetDockerCredentials returns the DockerCredentials field if non-nil, zero value otherwise.

### GetDockerCredentialsOk

`func (o *CommonsCreateRunnerImagePullSecretInput) GetDockerCredentialsOk() (*CommonsRunnerImagePullSecretDockerCredentials, bool)`

GetDockerCredentialsOk returns a tuple with the DockerCredentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDockerCredentials

`func (o *CommonsCreateRunnerImagePullSecretInput) SetDockerCredentials(v CommonsRunnerImagePullSecretDockerCredentials)`

SetDockerCredentials sets DockerCredentials field to given value.

### HasDockerCredentials

`func (o *CommonsCreateRunnerImagePullSecretInput) HasDockerCredentials() bool`

HasDockerCredentials returns a boolean if a field has been set.

### GetType

`func (o *CommonsCreateRunnerImagePullSecretInput) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsCreateRunnerImagePullSecretInput) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsCreateRunnerImagePullSecretInput) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


