# Entrance CLI User Documentation

The Entrance CLI is an interactive command-line tool for testing authentication flows against a real Entra External ID tenant. It exercises every flow in the Entrance library, prompting for input at each step.

## Building

From the repository root:

```bash
go build -o entrance-cli ./cmd/entrance-cli
```

## Running

```bash
entrance-cli \
  --tenant-subdomain contoso \
  --tenant-identifier contoso.onmicrosoft.com \
  --client-identifier 00000000-0000-0000-0000-000000000000 \
  --email user@contoso.com
```

## Flags

| Flag                   | Required | Default                          | Description                                              |
|------------------------|----------|----------------------------------|----------------------------------------------------------|
| `--tenant-subdomain`   | Yes      |                                  | Tenant subdomain, e.g. `contoso` for `contoso.ciamlogin.com` |
| `--tenant-identifier`  | Yes      |                                  | Tenant GUID or `{tenant}.onmicrosoft.com`. Use the GUID when SMS MFA is possible. |
| `--client-identifier`  | Yes      |                                  | Application (client) ID registered in Entra              |
| `--email`              | Yes      |                                  | Email address of the user to authenticate                |
| `--scopes`             | No       | `openid profile offline_access`  | Space-separated OAuth scopes                             |
| `--show-full-tokens`   | No       | `false`                          | Print full token values instead of truncating to 40 characters |

If any required flag is missing, the CLI prints usage information and exits with code 1.

## Interactive Menu

After startup, the CLI displays a menu:

```
Select an authentication flow:
  1. Email one-time passcode sign-in
  2. Password sign-in
  3. Password sign-in with multifactor authentication
  4. Password sign-in with registration
  5. Exit
Select a choice:
```

Enter the number of the desired flow. After each flow completes (whether it succeeds or encounters an error), the menu reappears so you can run another flow or exit.

## Authentication Flows

### Flow 1: Email One-Time Passcode Sign-In

This flow authenticates the user with a one-time passcode sent to their email address.

1. The CLI initiates a sign-in request and requests an email OTP challenge.
2. On success, it displays the challenge target (a masked version of the email address) and the expected passcode length:
   ```
   Challenge sent to: c***r@co**o.com
   Expected code length: 8
   ```
3. Check your email for the passcode, then enter it at the prompt:
   ```
   Enter the one-time passcode:
   ```
4. If the passcode is incorrect, you may retry up to 3 times. After 3 failed attempts the CLI prints `Maximum passcode attempts exceeded.` and returns to the menu.
5. On success, the CLI prints the token response.

### Flow 2: Password Sign-In

This flow authenticates the user with a password.

1. The CLI prompts for your password (input is hidden):
   ```
   Enter password:
   ```
2. It initiates a sign-in request and submits the password.
3. On success, the CLI prints the token response.

### Flow 3: Password Sign-In with Multifactor Authentication

This flow handles tenants that require MFA after password authentication.

1. The CLI prompts for your password (input is hidden):
   ```
   Enter password:
   ```
2. It initiates a sign-in request with the `MultifactorRequired` capability and submits the password.
3. If MFA is required, the CLI retrieves and displays your registered MFA methods:
   ```
     1. email (channel: email, hint: c***r@co**o.com)
     2. phone (channel: phone_sms, hint: +1 ***-***-1234)
   ```
4. Select a method by number:
   ```
   Select an MFA method (enter number):
   ```
5. The CLI sends the MFA challenge and displays the target and expected code length:
   ```
   Challenge sent to: c***r@co**o.com
   Expected code length: 8
   ```
6. Enter the MFA passcode:
   ```
   Enter the MFA one-time passcode:
   ```
7. On success, the CLI prints the token response.
8. If the tenant does not require MFA, the password step succeeds directly and tokens are printed without the MFA steps.

### Flow 4: Password Sign-In with Registration

This flow handles the case where the user must register a strong authentication method before completing sign-in.

1. The CLI prompts for your password (input is hidden):
   ```
   Enter password:
   ```
2. It initiates a sign-in request with the `RegistrationRequired` capability and submits the password.
3. If registration is required, the CLI retrieves and displays methods available for enrollment:
   ```
     1. email (channel: email, hint: )
     2. phone (channel: phone_sms, hint: )
   ```
4. Select a method by number:
   ```
   Select a registration method (enter number):
   ```
5. Enter the target value for the selected method:
   ```
   Enter the target value (e.g., email address or phone number):
   ```
6. The CLI sends a registration challenge. If the challenge type is `preverified` (e.g., the email address is already verified), OTP entry is skipped. Otherwise:
   ```
   Enter the registration one-time passcode:
   ```
7. After registration completes, the CLI exchanges the continuation token for security tokens and prints the token response.
8. If registration is not required, the password step succeeds directly and tokens are printed.

## Token Output

On successful authentication, the CLI prints:

```
Authentication succeeded.
  Access Token:  eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIs...
  ID Token:      eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIs...
  Refresh Token: 0.AVQAx8Bs7dFJ3EKZ1gQfGLOT0CYnHpECTm...
  Expires In:    3600
  Token Type:    Bearer
```

By default, token values longer than 40 characters are truncated with `...` appended. Pass `--show-full-tokens` to display the complete values.

## Error Messages

The CLI displays specific messages for each error condition:

| Condition                          | Message                                                |
|------------------------------------|--------------------------------------------------------|
| User does not exist                | `User not found.`                                      |
| Incorrect one-time passcode        | `Invalid passcode.`                                    |
| Incorrect password                 | `Invalid password.`                                    |
| Continuation token has expired     | `Session expired — please restart the flow.`           |
| Native auth not supported          | `Native auth unavailable — browser auth required.`     |
| SMS fraud protection triggered     | `Access denied (SMS fraud protection).`                |
| HTTP transport failure             | `HTTP error:` followed by the underlying error         |
| Other Entra API error              | `Authentication error:` followed by the error details  |

After any error, the CLI returns to the main menu so you can try again or select a different flow.

## Tenant Identifier Guidance

The `--tenant-identifier` flag accepts either a domain name (`contoso.onmicrosoft.com`) or a tenant GUID (`00000000-0000-0000-0000-000000000000`).

Use the **tenant GUID** form when testing flows that involve SMS-based MFA. The SMS OTP token endpoint requires the tenant ID form of the URL. For email-only flows, either form works.

## Exit Codes

| Code | Meaning                            |
|------|------------------------------------|
| 0    | Normal exit                        |
| 1    | Missing or invalid command-line flags |
