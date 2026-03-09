package request

import (
	"strings"
	"testing"
)

func TestBuildFormBodyWithEmptyMap(test *testing.T) {
	parameters := map[string]string{}

	result := BuildFormBody(parameters)

	if result != "" {
		test.Errorf("expected empty string, got %s", result)
	}
}

func TestBuildFormBodyWithSingleParameter(test *testing.T) {
	parameters := map[string]string{
		"client_id": "test-client",
	}

	result := BuildFormBody(parameters)

	expectedBody := "client_id=test-client"
	if result != expectedBody {
		test.Errorf("expected %s, got %s", expectedBody, result)
	}
}

func TestBuildFormBodyWithMultipleParameters(test *testing.T) {
	parameters := map[string]string{
		"client_id": "test-client",
		"username":  "user@contoso.com",
	}

	result := BuildFormBody(parameters)

	if !strings.Contains(result, "client_id=test-client") {
		test.Errorf("expected result to contain client_id=test-client, got %s", result)
	}
	if !strings.Contains(result, "username=user%40contoso.com") {
		test.Errorf("expected result to contain username=user%%40contoso.com, got %s", result)
	}
	if strings.Count(result, "&") != 1 {
		test.Errorf("expected exactly one ampersand separator, got %s", result)
	}
}

func TestBuildFormBodyWithSpecialCharacters(test *testing.T) {
	parameters := map[string]string{
		"password": "p@ss w0rd&more=yes",
	}

	result := BuildFormBody(parameters)

	expectedBody := "password=p%40ss+w0rd%26more%3Dyes"
	if result != expectedBody {
		test.Errorf("expected %s, got %s", expectedBody, result)
	}
}
