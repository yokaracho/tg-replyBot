package services

import (
	"context"

	"tg-replyBot/internal/ai"
	"tg-replyBot/internal/models"
	"tg-replyBot/pkg/logger"
)

type ReplyGenerator struct {
	aiProvider   ai.Provider
	styleManager *StyleManager
	logger       logger.Logger
}

func NewReplyGenerator(aiProvider ai.Provider, styleManager *StyleManager, logger logger.Logger) *ReplyGenerator {
	return &ReplyGenerator{
		aiProvider:   aiProvider,
		styleManager: styleManager,
		logger:       logger,
	}
}

func (rg *ReplyGenerator) GenerateReplies(ctx context.Context, request models.ReplyRequest) (models.ReplyResponse, error) {
	var replies []models.Reply

	styles := request.Styles
	if len(styles) == 0 {
		allStyles := rg.styleManager.GetAllStyles()
		for _, style := range allStyles {
			styles = append(styles, style.Name)
		}
	}

	for _, styleName := range styles {
		style, err := rg.styleManager.GetStyle(styleName)
		if err != nil {
			rg.logger.Error("Стиль не найден", "style", styleName, "error", err)
			continue
		}

		content, err := rg.aiProvider.GenerateReply(ctx, ai.Request{
			Message:         request.Message,
			ContextMessages: request.ContextMessages,
			Style:           style,
		})

		if err != nil {
			rg.logger.Error("Ошибка генерации ответа", "style", styleName, "error", err)
			continue
		}

		replies = append(replies, models.Reply{
			Style:   styleName,
			Content: content,
			Emoji:   style.Emoji,
		})
	}

	return models.ReplyResponse{Replies: replies}, nil
}
