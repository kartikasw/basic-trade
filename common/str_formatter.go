package common

import (
	"strings"
)

func FormatStrForFullTextSearch(keyword string) string {
	words := strings.Fields(keyword)
	for i, word := range words {
		words[i] = word + ":*"
	}
	formattedString := strings.Join(words, " & ")

	if formattedString == "" {
		return ":*"
	}

	return formattedString
}
