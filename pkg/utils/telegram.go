package utils

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func EscapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"*", "\\*",
		"_", "\\_",
		"`", "\\`",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
	)
	return replacer.Replace(text)
}

func FormatUserMention(user *tgbotapi.User) string {
	if user.UserName != "" {
		return fmt.Sprintf("@%s", user.UserName)
	}

	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}

	return fmt.Sprintf("[%s](tg://user?id=%d)", EscapeMarkdown(name), user.ID)
}

func CreateInlineKeyboard(buttons [][]InlineButton) tgbotapi.InlineKeyboardMarkup {
	var keyboard [][]tgbotapi.InlineKeyboardButton

	for _, row := range buttons {
		var keyboardRow []tgbotapi.InlineKeyboardButton
		for _, button := range row {
			keyboardRow = append(keyboardRow, tgbotapi.NewInlineKeyboardButtonData(button.Text, button.Data))
		}
		keyboard = append(keyboard, keyboardRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

type InlineButton struct {
	Text string
	Data string
}

func SplitLongMessage(text string, maxLength int) []string {
	if len(text) <= maxLength {
		return []string{text}
	}

	var messages []string
	words := strings.Fields(text)
	currentMessage := ""

	for _, word := range words {
		testMessage := currentMessage
		if testMessage != "" {
			testMessage += " "
		}
		testMessage += word

		if len(testMessage) > maxLength {
			if currentMessage != "" {
				messages = append(messages, currentMessage)
				currentMessage = word
			} else {
				// Слово слишком длинное, принудительно разбиваем
				for len(word) > maxLength {
					messages = append(messages, word[:maxLength])
					word = word[maxLength:]
				}
				currentMessage = word
			}
		} else {
			currentMessage = testMessage
		}
	}

	if currentMessage != "" {
		messages = append(messages, currentMessage)
	}

	return messages
}
