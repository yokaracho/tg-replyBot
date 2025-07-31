package memory

import (
	"context"
	"fmt"
	"sync"
	"time"

	"tg-replyBot/internal/models"
	"tg-replyBot/internal/storage"
)

type Storage struct {
	users    map[int64]*models.User
	contexts map[int64]*models.Context
	mutex    sync.RWMutex
	ttl      time.Duration
}

func New() storage.Storage {
	s := &Storage{
		users:    make(map[int64]*models.User),
		contexts: make(map[int64]*models.Context),
		ttl:      24 * time.Hour,
	}

	go s.startCleanup()

	return s
}

func (s *Storage) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[userID]
	if !exists {
		return nil, fmt.Errorf("пользователь не найден")
	}

	return user, nil
}

func (s *Storage) SaveUser(ctx context.Context, user *models.User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.Updated = time.Now()
	s.users[user.ID] = user
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.users, userID)
	return nil
}

func (s *Storage) GetContext(ctx context.Context, userID int64) (*models.Context, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	context, exists := s.contexts[userID]
	if !exists {
		return models.NewContext(userID), nil
	}

	return context, nil
}

func (s *Storage) SaveContext(ctx context.Context, context *models.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	context.Updated = time.Now()
	s.contexts[context.UserID] = context
	return nil
}

func (s *Storage) DeleteContext(ctx context.Context, userID int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.contexts, userID)
	return nil
}

func (s *Storage) Cleanup(ctx context.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.ttl)

	// Очистка пользователей
	for userID, user := range s.users {
		if user.Updated.Before(cutoff) {
			delete(s.users, userID)
		}
	}

	// Очистка контекстов
	for userID, context := range s.contexts {
		if context.Updated.Before(cutoff) {
			delete(s.contexts, userID)
		}
	}

	return nil
}

func (s *Storage) startCleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.Cleanup(context.Background())
		}
	}
}
