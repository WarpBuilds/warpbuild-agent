# CommonsRunnerImagePullSecretAWS

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessKeyId** | Pointer to **string** |  | [optional] 
**AwsEcrRepository** | Pointer to **string** | AWSECRRepository is the short name of the ecr repository For example, if the complete uri for an image is &lt;account_id&gt;.dkr.ecr.&lt;region&gt;.amazonaws.com/acme/customrunners:v1.5.0 The AWS ECR Repo is &#x60;acme/customrunners&#x60; | [optional] 
**Region** | Pointer to **string** |  | [optional] 
**SecretAccessKey** | Pointer to **string** |  | [optional] 

## Methods

### NewCommonsRunnerImagePullSecretAWS

`func NewCommonsRunnerImagePullSecretAWS() *CommonsRunnerImagePullSecretAWS`

NewCommonsRunnerImagePullSecretAWS instantiates a new CommonsRunnerImagePullSecretAWS object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommonsRunnerImagePullSecretAWSWithDefaults

`func NewCommonsRunnerImagePullSecretAWSWithDefaults() *CommonsRunnerImagePullSecretAWS`

NewCommonsRunnerImagePullSecretAWSWithDefaults instantiates a new CommonsRunnerImagePullSecretAWS object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessKeyId

`func (o *CommonsRunnerImagePullSecretAWS) GetAccessKeyId() string`

GetAccessKeyId returns the AccessKeyId field if non-nil, zero value otherwise.

### GetAccessKeyIdOk

`func (o *CommonsRunnerImagePullSecretAWS) GetAccessKeyIdOk() (*string, bool)`

GetAccessKeyIdOk returns a tuple with the AccessKeyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessKeyId

`func (o *CommonsRunnerImagePullSecretAWS) SetAccessKeyId(v string)`

SetAccessKeyId sets AccessKeyId field to given value.

### HasAccessKeyId

`func (o *CommonsRunnerImagePullSecretAWS) HasAccessKeyId() bool`

HasAccessKeyId returns a boolean if a field has been set.

### GetAwsEcrRepository

`func (o *CommonsRunnerImagePullSecretAWS) GetAwsEcrRepository() string`

GetAwsEcrRepository returns the AwsEcrRepository field if non-nil, zero value otherwise.

### GetAwsEcrRepositoryOk

`func (o *CommonsRunnerImagePullSecretAWS) GetAwsEcrRepositoryOk() (*string, bool)`

GetAwsEcrRepositoryOk returns a tuple with the AwsEcrRepository field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAwsEcrRepository

`func (o *CommonsRunnerImagePullSecretAWS) SetAwsEcrRepository(v string)`

SetAwsEcrRepository sets AwsEcrRepository field to given value.

### HasAwsEcrRepository

`func (o *CommonsRunnerImagePullSecretAWS) HasAwsEcrRepository() bool`

HasAwsEcrRepository returns a boolean if a field has been set.

### GetRegion

`func (o *CommonsRunnerImagePullSecretAWS) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *CommonsRunnerImagePullSecretAWS) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *CommonsRunnerImagePullSecretAWS) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *CommonsRunnerImagePullSecretAWS) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetSecretAccessKey

`func (o *CommonsRunnerImagePullSecretAWS) GetSecretAccessKey() string`

GetSecretAccessKey returns the SecretAccessKey field if non-nil, zero value otherwise.

### GetSecretAccessKeyOk

`func (o *CommonsRunnerImagePullSecretAWS) GetSecretAccessKeyOk() (*string, bool)`

GetSecretAccessKeyOk returns a tuple with the SecretAccessKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecretAccessKey

`func (o *CommonsRunnerImagePullSecretAWS) SetSecretAccessKey(v string)`

SetSecretAccessKey sets SecretAccessKey field to given value.

### HasSecretAccessKey

`func (o *CommonsRunnerImagePullSecretAWS) HasSecretAccessKey() bool`

HasSecretAccessKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


