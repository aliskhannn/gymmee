// Package service implements the core business logic.
package service

import (
	"context"
	"fmt"

	"github.com/aliskhannn/gymmee/internal/domain"
)

// ExerciseService manages operations related to gym exercises.
type ExerciseService struct {
	repo ExerciseRepository
}

// NewExerciseService creates a new instance of ExerciseService.
func NewExerciseService(repo ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

// GetAllAvailable retrieves both system-wide exercises and user-specific custom exercises.
func (s *ExerciseService) GetAllAvailable(ctx context.Context, userID int64) ([]domain.Exercise, error) {
	systemExs, err := s.repo.GetSystemExercises(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get system exercises: %w", err)
	}

	userExs, err := s.repo.GetUserExercises(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user exercises: %w", err)
	}

	// Merge both lists
	allExercises := append(systemExs, userExs...)
	return allExercises, nil
}