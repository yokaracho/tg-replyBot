package utils

import (
	"strings"
	"unicode/utf8"
)

func TruncateText(text string, maxLength int) string {
	if utf8.RuneCountInString(text) <= maxLength {
		return text
	}

	runes := []rune(text)
	if len(runes) <= maxLength-3 {
		return text
	}

	return string(runes[:maxLength-3]) + "..."
}

func CleanText(text string) string {
	text = strings.TrimSpace(text)

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return text
}

func IsQuestion(text string) bool {
	return strings.Contains(text, "?") ||
		strings.HasSuffix(text, "?")
}

func ContainsAny(text string, words []string) bool {
	textLower := strings.ToLower(text)
	for _, word := range words {
		if strings.Contains(textLower, strings.ToLower(word)) {
			return true
		}
	}
	return false
}
