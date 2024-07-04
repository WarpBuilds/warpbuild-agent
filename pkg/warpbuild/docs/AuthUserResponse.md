# AuthUserResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | **string** |  | 
<<<<<<< HEAD
**Organization** | [**V1Organization**](V1Organization.md) |  | 
**RefreshToken** | **string** |  | 
**ShouldShowVcsIntegration** | **bool** |  | 
=======
**IsDifferentOrg** | **bool** |  | 
**Organization** | [**V1Organization**](V1Organization.md) |  | 
**RefreshToken** | **string** |  | 
>>>>>>> prajjwal-warp-323
**User** | [**V1User**](V1User.md) |  | 

## Methods

### NewAuthUserResponse

<<<<<<< HEAD
`func NewAuthUserResponse(accessToken string, organization V1Organization, refreshToken string, shouldShowVcsIntegration bool, user V1User, ) *AuthUserResponse`
=======
`func NewAuthUserResponse(accessToken string, isDifferentOrg bool, organization V1Organization, refreshToken string, user V1User, ) *AuthUserResponse`
>>>>>>> prajjwal-warp-323

NewAuthUserResponse instantiates a new AuthUserResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthUserResponseWithDefaults

`func NewAuthUserResponseWithDefaults() *AuthUserResponse`

NewAuthUserResponseWithDefaults instantiates a new AuthUserResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *AuthUserResponse) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *AuthUserResponse) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *AuthUserResponse) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


<<<<<<< HEAD
=======
### GetIsDifferentOrg

`func (o *AuthUserResponse) GetIsDifferentOrg() bool`

GetIsDifferentOrg returns the IsDifferentOrg field if non-nil, zero value otherwise.

### GetIsDifferentOrgOk

`func (o *AuthUserResponse) GetIsDifferentOrgOk() (*bool, bool)`

GetIsDifferentOrgOk returns a tuple with the IsDifferentOrg field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsDifferentOrg

`func (o *AuthUserResponse) SetIsDifferentOrg(v bool)`

SetIsDifferentOrg sets IsDifferentOrg field to given value.


>>>>>>> prajjwal-warp-323
### GetOrganization

`func (o *AuthUserResponse) GetOrganization() V1Organization`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *AuthUserResponse) GetOrganizationOk() (*V1Organization, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *AuthUserResponse) SetOrganization(v V1Organization)`

SetOrganization sets Organization field to given value.


### GetRefreshToken

`func (o *AuthUserResponse) GetRefreshToken() string`

GetRefreshToken returns the RefreshToken field if non-nil, zero value otherwise.

### GetRefreshTokenOk

`func (o *AuthUserResponse) GetRefreshTokenOk() (*string, bool)`

GetRefreshTokenOk returns a tuple with the RefreshToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRefreshToken

`func (o *AuthUserResponse) SetRefreshToken(v string)`

SetRefreshToken sets RefreshToken field to given value.


<<<<<<< HEAD
### GetShouldShowVcsIntegration

`func (o *AuthUserResponse) GetShouldShowVcsIntegration() bool`

GetShouldShowVcsIntegration returns the ShouldShowVcsIntegration field if non-nil, zero value otherwise.

### GetShouldShowVcsIntegrationOk

`func (o *AuthUserResponse) GetShouldShowVcsIntegrationOk() (*bool, bool)`

GetShouldShowVcsIntegrationOk returns a tuple with the ShouldShowVcsIntegration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShouldShowVcsIntegration

`func (o *AuthUserResponse) SetShouldShowVcsIntegration(v bool)`

SetShouldShowVcsIntegration sets ShouldShowVcsIntegration field to given value.


=======
>>>>>>> prajjwal-warp-323
### GetUser

`func (o *AuthUserResponse) GetUser() V1User`

GetUser returns the User field if non-nil, zero value otherwise.

### GetUserOk

`func (o *AuthUserResponse) GetUserOk() (*V1User, bool)`

GetUserOk returns a tuple with the User field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUser

`func (o *AuthUserResponse) SetUser(v V1User)`

SetUser sets User field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


