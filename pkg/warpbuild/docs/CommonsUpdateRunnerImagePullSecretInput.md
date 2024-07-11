# CommonsUpdateRunnerImagePullSecretInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Aws** | Pointer to [**CommonsRunnerImagePullSecretAWS**](CommonsRunnerImagePullSecretAWS.md) |  | [optional] 
**DockerCredentials** | Pointer to [**CommonsRunnerImagePullSecretDockerCredentials**](CommonsRunnerImagePullSecretDockerCredentials.md) |  | [optional] 

## Methods

### NewCommonsUpdateRunnerImagePullSecretInput

`func NewCommonsUpdateRunnerImagePullSecretInput() *CommonsUpdateRunnerImagePullSecretInput`

NewCommonsUpdateRunnerImagePullSecretInput instantiates a new CommonsUpdateRunnerImagePullSecretInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsUpdateRunnerImagePullSecretInputWithDefaults

`func NewCommonsUpdateRunnerImagePullSecretInputWithDefaults() *CommonsUpdateRunnerImagePullSecretInput`

NewCommonsUpdateRunnerImagePullSecretInputWithDefaults instantiates a new CommonsUpdateRunnerImagePullSecretInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAws

`func (o *CommonsUpdateRunnerImagePullSecretInput) GetAws() CommonsRunnerImagePullSecretAWS`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *CommonsUpdateRunnerImagePullSecretInput) GetAwsOk() (*CommonsRunnerImagePullSecretAWS, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *CommonsUpdateRunnerImagePullSecretInput) SetAws(v CommonsRunnerImagePullSecretAWS)`

SetAws sets Aws field to given value.

### HasAws

`func (o *CommonsUpdateRunnerImagePullSecretInput) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetDockerCredentials

`func (o *CommonsUpdateRunnerImagePullSecretInput) GetDockerCredentials() CommonsRunnerImagePullSecretDockerCredentials`

GetDockerCredentials returns the DockerCredentials field if non-nil, zero value otherwise.

### GetDockerCredentialsOk

`func (o *CommonsUpdateRunnerImagePullSecretInput) GetDockerCredentialsOk() (*CommonsRunnerImagePullSecretDockerCredentials, bool)`

GetDockerCredentialsOk returns a tuple with the DockerCredentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDockerCredentials

`func (o *CommonsUpdateRunnerImagePullSecretInput) SetDockerCredentials(v CommonsRunnerImagePullSecretDockerCredentials)`

SetDockerCredentials sets DockerCredentials field to given value.

### HasDockerCredentials

`func (o *CommonsUpdateRunnerImagePullSecretInput) HasDockerCredentials() bool`

HasDockerCredentials returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


