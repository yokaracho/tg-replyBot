package services

import (
	"context"

	"tg-replyBot/internal/models"
	"tg-replyBot/internal/storage"
	"tg-replyBot/pkg/logger"
)

type ContextManager struct {
	storage storage.Storage
	logger  logger.Logger
}

func NewContextManager(storage storage.Storage, logger logger.Logger) *ContextManager {
	return &ContextManager{
		storage: storage,
		logger:  logger,
	}
}

func (cm *ContextManager) GetContext(ctx context.Context, userID int64) (*models.Context, error) {
	userContext, err := cm.storage.GetContext(ctx, userID)
	if err != nil {
		cm.logger.Error("Ошибка получения контекста", "user_id", userID, "error", err)
		return models.NewContext(userID), nil
	}

	return userContext, nil
}

func (cm *ContextManager) AddMessage(ctx context.Context, userID int64, message string) error {
	userContext, err := cm.GetContext(ctx, userID)
	if err != nil {
		return err
	}

	userContext.AddMessage(message)

	if err := cm.storage.SaveContext(ctx, userContext); err != nil {
		cm.logger.Error("Ошибка сохранения контекста", "user_id", userID, "error", err)
		return err
	}

	return nil
}

func (cm *ContextManager) SaveReplies(ctx context.Context, userID int64, replies models.Replies) error {
	userContext, err := cm.GetContext(ctx, userID)
	if err != nil {
		return err
	}

	userContext.Replies = replies

	if err := cm.storage.SaveContext(ctx, userContext); err != nil {
		cm.logger.Error("Ошибка сохранения ответов", "user_id", userID, "error", err)
		return err
	}

	return nil
}

func (cm *ContextManager) ClearContext(ctx context.Context, userID int64) error {
	if err := cm.storage.DeleteContext(ctx, userID); err != nil {
		cm.logger.Error("Ошибка очистки контекста", "user_id", userID, "error", err)
		return err
	}

	return nil
}
