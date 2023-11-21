# SwitchOrganizationResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccessToken** | **string** |  | 
**Organization** | [**V1Organization**](V1Organization.md) |  | 

## Methods

### NewSwitchOrganizationResponse

`func NewSwitchOrganizationResponse(accessToken string, organization V1Organization, ) *SwitchOrganizationResponse`

NewSwitchOrganizationResponse instantiates a new SwitchOrganizationResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchOrganizationResponseWithDefaults

`func NewSwitchOrganizationResponseWithDefaults() *SwitchOrganizationResponse`

NewSwitchOrganizationResponseWithDefaults instantiates a new SwitchOrganizationResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccessToken

`func (o *SwitchOrganizationResponse) GetAccessToken() string`

GetAccessToken returns the AccessToken field if non-nil, zero value otherwise.

### GetAccessTokenOk

`func (o *SwitchOrganizationResponse) GetAccessTokenOk() (*string, bool)`

GetAccessTokenOk returns a tuple with the AccessToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessToken

`func (o *SwitchOrganizationResponse) SetAccessToken(v string)`

SetAccessToken sets AccessToken field to given value.


### GetOrganization

`func (o *SwitchOrganizationResponse) GetOrganization() V1Organization`

GetOrganization returns the Organization field if non-nil, zero value otherwise.

### GetOrganizationOk

`func (o *SwitchOrganizationResponse) GetOrganizationOk() (*V1Organization, bool)`

GetOrganizationOk returns a tuple with the Organization field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganization

`func (o *SwitchOrganizationResponse) SetOrganization(v V1Organization)`

SetOrganization sets Organization field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


