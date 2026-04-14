// Package service implements the core business logic of the application.
package service

import (
	"context"
	"fmt"

	"github.com/aliskhannn/gym-log/internal/domain"
)

// UserRepository defines the data access methods required by the UserService.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// UserService coordinates business operations related to users.
type UserService struct {
	repo UserRepository
}

// NewUserService creates and returns a new instance of UserService.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetOrCreateUser retrieves a user by their Telegram ID.
// If the user does not exist in the database, it creates a new record.
func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, username *string) (*domain.User, error) {
	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by tg id: %w", err)
	}

	// User exists, return immediately
	if user != nil {
		return user, nil
	}

	// User does not exist, create a new one
	newUser := &domain.User{
		TelegramID: telegramID,
		Username:   username,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create new user: %w", err)
	}

	return newUser, nil
}
