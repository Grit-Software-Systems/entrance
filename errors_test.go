package entrance_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Grit-Software-Systems/entrance"
)

func TestAccessDeniedErrorMessage(testing *testing.T) {
	accessDeniedError := entrance.AccessDeniedError{}
	expectedMessage := "Access was denied. The request was blocked by fraud protection or administrator policy."
	if accessDeniedError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, accessDeniedError.Error())
	}
}

func TestAuthenticationErrorMessage(testing *testing.T) {
	authenticationError := entrance.AuthenticationError{}
	expectedMessage := "An unexpected authentication error occurred."
	if authenticationError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, authenticationError.Error())
	}
}

func TestExpiredTokenErrorMessage(testing *testing.T) {
	expiredTokenError := entrance.ExpiredTokenError{}
	expectedMessage := "The continuation token has expired. Please restart the authentication flow."
	if expiredTokenError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, expiredTokenError.Error())
	}
}

func TestInvalidPasscodeErrorMessage(testing *testing.T) {
	invalidPasscodeError := entrance.InvalidPasscodeError{}
	expectedMessage := "The one-time passcode is incorrect."
	if invalidPasscodeError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, invalidPasscodeError.Error())
	}
}

func TestInvalidPasswordErrorMessage(testing *testing.T) {
	invalidPasswordError := entrance.InvalidPasswordError{}
	expectedMessage := "The password is incorrect."
	if invalidPasswordError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, invalidPasswordError.Error())
	}
}

func TestMultifactorRequiredErrorMessage(testing *testing.T) {
	multifactorRequiredError := entrance.MultifactorRequiredError{}
	expectedMessage := "Multifactor authentication is required. Complete the additional verification step to proceed."
	if multifactorRequiredError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, multifactorRequiredError.Error())
	}
}

func TestRedirectRequiredErrorMessage(testing *testing.T) {
	redirectRequiredError := entrance.RedirectRequiredError{}
	expectedMessage := "Native authentication is not available. The client must fall back to browser-based authentication."
	if redirectRequiredError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, redirectRequiredError.Error())
	}
}

func TestRegistrationRequiredErrorMessage(testing *testing.T) {
	registrationRequiredError := entrance.RegistrationRequiredError{}
	expectedMessage := "A strong authentication method must be registered before sign-in can complete."
	if registrationRequiredError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, registrationRequiredError.Error())
	}
}

func TestUserNotFoundErrorMessage(testing *testing.T) {
	userNotFoundError := entrance.UserNotFoundError{}
	expectedMessage := "The user account was not found."
	if userNotFoundError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, userNotFoundError.Error())
	}
}

func TestRequestErrorMessage(testing *testing.T) {
	requestError := entrance.RequestError{
		Message:    "something went wrong",
		Underlying: nil,
	}
	expectedMessage := "something went wrong"
	if requestError.Error() != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, requestError.Error())
	}
}

func TestRequestErrorUnwrap(testing *testing.T) {
	underlyingError := fmt.Errorf("network timeout")
	requestError := entrance.RequestError{
		Message:    "request failed",
		Underlying: underlyingError,
	}
	unwrappedError := requestError.Unwrap()
	if unwrappedError != underlyingError {
		testing.Errorf("expected underlying error %v, got %v", underlyingError, unwrappedError)
	}
}

func TestRequestErrorUnwrapNil(testing *testing.T) {
	requestError := entrance.RequestError{
		Message:    "request failed",
		Underlying: nil,
	}
	unwrappedError := requestError.Unwrap()
	if unwrappedError != nil {
		testing.Errorf("expected nil, got %v", unwrappedError)
	}
}

func TestErrorsAsAccessDeniedError(testing *testing.T) {
	var targetError entrance.AccessDeniedError
	sourceError := entrance.AccessDeniedError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "access_denied"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match AccessDeniedError")
	}
}

func TestErrorsAsExpiredTokenError(testing *testing.T) {
	var targetError entrance.ExpiredTokenError
	sourceError := entrance.ExpiredTokenError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "expired_token"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match ExpiredTokenError")
	}
}

func TestErrorsAsInvalidPasscodeError(testing *testing.T) {
	var targetError entrance.InvalidPasscodeError
	sourceError := entrance.InvalidPasscodeError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "invalid_grant"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match InvalidPasscodeError")
	}
}

func TestErrorsAsInvalidPasswordError(testing *testing.T) {
	var targetError entrance.InvalidPasswordError
	sourceError := entrance.InvalidPasswordError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "invalid_grant"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match InvalidPasswordError")
	}
}

func TestErrorsAsMultifactorRequiredError(testing *testing.T) {
	var targetError entrance.MultifactorRequiredError
	sourceError := entrance.MultifactorRequiredError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "invalid_grant"},
		ContinuationToken:  "token123",
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match MultifactorRequiredError")
	}
	if targetError.ContinuationToken != "token123" {
		testing.Errorf("expected ContinuationToken %q, got %q", "token123", targetError.ContinuationToken)
	}
}

func TestErrorsAsRegistrationRequiredError(testing *testing.T) {
	var targetError entrance.RegistrationRequiredError
	sourceError := entrance.RegistrationRequiredError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "invalid_grant"},
		ContinuationToken:  "regtoken456",
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match RegistrationRequiredError")
	}
	if targetError.ContinuationToken != "regtoken456" {
		testing.Errorf("expected ContinuationToken %q, got %q", "regtoken456", targetError.ContinuationToken)
	}
}

func TestErrorsAsRedirectRequiredError(testing *testing.T) {
	var targetError entrance.RedirectRequiredError
	sourceError := entrance.RedirectRequiredError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "redirect"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match RedirectRequiredError")
	}
}

func TestErrorsAsUserNotFoundError(testing *testing.T) {
	var targetError entrance.UserNotFoundError
	sourceError := entrance.UserNotFoundError{
		AuthenticationError: entrance.AuthenticationError{ErrorCode: "user_not_found"},
	}
	isMatch := errors.As(sourceError, &targetError)
	if !isMatch {
		testing.Error("expected errors.As to match UserNotFoundError")
	}
}
