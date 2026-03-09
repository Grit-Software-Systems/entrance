package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/Grit-Software-Systems/entrance"
)

const (
	challengeTypePreverified = "preverified"
	maximumPasscodeAttempts   = 3

	multifactorMethodPrompt      = "Select an MFA method (enter number): "
	multifactorPasscodePrompt    = "Enter the MFA one-time passcode: "
	passcodePrompt               = "Enter the one-time passcode: "
	passwordPrompt               = "Enter password: "
	registrationMethodPrompt     = "Select a registration method (enter number): "
	registrationPasscodePrompt   = "Enter the registration one-time passcode: "
	registrationTargetPrompt     = "Enter the target value (e.g., email address or phone number): "
)

func attemptPasscodeRedemption(authenticationClient entrance.Client, requestContext context.Context, continuationToken string, showFullTokens bool) bool {
	var passcode string = readLine(passcodePrompt)
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemOneTimePasscode(requestContext, continuationToken, passcode)
	if redeemError != nil {
		var invalidPasscodeError entrance.InvalidPasscodeError
		if errors.As(redeemError, &invalidPasscodeError) {
			printError(redeemError)
			result := false
			return result
		}
		printError(redeemError)
		result := true
		return result
	}
	printTokenResponse(tokenResponse, showFullTokens)
	result := true
	return result
}

func handleMultifactorChallenge(authenticationClient entrance.Client, requestContext context.Context, continuationToken string, showFullTokens bool) {
	var introspectResponse entrance.IntrospectResponse
	var introspectError error
	introspectResponse, introspectError = authenticationClient.Introspect(requestContext, continuationToken)
	if introspectError != nil {
		printError(introspectError)
		return
	}
	printAuthenticationMethods(introspectResponse.Methods)
	var methodCount int = len(introspectResponse.Methods)
	var selectedIndex int = readChoice(multifactorMethodPrompt, 1, methodCount)
	var selectedMethod entrance.AuthenticationMethod = introspectResponse.Methods[selectedIndex-1]
	var challengeResponse entrance.ChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.ChallengeMultifactor(requestContext, introspectResponse.ContinuationToken, selectedMethod.Identifier)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	printChallengeDetails(challengeResponse.ChallengeTargetLabel, challengeResponse.CodeLength)
	var passcode string = readLine(multifactorPasscodePrompt)
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemMultifactorOneTimePasscode(requestContext, challengeResponse.ContinuationToken, passcode)
	if redeemError != nil {
		printError(redeemError)
		return
	}
	printTokenResponse(tokenResponse, showFullTokens)
}

func handleRegistrationChallenge(authenticationClient entrance.Client, requestContext context.Context, continuationToken string, showFullTokens bool) {
	var introspectResponse entrance.RegistrationIntrospectResponse
	var introspectError error
	introspectResponse, introspectError = authenticationClient.IntrospectRegistration(requestContext, continuationToken)
	if introspectError != nil {
		printError(introspectError)
		return
	}
	printAuthenticationMethods(introspectResponse.Methods)
	var methodCount int = len(introspectResponse.Methods)
	var selectedIndex int = readChoice(registrationMethodPrompt, 1, methodCount)
	var selectedMethod entrance.AuthenticationMethod = introspectResponse.Methods[selectedIndex-1]
	var targetValue string = readLine(registrationTargetPrompt)
	var challengeResponse entrance.RegistrationChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.ChallengeRegistration(requestContext, introspectResponse.ContinuationToken, selectedMethod.ChallengeChannel, targetValue)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	var registrationContinuationToken string = resolveRegistrationToken(authenticationClient, requestContext, challengeResponse)
	if registrationContinuationToken == "" {
		return
	}
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemContinuationToken(requestContext, registrationContinuationToken)
	if redeemError != nil {
		printError(redeemError)
		return
	}
	printTokenResponse(tokenResponse, showFullTokens)
}

func resolveRegistrationToken(authenticationClient entrance.Client, requestContext context.Context, challengeResponse entrance.RegistrationChallengeResponse) string {
	var challengeTypeString string = string(challengeResponse.ChallengeType)
	if challengeTypeString == challengeTypePreverified {
		result := challengeResponse.ContinuationToken
		return result
	}
	printChallengeDetails(challengeResponse.ChallengeTargetLabel, challengeResponse.CodeLength)
	var passcode string = readLine(registrationPasscodePrompt)
	var continueResponse entrance.RegistrationContinueResponse
	var continueError error
	continueResponse, continueError = authenticationClient.ContinueRegistration(requestContext, challengeResponse.ContinuationToken, passcode)
	if continueError != nil {
		printError(continueError)
		result := ""
		return result
	}
	result := continueResponse.ContinuationToken
	return result
}

func signInWithOneTimePasscode(authenticationClient entrance.Client, emailAddress string, showFullTokens bool) {
	var requestContext context.Context = context.Background()
	var challengeMethods []entrance.ChallengeMethod = []entrance.ChallengeMethod{
		entrance.ChallengeMethodOneTimePasscode,
		entrance.ChallengeMethodRedirect,
	}
	var initiateResponse entrance.InitiateResponse
	var initiateError error
	initiateResponse, initiateError = authenticationClient.Initiate(requestContext, emailAddress, challengeMethods, nil)
	if initiateError != nil {
		printError(initiateError)
		return
	}
	var challengeResponse entrance.ChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.Challenge(requestContext, initiateResponse.ContinuationToken, challengeMethods)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	printChallengeDetails(challengeResponse.ChallengeTargetLabel, challengeResponse.CodeLength)
	for attempt := 0; attempt < maximumPasscodeAttempts; attempt++ {
		var finished bool = attemptPasscodeRedemption(authenticationClient, requestContext, challengeResponse.ContinuationToken, showFullTokens)
		if finished {
			return
		}
	}
	fmt.Println("Maximum passcode attempts exceeded.")
}

func signInWithPassword(authenticationClient entrance.Client, emailAddress string, showFullTokens bool) {
	var password string = readPassword(passwordPrompt)
	var requestContext context.Context = context.Background()
	var challengeMethods []entrance.ChallengeMethod = []entrance.ChallengeMethod{
		entrance.ChallengeMethodPassword,
		entrance.ChallengeMethodRedirect,
	}
	var initiateResponse entrance.InitiateResponse
	var initiateError error
	initiateResponse, initiateError = authenticationClient.Initiate(requestContext, emailAddress, challengeMethods, nil)
	if initiateError != nil {
		printError(initiateError)
		return
	}
	var challengeResponse entrance.ChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.Challenge(requestContext, initiateResponse.ContinuationToken, challengeMethods)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemPassword(requestContext, challengeResponse.ContinuationToken, password)
	if redeemError != nil {
		printError(redeemError)
		return
	}
	printTokenResponse(tokenResponse, showFullTokens)
}

func signInWithPasswordAndMultifactor(authenticationClient entrance.Client, emailAddress string, showFullTokens bool) {
	var password string = readPassword(passwordPrompt)
	var requestContext context.Context = context.Background()
	var challengeMethods []entrance.ChallengeMethod = []entrance.ChallengeMethod{
		entrance.ChallengeMethodPassword,
		entrance.ChallengeMethodRedirect,
	}
	var capabilities []entrance.Capability = []entrance.Capability{
		entrance.CapabilityMultifactorRequired,
	}
	var initiateResponse entrance.InitiateResponse
	var initiateError error
	initiateResponse, initiateError = authenticationClient.Initiate(requestContext, emailAddress, challengeMethods, capabilities)
	if initiateError != nil {
		printError(initiateError)
		return
	}
	var challengeResponse entrance.ChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.Challenge(requestContext, initiateResponse.ContinuationToken, challengeMethods)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemPassword(requestContext, challengeResponse.ContinuationToken, password)
	if redeemError != nil {
		var multifactorError entrance.MultifactorRequiredError
		if errors.As(redeemError, &multifactorError) {
			handleMultifactorChallenge(authenticationClient, requestContext, multifactorError.ContinuationToken, showFullTokens)
			return
		}
		printError(redeemError)
		return
	}
	printTokenResponse(tokenResponse, showFullTokens)
}

func signInWithPasswordAndRegistration(authenticationClient entrance.Client, emailAddress string, showFullTokens bool) {
	var password string = readPassword(passwordPrompt)
	var requestContext context.Context = context.Background()
	var challengeMethods []entrance.ChallengeMethod = []entrance.ChallengeMethod{
		entrance.ChallengeMethodPassword,
		entrance.ChallengeMethodRedirect,
	}
	var capabilities []entrance.Capability = []entrance.Capability{
		entrance.CapabilityRegistrationRequired,
	}
	var initiateResponse entrance.InitiateResponse
	var initiateError error
	initiateResponse, initiateError = authenticationClient.Initiate(requestContext, emailAddress, challengeMethods, capabilities)
	if initiateError != nil {
		printError(initiateError)
		return
	}
	var challengeResponse entrance.ChallengeResponse
	var challengeError error
	challengeResponse, challengeError = authenticationClient.Challenge(requestContext, initiateResponse.ContinuationToken, challengeMethods)
	if challengeError != nil {
		printError(challengeError)
		return
	}
	var tokenResponse entrance.TokenResponse
	var redeemError error
	tokenResponse, redeemError = authenticationClient.RedeemPassword(requestContext, challengeResponse.ContinuationToken, password)
	if redeemError != nil {
		var registrationError entrance.RegistrationRequiredError
		if errors.As(redeemError, &registrationError) {
			handleRegistrationChallenge(authenticationClient, requestContext, registrationError.ContinuationToken, showFullTokens)
			return
		}
		printError(redeemError)
		return
	}
	printTokenResponse(tokenResponse, showFullTokens)
}
