package ai

import (
	"context"

	"tg-replyBot/internal/models"
)

type Provider interface {
	GenerateReply(ctx context.Context, request Request) (string, error)
}

type Request struct {
	Message         string
	ContextMessages []string
	Style           models.Style
}
