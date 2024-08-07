# Go API client for warpbuild

This is the docs for warp builds api for argonaut

## Overview
This API client was generated by the [OpenAPI Generator](https://openapi-generator.tech) project.  By using the [OpenAPI-spec](https://www.openapis.org/) from a remote server, you can easily generate an API client.

- API version: 0.4.0
- Package version: /Users/prashant/warpbuilds/warpbuild-agent/pkg/warpbuild
- Build package: org.openapitools.codegen.languages.GoClientCodegen
For more information, please visit [http://www.swagger.io/support](http://www.swagger.io/support)

## Installation

Install the following dependencies:

```shell
go get github.com/stretchr/testify/assert
go get golang.org/x/net/context
```

Put the package under your project folder and add the following in import:

```golang
import warpbuild "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
```

To use a proxy, set the environment variable `HTTP_PROXY`:

```golang
os.Setenv("HTTP_PROXY", "http://proxy_name:proxy_port")
```

## Configuration of Server URL

Default configuration comes with `Servers` field that contains server objects as defined in the OpenAPI specification.

### Select Server Configuration

For using other server than the one defined on index 0 set context value `sw.ContextServerIndex` of type `int`.

```golang
ctx := context.WithValue(context.Background(), warpbuild.ContextServerIndex, 1)
```

### Templated Server URL

Templated server URL is formatted using default variables from configuration or from context value `sw.ContextServerVariables` of type `map[string]string`.

```golang
ctx := context.WithValue(context.Background(), warpbuild.ContextServerVariables, map[string]string{
	"basePath": "v2",
})
```

Note, enum values are always validated and all unused variables are silently ignored.

### URLs Configuration per Operation

Each operation can use different server URL defined using `OperationServers` map in the `Configuration`.
An operation is uniquely identified by `"{classname}Service.{nickname}"` string.
Similar rules for overriding default operation server index and variables applies by using `sw.ContextOperationServerIndices` and `sw.ContextOperationServerVariables` context maps.

```golang
ctx := context.WithValue(context.Background(), warpbuild.ContextOperationServerIndices, map[string]int{
	"{classname}Service.{nickname}": 2,
})
ctx = context.WithValue(context.Background(), warpbuild.ContextOperationServerVariables, map[string]map[string]string{
	"{classname}Service.{nickname}": {
		"port": "8443",
	},
})
```

## Documentation for API Endpoints

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*V1AuthAPI* | [**AuthTokensGet**](docs/V1AuthAPI.md#authtokensget) | **Get** /auth/tokens | List user tokens
*V1AuthAPI* | [**AuthUser**](docs/V1AuthAPI.md#authuser) | **Post** /auth | Auth user
*V1AuthAPI* | [**AuthUsersGet**](docs/V1AuthAPI.md#authusersget) | **Get** /auth/users | List users
*V1AuthAPI* | [**GetAuthURL**](docs/V1AuthAPI.md#getauthurl) | **Get** /auth/login/{provider} | Get auth url
*V1AuthAPI* | [**GetMe**](docs/V1AuthAPI.md#getme) | **Get** /auth/me | Auth user
*V1AuthAPI* | [**Logout**](docs/V1AuthAPI.md#logout) | **Patch** /auth/logout | Logout
*V1AuthAPI* | [**RefreshToken**](docs/V1AuthAPI.md#refreshtoken) | **Patch** /auth/token/refresh | Refresh token
*V1AuthAPI* | [**SwitchOrganization**](docs/V1AuthAPI.md#switchorganization) | **Patch** /auth/switch | Switch organization
*V1BillingAPI* | [**PostUsageForInternalService**](docs/V1BillingAPI.md#postusageforinternalservice) | **Post** /billing/usage/internal | Post Usage for internal service
*V1DebuggerAPI* | [**DebugPublishEvent**](docs/V1DebuggerAPI.md#debugpublishevent) | **Post** /debugger/events/publish | Publish an event to the event bus
*V1InsightsIntegrationsAPI* | [**GitHubCallback**](docs/V1InsightsIntegrationsAPI.md#githubcallback) | **Post** /insights/integrations/github/callback | GitHub callback for insights
*V1JobsAPI* | [**GetCostSummary**](docs/V1JobsAPI.md#getcostsummary) | **Get** /jobs/cost-summary | GetCostSummary
*V1JobsAPI* | [**GetDaywiseCosts**](docs/V1JobsAPI.md#getdaywisecosts) | **Get** /jobs/daywise-costs | GetDaywiseCosts
*V1OrganizationAPI* | [**CreateOrganization**](docs/V1OrganizationAPI.md#createorganization) | **Post** /organization | Adds a new organisation for a current user
*V1OrganizationAPI* | [**GetOrganization**](docs/V1OrganizationAPI.md#getorganization) | **Get** /organization | Get organization details for the current organization. Current organization is figured from the authorization token
*V1OrganizationAPI* | [**ListOrgUsers**](docs/V1OrganizationAPI.md#listorgusers) | **Get** /organization/users | ListOrgUsers list the users for the current organization
*V1OrganizationAPI* | [**ListUserOrganizations**](docs/V1OrganizationAPI.md#listuserorganizations) | **Get** /organizations | ListUserOrganizations lists all the organization user has access to.
*V1OrganizationAPI* | [**UpdateOrganization**](docs/V1OrganizationAPI.md#updateorganization) | **Patch** /organization | Updates existing organization based on the fields provided.
*V1RunnerImagePullSecretsAPI* | [**CreateRunnerImagePullSecret**](docs/V1RunnerImagePullSecretsAPI.md#createrunnerimagepullsecret) | **Post** /runner-image-pull-secrets | Create a new runner image pull secret.
*V1RunnerImagePullSecretsAPI* | [**DeleteRunnerImagePullSecret**](docs/V1RunnerImagePullSecretsAPI.md#deleterunnerimagepullsecret) | **Delete** /runner-image-pull-secrets/{id} | Delete runner image pull secret details for the id.
*V1RunnerImagePullSecretsAPI* | [**GetRunnerImagePullSecret**](docs/V1RunnerImagePullSecretsAPI.md#getrunnerimagepullsecret) | **Get** /runner-image-pull-secrets/{id} | Get runner image pull secret details for the id.
*V1RunnerImagePullSecretsAPI* | [**ListRunnerImagePullSecrets**](docs/V1RunnerImagePullSecretsAPI.md#listrunnerimagepullsecrets) | **Get** /runner-image-pull-secrets | List all runner image pull secrets.
*V1RunnerImagePullSecretsAPI* | [**UpdateRunnerImagePullSecret**](docs/V1RunnerImagePullSecretsAPI.md#updaterunnerimagepullsecret) | **Put** /runner-image-pull-secrets/{id} | Update runner image pull secret details for the id.
*V1RunnerImageVersionsAPI* | [**DeleteRunnerImageVersion**](docs/V1RunnerImageVersionsAPI.md#deleterunnerimageversion) | **Delete** /runner-image-versions/{id} | Delete runner image version details for the id.
*V1RunnerImageVersionsAPI* | [**GetRunnerImageVersion**](docs/V1RunnerImageVersionsAPI.md#getrunnerimageversion) | **Get** /runner-image-versions/{id} | Get runner image version details for the id.
*V1RunnerImageVersionsAPI* | [**ListRunnerImageVersions**](docs/V1RunnerImageVersionsAPI.md#listrunnerimageversions) | **Get** /runner-image-versions | List all runner image versions.
*V1RunnerImageVersionsAPI* | [**UpdateRunnerImageVersion**](docs/V1RunnerImageVersionsAPI.md#updaterunnerimageversion) | **Patch** /runner-image-versions/{id} | Update runner image version details for the id.
*V1RunnerImagesAPI* | [**CreateRunnerImage**](docs/V1RunnerImagesAPI.md#createrunnerimage) | **Post** /runner-images | Create a new runner image.
*V1RunnerImagesAPI* | [**DeleteRunnerImage**](docs/V1RunnerImagesAPI.md#deleterunnerimage) | **Delete** /runner-images/{id} | Delete runner image details for the id.
*V1RunnerImagesAPI* | [**GetRunnerImage**](docs/V1RunnerImagesAPI.md#getrunnerimage) | **Get** /runner-images/{id} | Get runner image details for the id.
*V1RunnerImagesAPI* | [**ListRunnerImages**](docs/V1RunnerImagesAPI.md#listrunnerimages) | **Get** /runner-images | List all runner images.
*V1RunnerImagesAPI* | [**UpdateRunnerImage**](docs/V1RunnerImagesAPI.md#updaterunnerimage) | **Put** /runner-images/{id} | Update runner image details for the id.
*V1RunnerInstanceAPI* | [**AddRunnerInstance**](docs/V1RunnerInstanceAPI.md#addrunnerinstance) | **Post** /runners_instance | Add a new runner instance
*V1RunnerInstanceAPI* | [**GetRunnerInstanceAllocationDetails**](docs/V1RunnerInstanceAPI.md#getrunnerinstanceallocationdetails) | **Get** /runners_instance/{id}/allocation_details | Get runner instance allocation details for the id
*V1RunnerInstanceAPI* | [**GetRunnerInstancePresignedLogUploadURL**](docs/V1RunnerInstanceAPI.md#getrunnerinstancepresignedloguploadurl) | **Get** /runners_instance/{id}/presigned_log_upload_url | Gets a presigned url for uploading logs for a runner instance
*V1RunnerInstanceAPI* | [**GetRunnerLastJobProcessedMeta**](docs/V1RunnerInstanceAPI.md#getrunnerlastjobprocessedmeta) | **Get** /runner_instance/internal/{id}/last_job_processed_meta | Get runner last used job meta
*V1RunnerInstanceAPI* | [**RunnerInstanceCleanupHook**](docs/V1RunnerInstanceAPI.md#runnerinstancecleanuphook) | **Post** /runners_instance/{id}/cleanup_hook | Get runner instance allocation details for the id
*V1RunnersAPI* | [**ComputeCustomRunnerRate**](docs/V1RunnersAPI.md#computecustomrunnerrate) | **Post** /runners/cost/calculator | Get ComputeCustomRunnerRate details
*V1RunnersAPI* | [**DeleteRunner**](docs/V1RunnersAPI.md#deleterunner) | **Delete** /runners/{id} | delete runner for the id. Current organization is figured from the authorization token
*V1RunnersAPI* | [**GetRunner**](docs/V1RunnersAPI.md#getrunner) | **Get** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token
*V1RunnersAPI* | [**GetRunnerSetDefaultGroup**](docs/V1RunnersAPI.md#getrunnersetdefaultgroup) | **Get** /runners/default-group | Get default group for runner set
*V1RunnersAPI* | [**GetRunnersUsage**](docs/V1RunnersAPI.md#getrunnersusage) | **Get** /runners/usage | Get runtimes for runners of the organisation
*V1RunnersAPI* | [**ListRunners**](docs/V1RunnersAPI.md#listrunners) | **Get** /runners | ListRunners lists all the runners for an org.
*V1RunnersAPI* | [**SetRunnerSetDefaultGroup**](docs/V1RunnersAPI.md#setrunnersetdefaultgroup) | **Patch** /runners/default-group | Set default group for runner set
*V1RunnersAPI* | [**SetupRunner**](docs/V1RunnersAPI.md#setuprunner) | **Post** /runners | Adds a new runner for a current organization
*V1RunnersAPI* | [**UpdateRunner**](docs/V1RunnersAPI.md#updaterunner) | **Patch** /runners/{id} | Get runner details for the id. Current organization is figured from the authorization token
*V1SkuAPI* | [**GetSku**](docs/V1SkuAPI.md#getsku) | **Get** /sku/{id} | Get default group for runner set
*V1SkuAPI* | [**ListSku**](docs/V1SkuAPI.md#listsku) | **Get** /sku | ListAllSku lists all the runners sku for an org.
*V1SubscriptionsAPI* | [**DeleteCurrentSubscription**](docs/V1SubscriptionsAPI.md#deletecurrentsubscription) | **Delete** /subscription | Cancel Org current Subscription
*V1SubscriptionsAPI* | [**DeleteStripePaymentMethod**](docs/V1SubscriptionsAPI.md#deletestripepaymentmethod) | **Delete** /subscription/stripe/payment_method/{payment_method_id} | delete stripe setup intent payment method
*V1SubscriptionsAPI* | [**GetBillingInfo**](docs/V1SubscriptionsAPI.md#getbillinginfo) | **Get** /billing/info | Get Billing Info
*V1SubscriptionsAPI* | [**GetCustomerPortalUrl**](docs/V1SubscriptionsAPI.md#getcustomerportalurl) | **Post** /subscription/customer_portal_url | Get customer portal url
*V1SubscriptionsAPI* | [**GetSubscriptionDetails**](docs/V1SubscriptionsAPI.md#getsubscriptiondetails) | **Get** /subscription | Get Current Org Subscription Details
*V1SubscriptionsAPI* | [**InitateSubscriptionCheckout**](docs/V1SubscriptionsAPI.md#initatesubscriptioncheckout) | **Post** /billing/checkout | Initiate Checkout for subscription with PG
*V1SubscriptionsAPI* | [**InitiateSetupIntent**](docs/V1SubscriptionsAPI.md#initiatesetupintent) | **Post** /billing/setup_intent/init | Initiate Checkout for subscription with PG
*V1SubscriptionsAPI* | [**PostSetupIntent**](docs/V1SubscriptionsAPI.md#postsetupintent) | **Post** /billing/setup_intent/post_processor | Post Checkout processing for subscription with PG
*V1SubscriptionsAPI* | [**StripePaymentMethodDefault**](docs/V1SubscriptionsAPI.md#stripepaymentmethoddefault) | **Patch** /subscription/stripe/payment_method/{payment_method_id} | update stripe payment method to default
*V1SubscriptionsAPI* | [**SubscriptionPGWebhook**](docs/V1SubscriptionsAPI.md#subscriptionpgwebhook) | **Post** /subscription/{gateway}/webhook | S2S Webhook received from PG
*V1SubscriptionsAPI* | [**UpdateBillingInfo**](docs/V1SubscriptionsAPI.md#updatebillinginfo) | **Patch** /billing/info | Update Billing Info
*V1UiAPI* | [**GetBannerMessages**](docs/V1UiAPI.md#getbannermessages) | **Get** /ui/banner-messages | Get specific banner messages for UI/Org or all
*V1VcsAPI* | [**ApproveVCSIntegration**](docs/V1VcsAPI.md#approvevcsintegration) | **Put** /vcs/approve-integration | This handles the callback for approving an installation
*V1VcsAPI* | [**CreateVCSGitRepo**](docs/V1VcsAPI.md#createvcsgitrepo) | **Post** /vcs/repos | create vcs repo based on repo internal id
*V1VcsAPI* | [**CreateVCSIntegration**](docs/V1VcsAPI.md#createvcsintegration) | **Post** /vcs/integrations | Create a new vcs integration
*V1VcsAPI* | [**DeleteVCSIntegration**](docs/V1VcsAPI.md#deletevcsintegration) | **Delete** /vcs/integrations/{integration_id} | Delete an existing vcs integration
*V1VcsAPI* | [**GetVCSGitRepo**](docs/V1VcsAPI.md#getvcsgitrepo) | **Get** /vcs/repos/{id} | get vcs repo based on repo internal id
*V1VcsAPI* | [**ListVCSEntites**](docs/V1VcsAPI.md#listvcsentites) | **Get** /vcs/entities | Lists all vcs entities for vcs integration
*V1VcsAPI* | [**ListVCSIntegration**](docs/V1VcsAPI.md#listvcsintegration) | **Get** /vcs/integrations | Lists all vcs integration for provider
*V1VcsAPI* | [**ListVCSRepos**](docs/V1VcsAPI.md#listvcsrepos) | **Get** /vcs/repos | Lists all vcs repos for vcs integration
*V1VcsAPI* | [**ListVCSRunnerGroups**](docs/V1VcsAPI.md#listvcsrunnergroups) | **Post** /vcs/list-runner-groups | Lists all vcs runner groups
*V1VcsAPI* | [**UpdateVCSIntegration**](docs/V1VcsAPI.md#updatevcsintegration) | **Put** /vcs/integrations/{integration_id} | Update an existing vcs integration
*V1WorkflowsAPI* | [**GetPullRequestAuthURL**](docs/V1WorkflowsAPI.md#getpullrequestauthurl) | **Get** /workflows/pr-auth-url | Get auth url required for GH PR
*V1WorkflowsAPI* | [**ListWorkflows**](docs/V1WorkflowsAPI.md#listworkflows) | **Get** /workflows | Lists all workflows (workflows) for organization according to repo
*V1WorkflowsAPI* | [**PullWorkflows**](docs/V1WorkflowsAPI.md#pullworkflows) | **Patch** /workflows/pull | Pulls all workflows from the provider to the database
*V1WorkflowsAPI* | [**WarpWorkflows**](docs/V1WorkflowsAPI.md#warpworkflows) | **Patch** /workflows/warp | Warps workflows for organization according to given internal workflow ids


## Documentation For Models

 - [ApproveVCSIntegrationRequest](docs/ApproveVCSIntegrationRequest.md)
 - [AuthUserRequest](docs/AuthUserRequest.md)
 - [AuthUserResponse](docs/AuthUserResponse.md)
 - [CommonsAddRunnerInstanceInput](docs/CommonsAddRunnerInstanceInput.md)
 - [CommonsBalance](docs/CommonsBalance.md)
 - [CommonsBalanceDetails](docs/CommonsBalanceDetails.md)
 - [CommonsBannerMessage](docs/CommonsBannerMessage.md)
 - [CommonsBillingInfo](docs/CommonsBillingInfo.md)
 - [CommonsCacheCostSummary](docs/CommonsCacheCostSummary.md)
 - [CommonsContainerRunnerImage](docs/CommonsContainerRunnerImage.md)
 - [CommonsContainerRunnerImageUpdate](docs/CommonsContainerRunnerImageUpdate.md)
 - [CommonsContainerRunnerImageVersion](docs/CommonsContainerRunnerImageVersion.md)
 - [CommonsCostSummary](docs/CommonsCostSummary.md)
 - [CommonsCoupon](docs/CommonsCoupon.md)
 - [CommonsCreateRepoOptions](docs/CommonsCreateRepoOptions.md)
 - [CommonsCreateRunnerImageInput](docs/CommonsCreateRunnerImageInput.md)
 - [CommonsCreateRunnerImagePullSecretInput](docs/CommonsCreateRunnerImagePullSecretInput.md)
 - [CommonsDaywiseCost](docs/CommonsDaywiseCost.md)
 - [CommonsDaywiseRuntime](docs/CommonsDaywiseRuntime.md)
 - [CommonsGetPresignedLogUploadURLOutput](docs/CommonsGetPresignedLogUploadURLOutput.md)
 - [CommonsGithubRunnerApplicationDetails](docs/CommonsGithubRunnerApplicationDetails.md)
 - [CommonsInstanceSku](docs/CommonsInstanceSku.md)
 - [CommonsInstanceSkuProperties](docs/CommonsInstanceSkuProperties.md)
 - [CommonsInternalPostUsageInput](docs/CommonsInternalPostUsageInput.md)
 - [CommonsInternalPostUsageOutput](docs/CommonsInternalPostUsageOutput.md)
 - [CommonsJobRunnerInfo](docs/CommonsJobRunnerInfo.md)
 - [CommonsLastJobProcessedMeta](docs/CommonsLastJobProcessedMeta.md)
 - [CommonsListRunnerImagePullSecretsOutput](docs/CommonsListRunnerImagePullSecretsOutput.md)
 - [CommonsListRunnerImageVersionsOutput](docs/CommonsListRunnerImageVersionsOutput.md)
 - [CommonsListRunnerImagesOutput](docs/CommonsListRunnerImagesOutput.md)
 - [CommonsListTokensOptions](docs/CommonsListTokensOptions.md)
 - [CommonsListUsersOptions](docs/CommonsListUsersOptions.md)
 - [CommonsListUsersResponse](docs/CommonsListUsersResponse.md)
 - [CommonsListVCSRunnerGroupsInput](docs/CommonsListVCSRunnerGroupsInput.md)
 - [CommonsListVCSRunnerGroupsResponse](docs/CommonsListVCSRunnerGroupsResponse.md)
 - [CommonsOrganization](docs/CommonsOrganization.md)
 - [CommonsPaymentDetails](docs/CommonsPaymentDetails.md)
 - [CommonsPaymentMethod](docs/CommonsPaymentMethod.md)
 - [CommonsPostPaymentMethodSetupInput](docs/CommonsPostPaymentMethodSetupInput.md)
 - [CommonsProviderInstanceSkuMapping](docs/CommonsProviderInstanceSkuMapping.md)
 - [CommonsRateCalculationInput](docs/CommonsRateCalculationInput.md)
 - [CommonsRateCalculationOutput](docs/CommonsRateCalculationOutput.md)
 - [CommonsRepo](docs/CommonsRepo.md)
 - [CommonsReqCheckoutSession](docs/CommonsReqCheckoutSession.md)
 - [CommonsReqSetupIntentInit](docs/CommonsReqSetupIntentInit.md)
 - [CommonsResCheckoutSession](docs/CommonsResCheckoutSession.md)
 - [CommonsResSetupIntentInit](docs/CommonsResSetupIntentInit.md)
 - [CommonsRunner](docs/CommonsRunner.md)
 - [CommonsRunnerGroup](docs/CommonsRunnerGroup.md)
 - [CommonsRunnerImage](docs/CommonsRunnerImage.md)
 - [CommonsRunnerImageHook](docs/CommonsRunnerImageHook.md)
 - [CommonsRunnerImagePullSecret](docs/CommonsRunnerImagePullSecret.md)
 - [CommonsRunnerImagePullSecretAWS](docs/CommonsRunnerImagePullSecretAWS.md)
 - [CommonsRunnerImagePullSecretDockerCredentials](docs/CommonsRunnerImagePullSecretDockerCredentials.md)
 - [CommonsRunnerImageSettings](docs/CommonsRunnerImageSettings.md)
 - [CommonsRunnerImageVersion](docs/CommonsRunnerImageVersion.md)
 - [CommonsRunnerInfo](docs/CommonsRunnerInfo.md)
 - [CommonsRunnerInstance](docs/CommonsRunnerInstance.md)
 - [CommonsRunnerInstanceAllocationDetails](docs/CommonsRunnerInstanceAllocationDetails.md)
 - [CommonsRunnerInstanceConfiguration](docs/CommonsRunnerInstanceConfiguration.md)
 - [CommonsRunnerSetConfiguration](docs/CommonsRunnerSetConfiguration.md)
 - [CommonsRunnerSetDefaultGroup](docs/CommonsRunnerSetDefaultGroup.md)
 - [CommonsRunnersCostSummary](docs/CommonsRunnersCostSummary.md)
 - [CommonsRunnersUsage](docs/CommonsRunnersUsage.md)
 - [CommonsRunnerwiseRuntime](docs/CommonsRunnerwiseRuntime.md)
 - [CommonsSetRunnerSetDefaultGroupInput](docs/CommonsSetRunnerSetDefaultGroupInput.md)
 - [CommonsSetupRunnerInput](docs/CommonsSetupRunnerInput.md)
 - [CommonsStorage](docs/CommonsStorage.md)
 - [CommonsSubscriptionDetails](docs/CommonsSubscriptionDetails.md)
 - [CommonsUpcomingBill](docs/CommonsUpcomingBill.md)
 - [CommonsUpdateBillingInfoInput](docs/CommonsUpdateBillingInfoInput.md)
 - [CommonsUpdateContainerRunnerImageVersion](docs/CommonsUpdateContainerRunnerImageVersion.md)
 - [CommonsUpdateRunnerImageInput](docs/CommonsUpdateRunnerImageInput.md)
 - [CommonsUpdateRunnerImagePullSecretInput](docs/CommonsUpdateRunnerImagePullSecretInput.md)
 - [CommonsUpdateRunnerImageVersionInput](docs/CommonsUpdateRunnerImageVersionInput.md)
 - [CommonsUpdateRunnerInput](docs/CommonsUpdateRunnerInput.md)
 - [CommonsUserToken](docs/CommonsUserToken.md)
 - [CommonsVCSIntegrationLean](docs/CommonsVCSIntegrationLean.md)
 - [CommonsWarpbuildImage](docs/CommonsWarpbuildImage.md)
 - [CommonsWorkflow](docs/CommonsWorkflow.md)
 - [CommonsWorkflowStats](docs/CommonsWorkflowStats.md)
 - [CreateVCSIntegrationRequest](docs/CreateVCSIntegrationRequest.md)
 - [DebuggerPublishEventInput](docs/DebuggerPublishEventInput.md)
 - [InsightsCallbackInput](docs/InsightsCallbackInput.md)
 - [ListWorkflowsResponse](docs/ListWorkflowsResponse.md)
 - [MeResponse](docs/MeResponse.md)
 - [SwitchOrganizationRequest](docs/SwitchOrganizationRequest.md)
 - [SwitchOrganizationResponse](docs/SwitchOrganizationResponse.md)
 - [TokenRefreshRequest](docs/TokenRefreshRequest.md)
 - [TokenRefreshResponse](docs/TokenRefreshResponse.md)
 - [TypesGenericSuccessMessage](docs/TypesGenericSuccessMessage.md)
 - [UpdateOrganizationRequest](docs/UpdateOrganizationRequest.md)
 - [UpdateVCSIntegrationRequest](docs/UpdateVCSIntegrationRequest.md)
 - [UpdateVCSIntegrationResponse](docs/UpdateVCSIntegrationResponse.md)
 - [V1ListUsersForOrganizationResult](docs/V1ListUsersForOrganizationResult.md)
 - [V1Organization](docs/V1Organization.md)
 - [V1User](docs/V1User.md)
 - [VCSEntity](docs/VCSEntity.md)
 - [VCSIntegration](docs/VCSIntegration.md)
 - [WarpBuildAPIError](docs/WarpBuildAPIError.md)
 - [WarpWorkflowsRequest](docs/WarpWorkflowsRequest.md)
 - [WarpWorkflowsResponse](docs/WarpWorkflowsResponse.md)


## Documentation For Authorization


Authentication schemes defined for the API:
### JWTKeyAuth

- **Type**: API key
- **API key parameter name**: Authorization
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: Authorization and passed in as the auth context for each request.

Example

```golang
auth := context.WithValue(
		context.Background(),
		sw.ContextAPIKeys,
		map[string]sw.APIKey{
			"Authorization": {Key: "API_KEY_STRING"},
		},
	)
r, err := client.Service.Operation(auth, args)
```

### WarpBuildAdminSecretAuth

- **Type**: API key
- **API key parameter name**: X-Warp-Build-Admin-Secret
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: X-Warp-Build-Admin-Secret and passed in as the auth context for each request.

Example

```golang
auth := context.WithValue(
		context.Background(),
		sw.ContextAPIKeys,
		map[string]sw.APIKey{
			"X-Warp-Build-Admin-Secret": {Key: "API_KEY_STRING"},
		},
	)
r, err := client.Service.Operation(auth, args)
```

### WarpBuildServiceSecretAuth

- **Type**: API key
- **API key parameter name**: x-warpbuild-service-secret
- **Location**: HTTP header

Note, each API key must be added to a map of `map[string]APIKey` where the key is: x-warpbuild-service-secret and passed in as the auth context for each request.

Example

```golang
auth := context.WithValue(
		context.Background(),
		sw.ContextAPIKeys,
		map[string]sw.APIKey{
			"x-warpbuild-service-secret": {Key: "API_KEY_STRING"},
		},
	)
r, err := client.Service.Operation(auth, args)
```


## Documentation for Utility Methods

Due to the fact that model structure members are all pointers, this package contains
a number of utility functions to easily obtain pointers to values of basic types.
Each of these functions takes a value of the given basic type and returns a pointer to it:

* `PtrBool`
* `PtrInt`
* `PtrInt32`
* `PtrInt64`
* `PtrFloat`
* `PtrFloat32`
* `PtrFloat64`
* `PtrString`
* `PtrTime`

## Author

support@swagger.io

