package ai

import (
	"context"
	"fmt"
	"time"

	"tg-replyBot/pkg/logger"
)

type SmartFallback struct {
	ollama   Provider
	fallback Provider
	logger   logger.Logger
}

func NewSmartFallback(ollama Provider, fallback Provider, logger logger.Logger) Provider {
	return &SmartFallback{
		ollama:   ollama,
		fallback: fallback,
		logger:   logger,
	}
}

func (s *SmartFallback) GenerateReply(ctx context.Context, request Request) (string, error) {
	ollamaCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.logger.Debug("Пытаемся получить ответ от Ollama")

	reply, err := s.ollama.GenerateReply(ollamaCtx, request)
	if err == nil && reply != "" {
		s.logger.Debug("Получен ответ от Ollama", "reply_length", len(reply))
		return reply, nil
	}

	s.logger.Warn("Ollama недоступен, переключаемся на fallback", "error", err)

	fallbackReply, fallbackErr := s.fallback.GenerateReply(ctx, request)
	if fallbackErr != nil {
		s.logger.Error("Ошибка fallback провайдера", "error", fallbackErr)
		return "", fmt.Errorf("все провайдеры недоступны. Ollama: %w, Fallback: %v", err, fallbackErr)
	}

	s.logger.Debug("Получен ответ от fallback", "reply", fallbackReply)
	return fallbackReply, nil
}
