package entrance

import (
	"testing"
)

func TestEffectiveScopes_ReturnsDefaultWhenScopesIsEmpty(testing *testing.T) {
	configuration := Configuration{
		TenantSubdomain:  "contoso",
		TenantIdentifier: "contoso.onmicrosoft.com",
		ClientIdentifier: "00000000-0000-0000-0000-000000000000",
		Scopes:           "",
	}
	result := configuration.effectiveScopes()
	if result != defaultScopes {
		testing.Errorf("expected %q, got %q", defaultScopes, result)
	}
}

func TestEffectiveScopes_ReturnsCustomWhenScopesIsSet(testing *testing.T) {
	customScopes := "openid email"
	configuration := Configuration{
		TenantSubdomain:  "contoso",
		TenantIdentifier: "contoso.onmicrosoft.com",
		ClientIdentifier: "00000000-0000-0000-0000-000000000000",
		Scopes:           customScopes,
	}
	result := configuration.effectiveScopes()
	if result != customScopes {
		testing.Errorf("expected %q, got %q", customScopes, result)
	}
}
