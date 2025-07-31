package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-replyBot/internal/models"
)

type Handlers struct {
	bot *Bot
}

func NewHandlers(bot *Bot) *Handlers {
	return &Handlers{bot: bot}
}

func (h *Handlers) HandleMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID

	h.bot.logger.Info("Получено сообщение",
		"user_id", userID,
		"username", message.From.UserName,
		"text", message.Text)

	switch message.Command() {
	case "start":
		h.handleStart(chatID)
	case "help":
		h.handleHelp(chatID)
	case "styles":
		h.handleStyles(chatID)
	case "clear":
		h.handleClear(chatID, userID)
	default:
		if message.Text != "" {
			h.handleTextMessage(message)
		} else {
			h.bot.SendMessage(chatID, "❌ Пожалуйста, отправьте текстовое сообщение")
		}
	}
}

func (h *Handlers) handleStart(chatID int64) {
	welcomeText := `🤖 Привет! Я бот-помощник для ответов на сообщения.

📝 **Как пользоваться:**
1. Просто перешлите мне сообщение, на которое нужно ответить
2. Или напишите текст сообщения
3. Я предложу несколько вариантов ответов в разных стилях

⚙️ **Команды:**
/help - помощь
/clear - очистить контекст разговора
/styles - посмотреть доступные стили ответов

Попробуйте прислать мне любое сообщение!`

	h.bot.SendMessage(chatID, welcomeText)
}

func (h *Handlers) handleHelp(chatID int64) {
	helpText := `📚 **Подробная инструкция:**

1️⃣ **Отправьте сообщение** - любым способом:
   • Перешлите сообщение от другого человека
   • Скопируйте и вставьте текст
   • Напишите своими словами суть

2️⃣ **Получите варианты ответов** в разных стилях:
   • Дружелюбный - как с близким другом
   • Официальный - деловой стиль
   • Краткий - по существу
   • Развернутый - подробное объяснение
   • С юмором - легкий и непринужденный
   • Эмпатичный - с пониманием и поддержкой
   • Мотивирующий - воодушевляющий

3️⃣ **Выберите подходящий** и отправьте собеседнику!

💡 **Советы:**
- Бот запоминает контекст разговора для лучших предложений
- Используйте /clear для сброса контекста
- Чем подробнее опишете ситуацию, тем лучше будут предложения`

	h.bot.SendMessage(chatID, helpText)
}

func (h *Handlers) handleStyles(chatID int64) {
	stylesText := "🎨 **Доступные стили ответов:**\n\n"

	for _, style := range models.DefaultStyles {
		stylesText += fmt.Sprintf("%s **%s**\n%s\n\n",
			style.Emoji,
			style.DisplayName,
			style.Description)
	}

	h.bot.SendMessage(chatID, stylesText)
}

func (h *Handlers) handleClear(chatID, userID int64) {
	ctx := context.Background()
	err := h.bot.contextManager.ClearContext(ctx, userID)
	if err != nil {
		h.bot.logger.Error("Ошибка очистки контекста", "user_id", userID, "error", err)
		h.bot.SendMessage(chatID, "❌ Произошла ошибка при очистке контекста")
		return
	}

	h.bot.SendMessage(chatID, "🗑️ Контекст разговора очищен!")
}

func (h *Handlers) handleTextMessage(message *tgbotapi.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID
	messageText := message.Text

	ctx := context.Background()

	// Добавляем сообщение в контекст
	err := h.bot.contextManager.AddMessage(ctx, userID, messageText)
	if err != nil {
		h.bot.logger.Error("Ошибка добавления сообщения в контекст", "error", err)
	}

	// Показываем индикатор "печатает"
	h.bot.SendChatAction(chatID, tgbotapi.ChatTyping)

	// Получаем контекст для генерации ответов
	userContext, err := h.bot.contextManager.GetContext(ctx, userID)
	if err != nil {
		h.bot.logger.Error("Ошибка получения контекста", "error", err)
		h.bot.SendMessage(chatID, "❌ Произошла ошибка")
		return
	}

	// Генерируем ответы
	request := models.ReplyRequest{
		Message:         messageText,
		ContextMessages: userContext.GetRecentMessages(5),
	}

	response, err := h.bot.replyGenerator.GenerateReplies(ctx, request)
	if err != nil {
		h.bot.logger.Error("Ошибка генерации ответов", "error", err)
		h.bot.SendMessage(chatID, "❌ Произошла ошибка при генерации ответов")
		return
	}

	if len(response.Replies) == 0 {
		h.bot.SendMessage(chatID, "❌ Не удалось сгенерировать ответы")
		return
	}

	// Сохраняем ответы в контекст
	replies := make(models.Replies)
	for _, reply := range response.Replies {
		replies[reply.Style] = reply.Content
	}
	h.bot.contextManager.SaveReplies(ctx, userID, replies)

	// Создаем клавиатуру и отправляем ответ
	keyboard := h.createMainKeyboard(response.Replies[:min(3, len(response.Replies))])
	responseText := h.formatResponse(messageText, response.Replies[:min(3, len(response.Replies))])

	h.bot.SendMessageWithKeyboard(chatID, responseText, keyboard)
}

func (h *Handlers) createMainKeyboard(replies []models.Reply) tgbotapi.InlineKeyboardMarkup {
	var buttons [][]tgbotapi.InlineKeyboardButton

	// Первые 3 стиля
	for _, reply := range replies {
		button := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", reply.Emoji, models.DefaultStyles[getStyleIndex(reply.Style)].DisplayName),
			fmt.Sprintf("style_%s", reply.Style),
		)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{button})
	}

	// Кнопка "Все стили"
	allStylesButton := tgbotapi.NewInlineKeyboardButtonData("📝 Все стили", "all_styles")
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{allStylesButton})

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func (h *Handlers) formatResponse(originalMessage string, replies []models.Reply) string {
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

func getStyleIndex(styleName string) int {
	for i, style := range models.DefaultStyles {
		if style.Name == styleName {
			return i
		}
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
