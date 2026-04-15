// Package service implements the core business logic of the application.
package service

import (
	"context"
	"fmt"

	"github.com/aliskhannn/gymmee/internal/domain"
)

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

// UpdateUser updates the user's profile information.
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
