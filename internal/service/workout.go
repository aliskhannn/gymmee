// Package service implements the core business logic.
package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aliskhannn/gym-log/internal/domain"
	"github.com/aliskhannn/gym-log/pkg/calculator"
)

// WorkoutService coordinates business operations related to workouts.
type WorkoutService struct {
	workoutRepo  WorkoutRepository
	exerciseRepo ExerciseRepository
}

// NewWorkoutService creates a new instance of WorkoutService.
func NewWorkoutService(wRepo WorkoutRepository, eRepo ExerciseRepository) *WorkoutService {
	return &WorkoutService{
		workoutRepo:  wRepo,
		exerciseRepo: eRepo,
	}
}

// SetHintResult represents the context hint for the frontend.
type SetHintResult struct {
	LastWeight     float64                       `json:"last_weight"`
	LastReps       int                           `json:"last_reps"`
	PlatesRequired []calculator.PlateRequirement `json:"plates_required"`
}

// GetExerciseHint returns the user's last performance on an exercise and the plates needed.
func (s *WorkoutService) GetExerciseHint(ctx context.Context, user *domain.User, exerciseID int64) (*SetHintResult, error) {
	stats, err := s.workoutRepo.GetLastSetStats(ctx, user.ID, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get last set stats: %w", err)
	}

	if stats == nil {
		return nil, nil
	}

	var availablePlates []float64
	if err := json.Unmarshal([]byte(user.AvailablePlates), &availablePlates); err != nil {
		availablePlates = []float64{25, 20, 15, 10, 5, 2.5, 1.25}
	}

	plates := calculator.CalculatePlates(stats.Weight, user.BarbellWeight, availablePlates)

	return &SetHintResult{
		LastWeight:     stats.Weight,
		LastReps:       stats.Reps,
		PlatesRequired: plates,
	}, nil
}

// StartSession initiates a new workout session for the user.
func (s *WorkoutService) StartSession(ctx context.Context, userID int64, planDayID *int64) (*domain.WorkoutSession, error) {
	active, err := s.workoutRepo.GetActiveSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return active, nil
	}

	session := &domain.WorkoutSession{
		UserID:    userID,
		PlanDayID: planDayID,
	}

	if err := s.workoutRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}
