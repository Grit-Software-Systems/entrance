package entrance

import (
	"context"

	"github.com/Grit-Software-Systems/entrance/internal/request"
)

const (
	grantTypeMultifactorOneTimePasscode = "mfa_oob"

	formParameterMethodIdentifier = "id"
)

func (client Client) ChallengeMultifactor(
	requestContext context.Context,
	continuationToken string,
	methodIdentifier string,
) (ChallengeResponse, error) {
	challengeTypeValue := buildChallengeTypeValue([]ChallengeMethod{
		ChallengeMethodOneTimePasscode,
		ChallengeMethodRedirect,
	})
	formParameters := map[string]string{
		formParameterChallengeType:     challengeTypeValue,
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterMethodIdentifier:  methodIdentifier,
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

func (client Client) Introspect(
	requestContext context.Context,
	continuationToken string,
) (IntrospectResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathIntrospect,
	)
	var introspectResponse IntrospectResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &introspectResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return IntrospectResponse{}, classifiedError
	}
	return introspectResponse, nil
}

func (client Client) RedeemMultifactorOneTimePasscode(
	requestContext context.Context,
	continuationToken string,
	passcode string,
) (TokenResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterGrantType:         grantTypeMultifactorOneTimePasscode,
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

