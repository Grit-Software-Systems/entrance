package request

import (
	"testing"
)

func TestBuildEndpointUrlWithDefaultOnmicrosoftForm(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathInitiate)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/oauth2/v2.0/initiate"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithGuidTenantIdentifier(test *testing.T) {
	tenantIdentifier := "00000000-0000-0000-0000-000000000000"
	result := BuildEndpointUrl("contoso", tenantIdentifier, PathToken)

	expectedUrl := "https://contoso.ciamlogin.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithChallengePath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathChallenge)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/oauth2/v2.0/challenge"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithIntrospectPath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathIntrospect)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/oauth2/v2.0/introspect"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithRegistrationChallengePath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathRegistrationChallenge)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/register/v1.0/challenge"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithRegistrationContinuePath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathRegistrationContinue)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/register/v1.0/continue"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithRegistrationIntrospectPath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathRegistrationIntrospect)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/register/v1.0/introspect"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}

func TestBuildEndpointUrlWithTokenPath(test *testing.T) {
	result := BuildEndpointUrl("contoso", "", PathToken)

	expectedUrl := "https://contoso.ciamlogin.com/contoso.onmicrosoft.com/oauth2/v2.0/token"
	if result != expectedUrl {
		test.Errorf("expected %s, got %s", expectedUrl, result)
	}
}
