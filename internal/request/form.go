package request

import (
	"net/url"
	"strings"
)

func BuildFormBody(parameters map[string]string) string {
	var encodedPairs []string
	for parameterName, parameterValue := range parameters {
		encodedName := url.QueryEscape(parameterName)
		encodedValue := url.QueryEscape(parameterValue)
		encodedPair := encodedName + "=" + encodedValue
		encodedPairs = append(encodedPairs, encodedPair)
	}
	result := strings.Join(encodedPairs, "&")
	return result
}
