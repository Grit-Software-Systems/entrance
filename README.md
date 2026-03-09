# Entrance

A Go client library for [Entra External ID native authentication](https://learn.microsoft.com/en-us/entra/identity-platform/reference-native-authentication-api). Entrance wraps the multi-step OAuth 2.0 native authentication API behind a type-safe, debuggable interface supporting email OTP sign-in, password sign-in, multifactor authentication, and strong authentication method registration.

## Installation

```bash
go get github.com/Grit-Software-Systems/entrance
```

## Usage

### Create a Client

```go
configuration := entrance.Configuration{
    TenantSubdomain:  "contoso",
    TenantIdentifier: "contoso.onmicrosoft.com",
    ClientIdentifier: "00000000-0000-0000-0000-000000000000",
}
authenticationClient := entrance.NewClient(configuration)
```

To inject a custom `*http.Client` (for testing or proxy configuration):

```go
authenticationClient := entrance.NewClientWithHttpClient(configuration, customHttpClient)
```

### Sign In with Email One-Time Passcode

```go
initiateResponse, initiateError := authenticationClient.Initiate(
    requestContext, emailAddress,
    []entrance.ChallengeMethod{entrance.ChallengeMethodOneTimePasscode, entrance.ChallengeMethodRedirect},
    nil,
)

challengeResponse, challengeError := authenticationClient.Challenge(
    requestContext, initiateResponse.ContinuationToken,
    []entrance.ChallengeMethod{entrance.ChallengeMethodOneTimePasscode, entrance.ChallengeMethodRedirect},
)
// challengeResponse.ChallengeTargetLabel contains the masked email, e.g. "c***r@co**o.com"
// challengeResponse.CodeLength contains the expected OTP length
// Prompt the user for the one-time passcode...

tokenResponse, tokenError := authenticationClient.RedeemOneTimePasscode(
    requestContext, challengeResponse.ContinuationToken, userEnteredPasscode,
)
// tokenResponse.AccessToken, tokenResponse.IdToken, tokenResponse.RefreshToken are now available
```

### Sign In with Password

```go
initiateResponse, initiateError := authenticationClient.Initiate(
    requestContext, emailAddress,
    []entrance.ChallengeMethod{entrance.ChallengeMethodPassword, entrance.ChallengeMethodRedirect},
    nil,
)

challengeResponse, challengeError := authenticationClient.Challenge(
    requestContext, initiateResponse.ContinuationToken,
    []entrance.ChallengeMethod{entrance.ChallengeMethodPassword, entrance.ChallengeMethodRedirect},
)

tokenResponse, tokenError := authenticationClient.RedeemPassword(
    requestContext, challengeResponse.ContinuationToken, password,
)
```

### Password Sign-In with Multifactor Authentication

When a tenant requires MFA, `RedeemPassword` returns a `MultifactorRequiredError`. Include `CapabilityMultifactorRequired` in the `Initiate` call to opt in to native MFA handling. Pass the tenant GUID as `TenantIdentifier` when SMS MFA is possible, because the SMS OTP token endpoint requires the tenant ID form of the URL.

```go
configuration := entrance.Configuration{
    TenantSubdomain:  "contoso",
    TenantIdentifier: "00000000-0000-0000-0000-000000000000",
    ClientIdentifier: "11111111-1111-1111-1111-111111111111",
}
authenticationClient := entrance.NewClient(configuration)

initiateResponse, initiateError := authenticationClient.Initiate(
    requestContext, emailAddress,
    []entrance.ChallengeMethod{entrance.ChallengeMethodPassword, entrance.ChallengeMethodRedirect},
    []entrance.Capability{entrance.CapabilityMultifactorRequired},
)

challengeResponse, challengeError := authenticationClient.Challenge(
    requestContext, initiateResponse.ContinuationToken,
    []entrance.ChallengeMethod{entrance.ChallengeMethodPassword, entrance.ChallengeMethodRedirect},
)

tokenResponse, passwordError := authenticationClient.RedeemPassword(
    requestContext, challengeResponse.ContinuationToken, password,
)

var multifactorError entrance.MultifactorRequiredError
if errors.As(passwordError, &multifactorError) {
    introspectResponse, introspectError := authenticationClient.Introspect(
        requestContext, multifactorError.ContinuationToken,
    )
    // introspectResponse.Methods lists the user's registered MFA methods.
    // Present them to the user and let them choose one...

    mfaChallengeResponse, mfaChallengeError := authenticationClient.ChallengeMultifactor(
        requestContext, introspectResponse.ContinuationToken, selectedMethod.Identifier,
    )
    // Prompt the user for the MFA one-time passcode...

    tokenResponse, tokenError := authenticationClient.RedeemMultifactorOneTimePasscode(
        requestContext, mfaChallengeResponse.ContinuationToken, mfaPasscode,
        selectedMethod.ChallengeChannel,
    )
}
```

### Password Sign-In with Registration Required

When the user has no registered strong authentication method, `RedeemPassword` returns a `RegistrationRequiredError`. Include `CapabilityRegistrationRequired` in the `Initiate` call to opt in to native registration handling.

```go
initiateResponse, initiateError := authenticationClient.Initiate(
    requestContext, emailAddress,
    []entrance.ChallengeMethod{entrance.ChallengeMethodPassword, entrance.ChallengeMethodRedirect},
    []entrance.Capability{entrance.CapabilityRegistrationRequired},
)

// ... Challenge and RedeemPassword as above ...

var registrationError entrance.RegistrationRequiredError
if errors.As(passwordError, &registrationError) {
    registrationIntrospectResponse, introspectError := authenticationClient.IntrospectRegistration(
        requestContext, registrationError.ContinuationToken,
    )
    // registrationIntrospectResponse.Methods lists methods available for enrollment.
    // Present them to the user...

    registrationChallengeResponse, challengeError := authenticationClient.ChallengeRegistration(
        requestContext, registrationIntrospectResponse.ContinuationToken,
        "email", userEmailAddress,
    )
    // If registrationChallengeResponse.ChallengeType is "preverified", skip OTP entry.
    // Otherwise, prompt the user for the one-time passcode...

    registrationContinueResponse, continueError := authenticationClient.ContinueRegistration(
        requestContext, registrationChallengeResponse.ContinuationToken, registrationPasscode,
    )

    tokenResponse, tokenError := authenticationClient.RedeemContinuationToken(
        requestContext, registrationContinueResponse.ContinuationToken,
    )
}
```

### Handling Redirect Fallback

If the server cannot fulfill the requested challenge type natively, `Challenge` returns a `RedirectRequiredError`. The `ChallengeMethodRedirect` value must always be included in challenge type lists as a required fallback.

```go
challengeResponse, challengeError := authenticationClient.Challenge(
    requestContext, initiateResponse.ContinuationToken,
    []entrance.ChallengeMethod{entrance.ChallengeMethodOneTimePasscode, entrance.ChallengeMethodRedirect},
)

var redirectError entrance.RedirectRequiredError
if errors.As(challengeError, &redirectError) {
    // Native authentication is not available. Fall back to browser-based authentication.
}
```

## Configuration

| Field              | Description                                                                                       |
|--------------------|---------------------------------------------------------------------------------------------------|
| `TenantSubdomain`  | The tenant subdomain, e.g. `"contoso"` for `contoso.ciamlogin.com`                                |
| `TenantIdentifier` | The tenant GUID or `"{tenant}.onmicrosoft.com"`. Use the GUID when SMS MFA is possible.           |
| `ClientIdentifier` | The application (client) ID registered in Entra                                                   |
| `Scopes`           | Space-separated OAuth scopes. Defaults to `"openid profile offline_access"` when empty.           |

## Error Handling

All API errors are returned as typed structs implementing the `error` interface. Use `errors.As` to inspect specific error types.

| Error Type                   | Condition                                                        |
|------------------------------|------------------------------------------------------------------|
| `UserNotFoundError`          | The username does not exist                                      |
| `InvalidPasscodeError`       | The submitted one-time passcode is incorrect                     |
| `InvalidPasswordError`       | The submitted password is incorrect                              |
| `ExpiredTokenError`          | A continuation token has expired                                 |
| `RedirectRequiredError`      | Native authentication is unavailable; fall back to browser auth  |
| `MultifactorRequiredError`   | MFA is required; contains the continuation token to proceed      |
| `RegistrationRequiredError`  | Strong auth method registration is required; contains the token  |
| `AccessDeniedError`          | SMS fraud protection blocked the request                         |
| `RequestError`               | HTTP transport failure or unexpected response                    |
| `AuthenticationError`        | Base error for any other Entra API error                         |

`MultifactorRequiredError` and `RegistrationRequiredError` carry a `ContinuationToken` field needed to proceed with the MFA or registration flow.

`RequestError` wraps the underlying transport error, accessible via `errors.Unwrap`.

## Challenge Methods

| Constant                          | Value        | Description                              |
|-----------------------------------|--------------|------------------------------------------|
| `ChallengeMethodOneTimePasscode`  | `"oob"`      | Out-of-band one-time passcode via email  |
| `ChallengeMethodPassword`         | `"password"` | Password-based authentication            |
| `ChallengeMethodRedirect`         | `"redirect"` | Fallback to browser auth (always required) |

## Capabilities

| Constant                          | Value                    | Description                                  |
|-----------------------------------|--------------------------|----------------------------------------------|
| `CapabilityMultifactorRequired`   | `"mfa_required"`         | The app can handle the native MFA flow       |
| `CapabilityRegistrationRequired`  | `"registration_required"`| The app can drive the registration UI        |

## API Reference

### Client Methods

#### Sign-In

- **`Initiate`** — Starts a sign-in flow for the given username.
- **`Challenge`** — Requests the server to issue a challenge using the specified method.
- **`RedeemOneTimePasscode`** — Submits an email one-time passcode to obtain tokens.
- **`RedeemPassword`** — Submits a password to obtain tokens.

#### Multifactor Authentication

- **`Introspect`** — Lists the user's registered MFA methods.
- **`ChallengeMultifactor`** — Sends an MFA challenge to the selected method.
- **`RedeemMultifactorOneTimePasscode`** — Submits the MFA one-time passcode to obtain tokens. Automatically uses the tenant-ID token endpoint for SMS OTP.

#### Strong Authentication Method Registration

- **`IntrospectRegistration`** — Lists methods available for enrollment.
- **`ChallengeRegistration`** — Sends a challenge to the chosen registration target.
- **`ContinueRegistration`** — Submits the OTP to complete method registration.
- **`RedeemContinuationToken`** — Exchanges a continuation token for security tokens after registration completes.

## Architecture

The library exposes step-by-step methods rather than a high-level orchestrator. Each API call is a separate method, giving the caller full control over flow, error handling, and UI between steps.

Continuation tokens thread the multi-step flow together. Each response returns a new token that must be passed to the next request. Previous tokens become invalid after each step.

### Package Layout

```
entrance/
├── client.go              Client struct and constructors
├── configuration.go       Configuration struct
├── errors.go              Typed error structs
├── sign_in.go             Sign-in methods (Initiate, Challenge, Redeem*)
├── multifactor.go         MFA methods (Introspect, ChallengeMultifactor, RedeemMultifactor*)
├── registration.go        Registration methods (IntrospectRegistration, Challenge*, Continue*, RedeemContinuationToken)
├── types.go               Public request/response types and constants
├── locales/
│   └── en.json            English translations (embedded at compile time)
└── internal/
    ├── global/
    │   ├── catalog.go     Loads translations into golang.org/x/text catalog
    │   └── messages.go    Message key constants
    └── request/
        ├── endpoint.go    URL construction
        ├── form.go        Form-encoded request builder
        └── sender.go      HTTP POST and JSON response parsing
```

## License

See [LICENSE](LICENSE) for details.
