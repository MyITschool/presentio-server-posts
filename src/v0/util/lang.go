package util

import (
	"presentio-server-posts/src/v0/deps/detector/langdet/langdetdef"
)

func GetLang(text string) string {
	detector := langdetdef.NewWithDefaultLanguages()

	lang := detector.GetClosestLanguage(text)

	if lang == "hebrew" || lang == "undefined" {
		lang = "english" // because there is no default configuration for `hebrew` in postgres
	}

	return lang
}
