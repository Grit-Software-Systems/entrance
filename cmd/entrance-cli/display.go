package main

import (
	"errors"
	"fmt"

	"github.com/Grit-Software-Systems/entrance"
)

const (
	tokenDisplayLength = 40
)

func printAuthenticationMethods(methods []entrance.AuthenticationMethod) {
	for index := 0; index < len(methods); index++ {
		var method entrance.AuthenticationMethod = methods[index]
		fmt.Printf("  %d. %s %s\n", index+1, method.ChallengeChannel, method.LoginHint)
	}
}

func printChallengeDetails(challengeTargetLabel string, codeLength int) {
	fmt.Printf("Challenge sent to: %s\n", challengeTargetLabel)
	fmt.Printf("Expected code length: %d\n", codeLength)
}

func printError(errorValue error) {
	var recognized bool = printRecognizedError(errorValue)
	if recognized {
		return
	}
	var requestError entrance.RequestError
	if errors.As(errorValue, &requestError) {
		fmt.Printf("HTTP error: %s\n", requestError.Error())
		return
	}
	var authenticationError entrance.AuthenticationError
	if errors.As(errorValue, &authenticationError) {
		fmt.Printf("Authentication error: %s\n", authenticationError.Error())
		fmt.Printf("  Error code: %s\n", authenticationError.ErrorCode)
		fmt.Printf("  Suberror: %s\n", authenticationError.SubErrorCode)
		fmt.Printf("  Description: %s\n", authenticationError.ErrorDescription)
		return
	}
	fmt.Printf("Error: %s\n", errorValue.Error())
}

func printRecognizedError(errorValue error) bool {
	var accessDeniedError entrance.AccessDeniedError
	if errors.As(errorValue, &accessDeniedError) {
		fmt.Println("Access denied (SMS fraud protection).")
		result := true
		return result
	}
	var expiredTokenError entrance.ExpiredTokenError
	if errors.As(errorValue, &expiredTokenError) {
		fmt.Println("Session expired — please restart the flow.")
		result := true
		return result
	}
	var invalidPasscodeError entrance.InvalidPasscodeError
	if errors.As(errorValue, &invalidPasscodeError) {
		fmt.Println("Invalid passcode.")
		result := true
		return result
	}
	var invalidPasswordError entrance.InvalidPasswordError
	if errors.As(errorValue, &invalidPasswordError) {
		fmt.Println("Invalid password.")
		result := true
		return result
	}
	var redirectRequiredError entrance.RedirectRequiredError
	if errors.As(errorValue, &redirectRequiredError) {
		fmt.Println("Native auth unavailable — browser auth required.")
		result := true
		return result
	}
	var userNotFoundError entrance.UserNotFoundError
	if errors.As(errorValue, &userNotFoundError) {
		fmt.Println("User not found.")
		result := true
		return result
	}
	result := false
	return result
}

func printTokenResponse(tokenResponse entrance.TokenResponse, showFullTokens bool) {
	var accessTokenDisplay string = truncateToken(tokenResponse.AccessToken, showFullTokens)
	var idTokenDisplay string = truncateToken(tokenResponse.IdToken, showFullTokens)
	var refreshTokenDisplay string = truncateToken(tokenResponse.RefreshToken, showFullTokens)

	fmt.Println("Authentication succeeded.")
	fmt.Printf("  Access Token:  %s\n", accessTokenDisplay)
	fmt.Printf("  ID Token:      %s\n", idTokenDisplay)
	fmt.Printf("  Refresh Token: %s\n", refreshTokenDisplay)
	fmt.Printf("  Expires In:    %d\n", tokenResponse.ExpiresIn)
	fmt.Printf("  Token Type:    %s\n", tokenResponse.TokenType)
}

func truncateToken(tokenValue string, showFullTokens bool) string {
	if showFullTokens {
		result := tokenValue
		return result
	}
	if len(tokenValue) <= tokenDisplayLength {
		result := tokenValue
		return result
	}
	result := tokenValue[:tokenDisplayLength] + "..."
	return result
}
