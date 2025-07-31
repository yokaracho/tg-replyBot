package bot

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tg-replyBot/internal/config"
	"tg-replyBot/internal/services"
	"tg-replyBot/pkg/logger"
)

type Bot struct {
	api             *tgbotapi.BotAPI
	config          config.TelegramConfig
	contextManager  *services.ContextManager
	replyGenerator  *services.ReplyGenerator
	logger          logger.Logger
	handlers        *Handlers
	callbackHandler *CallbackHandler
}

func New(
	config config.TelegramConfig,
	contextManager *services.ContextManager,
	replyGenerator *services.ReplyGenerator,
	logger logger.Logger,
) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Telegram API: %w", err)
	}

	api.Debug = config.Debug
	logger.Info("Авторизован как", "username", api.Self.UserName)

	bot := &Bot{
		api:            api,
		config:         config,
		contextManager: contextManager,
		replyGenerator: replyGenerator,
		logger:         logger,
	}

	// Инициализируем обработчики
	bot.handlers = NewHandlers(bot)
	bot.callbackHandler = NewCallbackHandler(bot)

	return bot, nil
}

func (b *Bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.config.Timeout

	updates := b.api.GetUpdatesChan(u)

	b.logger.Info("Бот начал получение обновлений")

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("Получен сигнал остановки")
			return ctx.Err()
		case update := <-updates:
			go b.handleUpdate(update)
		}
	}
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
	b.logger.Info("Получение обновлений остановлено")
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("Panic в обработке обновления", "panic", r)
		}
	}()

	if update.Message != nil {
		b.handlers.HandleMessage(update.Message)
	} else if update.CallbackQuery != nil {
		b.callbackHandler.HandleCallback(update.CallbackQuery)
	}
}

func (b *Bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) SendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) EditMessage(chatID int64, messageID int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, text)
	editMsg.ParseMode = "Markdown"
	if keyboard != nil {
		editMsg.ReplyMarkup = keyboard
	}
	_, err := b.api.Send(editMsg)
	return err
}

func (b *Bot) SendChatAction(chatID int64, action string) error {
	chatAction := tgbotapi.NewChatAction(chatID, action)
	_, err := b.api.Send(chatAction)
	return err
}
