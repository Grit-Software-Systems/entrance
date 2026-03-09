package entrance

import (
	"golang.org/x/text/language"

	"github.com/Grit-Software-Systems/entrance/internal/global"
)

type AccessDeniedError struct {
	AuthenticationError
}

func (accessDeniedError AccessDeniedError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageAccessDenied)
	return result
}

type AuthenticationError struct {
	ErrorCode        string
	SubErrorCode     string
	ErrorDescription string
	CorrelationId    string
	Timestamp        string
}

func (authenticationError AuthenticationError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageUnexpectedError)
	return result
}

type ExpiredTokenError struct {
	AuthenticationError
}

func (expiredTokenError ExpiredTokenError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageExpiredToken)
	return result
}

type InvalidPasscodeError struct {
	AuthenticationError
}

func (invalidPasscodeError InvalidPasscodeError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageInvalidPasscode)
	return result
}

type InvalidPasswordError struct {
	AuthenticationError
}

func (invalidPasswordError InvalidPasswordError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageInvalidPassword)
	return result
}

type MultifactorRequiredError struct {
	AuthenticationError
	ContinuationToken string
}

func (multifactorRequiredError MultifactorRequiredError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageMultifactorRequired)
	return result
}

type RedirectRequiredError struct {
	AuthenticationError
}

func (redirectRequiredError RedirectRequiredError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageRedirectRequired)
	return result
}

type RegistrationRequiredError struct {
	AuthenticationError
	ContinuationToken string
}

func (registrationRequiredError RegistrationRequiredError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageRegistrationRequired)
	return result
}

type RequestError struct {
	Message    string
	Underlying error
}

func (requestError RequestError) Error() string {
	result := requestError.Message
	return result
}

func (requestError RequestError) Unwrap() error {
	result := requestError.Underlying
	return result
}

type UserNotFoundError struct {
	AuthenticationError
}

func (userNotFoundError UserNotFoundError) Error() string {
	result := global.TranslateMessage(language.English, global.MessageUserNotFound)
	return result
}
