package entrance

import (
	"context"
	"strings"

	"github.com/Grit-Software-Systems/entrance/internal/request"
)

const (
	formParameterCapabilities      = "capabilities"
	formParameterChallengeType     = "challenge_type"
	formParameterClientIdentifier  = "client_id"
	formParameterContinuationToken = "continuation_token"
	formParameterGrantType         = "grant_type"
	formParameterOneTimePasscode   = "oob"
	formParameterPassword          = "password"
	formParameterScope             = "scope"
	formParameterUsername          = "username"

	grantTypeOneTimePasscode = "oob"
	grantTypePassword        = "password"
)

func buildCapabilitiesValue(capabilities []Capability) string {
	var values []string
	for _, item := range capabilities {
		values = append(values, string(item))
	}
	result := strings.Join(values, " ")
	return result
}

func buildChallengeTypeValue(challengeTypes []ChallengeMethod) string {
	var values []string
	for _, item := range challengeTypes {
		values = append(values, string(item))
	}
	result := strings.Join(values, " ")
	return result
}

func (client Client) Challenge(
	requestContext context.Context,
	continuationToken string,
	challengeTypes []ChallengeMethod,
) (ChallengeResponse, error) {
	challengeTypeValue := buildChallengeTypeValue(challengeTypes)
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterChallengeType:     challengeTypeValue,
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathChallenge,
	)
	var challengeResponse ChallengeResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &challengeResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return ChallengeResponse{}, classifiedError
	}
	if challengeResponse.ChallengeType == ChallengeMethodRedirect {
		result := RedirectRequiredError{}
		return ChallengeResponse{}, result
	}
	return challengeResponse, nil
}

func (client Client) Initiate(
	requestContext context.Context,
	username string,
	challengeTypes []ChallengeMethod,
	capabilities []Capability,
) (InitiateResponse, error) {
	challengeTypeValue := buildChallengeTypeValue(challengeTypes)
	formParameters := map[string]string{
		formParameterClientIdentifier: client.configuration.ClientIdentifier,
		formParameterUsername:         username,
		formParameterChallengeType:    challengeTypeValue,
	}
	if len(capabilities) > 0 {
		capabilitiesValue := buildCapabilitiesValue(capabilities)
		formParameters[formParameterCapabilities] = capabilitiesValue
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathInitiate,
	)
	var initiateResponse InitiateResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &initiateResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return InitiateResponse{}, classifiedError
	}
	return initiateResponse, nil
}

func (client Client) RedeemOneTimePasscode(
	requestContext context.Context,
	continuationToken string,
	passcode string,
) (TokenResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterGrantType:         grantTypeOneTimePasscode,
		formParameterOneTimePasscode:   passcode,
		formParameterScope:             client.configuration.effectiveScopes(),
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathToken,
	)
	var tokenResponse TokenResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &tokenResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return TokenResponse{}, classifiedError
	}
	return tokenResponse, nil
}

func (client Client) RedeemPassword(
	requestContext context.Context,
	continuationToken string,
	password string,
) (TokenResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterGrantType:         grantTypePassword,
		formParameterPassword:          password,
		formParameterScope:             client.configuration.effectiveScopes(),
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathToken,
	)
	var tokenResponse TokenResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &tokenResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return TokenResponse{}, classifiedError
	}
	return tokenResponse, nil
}
