# VCSIntegration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AccountId** | Pointer to **string** |  | [optional] 
**AccountOwner** | Pointer to **string** |  | [optional] 
**ApplicationId** | Pointer to **string** |  | [optional] 
**ConfigurationUrl** | Pointer to **string** | ConfigurationURL  Github - Is the installation&#39;s settings page on the github account. - This field is only populated after the github account is successfully connected.  Gitlab - &lt;no-equivalent-exists&gt; | [optional] 
**ConnectionStatus** | **string** |  | 
**CreatedAt** | **string** |  | 
**Id** | **string** |  | 
**InstallationId** | Pointer to **string** |  | [optional] 
**IntegrationType** | **string** |  | 
**IntegrationUrl** | **string** |  | 
**OrganizationId** | Pointer to **string** |  | [optional] 
**Provider** | **string** |  | 
**TargetType** | **string** |  | 
**UpdatedAt** | **string** |  | 
**VcsUrl** | Pointer to **string** |  | [optional] 

## Methods

### NewVCSIntegration

`func NewVCSIntegration(connectionStatus string, createdAt string, id string, integrationType string, integrationUrl string, provider string, targetType string, updatedAt string, ) *VCSIntegration`

NewVCSIntegration instantiates a new VCSIntegration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVCSIntegrationWithDefaults

`func NewVCSIntegrationWithDefaults() *VCSIntegration`

NewVCSIntegrationWithDefaults instantiates a new VCSIntegration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAccountId

`func (o *VCSIntegration) GetAccountId() string`

GetAccountId returns the AccountId field if non-nil, zero value otherwise.

### GetAccountIdOk

`func (o *VCSIntegration) GetAccountIdOk() (*string, bool)`

GetAccountIdOk returns a tuple with the AccountId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountId

`func (o *VCSIntegration) SetAccountId(v string)`

SetAccountId sets AccountId field to given value.

### HasAccountId

`func (o *VCSIntegration) HasAccountId() bool`

HasAccountId returns a boolean if a field has been set.

### GetAccountOwner

`func (o *VCSIntegration) GetAccountOwner() string`

GetAccountOwner returns the AccountOwner field if non-nil, zero value otherwise.

### GetAccountOwnerOk

`func (o *VCSIntegration) GetAccountOwnerOk() (*string, bool)`

GetAccountOwnerOk returns a tuple with the AccountOwner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccountOwner

`func (o *VCSIntegration) SetAccountOwner(v string)`

SetAccountOwner sets AccountOwner field to given value.

### HasAccountOwner

`func (o *VCSIntegration) HasAccountOwner() bool`

HasAccountOwner returns a boolean if a field has been set.

### GetApplicationId

`func (o *VCSIntegration) GetApplicationId() string`

GetApplicationId returns the ApplicationId field if non-nil, zero value otherwise.

### GetApplicationIdOk

`func (o *VCSIntegration) GetApplicationIdOk() (*string, bool)`

GetApplicationIdOk returns a tuple with the ApplicationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetApplicationId

`func (o *VCSIntegration) SetApplicationId(v string)`

SetApplicationId sets ApplicationId field to given value.

### HasApplicationId

`func (o *VCSIntegration) HasApplicationId() bool`

HasApplicationId returns a boolean if a field has been set.

### GetConfigurationUrl

`func (o *VCSIntegration) GetConfigurationUrl() string`

GetConfigurationUrl returns the ConfigurationUrl field if non-nil, zero value otherwise.

### GetConfigurationUrlOk

`func (o *VCSIntegration) GetConfigurationUrlOk() (*string, bool)`

GetConfigurationUrlOk returns a tuple with the ConfigurationUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigurationUrl

`func (o *VCSIntegration) SetConfigurationUrl(v string)`

SetConfigurationUrl sets ConfigurationUrl field to given value.

### HasConfigurationUrl

`func (o *VCSIntegration) HasConfigurationUrl() bool`

HasConfigurationUrl returns a boolean if a field has been set.

### GetConnectionStatus

`func (o *VCSIntegration) GetConnectionStatus() string`

GetConnectionStatus returns the ConnectionStatus field if non-nil, zero value otherwise.

### GetConnectionStatusOk

`func (o *VCSIntegration) GetConnectionStatusOk() (*string, bool)`

GetConnectionStatusOk returns a tuple with the ConnectionStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionStatus

`func (o *VCSIntegration) SetConnectionStatus(v string)`

SetConnectionStatus sets ConnectionStatus field to given value.


### GetCreatedAt

`func (o *VCSIntegration) GetCreatedAt() string`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *VCSIntegration) GetCreatedAtOk() (*string, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *VCSIntegration) SetCreatedAt(v string)`

SetCreatedAt sets CreatedAt field to given value.


### GetId

`func (o *VCSIntegration) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *VCSIntegration) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *VCSIntegration) SetId(v string)`

SetId sets Id field to given value.


### GetInstallationId

`func (o *VCSIntegration) GetInstallationId() string`

GetInstallationId returns the InstallationId field if non-nil, zero value otherwise.

### GetInstallationIdOk

`func (o *VCSIntegration) GetInstallationIdOk() (*string, bool)`

GetInstallationIdOk returns a tuple with the InstallationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationId

`func (o *VCSIntegration) SetInstallationId(v string)`

SetInstallationId sets InstallationId field to given value.

### HasInstallationId

`func (o *VCSIntegration) HasInstallationId() bool`

HasInstallationId returns a boolean if a field has been set.

### GetIntegrationType

`func (o *VCSIntegration) GetIntegrationType() string`

GetIntegrationType returns the IntegrationType field if non-nil, zero value otherwise.

### GetIntegrationTypeOk

`func (o *VCSIntegration) GetIntegrationTypeOk() (*string, bool)`

GetIntegrationTypeOk returns a tuple with the IntegrationType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrationType

`func (o *VCSIntegration) SetIntegrationType(v string)`

SetIntegrationType sets IntegrationType field to given value.


### GetIntegrationUrl

`func (o *VCSIntegration) GetIntegrationUrl() string`

GetIntegrationUrl returns the IntegrationUrl field if non-nil, zero value otherwise.

### GetIntegrationUrlOk

`func (o *VCSIntegration) GetIntegrationUrlOk() (*string, bool)`

GetIntegrationUrlOk returns a tuple with the IntegrationUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegrationUrl

`func (o *VCSIntegration) SetIntegrationUrl(v string)`

SetIntegrationUrl sets IntegrationUrl field to given value.


### GetOrganizationId

`func (o *VCSIntegration) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *VCSIntegration) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *VCSIntegration) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.

### HasOrganizationId

`func (o *VCSIntegration) HasOrganizationId() bool`

HasOrganizationId returns a boolean if a field has been set.

### GetProvider

`func (o *VCSIntegration) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *VCSIntegration) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *VCSIntegration) SetProvider(v string)`

SetProvider sets Provider field to given value.


### GetTargetType

`func (o *VCSIntegration) GetTargetType() string`

GetTargetType returns the TargetType field if non-nil, zero value otherwise.

### GetTargetTypeOk

`func (o *VCSIntegration) GetTargetTypeOk() (*string, bool)`

GetTargetTypeOk returns a tuple with the TargetType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetType

`func (o *VCSIntegration) SetTargetType(v string)`

SetTargetType sets TargetType field to given value.


### GetUpdatedAt

`func (o *VCSIntegration) GetUpdatedAt() string`

GetUpdatedAt returns the UpdatedAt field if non-nil, zero value otherwise.

### GetUpdatedAtOk

`func (o *VCSIntegration) GetUpdatedAtOk() (*string, bool)`

GetUpdatedAtOk returns a tuple with the UpdatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdatedAt

`func (o *VCSIntegration) SetUpdatedAt(v string)`

SetUpdatedAt sets UpdatedAt field to given value.


### GetVcsUrl

`func (o *VCSIntegration) GetVcsUrl() string`

GetVcsUrl returns the VcsUrl field if non-nil, zero value otherwise.

### GetVcsUrlOk

`func (o *VCSIntegration) GetVcsUrlOk() (*string, bool)`

GetVcsUrlOk returns a tuple with the VcsUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVcsUrl

`func (o *VCSIntegration) SetVcsUrl(v string)`

SetVcsUrl sets VcsUrl field to given value.

### HasVcsUrl

`func (o *VCSIntegration) HasVcsUrl() bool`

HasVcsUrl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


