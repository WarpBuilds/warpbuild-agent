# UpdateVCSIntegrationRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **string** | Code is  - &#39;code&#39; param in the callback from gitlab | [optional] 
**Id** | **string** |  | 
**InstallationId** | Pointer to **string** | InstallationID  - &#39;installation_id&#39; param from github installation | [optional] 
**SetupAction** | Pointer to **string** | SetupAction  - &#39;setup_action&#39; param from github installation | [optional] 

## Methods

### NewUpdateVCSIntegrationRequest

`func NewUpdateVCSIntegrationRequest(id string, ) *UpdateVCSIntegrationRequest`

NewUpdateVCSIntegrationRequest instantiates a new UpdateVCSIntegrationRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateVCSIntegrationRequestWithDefaults

`func NewUpdateVCSIntegrationRequestWithDefaults() *UpdateVCSIntegrationRequest`

NewUpdateVCSIntegrationRequestWithDefaults instantiates a new UpdateVCSIntegrationRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *UpdateVCSIntegrationRequest) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *UpdateVCSIntegrationRequest) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *UpdateVCSIntegrationRequest) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *UpdateVCSIntegrationRequest) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetId

`func (o *UpdateVCSIntegrationRequest) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *UpdateVCSIntegrationRequest) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *UpdateVCSIntegrationRequest) SetId(v string)`

SetId sets Id field to given value.


### GetInstallationId

`func (o *UpdateVCSIntegrationRequest) GetInstallationId() string`

GetInstallationId returns the InstallationId field if non-nil, zero value otherwise.

### GetInstallationIdOk

`func (o *UpdateVCSIntegrationRequest) GetInstallationIdOk() (*string, bool)`

GetInstallationIdOk returns a tuple with the InstallationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationId

`func (o *UpdateVCSIntegrationRequest) SetInstallationId(v string)`

SetInstallationId sets InstallationId field to given value.

### HasInstallationId

`func (o *UpdateVCSIntegrationRequest) HasInstallationId() bool`

HasInstallationId returns a boolean if a field has been set.

### GetSetupAction

`func (o *UpdateVCSIntegrationRequest) GetSetupAction() string`

GetSetupAction returns the SetupAction field if non-nil, zero value otherwise.

### GetSetupActionOk

`func (o *UpdateVCSIntegrationRequest) GetSetupActionOk() (*string, bool)`

GetSetupActionOk returns a tuple with the SetupAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSetupAction

`func (o *UpdateVCSIntegrationRequest) SetSetupAction(v string)`

SetSetupAction sets SetupAction field to given value.

### HasSetupAction

`func (o *UpdateVCSIntegrationRequest) HasSetupAction() bool`

HasSetupAction returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


