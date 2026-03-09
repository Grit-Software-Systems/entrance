package entrance

import (
	"errors"
	"fmt"
	"testing"
)

func TestClassifyErrorUserNotFound(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:        "user_not_found",
		ErrorDescription: "User does not exist",
		CorrelationId:    "abc-123",
		Timestamp:        "2024-01-01",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var userNotFoundError UserNotFoundError
	isMatch := errors.As(classifiedError, &userNotFoundError)
	if !isMatch {
		testing.Errorf("expected UserNotFoundError, got %T", classifiedError)
	}
}

func TestClassifyErrorInvalidPasscode(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "invalid_grant",
		SubErrorCode: "invalid_oob_value",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var invalidPasscodeError InvalidPasscodeError
	isMatch := errors.As(classifiedError, &invalidPasscodeError)
	if !isMatch {
		testing.Errorf("expected InvalidPasscodeError, got %T", classifiedError)
	}
}

func TestClassifyErrorInvalidPassword(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "invalid_grant",
		SubErrorCode: "password_is_invalid",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var invalidPasswordError InvalidPasswordError
	isMatch := errors.As(classifiedError, &invalidPasswordError)
	if !isMatch {
		testing.Errorf("expected InvalidPasswordError, got %T", classifiedError)
	}
}

func TestClassifyErrorExpiredToken(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode: "expired_token",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var expiredTokenError ExpiredTokenError
	isMatch := errors.As(classifiedError, &expiredTokenError)
	if !isMatch {
		testing.Errorf("expected ExpiredTokenError, got %T", classifiedError)
	}
}

func TestClassifyErrorMultifactorRequired(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "invalid_grant",
		SubErrorCode: "mfa_required",
	}
	continuationToken := "mfa-continuation-token"
	classifiedError := classifyByErrorAndSubError(authenticationError, continuationToken)
	var multifactorRequiredError MultifactorRequiredError
	isMatch := errors.As(classifiedError, &multifactorRequiredError)
	if !isMatch {
		testing.Errorf("expected MultifactorRequiredError, got %T", classifiedError)
	}
	if multifactorRequiredError.ContinuationToken != continuationToken {
		testing.Errorf("expected ContinuationToken %q, got %q", continuationToken, multifactorRequiredError.ContinuationToken)
	}
}

func TestClassifyErrorRegistrationRequired(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "invalid_grant",
		SubErrorCode: "registration_required",
	}
	continuationToken := "registration-continuation-token"
	classifiedError := classifyByErrorAndSubError(authenticationError, continuationToken)
	var registrationRequiredError RegistrationRequiredError
	isMatch := errors.As(classifiedError, &registrationRequiredError)
	if !isMatch {
		testing.Errorf("expected RegistrationRequiredError, got %T", classifiedError)
	}
	if registrationRequiredError.ContinuationToken != continuationToken {
		testing.Errorf("expected ContinuationToken %q, got %q", continuationToken, registrationRequiredError.ContinuationToken)
	}
}

func TestClassifyErrorAccessDeniedBlockedByAdmin(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "access_denied",
		SubErrorCode: "provider_blocked_by_admin",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var accessDeniedError AccessDeniedError
	isMatch := errors.As(classifiedError, &accessDeniedError)
	if !isMatch {
		testing.Errorf("expected AccessDeniedError, got %T", classifiedError)
	}
}

func TestClassifyErrorAccessDeniedBlockedByReputation(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:    "access_denied",
		SubErrorCode: "provider_blocked_by_rep",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var accessDeniedError AccessDeniedError
	isMatch := errors.As(classifiedError, &accessDeniedError)
	if !isMatch {
		testing.Errorf("expected AccessDeniedError, got %T", classifiedError)
	}
}

func TestClassifyErrorUnrecognizedFallsBackToAuthenticationError(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:        "some_unknown_error",
		ErrorDescription: "Something unexpected",
		CorrelationId:    "xyz-789",
		Timestamp:        "2024-06-15",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var matchedError AuthenticationError
	isMatch := errors.As(classifiedError, &matchedError)
	if !isMatch {
		testing.Errorf("expected AuthenticationError, got %T", classifiedError)
	}
	if matchedError.ErrorCode != "some_unknown_error" {
		testing.Errorf("expected ErrorCode %q, got %q", "some_unknown_error", matchedError.ErrorCode)
	}
}

func TestClassifyErrorPreservesAllPayloadFields(testing *testing.T) {
	authenticationError := AuthenticationError{
		ErrorCode:        "user_not_found",
		SubErrorCode:     "",
		ErrorDescription: "The user was not found",
		CorrelationId:    "corr-id-456",
		Timestamp:        "2024-03-20T10:00:00Z",
	}
	classifiedError := classifyByErrorAndSubError(authenticationError, "")
	var userNotFoundError UserNotFoundError
	isMatch := errors.As(classifiedError, &userNotFoundError)
	if !isMatch {
		testing.Fatalf("expected UserNotFoundError, got %T", classifiedError)
	}
	if userNotFoundError.ErrorCode != "user_not_found" {
		testing.Errorf("expected ErrorCode %q, got %q", "user_not_found", userNotFoundError.ErrorCode)
	}
	if userNotFoundError.ErrorDescription != "The user was not found" {
		testing.Errorf("expected ErrorDescription %q, got %q", "The user was not found", userNotFoundError.ErrorDescription)
	}
	if userNotFoundError.CorrelationId != "corr-id-456" {
		testing.Errorf("expected CorrelationId %q, got %q", "corr-id-456", userNotFoundError.CorrelationId)
	}
	if userNotFoundError.Timestamp != "2024-03-20T10:00:00Z" {
		testing.Errorf("expected Timestamp %q, got %q", "2024-03-20T10:00:00Z", userNotFoundError.Timestamp)
	}
}

func TestClassifyErrorWrapsNonPayloadErrorAsRequestError(testing *testing.T) {
	originalError := fmt.Errorf("network timeout")
	classifiedError := classifyError(originalError)
	var requestError RequestError
	isMatch := errors.As(classifiedError, &requestError)
	if !isMatch {
		testing.Errorf("expected RequestError, got %T", classifiedError)
	}
	unwrappedError := requestError.Unwrap()
	if unwrappedError != originalError {
		testing.Errorf("expected unwrapped error to be the original error")
	}
}
