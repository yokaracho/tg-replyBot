package bot

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-replyBot/internal/models"
)

type CallbackHandler struct {
	bot *Bot
}

func NewCallbackHandler(bot *Bot) *CallbackHandler {
	return &CallbackHandler{bot: bot}
}

func (ch *CallbackHandler) HandleCallback(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	ch.bot.logger.Info("Получен callback",
		"user_id", userID,
		"data", callback.Data)

	// Отвечаем на callback
	callbackResponse := tgbotapi.NewCallback(callback.ID, "")
	ch.bot.api.Request(callbackResponse)

	ctx := context.Background()

	userContext, err := ch.bot.contextManager.GetContext(ctx, userID)
	if err != nil {
		ch.bot.logger.Error("Ошибка получения контекста", "error", err)
		ch.bot.EditMessage(chatID, messageID, "❌ Данные не найдены. Отправьте сообщение заново.", nil)
		return
	}

	switch callback.Data {
	case "all_styles":
		ch.showAllStyles(chatID, messageID, userContext)
	case "back_to_main":
		ch.showMainStyles(chatID, messageID, userContext)
	default:
		if strings.HasPrefix(callback.Data, "style_") {
			styleName := strings.TrimPrefix(callback.Data, "style_")
			ch.showSpecificStyle(chatID, messageID, userContext, styleName)
		}
	}
}

func (ch *CallbackHandler) showAllStyles(chatID int64, messageID int, userContext *models.Context) {
	if userContext.LastMessage == "" || len(userContext.Replies) == 0 {
		ch.bot.EditMessage(chatID, messageID, "❌ Данные не найдены. Отправьте сообщение заново.", nil)
		return
	}

	// Формируем все ответы
	var allReplies []models.Reply
	for _, style := range models.DefaultStyles {
		if reply, exists := userContext.Replies[style.Name]; exists {
			allReplies = append(allReplies, models.Reply{
				Style:   style.Name,
				Content: reply,
				Emoji:   style.Emoji,
			})
		}
	}

	responseText := ch.formatAllStylesResponse(userContext.LastMessage, allReplies)

	// Кнопка "Назад"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_main"),
		),
	)

	ch.bot.EditMessage(chatID, messageID, responseText, &keyboard)
}

func (ch *CallbackHandler) showMainStyles(chatID int64, messageID int, userContext *models.Context) {
	if userContext.LastMessage == "" || len(userContext.Replies) == 0 {
		ch.bot.EditMessage(chatID, messageID, "❌ Данные не найдены. Отправьте сообщение заново.", nil)
		return
	}

	// Формируем первые 3 ответа
	var mainReplies []models.Reply
	count := 0
	for _, style := range models.DefaultStyles {
		if count >= 3 {
			break
		}
		if reply, exists := userContext.Replies[style.Name]; exists {
			mainReplies = append(mainReplies, models.Reply{
				Style:   style.Name,
				Content: reply,
				Emoji:   style.Emoji,
			})
			count++
		}
	}

	responseText := ch.formatMainStylesResponse(userContext.LastMessage, mainReplies)
	keyboard := ch.createMainKeyboard(mainReplies)

	ch.bot.EditMessage(chatID, messageID, responseText, &keyboard)
}

func (ch *CallbackHandler) showSpecificStyle(chatID int64, messageID int, userContext *models.Context, styleName string) {
	reply, exists := userContext.Replies[styleName]
	if !exists {
		ch.bot.EditMessage(chatID, messageID, "❌ Стиль не найден.", nil)
		return
	}

	// Находим информацию о стиле
	var style models.Style
	found := false
	for _, s := range models.DefaultStyles {
		if s.Name == styleName {
			style = s
			found = true
			break
		}
	}

	if !found {
		ch.bot.EditMessage(chatID, messageID, "❌ Стиль не найден.", nil)
		return
	}

	displayMessage := userContext.LastMessage
	if len(displayMessage) > 100 {
		displayMessage = displayMessage[:100] + "..."
	}

	responseText := fmt.Sprintf("💬 **Сообщение:** %s\n\n%s **%s:**\n%s",
		displayMessage,
		style.Emoji,
		style.DisplayName,
		reply)

	// Кнопки навигации
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔙 Назад", "back_to_main"),
			tgbotapi.NewInlineKeyboardButtonData("📝 Все стили", "all_styles"),
		),
	)

	ch.bot.EditMessage(chatID, messageID, responseText, &keyboard)
}

func (ch *CallbackHandler) formatAllStylesResponse(originalMessage string, replies []models.Reply) string {
	displayMessage := originalMessage
	if len(displayMessage) > 100 {
		displayMessage = displayMessage[:100] + "..."
	}

	responseText := fmt.Sprintf("💬 **Сообщение:** %s\n\n🎯 **Все варианты ответов:**\n\n", displayMessage)

	for _, reply := range replies {
		styleDisplayName := models.DefaultStyles[getStyleIndex(reply.Style)].DisplayName
		responseText += fmt.Sprintf("%s **%s:**\n%s\n\n",
			reply.Emoji,
			styleDisplayName,
			reply.Content)
	}

	return responseText
}

func (ch *CallbackHandler) formatMainStylesResponse(originalMessage string, replies []models.Reply) string {
	displayMessage := originalMessage
	if len(displayMessage) > 100 {
		displayMessage = displayMessage[:100] + "..."
	}

	responseText := fmt.Sprintf("💬 **Сообщение:** %s\n\n🎯 **Варианты ответов:**\n\n", displayMessage)

	for _, reply := range replies {
		styleDisplayName := models.DefaultStyles[getStyleIndex(reply.Style)].DisplayName
		responseText += fmt.Sprintf("%s **%s:**\n%s\n\n",
			reply.Emoji,
			styleDisplayName,
			reply.Content)
	}

	return responseText
}

func (ch *CallbackHandler) createMainKeyboard(replies []models.Reply) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton

	for _, reply := range replies {
		styleDisplayName := models.DefaultStyles[getStyleIndex(reply.Style)].DisplayName
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", reply.Emoji, styleDisplayName),
			fmt.Sprintf("style_%s", reply.Style),
		)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}

	allStylesButton := tgbotapi.NewInlineKeyboardButtonData("📝 Все стили", "all_styles")
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{allStylesButton})

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}
