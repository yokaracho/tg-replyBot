package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"tg-replyBot/internal/models"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockStorage) SaveUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStorage) DeleteUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStorage) GetContext(ctx context.Context, userID int64) (*models.Context, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return models.NewContext(userID), args.Error(1)
	}
	return args.Get(0).(*models.Context), args.Error(1)
}

func (m *MockStorage) SaveContext(ctx context.Context, context *models.Context) error {
	args := m.Called(ctx, context)
	return args.Error(0)
}

func (m *MockStorage) DeleteContext(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStorage) Cleanup(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
