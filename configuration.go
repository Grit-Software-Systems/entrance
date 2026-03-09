package entrance

const (
	defaultScopes = "openid profile offline_access"
)

type Configuration struct {
	TenantSubdomain  string
	TenantIdentifier string
	ClientIdentifier string
	Scopes           string
}

func (configuration Configuration) effectiveScopes() string {
	if configuration.Scopes != "" {
		result := configuration.Scopes
		return result
	}
	result := defaultScopes
	return result
}
