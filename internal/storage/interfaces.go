package storage

import (
	"context"
	"tg-replyBot/internal/models"
)

type Storage interface {
	// Пользователи
	GetUser(ctx context.Context, userID int64) (*models.User, error)
	SaveUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int64) error

	// Контекст
	GetContext(ctx context.Context, userID int64) (*models.Context, error)
	SaveContext(ctx context.Context, context *models.Context) error
	DeleteContext(ctx context.Context, userID int64) error

	// Очистка
	Cleanup(ctx context.Context) error
}
