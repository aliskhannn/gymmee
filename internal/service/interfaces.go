package service

import (
	"context"
	"time"

	"github.com/aliskhannn/gym-log/internal/domain"
	"github.com/aliskhannn/gym-log/internal/repository/sqlite"
)

// UserRepository defines the data access methods required by the UserService.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

// HabitRepository defines the data access methods for habits.
type HabitRepository interface {
	Create(ctx context.Context, habit *domain.Habit) error
	GetUserHabits(ctx context.Context, userID int64) ([]domain.Habit, error)
	GetDailyLogs(ctx context.Context, userID int64, date time.Time) ([]domain.HabitLog, error)
	ToggleLog(ctx context.Context, habitID int64, date time.Time, completed bool) error
}

// WorkoutRepository defines the data access methods for workouts.
type WorkoutRepository interface {
	CreateSession(ctx context.Context, session *domain.WorkoutSession) error
	GetActiveSession(ctx context.Context, userID int64) (*domain.WorkoutSession, error)
	FinishSession(ctx context.Context, sessionID int64) error
	AddSet(ctx context.Context, set *domain.WorkoutSet) error
	GetLastSetStats(ctx context.Context, userID, exerciseID int64) (*sqlite.LastSetStats, error)
	GetHistory(ctx context.Context, userID int64) ([]domain.WorkoutSession, error) // <-- Добавили
}

// ExerciseRepository defines the data access methods for exercises.
type ExerciseRepository interface {
	GetSystemExercises(ctx context.Context) ([]domain.Exercise, error)
	GetUserExercises(ctx context.Context, userID int64) ([]domain.Exercise, error)
	Create(ctx context.Context, exercise *domain.Exercise) error
}
