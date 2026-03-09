package entrance

import (
	_ "embed"

	"golang.org/x/text/language"

	"github.com/Grit-Software-Systems/entrance/internal/global"
)

//go:embed locales/en.json
var englishTranslationsContent []byte

func init() {
	global.LoadTranslations(language.English, englishTranslationsContent)
}
