package entrance

import (
	"context"

	"github.com/Grit-Software-Systems/entrance/internal/request"
)

const (
	formParameterChallengeChannel = "challenge_channel"
	formParameterChallengeTarget  = "challenge_target"

	grantTypeContinuationToken = "continuation_token"
)

func (client Client) ChallengeRegistration(
	requestContext context.Context,
	continuationToken string,
	challengeChannel string,
	challengeTarget string,
) (RegistrationChallengeResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterChallengeType:     string(ChallengeMethodOneTimePasscode),
		formParameterChallengeChannel:  challengeChannel,
		formParameterChallengeTarget:   challengeTarget,
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathRegistrationChallenge,
	)
	var registrationChallengeResponse RegistrationChallengeResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &registrationChallengeResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return RegistrationChallengeResponse{}, classifiedError
	}
	return registrationChallengeResponse, nil
}

func (client Client) ContinueRegistration(
	requestContext context.Context,
	continuationToken string,
	passcode string,
) (RegistrationContinueResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterOneTimePasscode:   passcode,
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathRegistrationContinue,
	)
	var registrationContinueResponse RegistrationContinueResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &registrationContinueResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return RegistrationContinueResponse{}, classifiedError
	}
	return registrationContinueResponse, nil
}

func (client Client) IntrospectRegistration(
	requestContext context.Context,
	continuationToken string,
) (RegistrationIntrospectResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
	}
	formBody := request.BuildFormBody(formParameters)
	endpointUrl := request.BuildEndpointUrl(
		client.configuration.TenantSubdomain,
		client.configuration.TenantIdentifier,
		request.PathRegistrationIntrospect,
	)
	var registrationIntrospectResponse RegistrationIntrospectResponse
	sendError := request.SendRequest(
		requestContext, client.httpClient, endpointUrl, formBody, &registrationIntrospectResponse,
	)
	if sendError != nil {
		classifiedError := classifyError(sendError)
		return RegistrationIntrospectResponse{}, classifiedError
	}
	return registrationIntrospectResponse, nil
}

func (client Client) RedeemContinuationToken(
	requestContext context.Context,
	continuationToken string,
) (TokenResponse, error) {
	formParameters := map[string]string{
		formParameterClientIdentifier:  client.configuration.ClientIdentifier,
		formParameterContinuationToken: continuationToken,
		formParameterGrantType:         grantTypeContinuationToken,
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
