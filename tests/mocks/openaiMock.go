package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"tg-replyBot/internal/ai"
)

type MockAIProvider struct {
	mock.Mock
}

func (m *MockAIProvider) GenerateReply(ctx context.Context, request ai.Request) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}
