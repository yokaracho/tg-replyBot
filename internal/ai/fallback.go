package ai

import (
	"context"
	"strings"

	"tg-replyBot/pkg/logger"
)

type Fallback struct {
	logger logger.Logger
}

func NewFallback(logger logger.Logger) Provider {
	return &Fallback{
		logger: logger,
	}
}

func (f *Fallback) GenerateReply(ctx context.Context, request Request) (string, error) {
	messageLower := strings.ToLower(request.Message)
	styleName := request.Style.Name

	if strings.Contains(messageLower, "спасибо") || strings.Contains(messageLower, "благодарю") {
		return f.getThankYouReply(styleName), nil
	} else if strings.Contains(messageLower, "извини") || strings.Contains(messageLower, "прости") {
		return f.getApologyReply(styleName), nil
	} else if strings.Contains(request.Message, "?") {
		return f.getQuestionReply(styleName), nil
	}

	return f.getGenericReply(styleName), nil
}

func (f *Fallback) getThankYouReply(styleName string) string {
	replies := map[string]string{
		"friendly":   "Пожалуйста! Всегда рад помочь! 😊",
		"formal":     "Не за что. Обращайтесь при необходимости.",
		"brief":      "👍",
		"detailed":   "Пожалуйста! Я очень рад, что смог помочь. Если возникнут еще вопросы, не стесняйтесь обращаться.",
		"humorous":   "Да не за что! Я же не зверь какой-то 😄",
		"empathetic": "Всегда пожалуйста! Очень приятно помогать хорошим людям ❤️",
		"motivating": "С удовольствием! Продолжай в том же духе! 💪",
	}

	if reply, exists := replies[styleName]; exists {
		return reply
	}
	return "Пожалуйста!"
}

func (f *Fallback) getApologyReply(styleName string) string {
	replies := map[string]string{
		"friendly":   "Да ладно, все нормально! Не переживай 😊",
		"formal":     "Принято. Вопрос закрыт.",
		"brief":      "Окей",
		"detailed":   "Все в порядке, такое случается с каждым. Главное, что мы это обсудили.",
		"humorous":   "Да забей, я не обидчивый! 😄",
		"empathetic": "Не стоит переживать, все мы люди. Ты хороший человек ❤️",
		"motivating": "Это показывает твою зрелость! Признавать ошибки - это сила! 💪",
	}

	if reply, exists := replies[styleName]; exists {
		return reply
	}
	return "Все нормально!"
}

func (f *Fallback) getQuestionReply(styleName string) string {
	replies := map[string]string{
		"friendly":   "Хороший вопрос! Давай разберемся вместе 😊",
		"formal":     "Рассмотрю ваш вопрос и предоставлю информацию.",
		"brief":      "Разберусь",
		"detailed":   "Это интересный вопрос, который требует детального рассмотрения. Давайте обсудим все аспекты.",
		"humorous":   "О, это тот самый вопрос на миллион! 😄 Сейчас подумаем...",
		"empathetic": "Понимаю, что этот вопрос важен для тебя. Давай вместе найдем ответ ❤️",
		"motivating": "Отличный вопрос! Любознательность - путь к успеху! 💪",
	}

	if reply, exists := replies[styleName]; exists {
		return reply
	}
	return "Интересный вопрос!"
}

func (f *Fallback) getGenericReply(styleName string) string {
	replies := map[string]string{
		"friendly":   "Понял тебя! Звучит интересно 😊",
		"formal":     "Принято к сведению. Благодарю за информацию.",
		"brief":      "Ясно",
		"detailed":   "Спасибо, что поделился этой информацией. Это действительно важно обсудить.",
		"humorous":   "Ага, попал в точку! 😄",
		"empathetic": "Слышу тебя. Спасибо, что доверяешь мне свои мысли ❤️",
		"motivating": "Круто! Продолжай делиться своими идеями! 💪",
	}

	if reply, exists := replies[styleName]; exists {
		return reply
	}
	return "Понятно!"
}
