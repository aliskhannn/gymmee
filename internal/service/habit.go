package service

import (
	"context"
	"time"

	"github.com/aliskhannn/gymmee/internal/domain"
)

// HabitWithStatus combines a habit with its completion status for a specific day.
type HabitWithStatus struct {
	domain.Habit
	Completed bool `json:"completed"`
}

// HabitService manages user habits and tracking.
type HabitService struct {
	repo HabitRepository
}

// NewHabitService creates a new instance of HabitService.
func NewHabitService(repo HabitRepository) *HabitService {
	return &HabitService{repo: repo}
}

// GetDailyHabits returns all user habits with their completion status for the given date.
func (s *HabitService) GetDailyHabits(ctx context.Context, userID int64, date time.Time) ([]HabitWithStatus, error) {
	habits, err := s.repo.GetUserHabits(ctx, userID)
	if err != nil {
		return nil, err
	}

	logs, err := s.repo.GetDailyLogs(ctx, userID, date)
	if err != nil {
		return nil, err
	}
	
	statusMap := make(map[int64]bool)
	for _, log := range logs {
		statusMap[log.HabitID] = log.Completed
	}

	result := make([]HabitWithStatus, len(habits))
	for i, h := range habits {
		result[i] = HabitWithStatus{
			Habit:     h,
			Completed: statusMap[h.ID],
		}
	}

	return result, nil
}

// ToggleHabit updates the status of a habit for today.
func (s *HabitService) ToggleHabit(ctx context.Context, habitID int64, completed bool) error {
	return s.repo.ToggleLog(ctx, habitID, time.Now(), completed)
}

// CreateHabit registers a new habit for tracking.
func (s *HabitService) CreateHabit(ctx context.Context, userID int64, name string) (*domain.Habit, error) {
	habit := &domain.Habit{
		UserID: userID,
		Name:   name,
	}
	if err := s.repo.Create(ctx, habit); err != nil {
		return nil, err
	}
	return habit, nil
}
