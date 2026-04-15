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

type dbHabit struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func toDBHabit(d *domain.Habit) *dbHabit {
	return &dbHabit{
		ID:        d.ID,
		UserID:    d.UserID,
		Name:      d.Name,
		CreatedAt: d.CreatedAt,
	}
}

func toDomainHabit(dbH *dbHabit) domain.Habit {
	return domain.Habit{
		ID:        dbH.ID,
		UserID:    dbH.UserID,
		Name:      dbH.Name,
		CreatedAt: dbH.CreatedAt,
	}
}

type dbHabitLog struct {
	ID        int64     `db:"id"`
	HabitID   int64     `db:"habit_id"`
	Date      time.Time `db:"date"`
	Completed bool      `db:"completed"`
}

func toDomainHabitLog(dbHL *dbHabitLog) domain.HabitLog {
	return domain.HabitLog{
		ID:        dbHL.ID,
		HabitID:   dbHL.HabitID,
		Date:      dbHL.Date,
		Completed: dbHL.Completed,
	}
}

// NewHabitRepository creates a new instance of HabitRepository.
func NewHabitRepository(db *sqlx.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

// Create adds a new habit to the user's tracker.
func (r *HabitRepository) Create(ctx context.Context, habit *domain.Habit) error {
	dbh := toDBHabit(habit)
	query := `INSERT INTO habits (user_id, name) VALUES (:user_id, :name) RETURNING id, created_at`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if err := stmt.QueryRowxContext(ctx, dbh).Scan(&habit.ID, &habit.CreatedAt); err != nil {
		return err
	}
	return nil
}

// GetUserHabits retrieves all habits for a specific user.
func (r *HabitRepository) GetUserHabits(ctx context.Context, userID int64) ([]domain.Habit, error) {
	query := `SELECT * FROM habits WHERE user_id = $1 ORDER BY created_at ASC`
	var dbHabits []dbHabit
	if err := r.db.SelectContext(ctx, &dbHabits, query, userID); err != nil {
		return nil, err
	}

	habits := make([]domain.Habit, 0, len(dbHabits))
	for _, h := range dbHabits {
		habits = append(habits, toDomainHabit(&h))
	}
	return habits, nil
}

// GetDailyLogs retrieves completion status for all habits on a specific date.
func (r *HabitRepository) GetDailyLogs(ctx context.Context, userID int64, date time.Time) ([]domain.HabitLog, error) {
	dateStr := date.Format("2006-01-02")
	query := `
       SELECT hl.* FROM habit_logs hl
       JOIN habits h ON hl.habit_id = h.id
       WHERE h.user_id = $1 AND hl.date = $2
    `
	var dbLogs []dbHabitLog
	if err := r.db.SelectContext(ctx, &dbLogs, query, userID, dateStr); err != nil {
		return nil, err
	}

	logs := make([]domain.HabitLog, 0, len(dbLogs))
	for _, l := range dbLogs {
		logs = append(logs, toDomainHabitLog(&l))
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
