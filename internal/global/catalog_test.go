package global_test

import (
	"encoding/json"
	"os"
	"testing"

	"golang.org/x/text/language"

	"github.com/Grit-Software-Systems/entrance/internal/global"
)

func TestMain(main *testing.M) {
	translationsContent, readError := os.ReadFile("../../locales/en.json")
	if readError != nil {
		panic(readError)
	}
	global.LoadTranslations(language.English, translationsContent)
	exitCode := main.Run()
	os.Exit(exitCode)
}

func TestTranslateMessageAccessDenied(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageAccessDenied)
	expectedMessage := "Access was denied. The request was blocked by fraud protection or administrator policy."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageExpiredToken(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageExpiredToken)
	expectedMessage := "The continuation token has expired. Please restart the authentication flow."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageHttpRequestFailed(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageHttpRequestFailed)
	expectedMessage := "The HTTP request to the authentication server failed."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageInvalidPasscode(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageInvalidPasscode)
	expectedMessage := "The one-time passcode is incorrect."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageInvalidPassword(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageInvalidPassword)
	expectedMessage := "The password is incorrect."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageMultifactorRequired(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageMultifactorRequired)
	expectedMessage := "Multifactor authentication is required. Complete the additional verification step to proceed."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageRedirectRequired(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageRedirectRequired)
	expectedMessage := "Native authentication is not available. The client must fall back to browser-based authentication."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageRegistrationRequired(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageRegistrationRequired)
	expectedMessage := "A strong authentication method must be registered before sign-in can complete."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageResponseParsingFailed(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageResponseParsingFailed)
	expectedMessage := "The response from the authentication server could not be parsed."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageUnexpectedError(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageUnexpectedError)
	expectedMessage := "An unexpected authentication error occurred."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestTranslateMessageUserNotFound(testing *testing.T) {
	translatedMessage := global.TranslateMessage(language.English, global.MessageUserNotFound)
	expectedMessage := "The user account was not found."
	if translatedMessage != expectedMessage {
		testing.Errorf("expected %q, got %q", expectedMessage, translatedMessage)
	}
}

func TestAllMessageKeysHaveTranslations(testing *testing.T) {
	translationsContent, readError := os.ReadFile("../../locales/en.json")
	if readError != nil {
		testing.Fatalf("failed to read en.json: %v", readError)
	}
	var translations map[string]string
	unmarshalError := json.Unmarshal(translationsContent, &translations)
	if unmarshalError != nil {
		testing.Fatalf("failed to unmarshal en.json: %v", unmarshalError)
	}
	messageKeys := []string{
		global.MessageAccessDenied,
		global.MessageExpiredToken,
		global.MessageHttpRequestFailed,
		global.MessageInvalidPasscode,
		global.MessageInvalidPassword,
		global.MessageMultifactorRequired,
		global.MessageRedirectRequired,
		global.MessageRegistrationRequired,
		global.MessageResponseParsingFailed,
		global.MessageUnexpectedError,
		global.MessageUserNotFound,
	}
	for _, item := range messageKeys {
		translatedText, exists := translations[item]
		if !exists {
			testing.Errorf("message key %q is missing from en.json", item)
		}
		if translatedText == "" {
			testing.Errorf("message key %q has an empty translation in en.json", item)
		}
	}
}
