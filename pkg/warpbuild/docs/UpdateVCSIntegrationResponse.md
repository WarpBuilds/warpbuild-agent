# UpdateVCSIntegrationResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NewSession** | Pointer to [**AuthUserResponse**](AuthUserResponse.md) |  | [optional] 
**VcsIntegration** | Pointer to [**VCSIntegration**](VCSIntegration.md) |  | [optional] 

## Methods

### NewUpdateVCSIntegrationResponse

`func NewUpdateVCSIntegrationResponse() *UpdateVCSIntegrationResponse`

NewUpdateVCSIntegrationResponse instantiates a new UpdateVCSIntegrationResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUpdateVCSIntegrationResponseWithDefaults

`func NewUpdateVCSIntegrationResponseWithDefaults() *UpdateVCSIntegrationResponse`

NewUpdateVCSIntegrationResponseWithDefaults instantiates a new UpdateVCSIntegrationResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNewSession

`func (o *UpdateVCSIntegrationResponse) GetNewSession() AuthUserResponse`

GetNewSession returns the NewSession field if non-nil, zero value otherwise.

### GetNewSessionOk

`func (o *UpdateVCSIntegrationResponse) GetNewSessionOk() (*AuthUserResponse, bool)`

GetNewSessionOk returns a tuple with the NewSession field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNewSession

`func (o *UpdateVCSIntegrationResponse) SetNewSession(v AuthUserResponse)`

SetNewSession sets NewSession field to given value.

### HasNewSession

`func (o *UpdateVCSIntegrationResponse) HasNewSession() bool`

HasNewSession returns a boolean if a field has been set.

### GetVcsIntegration

`func (o *UpdateVCSIntegrationResponse) GetVcsIntegration() VCSIntegration`

GetVcsIntegration returns the VcsIntegration field if non-nil, zero value otherwise.

### GetVcsIntegrationOk

`func (o *UpdateVCSIntegrationResponse) GetVcsIntegrationOk() (*VCSIntegration, bool)`

GetVcsIntegrationOk returns a tuple with the VcsIntegration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsIntegration

`func (o *UpdateVCSIntegrationResponse) SetVcsIntegration(v VCSIntegration)`

SetVcsIntegration sets VcsIntegration field to given value.

### HasVcsIntegration

`func (o *UpdateVCSIntegrationResponse) HasVcsIntegration() bool`

HasVcsIntegration returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


