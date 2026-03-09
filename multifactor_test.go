package entrance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIntrospectReturnsIntrospectResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"introspect-token-789",
			"methods":[
				{"id":"email-method-1","challenge_type":"oob","challenge_channel":"email","login_hint":"c***r@co**o.com"},
				{"id":"sms-method-2","challenge_type":"oob","challenge_channel":"sms","login_hint":"+1***56"}
			]
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	introspectResponse, introspectError := testClient.Introspect(
		context.Background(),
		"mfa-continuation-token",
	)

	if introspectError != nil {
		testing.Fatalf("expected no error, got %v", introspectError)
	}
	if introspectResponse.ContinuationToken != "introspect-token-789" {
		testing.Errorf("expected ContinuationToken %q, got %q", "introspect-token-789", introspectResponse.ContinuationToken)
	}
	if len(introspectResponse.Methods) != 2 {
		testing.Fatalf("expected 2 methods, got %d", len(introspectResponse.Methods))
	}
	firstMethod := introspectResponse.Methods[0]
	if firstMethod.Identifier != "email-method-1" {
		testing.Errorf("expected first method Identifier %q, got %q", "email-method-1", firstMethod.Identifier)
	}
	if firstMethod.ChallengeType != ChallengeMethodOneTimePasscode {
		testing.Errorf("expected first method ChallengeType %q, got %q", ChallengeMethodOneTimePasscode, firstMethod.ChallengeType)
	}
	if firstMethod.ChallengeChannel != "email" {
		testing.Errorf("expected first method ChallengeChannel %q, got %q", "email", firstMethod.ChallengeChannel)
	}
	if firstMethod.LoginHint != "c***r@co**o.com" {
		testing.Errorf("expected first method LoginHint %q, got %q", "c***r@co**o.com", firstMethod.LoginHint)
	}
	secondMethod := introspectResponse.Methods[1]
	if secondMethod.Identifier != "sms-method-2" {
		testing.Errorf("expected second method Identifier %q, got %q", "sms-method-2", secondMethod.Identifier)
	}
	if secondMethod.ChallengeChannel != "sms" {
		testing.Errorf("expected second method ChallengeChannel %q, got %q", "sms", secondMethod.ChallengeChannel)
	}
}

func TestIntrospectReturnsExpiredTokenErrorOnExpiredToken(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"expired_token","error_description":"Continuation token has expired"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, introspectError := testClient.Introspect(
		context.Background(),
		"expired-token",
	)

	var expiredTokenError ExpiredTokenError
	isMatch := errors.As(introspectError, &expiredTokenError)
	if !isMatch {
		testing.Errorf("expected ExpiredTokenError, got %T: %v", introspectError, introspectError)
	}
}

func TestChallengeMultifactorReturnsChallengeResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"mfa-challenge-token-456",
			"challenge_type":"oob",
			"challenge_channel":"email",
			"code_length":8,
			"challenge_target_label":"c***r@co**o.com"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	challengeResponse, challengeError := testClient.ChallengeMultifactor(
		context.Background(),
		"introspect-token-789",
		"email-method-1",
	)

	if challengeError != nil {
		testing.Fatalf("expected no error, got %v", challengeError)
	}
	if challengeResponse.ContinuationToken != "mfa-challenge-token-456" {
		testing.Errorf("expected ContinuationToken %q, got %q", "mfa-challenge-token-456", challengeResponse.ContinuationToken)
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

func TestChallengeMultifactorSendsMethodIdentifier(testing *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		bodyBytes := make([]byte, httpRequest.ContentLength)
		httpRequest.Body.Read(bodyBytes)
		capturedBody = string(bodyBytes)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"token",
			"challenge_type":"oob",
			"challenge_channel":"email",
			"code_length":8,
			"challenge_target_label":"masked"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	testClient.ChallengeMultifactor(
		context.Background(),
		"introspect-token",
		"email-method-1",
	)

	if capturedBody == "" {
		testing.Fatal("expected request body to be captured")
	}
	parsedForm := parseFormBody(capturedBody)
	methodIdentifierValue := parsedForm["id"]
	if methodIdentifierValue != "email-method-1" {
		testing.Errorf("expected id %q, got %q", "email-method-1", methodIdentifierValue)
	}
}

func TestChallengeMultifactorReturnsRedirectRequiredErrorOnRedirect(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"token","challenge_type":"redirect"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, challengeError := testClient.ChallengeMultifactor(
		context.Background(),
		"introspect-token",
		"method-1",
	)

	var redirectRequiredError RedirectRequiredError
	isMatch := errors.As(challengeError, &redirectRequiredError)
	if !isMatch {
		testing.Errorf("expected RedirectRequiredError, got %T: %v", challengeError, challengeError)
	}
}

func TestRedeemMultifactorOneTimePasscodeReturnsTokenResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"mfa-access-token",
			"id_token":"mfa-id-token",
			"refresh_token":"mfa-refresh-token",
			"expires_in":3600
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	tokenResponse, tokenError := testClient.RedeemMultifactorOneTimePasscode(
		context.Background(),
		"mfa-challenge-token-456",
		"12345678",
	)

	if tokenError != nil {
		testing.Fatalf("expected no error, got %v", tokenError)
	}
	if tokenResponse.TokenType != "Bearer" {
		testing.Errorf("expected TokenType %q, got %q", "Bearer", tokenResponse.TokenType)
	}
	if tokenResponse.AccessToken != "mfa-access-token" {
		testing.Errorf("expected AccessToken %q, got %q", "mfa-access-token", tokenResponse.AccessToken)
	}
	if tokenResponse.IdToken != "mfa-id-token" {
		testing.Errorf("expected IdToken %q, got %q", "mfa-id-token", tokenResponse.IdToken)
	}
	if tokenResponse.RefreshToken != "mfa-refresh-token" {
		testing.Errorf("expected RefreshToken %q, got %q", "mfa-refresh-token", tokenResponse.RefreshToken)
	}
	if tokenResponse.ExpiresIn != 3600 {
		testing.Errorf("expected ExpiresIn 3600, got %d", tokenResponse.ExpiresIn)
	}
}

func TestRedeemMultifactorOneTimePasscodeReturnsInvalidPasscodeErrorOnInvalidCode(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"invalid_grant","suberror":"invalid_oob_value","error_description":"Invalid OTP"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemMultifactorOneTimePasscode(
		context.Background(),
		"mfa-challenge-token",
		"wrong-code",
	)

	var invalidPasscodeError InvalidPasscodeError
	isMatch := errors.As(tokenError, &invalidPasscodeError)
	if !isMatch {
		testing.Errorf("expected InvalidPasscodeError, got %T: %v", tokenError, tokenError)
	}
}

func TestRedeemMultifactorOneTimePasscodeUsesTenantIdentifier(testing *testing.T) {
	var capturedUrl string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		capturedUrl = httpRequest.URL.Path
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"token",
			"id_token":"token",
			"refresh_token":"token",
			"expires_in":3600
		}`))
	}))
	defer server.Close()

	tenantIdentifier := "00000000-0000-0000-0000-000000000000"
	testConfiguration := Configuration{
		TenantSubdomain:  "testserver",
		TenantIdentifier: tenantIdentifier,
		ClientIdentifier: "test-client-id",
	}
	testHttpClient := &http.Client{
		Transport: &redirectTransport{
			targetBaseUrl: server.URL,
		},
	}
	testClient := NewClientWithHttpClient(testConfiguration, testHttpClient)
	testClient.RedeemMultifactorOneTimePasscode(
		context.Background(),
		"mfa-token",
		"12345678",
	)

	expectedPathSegment := tenantIdentifier + "/oauth2/v2.0/token"
	if !strings.Contains(capturedUrl, expectedPathSegment) {
		testing.Errorf("expected URL path to contain %q, got %q", expectedPathSegment, capturedUrl)
	}
}

func TestRedeemMultifactorOneTimePasscodeSendsCorrectGrantType(testing *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		bodyBytes := make([]byte, httpRequest.ContentLength)
		httpRequest.Body.Read(bodyBytes)
		capturedBody = string(bodyBytes)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"token",
			"id_token":"token",
			"refresh_token":"token",
			"expires_in":3600
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	testClient.RedeemMultifactorOneTimePasscode(
		context.Background(),
		"mfa-token",
		"12345678",
	)

	if capturedBody == "" {
		testing.Fatal("expected request body to be captured")
	}
	parsedForm := parseFormBody(capturedBody)
	grantTypeValue := parsedForm["grant_type"]
	if grantTypeValue != "mfa_oob" {
		testing.Errorf("expected grant_type %q, got %q", "mfa_oob", grantTypeValue)
	}
}
