package entrance

import (
	"net/http"
	"testing"
)

func TestNewClient_CreatesClientWithNonNilHttpClient(testing *testing.T) {
	configuration := Configuration{
		TenantSubdomain:  "contoso",
		TenantIdentifier: "contoso.onmicrosoft.com",
		ClientIdentifier: "00000000-0000-0000-0000-000000000000",
	}
	result := NewClient(configuration)
	if result.httpClient == nil {
		testing.Error("expected httpClient to be non-nil")
	}
}

func TestNewClient_StoresConfiguration(testing *testing.T) {
	configuration := Configuration{
		TenantSubdomain:  "contoso",
		TenantIdentifier: "contoso.onmicrosoft.com",
		ClientIdentifier: "00000000-0000-0000-0000-000000000000",
	}
	result := NewClient(configuration)
	if result.configuration.TenantSubdomain != configuration.TenantSubdomain {
		testing.Errorf("expected TenantSubdomain %q, got %q",
			configuration.TenantSubdomain, result.configuration.TenantSubdomain)
	}
	if result.configuration.ClientIdentifier != configuration.ClientIdentifier {
		testing.Errorf("expected ClientIdentifier %q, got %q",
			configuration.ClientIdentifier, result.configuration.ClientIdentifier)
	}
}

func TestNewClientWithHttpClient_UsesProvidedClient(testing *testing.T) {
	configuration := Configuration{
		TenantSubdomain:  "contoso",
		TenantIdentifier: "contoso.onmicrosoft.com",
		ClientIdentifier: "00000000-0000-0000-0000-000000000000",
	}
	customHttpClient := &http.Client{}
	result := NewClientWithHttpClient(configuration, customHttpClient)
	if result.httpClient != customHttpClient {
		testing.Error("expected httpClient to be the provided custom client")
	}
}
