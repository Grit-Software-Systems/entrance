package request

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	contentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
	contentTypeHeader         = "Content-Type"
)

type ErrorPayload struct {
	ErrorCode         string `json:"error"`
	SubErrorCode      string `json:"suberror"`
	ErrorDescription  string `json:"error_description"`
	CorrelationId     string `json:"correlation_id"`
	Timestamp         string `json:"timestamp"`
	ContinuationToken string `json:"continuation_token"`
}

func (errorPayload *ErrorPayload) Error() string {
	result := fmt.Sprintf(
		"authentication error: %s (suberror: %s, description: %s)",
		errorPayload.ErrorCode,
		errorPayload.SubErrorCode,
		errorPayload.ErrorDescription,
	)
	return result
}

func SendRequest(
	requestContext context.Context,
	httpClient *http.Client,
	endpointUrl string,
	formBody string,
	successTarget interface{},
) error {
	bodyReader := strings.NewReader(formBody)
	httpRequest, requestCreationError := http.NewRequestWithContext(
		requestContext, http.MethodPost, endpointUrl, bodyReader,
	)
	if requestCreationError != nil {
		return fmt.Errorf("%w", requestCreationError)
	}
	httpRequest.Header.Set(contentTypeHeader, contentTypeFormUrlEncoded)
	httpResponse, requestError := httpClient.Do(httpRequest)
	if requestError != nil {
		return fmt.Errorf("%w", requestError)
	}
	defer httpResponse.Body.Close()
	responseBody, readError := io.ReadAll(httpResponse.Body)
	if readError != nil {
		return fmt.Errorf("%w", readError)
	}
	if httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300 {
		result := json.Unmarshal(responseBody, successTarget)
		return result
	}
	var errorPayload ErrorPayload
	unmarshalError := json.Unmarshal(responseBody, &errorPayload)
	if unmarshalError != nil {
		return fmt.Errorf("failed to parse error response: %w", unmarshalError)
	}
	result := &errorPayload
	return result
}
