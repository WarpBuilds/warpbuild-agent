# MeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
<<<<<<< HEAD
=======
**Extras** | Pointer to **map[string]string** |  | [optional] 
>>>>>>> prajjwal-warp-323
**Organization** | [**V1Organization**](V1Organization.md) |  | 
**User** | [**V1User**](V1User.md) |  | 
**VcsIntegration** | [**CommonsVCSIntegrationLean**](CommonsVCSIntegrationLean.md) |  | 

## Methods

### NewMeResponse

`func NewMeResponse(organization V1Organization, user V1User, vcsIntegration CommonsVCSIntegrationLean, ) *MeResponse`

NewMeResponse instantiates a new MeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMeResponseWithDefaults

`func NewMeResponseWithDefaults() *MeResponse`

NewMeResponseWithDefaults instantiates a new MeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

<<<<<<< HEAD
=======
### GetExtras

`func (o *MeResponse) GetExtras() map[string]string`

GetExtras returns the Extras field if non-nil, zero value otherwise.

### GetExtrasOk

`func (o *MeResponse) GetExtrasOk() (*map[string]string, bool)`

GetExtrasOk returns a tuple with the Extras field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtras

`func (o *MeResponse) SetExtras(v map[string]string)`

SetExtras sets Extras field to given value.

### HasExtras

`func (o *MeResponse) HasExtras() bool`

HasExtras returns a boolean if a field has been set.

>>>>>>> prajjwal-warp-323
### GetOrganization

`func (o *MeResponse) GetOrganization() V1Organization`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *MeResponse) GetOrganizationOk() (*V1Organization, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *MeResponse) SetOrganization(v V1Organization)`

SetOrganization sets Organization field to given value.


### GetUser

`func (o *MeResponse) GetUser() V1User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *MeResponse) GetUserOk() (*V1User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *MeResponse) SetUser(v V1User)`

SetUser sets User field to given value.


### GetVcsIntegration

`func (o *MeResponse) GetVcsIntegration() CommonsVCSIntegrationLean`

GetVcsIntegration returns the VcsIntegration field if non-nil, zero value otherwise.

### GetVcsIntegrationOk

`func (o *MeResponse) GetVcsIntegrationOk() (*CommonsVCSIntegrationLean, bool)`

GetVcsIntegrationOk returns a tuple with the VcsIntegration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsIntegration

`func (o *MeResponse) SetVcsIntegration(v CommonsVCSIntegrationLean)`

SetVcsIntegration sets VcsIntegration field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


