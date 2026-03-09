package entrance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntrospectRegistrationReturnsRegistrationIntrospectResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"registration-introspect-token-123",
			"methods":[
				{"id":"email","challenge_type":"oob","challenge_channel":"email","login_hint":"caseyjensen@contoso.com"},
				{"id":"sms","challenge_type":"oob","challenge_channel":"sms"}
			]
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	introspectResponse, introspectError := testClient.IntrospectRegistration(
		context.Background(),
		"registration-continuation-token",
	)

	if introspectError != nil {
		testing.Fatalf("expected no error, got %v", introspectError)
	}
	if introspectResponse.ContinuationToken != "registration-introspect-token-123" {
		testing.Errorf("expected ContinuationToken %q, got %q", "registration-introspect-token-123", introspectResponse.ContinuationToken)
	}
	if len(introspectResponse.Methods) != 2 {
		testing.Fatalf("expected 2 methods, got %d", len(introspectResponse.Methods))
	}
	firstMethod := introspectResponse.Methods[0]
	if firstMethod.Identifier != "email" {
		testing.Errorf("expected first method Identifier %q, got %q", "email", firstMethod.Identifier)
	}
	if firstMethod.ChallengeType != ChallengeMethodOneTimePasscode {
		testing.Errorf("expected first method ChallengeType %q, got %q", ChallengeMethodOneTimePasscode, firstMethod.ChallengeType)
	}
	if firstMethod.ChallengeChannel != "email" {
		testing.Errorf("expected first method ChallengeChannel %q, got %q", "email", firstMethod.ChallengeChannel)
	}
	if firstMethod.LoginHint != "caseyjensen@contoso.com" {
		testing.Errorf("expected first method LoginHint %q, got %q", "caseyjensen@contoso.com", firstMethod.LoginHint)
	}
	secondMethod := introspectResponse.Methods[1]
	if secondMethod.Identifier != "sms" {
		testing.Errorf("expected second method Identifier %q, got %q", "sms", secondMethod.Identifier)
	}
	if secondMethod.ChallengeChannel != "sms" {
		testing.Errorf("expected second method ChallengeChannel %q, got %q", "sms", secondMethod.ChallengeChannel)
	}
}

func TestIntrospectRegistrationReturnsExpiredTokenErrorOnExpiredToken(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"expired_token","error_description":"Continuation token has expired"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, introspectError := testClient.IntrospectRegistration(
		context.Background(),
		"expired-token",
	)

	var expiredTokenError ExpiredTokenError
	isMatch := errors.As(introspectError, &expiredTokenError)
	if !isMatch {
		testing.Errorf("expected ExpiredTokenError, got %T: %v", introspectError, introspectError)
	}
}

func TestChallengeRegistrationReturnsRegistrationChallengeResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"registration-challenge-token-456",
			"challenge_type":"oob",
			"challenge_channel":"email",
			"code_length":8,
			"challenge_target_label":"c***r@co**o.com"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	challengeResponse, challengeError := testClient.ChallengeRegistration(
		context.Background(),
		"registration-introspect-token-123",
		"email",
		"caseyjensen@contoso.com",
	)

	if challengeError != nil {
		testing.Fatalf("expected no error, got %v", challengeError)
	}
	if challengeResponse.ContinuationToken != "registration-challenge-token-456" {
		testing.Errorf("expected ContinuationToken %q, got %q", "registration-challenge-token-456", challengeResponse.ContinuationToken)
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

func TestChallengeRegistrationReturnsPreverifiedChallengeType(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"continuation_token":"registration-preverified-token-789",
			"challenge_type":"preverified",
			"challenge_channel":"email"
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	challengeResponse, challengeError := testClient.ChallengeRegistration(
		context.Background(),
		"registration-introspect-token-123",
		"email",
		"caseyjensen@contoso.com",
	)

	if challengeError != nil {
		testing.Fatalf("expected no error, got %v", challengeError)
	}
	if challengeResponse.ContinuationToken != "registration-preverified-token-789" {
		testing.Errorf("expected ContinuationToken %q, got %q", "registration-preverified-token-789", challengeResponse.ContinuationToken)
	}
	if challengeResponse.ChallengeType != "preverified" {
		testing.Errorf("expected ChallengeType %q, got %q", "preverified", challengeResponse.ChallengeType)
	}
}

func TestChallengeRegistrationSendsCorrectFormParameters(testing *testing.T) {
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
	testClient.ChallengeRegistration(
		context.Background(),
		"registration-token",
		"sms",
		"+15551234567",
	)

	if capturedBody == "" {
		testing.Fatal("expected request body to be captured")
	}
	parsedForm := parseFormBody(capturedBody)
	challengeTypeValue := parsedForm["challenge_type"]
	if challengeTypeValue != "oob" {
		testing.Errorf("expected challenge_type %q, got %q", "oob", challengeTypeValue)
	}
	challengeChannelValue := parsedForm["challenge_channel"]
	if challengeChannelValue != "sms" {
		testing.Errorf("expected challenge_channel %q, got %q", "sms", challengeChannelValue)
	}
	challengeTargetValue := parsedForm["challenge_target"]
	if challengeTargetValue != "+15551234567" {
		testing.Errorf("expected challenge_target %q, got %q", "+15551234567", challengeTargetValue)
	}
}

func TestContinueRegistrationReturnsRegistrationContinueResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"registration-continue-token-101"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	continueResponse, continueError := testClient.ContinueRegistration(
		context.Background(),
		"registration-challenge-token-456",
		"12345678",
	)

	if continueError != nil {
		testing.Fatalf("expected no error, got %v", continueError)
	}
	if continueResponse.ContinuationToken != "registration-continue-token-101" {
		testing.Errorf("expected ContinuationToken %q, got %q", "registration-continue-token-101", continueResponse.ContinuationToken)
	}
}

func TestContinueRegistrationReturnsInvalidPasscodeErrorOnWrongPasscode(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"invalid_grant","suberror":"invalid_oob_value","error_description":"Invalid OTP"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, continueError := testClient.ContinueRegistration(
		context.Background(),
		"registration-challenge-token-456",
		"wrong-code",
	)

	var invalidPasscodeError InvalidPasscodeError
	isMatch := errors.As(continueError, &invalidPasscodeError)
	if !isMatch {
		testing.Errorf("expected InvalidPasscodeError, got %T: %v", continueError, continueError)
	}
}

func TestRedeemContinuationTokenReturnsTokenResponseOnSuccess(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{
			"token_type":"Bearer",
			"access_token":"registration-access-token",
			"id_token":"registration-id-token",
			"refresh_token":"registration-refresh-token",
			"expires_in":3600
		}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	tokenResponse, tokenError := testClient.RedeemContinuationToken(
		context.Background(),
		"registration-continue-token-101",
	)

	if tokenError != nil {
		testing.Fatalf("expected no error, got %v", tokenError)
	}
	if tokenResponse.TokenType != "Bearer" {
		testing.Errorf("expected TokenType %q, got %q", "Bearer", tokenResponse.TokenType)
	}
	if tokenResponse.AccessToken != "registration-access-token" {
		testing.Errorf("expected AccessToken %q, got %q", "registration-access-token", tokenResponse.AccessToken)
	}
	if tokenResponse.IdToken != "registration-id-token" {
		testing.Errorf("expected IdToken %q, got %q", "registration-id-token", tokenResponse.IdToken)
	}
	if tokenResponse.RefreshToken != "registration-refresh-token" {
		testing.Errorf("expected RefreshToken %q, got %q", "registration-refresh-token", tokenResponse.RefreshToken)
	}
	if tokenResponse.ExpiresIn != 3600 {
		testing.Errorf("expected ExpiresIn 3600, got %d", tokenResponse.ExpiresIn)
	}
}

func TestRedeemContinuationTokenSendsCorrectGrantType(testing *testing.T) {
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
	testClient.RedeemContinuationToken(
		context.Background(),
		"registration-continue-token",
	)

	if capturedBody == "" {
		testing.Fatal("expected request body to be captured")
	}
	parsedForm := parseFormBody(capturedBody)
	grantTypeValue := parsedForm["grant_type"]
	if grantTypeValue != "continuation_token" {
		testing.Errorf("expected grant_type %q, got %q", "continuation_token", grantTypeValue)
	}
}

func TestRedeemContinuationTokenReturnsAccessDeniedErrorOnFraudProtection(testing *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"access_denied","suberror":"provider_blocked_by_rep","error_description":"Blocked by fraud protection"}`))
	}))
	defer server.Close()

	testClient := createTestClient(server)
	_, tokenError := testClient.RedeemContinuationToken(
		context.Background(),
		"registration-continue-token",
	)

	var accessDeniedError AccessDeniedError
	isMatch := errors.As(tokenError, &accessDeniedError)
	if !isMatch {
		testing.Errorf("expected AccessDeniedError, got %T: %v", tokenError, tokenError)
	}
}
