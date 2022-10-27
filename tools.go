package wordgenerator

import (
	"strings"

	stringTools "github.com/olehmushka/golang-toolkit/string_tools"
)

func hasInsideStrWordsLessThan(str string, min int) bool {
	for _, word := range strings.Split(str, " ") {
		if len(word) < min {
			return true
		}
	}
	return false
}

func makeInsideWordsCapitalized(str string) string {
	var out string
	for _, word := range strings.Split(str, " ") {
		out += stringTools.Capitalize(word)
	}
	return out
}
