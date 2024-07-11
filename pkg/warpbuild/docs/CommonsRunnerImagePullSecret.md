# CommonsRunnerImagePullSecret

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alias** | **string** |  | 
**Aws** | Pointer to [**CommonsRunnerImagePullSecretAWS**](CommonsRunnerImagePullSecretAWS.md) |  | [optional] 
**CreatedAt** | **string** |  | 
**DockerCredentials** | Pointer to [**CommonsRunnerImagePullSecretDockerCredentials**](CommonsRunnerImagePullSecretDockerCredentials.md) |  | [optional] 
**Id** | **string** |  | 
**OrganizationId** | **string** |  | 
**Type** | **string** |  | 
**UpdatedAt** | **string** |  | 

## Methods

### NewCommonsRunnerImagePullSecret

`func NewCommonsRunnerImagePullSecret(alias string, createdAt string, id string, organizationId string, type_ string, updatedAt string, ) *CommonsRunnerImagePullSecret`

NewCommonsRunnerImagePullSecret instantiates a new CommonsRunnerImagePullSecret object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImagePullSecretWithDefaults

`func NewCommonsRunnerImagePullSecretWithDefaults() *CommonsRunnerImagePullSecret`

NewCommonsRunnerImagePullSecretWithDefaults instantiates a new CommonsRunnerImagePullSecret object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlias

`func (o *CommonsRunnerImagePullSecret) GetAlias() string`

GetAlias returns the Alias field if non-nil, zero value otherwise.

### GetAliasOk

`func (o *CommonsRunnerImagePullSecret) GetAliasOk() (*string, bool)`

GetAliasOk returns a tuple with the Alias field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlias

`func (o *CommonsRunnerImagePullSecret) SetAlias(v string)`

SetAlias sets Alias field to given value.


### GetAws

`func (o *CommonsRunnerImagePullSecret) GetAws() CommonsRunnerImagePullSecretAWS`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *CommonsRunnerImagePullSecret) GetAwsOk() (*CommonsRunnerImagePullSecretAWS, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *CommonsRunnerImagePullSecret) SetAws(v CommonsRunnerImagePullSecretAWS)`

SetAws sets Aws field to given value.

### HasAws

`func (o *CommonsRunnerImagePullSecret) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetCreatedAt

`func (o *CommonsRunnerImagePullSecret) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *CommonsRunnerImagePullSecret) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *CommonsRunnerImagePullSecret) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetDockerCredentials

`func (o *CommonsRunnerImagePullSecret) GetDockerCredentials() CommonsRunnerImagePullSecretDockerCredentials`

GetDockerCredentials returns the DockerCredentials field if non-nil, zero value otherwise.

### GetDockerCredentialsOk

`func (o *CommonsRunnerImagePullSecret) GetDockerCredentialsOk() (*CommonsRunnerImagePullSecretDockerCredentials, bool)`

GetDockerCredentialsOk returns a tuple with the DockerCredentials field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDockerCredentials

`func (o *CommonsRunnerImagePullSecret) SetDockerCredentials(v CommonsRunnerImagePullSecretDockerCredentials)`

SetDockerCredentials sets DockerCredentials field to given value.

### HasDockerCredentials

`func (o *CommonsRunnerImagePullSecret) HasDockerCredentials() bool`

HasDockerCredentials returns a boolean if a field has been set.

### GetId

`func (o *CommonsRunnerImagePullSecret) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *CommonsRunnerImagePullSecret) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *CommonsRunnerImagePullSecret) SetId(v string)`

SetId sets Id field to given value.


### GetOrganizationId

`func (o *CommonsRunnerImagePullSecret) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *CommonsRunnerImagePullSecret) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *CommonsRunnerImagePullSecret) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.


### GetType

`func (o *CommonsRunnerImagePullSecret) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *CommonsRunnerImagePullSecret) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *CommonsRunnerImagePullSecret) SetType(v string)`

SetType sets Type field to given value.


### GetUpdatedAt

`func (o *CommonsRunnerImagePullSecret) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *CommonsRunnerImagePullSecret) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *CommonsRunnerImagePullSecret) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


