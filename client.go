package entrance

import (
	"net/http"
)

type Client struct {
	configuration Configuration
	httpClient    *http.Client
}

func NewClient(configuration Configuration) Client {
	result := Client{
		configuration: configuration,
		httpClient:    http.DefaultClient,
	}
	return result
}

func NewClientWithHttpClient(configuration Configuration, httpClient *http.Client) Client {
	result := Client{
		configuration: configuration,
		httpClient:    httpClient,
	}
	return result
}
