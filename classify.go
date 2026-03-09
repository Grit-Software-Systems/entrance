package entrance

import (
	"errors"

	"github.com/Grit-Software-Systems/entrance/internal/request"
)

const (
	errorCodeAccessDenied = "access_denied"
	errorCodeExpiredToken = "expired_token"
	errorCodeInvalidGrant = "invalid_grant"
	errorCodeUserNotFound = "user_not_found"

	subErrorCodeInvalidOneTimePasscodeValue = "invalid_oob_value"
	subErrorCodeMultifactorRequired         = "mfa_required"
	subErrorCodePasswordIsInvalid           = "password_is_invalid"
	subErrorCodeProviderBlockedByAdmin      = "provider_blocked_by_admin"
	subErrorCodeProviderBlockedByReputation = "provider_blocked_by_rep"
	subErrorCodeRegistrationRequired        = "registration_required"
)

func classifyAccessDenied(authenticationError AuthenticationError) error {
	switch authenticationError.SubErrorCode {
	case subErrorCodeProviderBlockedByAdmin, subErrorCodeProviderBlockedByReputation:
		result := AccessDeniedError{AuthenticationError: authenticationError}
		return result
	default:
		return authenticationError
	}
}

func classifyByErrorAndSubError(authenticationError AuthenticationError, continuationToken string) error {
	switch authenticationError.ErrorCode {
	case errorCodeAccessDenied:
		result := classifyAccessDenied(authenticationError)
		return result
	case errorCodeExpiredToken:
		result := ExpiredTokenError{AuthenticationError: authenticationError}
		return result
	case errorCodeInvalidGrant:
		result := classifyInvalidGrant(authenticationError, continuationToken)
		return result
	case errorCodeUserNotFound:
		result := UserNotFoundError{AuthenticationError: authenticationError}
		return result
	default:
		return authenticationError
	}
}

func classifyError(sendError error) error {
	var errorPayload *request.ErrorPayload
	isPayloadError := errors.As(sendError, &errorPayload)
	if !isPayloadError {
		result := RequestError{
			Message:    sendError.Error(),
			Underlying: sendError,
		}
		return result
	}
	authenticationError := AuthenticationError{
		ErrorCode:        errorPayload.ErrorCode,
		SubErrorCode:     errorPayload.SubErrorCode,
		ErrorDescription: errorPayload.ErrorDescription,
		CorrelationId:    errorPayload.CorrelationId,
		Timestamp:        errorPayload.Timestamp,
	}
	result := classifyByErrorAndSubError(authenticationError, errorPayload.ContinuationToken)
	return result
}

func classifyInvalidGrant(authenticationError AuthenticationError, continuationToken string) error {
	switch authenticationError.SubErrorCode {
	case subErrorCodeInvalidOneTimePasscodeValue:
		result := InvalidPasscodeError{AuthenticationError: authenticationError}
		return result
	case subErrorCodeMultifactorRequired:
		result := MultifactorRequiredError{
			AuthenticationError: authenticationError,
			ContinuationToken:   continuationToken,
		}
		return result
	case subErrorCodePasswordIsInvalid:
		result := InvalidPasswordError{AuthenticationError: authenticationError}
		return result
	case subErrorCodeRegistrationRequired:
		result := RegistrationRequiredError{
			AuthenticationError: authenticationError,
			ContinuationToken:   continuationToken,
		}
		return result
	default:
		return authenticationError
	}
}
