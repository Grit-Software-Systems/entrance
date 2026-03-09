package global

import (
	"encoding/json"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func LoadTranslations(languageTag language.Tag, translationsContent []byte) {
	var translations map[string]string
	unmarshalError := json.Unmarshal(translationsContent, &translations)
	if unmarshalError != nil {
		panic(unmarshalError)
	}
	for messageKey, translatedText := range translations {
		message.SetString(languageTag, messageKey, translatedText)
	}
}

func TranslateMessage(languageTag language.Tag, messageKey string) string {
	printer := message.NewPrinter(languageTag)
	result := printer.Sprintf(messageKey)
	return result
}
