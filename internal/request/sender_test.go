package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendRequestWithSuccessfulResponse(test *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"continuation_token":"test-token-123"}`))
	}))
	defer server.Close()

	type successResponse struct {
		ContinuationToken string `json:"continuation_token"`
	}
	var target successResponse

	sendError := SendRequest(context.Background(), server.Client(), server.URL, "", &target)

	if sendError != nil {
		test.Fatalf("expected no error, got %v", sendError)
	}
	if target.ContinuationToken != "test-token-123" {
		test.Errorf("expected continuation_token test-token-123, got %s", target.ContinuationToken)
	}
}

func TestSendRequestWithErrorResponse(test *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte(`{"error":"user_not_found","error_description":"User not found"}`))
	}))
	defer server.Close()

	var target map[string]interface{}

	sendError := SendRequest(context.Background(), server.Client(), server.URL, "", &target)

	if sendError == nil {
		test.Fatal("expected an error, got nil")
	}
	errorMessage := sendError.Error()
	if !strings.Contains(errorMessage, "user_not_found") {
		test.Errorf("expected error to contain user_not_found, got %s", errorMessage)
	}
}

func TestSendRequestWithNetworkFailure(test *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {}))
	server.Close()

	var target map[string]interface{}

	sendError := SendRequest(context.Background(), server.Client(), server.URL, "", &target)

	if sendError == nil {
		test.Fatal("expected an error for network failure, got nil")
	}
}

func TestSendRequestWithNonJsonBody(test *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "text/plain")
		responseWriter.WriteHeader(http.StatusBadRequest)
		responseWriter.Write([]byte("not json"))
	}))
	defer server.Close()

	var target map[string]interface{}

	sendError := SendRequest(context.Background(), server.Client(), server.URL, "", &target)

	if sendError == nil {
		test.Fatal("expected an error for non-JSON body, got nil")
	}
	errorMessage := sendError.Error()
	if !strings.Contains(errorMessage, "parse") {
		test.Errorf("expected error about parsing, got %s", errorMessage)
	}
}

func TestSendRequestSetsContentTypeHeader(test *testing.T) {
	var capturedContentType string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		capturedContentType = request.Header.Get("Content-Type")
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{}`))
	}))
	defer server.Close()

	var target map[string]interface{}

	SendRequest(context.Background(), server.Client(), server.URL, "", &target)

	if capturedContentType != "application/x-www-form-urlencoded" {
		test.Errorf("expected Content-Type application/x-www-form-urlencoded, got %s", capturedContentType)
	}
}

func TestSendRequestSendsFormBody(test *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		bodyBytes := make([]byte, request.ContentLength)
		request.Body.Read(bodyBytes)
		capturedBody = string(bodyBytes)
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{}`))
	}))
	defer server.Close()

	var target map[string]interface{}
	formBody := "client_id=test-client&username=user%40contoso.com"

	SendRequest(context.Background(), server.Client(), server.URL, formBody, &target)

	if capturedBody != formBody {
		test.Errorf("expected body %s, got %s", formBody, capturedBody)
	}
}
