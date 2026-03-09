package request

import (
	"net/url"
)

const (
	hostSuffix   = ".ciamlogin.com"
	tenantSuffix = ".onmicrosoft.com"
	urlScheme    = "https"
)

const (
	PathChallenge              = "/oauth2/v2.0/challenge"
	PathInitiate               = "/oauth2/v2.0/initiate"
	PathIntrospect             = "/oauth2/v2.0/introspect"
	PathRegistrationChallenge  = "/register/v1.0/challenge"
	PathRegistrationContinue   = "/register/v1.0/continue"
	PathRegistrationIntrospect = "/register/v1.0/introspect"
	PathToken                  = "/oauth2/v2.0/token"
)

func BuildEndpointUrl(tenantSubdomain string, tenantIdentifier string, path string) string {
	var endpointUrl url.URL
	endpointUrl.Scheme = urlScheme
	endpointUrl.Host = tenantSubdomain + hostSuffix
	endpointUrl.Path = buildTenantPath(tenantSubdomain, tenantIdentifier, path)
	result := endpointUrl.String()
	return result
}

func buildTenantPath(tenantSubdomain string, tenantIdentifier string, path string) string {
	var tenantSegment string
	if tenantIdentifier != "" {
		tenantSegment = tenantIdentifier
	} else {
		tenantSegment = tenantSubdomain + tenantSuffix
	}
	result := "/" + tenantSegment + path
	return result
}
