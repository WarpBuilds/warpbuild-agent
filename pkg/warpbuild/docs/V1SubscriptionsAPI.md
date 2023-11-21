# \V1SubscriptionsAPI

All URIs are relative to *https://backend.warpbuild.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DeleteCurrentSubscription**](V1SubscriptionsAPI.md#DeleteCurrentSubscription) | **Delete** /subscription | Cancel Org current Subscription
[**DeleteStripePaymentMethod**](V1SubscriptionsAPI.md#DeleteStripePaymentMethod) | **Delete** /subscription/stripe/payment_method/{payment_method_id} | delete stripe setup intent payment method
[**GetSubscriptionDetails**](V1SubscriptionsAPI.md#GetSubscriptionDetails) | **Get** /subscription | Get Current Org Subscription Details
[**InitateSubscriptionCheckout**](V1SubscriptionsAPI.md#InitateSubscriptionCheckout) | **Post** /billing/checkout | Initiate Checkout for subscription with PG
[**InitiateSetupIntent**](V1SubscriptionsAPI.md#InitiateSetupIntent) | **Post** /billing/setup_intent/init | Initiate Checkout for subscription with PG
[**StripePaymentMethodDefault**](V1SubscriptionsAPI.md#StripePaymentMethodDefault) | **Patch** /subscription/stripe/payment_method/{payment_method_id} | update stripe payment method to default
[**SubscriptionPGWebhook**](V1SubscriptionsAPI.md#SubscriptionPGWebhook) | **Post** /subscription/{gateway}/webhook | S2S Webhook received from PG



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
    resp, r, err := apiClient.V1SubscriptionsAPI.DeleteCurrentSubscription(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.DeleteCurrentSubscription``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteCurrentSubscription`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.DeleteCurrentSubscription`: %v\n", resp)
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
    resp, r, err := apiClient.V1SubscriptionsAPI.DeleteStripePaymentMethod(context.Background(), paymentMethodId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.DeleteStripePaymentMethod``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DeleteStripePaymentMethod`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.DeleteStripePaymentMethod`: %v\n", resp)
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
    resp, r, err := apiClient.V1SubscriptionsAPI.GetSubscriptionDetails(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.GetSubscriptionDetails``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetSubscriptionDetails`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.GetSubscriptionDetails`: %v\n", resp)
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
    resp, r, err := apiClient.V1SubscriptionsAPI.InitateSubscriptionCheckout(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.InitateSubscriptionCheckout``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InitateSubscriptionCheckout`: CommonsResCheckoutSession
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.InitateSubscriptionCheckout`: %v\n", resp)
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
    resp, r, err := apiClient.V1SubscriptionsAPI.InitiateSetupIntent(context.Background()).Body(body).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.InitiateSetupIntent``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `InitiateSetupIntent`: CommonsResSetupIntentInit
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.InitiateSetupIntent`: %v\n", resp)
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
    resp, r, err := apiClient.V1SubscriptionsAPI.StripePaymentMethodDefault(context.Background(), paymentMethodId).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.StripePaymentMethodDefault``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `StripePaymentMethodDefault`: CommonsSubscriptionDetails
    fmt.Fprintf(os.Stdout, "Response from `V1SubscriptionsAPI.StripePaymentMethodDefault`: %v\n", resp)
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
    r, err := apiClient.V1SubscriptionsAPI.SubscriptionPGWebhook(context.Background(), gateway).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `V1SubscriptionsAPI.SubscriptionPGWebhook``: %v\n", err)
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

