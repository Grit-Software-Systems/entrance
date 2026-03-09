package entrance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestInitiateReturnsInitiateResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"initiate-token-123"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	initiateResponse, initiateError := testClient.Initiate(
		context.Background(),
		"user@contoso.com",
		[]ChallengeMethod{ChallengeMethodOneTimePasscode, ChallengeMethodRedirect},
		nil,
	)

	if initiateError != nil {
		testing.Fatalf("expected no error, got %v", initiateError)
	}
	if initiateResponse.ContinuationToken != "initiate-token-123" {
		testing.Errorf("expected ContinuationToken %q, got %q", "initiate-token-123", initiateResponse.ContinuationToken)
	}
}

func TestInitiateReturnsUserNotFoundErrorOnUnknownUser(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"user_not_found","error_description":"User does not exist"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, initiateError := testClient.Initiate(
		context.Background(),
		"unknown@contoso.com",
		[]ChallengeMethod{ChallengeMethodOneTimePasscode, ChallengeMethodRedirect},
		nil,
	)

	var userNotFoundError UserNotFoundError
	isMatch := errors.As(initiateError, &userNotFoundError)
	if !isMatch {
		testing.Errorf("expected UserNotFoundError, got %T: %v", initiateError, initiateError)
	}
}

func TestInitiateSendsCapabilitiesWhenProvided(testing *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		bodyBytes := make([]byte, httpRequest.ContentLength)
		httpRequest.Body.Read(bodyBytes)
		capturedBody = string(bodyBytes)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"token"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	testClient.Initiate(
		context.Background(),
		"user@contoso.com",
		[]ChallengeMethod{ChallengeMethodPassword, ChallengeMethodRedirect},
		[]Capability{CapabilityMultifactorRequired},
	)

	if capturedBody == "" {
		testing.Fatal("expected request body to be captured")
	}
	parsedForm := parseFormBody(capturedBody)
	capabilitiesValue := parsedForm["capabilities"]
	if capabilitiesValue != "mfa_required" {
		testing.Errorf("expected capabilities %q, got %q", "mfa_required", capabilitiesValue)
	}
}

func TestChallengeReturnsChallengeResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"challenge-token-456",
			"challenge_type":"oob",
			"challenge_channel":"email",
			"code_length":8,
			"challenge_target_label":"c***r@co**o.com"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	challengeResponse, challengeError := testClient.Challenge(
		context.Background(),
		"initiate-token-123",
		[]ChallengeMethod{ChallengeMethodOneTimePasscode, ChallengeMethodRedirect},
	)

	if challengeError != nil {
		testing.Fatalf("expected no error, got %v", challengeError)
	}
	if challengeResponse.ContinuationToken != "challenge-token-456" {
		testing.Errorf("expected ContinuationToken %q, got %q", "challenge-token-456", challengeResponse.ContinuationToken)
	}
	if challengeResponse.ChallengeType != ChallengeMethodOneTimePasscode {
		testing.Errorf("expected ChallengeType %q, got %q", ChallengeMethodOneTimePasscode, challengeResponse.ChallengeType)
	}
	if challengeResponse.ChallengeChannel != "email" {
		testing.Errorf("expected ChallengeChannel %q, got %q", "email", challengeResponse.ChallengeChannel)
	}
	if challengeResponse.CodeLength != 8 {
		testing.Errorf("expected CodeLength 8, got %d", challengeResponse.CodeLength)
	}
	if challengeResponse.ChallengeTargetLabel != "c***r@co**o.com" {
		testing.Errorf("expected ChallengeTargetLabel %q, got %q", "c***r@co**o.com", challengeResponse.ChallengeTargetLabel)
	}
}

func TestChallengeReturnsRedirectRequiredErrorOnRedirect(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"token","challenge_type":"redirect"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, challengeError := testClient.Challenge(
		context.Background(),
		"initiate-token-123",
		[]ChallengeMethod{ChallengeMethodOneTimePasscode, ChallengeMethodRedirect},
	)

	var redirectRequiredError RedirectRequiredError
	isMatch := errors.As(challengeError, &redirectRequiredError)
	if !isMatch {
		testing.Errorf("expected RedirectRequiredError, got %T: %v", challengeError, challengeError)
	}
}

func TestRedeemOneTimePasscodeReturnsTokenResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"access-token-value",
			"id_token":"id-token-value",
			"refresh_token":"refresh-token-value",
			"expires_in":3600
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	tokenResponse, tokenError := testClient.RedeemOneTimePasscode(
		context.Background(),
		"challenge-token-456",
		"12345678",
	)

	if tokenError != nil {
		testing.Fatalf("expected no error, got %v", tokenError)
	}
	if tokenResponse.TokenType != "Bearer" {
		testing.Errorf("expected TokenType %q, got %q", "Bearer", tokenResponse.TokenType)
	}
	if tokenResponse.AccessToken != "access-token-value" {
		testing.Errorf("expected AccessToken %q, got %q", "access-token-value", tokenResponse.AccessToken)
	}
	if tokenResponse.IdToken != "id-token-value" {
		testing.Errorf("expected IdToken %q, got %q", "id-token-value", tokenResponse.IdToken)
	}
	if tokenResponse.RefreshToken != "refresh-token-value" {
		testing.Errorf("expected RefreshToken %q, got %q", "refresh-token-value", tokenResponse.RefreshToken)
	}
	if tokenResponse.ExpiresIn != 3600 {
		testing.Errorf("expected ExpiresIn 3600, got %d", tokenResponse.ExpiresIn)
	}
}

func TestRedeemOneTimePasscodeReturnsInvalidPasscodeErrorOnInvalidCode(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"invalid_grant","suberror":"invalid_oob_value","error_description":"Invalid OTP"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemOneTimePasscode(
		context.Background(),
		"challenge-token-456",
		"wrong-code",
	)

	var invalidPasscodeError InvalidPasscodeError
	isMatch := errors.As(tokenError, &invalidPasscodeError)
	if !isMatch {
		testing.Errorf("expected InvalidPasscodeError, got %T: %v", tokenError, tokenError)
	}
}

func TestRedeemPasswordReturnsTokenResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"password-access-token",
			"id_token":"password-id-token",
			"refresh_token":"password-refresh-token",
			"expires_in":7200
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	tokenResponse, tokenError := testClient.RedeemPassword(
		context.Background(),
		"challenge-token-789",
		"correct-password",
	)

	if tokenError != nil {
		testing.Fatalf("expected no error, got %v", tokenError)
	}
	if tokenResponse.AccessToken != "password-access-token" {
		testing.Errorf("expected AccessToken %q, got %q", "password-access-token", tokenResponse.AccessToken)
	}
}

func TestRedeemPasswordReturnsMultifactorRequiredErrorWhenMfaRequired(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{
			"error":"invalid_grant",
			"suberror":"mfa_required",
			"error_description":"MFA is required",
			"continuation_token":"mfa-continuation-token"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemPassword(
		context.Background(),
		"challenge-token-789",
		"correct-password",
	)

	var multifactorRequiredError MultifactorRequiredError
	isMatch := errors.As(tokenError, &multifactorRequiredError)
	if !isMatch {
		testing.Fatalf("expected MultifactorRequiredError, got %T: %v", tokenError, tokenError)
	}
	if multifactorRequiredError.ContinuationToken != "mfa-continuation-token" {
		testing.Errorf("expected ContinuationToken %q, got %q", "mfa-continuation-token", multifactorRequiredError.ContinuationToken)
	}
}

func TestRedeemPasswordReturnsRegistrationRequiredErrorWhenRegistrationRequired(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{
			"error":"invalid_grant",
			"suberror":"registration_required",
			"error_description":"Registration is required",
			"continuation_token":"registration-continuation-token"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemPassword(
		context.Background(),
		"challenge-token-789",
		"correct-password",
	)

	var registrationRequiredError RegistrationRequiredError
	isMatch := errors.As(tokenError, &registrationRequiredError)
	if !isMatch {
		testing.Fatalf("expected RegistrationRequiredError, got %T: %v", tokenError, tokenError)
	}
	if registrationRequiredError.ContinuationToken != "registration-continuation-token" {
		testing.Errorf("expected ContinuationToken %q, got %q", "registration-continuation-token", registrationRequiredError.ContinuationToken)
	}
}

func TestRedeemPasswordReturnsInvalidPasswordErrorOnWrongPassword(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"invalid_grant","suberror":"password_is_invalid","error_description":"Wrong password"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemPassword(
		context.Background(),
		"challenge-token-789",
		"wrong-password",
	)

	var invalidPasswordError InvalidPasswordError
	isMatch := errors.As(tokenError, &invalidPasswordError)
	if !isMatch {
		testing.Errorf("expected InvalidPasswordError, got %T: %v", tokenError, tokenError)
	}
}

type redirectTransport struct {
	targetBaseUrl string
}

func (transport *redirectTransport) RoundTrip(originalRequest *http.Request) (*http.Response, error) {
	redirectedUrl := transport.targetBaseUrl + originalRequest.URL.Path
	parsedUrl, parseError := url.Parse(redirectedUrl)
	if parseError != nil {
		return nil, parseError
	}
	redirectedRequest := originalRequest.Clone(originalRequest.Context())
	redirectedRequest.URL = parsedUrl
	result, roundTripError := http.DefaultTransport.RoundTrip(redirectedRequest)
	return result, roundTripError
}

func createTestClient(server *httptest.Server) Client {
	testConfiguration := Configuration{
		TenantSubdomain:  "testserver",
		ClientIdentifier: "test-client-id",
	}
	testHttpClient := &http.Client{
		Transport: &redirectTransport{
			targetBaseUrl: server.URL,
		},
	}
	result := NewClientWithHttpClient(testConfiguration, testHttpClient)
	return result
}

func parseFormBody(formBody string) map[string]string {
	parsedValues, parseError := url.ParseQuery(formBody)
	if parseError != nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for parameterName, parameterValues := range parsedValues {
		if len(parameterValues) > 0 {
			result[parameterName] = parameterValues[0]
		}
	}
	return result
}
