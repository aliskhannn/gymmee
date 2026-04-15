// Package sqlite provides SQLite implementations for data access.
package sqlite

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/aliskhannn/gymmee/internal/domain"
)

// HabitRepository handles database operations for user habits and their logs.
type HabitRepository struct {
	db *sqlx.DB
}

// NewHabitRepository creates a new instance of HabitRepository.
func NewHabitRepository(db *sqlx.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

// Create adds a new habit to the user's tracker.
func (r *HabitRepository) Create(ctx context.Context, habit *domain.Habit) error {
	query := `INSERT INTO habits (user_id, name) VALUES (:user_id, :name) RETURNING id, created_at`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, habit).Scan(&habit.ID, &habit.CreatedAt)
}

// GetUserHabits retrieves all habits for a specific user.
func (r *HabitRepository) GetUserHabits(ctx context.Context, userID int64) ([]domain.Habit, error) {
	query := `SELECT * FROM habits WHERE user_id = $1 ORDER BY created_at ASC`
	var habits []domain.Habit
	if err := r.db.SelectContext(ctx, &habits, query, userID); err != nil {
		return nil, err
	}
	return habits, nil
}

// GetDailyLogs retrieves completion status for all habits on a specific date.
func (r *HabitRepository) GetDailyLogs(ctx context.Context, userID int64, date time.Time) ([]domain.HabitLog, error) {
	// Приводим время к формату YYYY-MM-DD для SQLite
	dateStr := date.Format("2006-01-02")
	query := `
		SELECT hl.* FROM habit_logs hl
		JOIN habits h ON hl.habit_id = h.id
		WHERE h.user_id = $1 AND hl.date = $2
	`
	var logs []domain.HabitLog
	if err := r.db.SelectContext(ctx, &logs, query, userID, dateStr); err != nil {
		return nil, err
	}
	return logs, nil
}

// ToggleLog flips the completion status of a habit for a specific date.
func (r *HabitRepository) ToggleLog(ctx context.Context, habitID int64, date time.Time, completed bool) error {
	dateStr := date.Format("2006-01-02")
	query := `
		INSERT INTO habit_logs (habit_id, date, completed)
		VALUES ($1, $2, $3)
		ON CONFLICT(habit_id, date) DO UPDATE SET completed = EXCLUDED.completed
	`
	_, err := r.db.ExecContext(ctx, query, habitID, dateStr, completed)
	return err
}
