# \V1SubscriptionsApi

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteCurrentSubscription**](V1SubscriptionsApi.md#DeleteCurrentSubscription) | **Delete** /subscription | Cancel Org current Subscription
[**DeleteStripePaymentMethod**](V1SubscriptionsApi.md#DeleteStripePaymentMethod) | **Delete** /subscription/stripe/payment_method/{payment_method_id} | delete stripe setup intent payment method
[**GetBillingInfo**](V1SubscriptionsApi.md#GetBillingInfo) | **Get** /billing/info | Get Billing Info
[**GetCustomerPortalUrl**](V1SubscriptionsApi.md#GetCustomerPortalUrl) | **Post** /subscription/customer_portal_url | Get customer portal url
[**GetSubscriptionDetails**](V1SubscriptionsApi.md#GetSubscriptionDetails) | **Get** /subscription | Get Current Org Subscription Details
[**InitateSubscriptionCheckout**](V1SubscriptionsApi.md#InitateSubscriptionCheckout) | **Post** /billing/checkout | Initiate Checkout for subscription with PG
[**InitiateSetupIntent**](V1SubscriptionsApi.md#InitiateSetupIntent) | **Post** /billing/setup_intent/init | Initiate Checkout for subscription with PG
[**PostSetupIntent**](V1SubscriptionsApi.md#PostSetupIntent) | **Post** /billing/setup_intent/post_processor | Post Checkout processing for subscription with PG
[**StripePaymentMethodDefault**](V1SubscriptionsApi.md#StripePaymentMethodDefault) | **Patch** /subscription/stripe/payment_method/{payment_method_id} | update stripe payment method to default
[**SubscriptionPGWebhook**](V1SubscriptionsApi.md#SubscriptionPGWebhook) | **Post** /subscription/{gateway}/webhook | S2S Webhook received from PG
[**UpdateBillingInfo**](V1SubscriptionsApi.md#UpdateBillingInfo) | **Patch** /billing/info | Update Billing Info



## DeleteCurrentSubscription

> CommonsSubscriptionDetails DeleteCurrentSubscription(ctx).Execute()

Cancel Org current Subscription

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.DeleteCurrentSubscription(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.DeleteCurrentSubscription``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteCurrentSubscription`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.DeleteCurrentSubscription`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteCurrentSubscriptionRequest struct via the builder pattern


### Return type

[**CommonsSubscriptionDetails**](CommonsSubscriptionDetails.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteStripePaymentMethod

> CommonsSubscriptionDetails DeleteStripePaymentMethod(ctx, paymentMethodId).Execute()

delete stripe setup intent payment method

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    paymentMethodId := "paymentMethodId_example" // string | ID for the stripe payment method

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.DeleteStripePaymentMethod(context.Background(), paymentMethodId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.DeleteStripePaymentMethod``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteStripePaymentMethod`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.DeleteStripePaymentMethod`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**paymentMethodId** | **string** | ID for the stripe payment method | 

### Other Parameters

Other parameters are passed through a pointer to a apiDeleteStripePaymentMethodRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsSubscriptionDetails**](CommonsSubscriptionDetails.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetBillingInfo

> CommonsBillingInfo GetBillingInfo(ctx).Execute()

Get Billing Info

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.GetBillingInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.GetBillingInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetBillingInfo`: CommonsBillingInfo
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.GetBillingInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetBillingInfoRequest struct via the builder pattern


### Return type

[**CommonsBillingInfo**](CommonsBillingInfo.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetCustomerPortalUrl

> string GetCustomerPortalUrl(ctx).Execute()

Get customer portal url

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.GetCustomerPortalUrl(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.GetCustomerPortalUrl``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetCustomerPortalUrl`: string
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.GetCustomerPortalUrl`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetCustomerPortalUrlRequest struct via the builder pattern


### Return type

**string**

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSubscriptionDetails

> CommonsSubscriptionDetails GetSubscriptionDetails(ctx).Execute()

Get Current Org Subscription Details

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.GetSubscriptionDetails(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.GetSubscriptionDetails``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetSubscriptionDetails`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.GetSubscriptionDetails`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiGetSubscriptionDetailsRequest struct via the builder pattern


### Return type

[**CommonsSubscriptionDetails**](CommonsSubscriptionDetails.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InitateSubscriptionCheckout

> CommonsResCheckoutSession InitateSubscriptionCheckout(ctx).Body(body).Execute()

Initiate Checkout for subscription with PG

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    body := *openapiclient.NewCommonsReqCheckoutSession("CancelUrl_example", "SuccessUrl_example") // CommonsReqCheckoutSession | initiate checkout session input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.InitateSubscriptionCheckout(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.InitateSubscriptionCheckout``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InitateSubscriptionCheckout`: CommonsResCheckoutSession
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.InitateSubscriptionCheckout`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInitateSubscriptionCheckoutRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsReqCheckoutSession**](CommonsReqCheckoutSession.md) | initiate checkout session input | 

### Return type

[**CommonsResCheckoutSession**](CommonsResCheckoutSession.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## InitiateSetupIntent

> CommonsResSetupIntentInit InitiateSetupIntent(ctx).Body(body).Execute()

Initiate Checkout for subscription with PG

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    body := *openapiclient.NewCommonsReqSetupIntentInit("CancelUrl_example", "SuccessUrl_example") // CommonsReqSetupIntentInit | initiate setup intent session input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.InitiateSetupIntent(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.InitiateSetupIntent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InitiateSetupIntent`: CommonsResSetupIntentInit
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.InitiateSetupIntent`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiInitiateSetupIntentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsReqSetupIntentInit**](CommonsReqSetupIntentInit.md) | initiate setup intent session input | 

### Return type

[**CommonsResSetupIntentInit**](CommonsResSetupIntentInit.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## PostSetupIntent

> CommonsPostPaymentMethodSetupInput PostSetupIntent(ctx).Body(body).Execute()

Post Checkout processing for subscription with PG

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    body := *openapiclient.NewCommonsPostPaymentMethodSetupInput() // CommonsPostPaymentMethodSetupInput | post setup intent session input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.PostSetupIntent(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.PostSetupIntent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `PostSetupIntent`: CommonsPostPaymentMethodSetupInput
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.PostSetupIntent`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiPostSetupIntentRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsPostPaymentMethodSetupInput**](CommonsPostPaymentMethodSetupInput.md) | post setup intent session input | 

### Return type

[**CommonsPostPaymentMethodSetupInput**](CommonsPostPaymentMethodSetupInput.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## StripePaymentMethodDefault

> CommonsSubscriptionDetails StripePaymentMethodDefault(ctx, paymentMethodId).Execute()

update stripe payment method to default

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    paymentMethodId := "paymentMethodId_example" // string | ID for the stripe payment method

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.StripePaymentMethodDefault(context.Background(), paymentMethodId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.StripePaymentMethodDefault``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `StripePaymentMethodDefault`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.StripePaymentMethodDefault`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**paymentMethodId** | **string** | ID for the stripe payment method | 

### Other Parameters

Other parameters are passed through a pointer to a apiStripePaymentMethodDefaultRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**CommonsSubscriptionDetails**](CommonsSubscriptionDetails.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SubscriptionPGWebhook

> SubscriptionPGWebhook(ctx, gateway).Execute()

S2S Webhook received from PG

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    gateway := "gateway_example" // string | gateway name, current only stripe

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.V1SubscriptionsApi.SubscriptionPGWebhook(context.Background(), gateway).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.SubscriptionPGWebhook``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**gateway** | **string** | gateway name, current only stripe | 

### Other Parameters

Other parameters are passed through a pointer to a apiSubscriptionPGWebhookRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateBillingInfo

> CommonsBillingInfo UpdateBillingInfo(ctx).Body(body).Execute()

Update Billing Info

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID/warpbuild"
)

func main() {
    body := *openapiclient.NewCommonsUpdateBillingInfoInput() // CommonsUpdateBillingInfoInput | billing info update input

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.V1SubscriptionsApi.UpdateBillingInfo(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsApi.UpdateBillingInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `UpdateBillingInfo`: CommonsBillingInfo
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsApi.UpdateBillingInfo`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUpdateBillingInfoRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**CommonsUpdateBillingInfoInput**](CommonsUpdateBillingInfoInput.md) | billing info update input | 

### Return type

[**CommonsBillingInfo**](CommonsBillingInfo.md)

### Authorization

[JWTKeyAuth](../README.md#JWTKeyAuth)

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

