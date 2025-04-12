package tools

import (
	"strings"
	"unicode"
)

func UcWords(s string) string {
	return strings.Join(capitalizeWords(strings.Fields(s)), " ")
}

func capitalizeWords(words []string) []string {
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return words
}
