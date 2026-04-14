package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/aliskhannn/gym-log/internal/domain"
)

// WorkoutRepository is the SQLite implementation for workout sessions and sets.
type WorkoutRepository struct {
	db *sqlx.DB
}

// NewWorkoutRepository creates a new instance of WorkoutRepository.
func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

// CreateSession starts a new workout session for a user.
func (r *WorkoutRepository) CreateSession(ctx context.Context, session *domain.WorkoutSession) error {
	query := `
		INSERT INTO workout_sessions (user_id, plan_day_id, started_at)
		VALUES (:user_id, :plan_day_id, CURRENT_TIMESTAMP)
		RETURNING id, started_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, session).Scan(&session.ID, &session.StartedAt)
}

// GetActiveSession returns an ongoing workout session (where ended_at IS NULL).
func (r *WorkoutRepository) GetActiveSession(ctx context.Context, userID int64) (*domain.WorkoutSession, error) {
	query := `SELECT * FROM workout_sessions WHERE user_id = $1 AND ended_at IS NULL LIMIT 1`
	var session domain.WorkoutSession

	err := r.db.GetContext(ctx, &session, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}

// FinishSession marks a workout session as completed by setting ended_at.
func (r *WorkoutRepository) FinishSession(ctx context.Context, sessionID int64) error {
	query := `UPDATE workout_sessions SET ended_at = CURRENT_TIMESTAMP WHERE id = $1 AND ended_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// AddSet logs a new set (weight and reps) for an exercise in a specific session.
func (r *WorkoutRepository) AddSet(ctx context.Context, set *domain.WorkoutSet) error {
	query := `
		INSERT INTO workout_sets (workout_session_id, exercise_id, weight, reps)
		VALUES (:workout_session_id, :exercise_id, :weight, :reps)
		RETURNING id, created_at
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, set).Scan(&set.ID, &set.CreatedAt)
}

// LastSetStats is a DTO used internally to scan the specific query result.
type LastSetStats struct {
	Weight float64 `db:"weight"`
	Reps   int     `db:"reps"`
}

// GetLastSetStats retrieves the weight and reps from the user's most recent interaction with an exercise.
// This is the core logic for the "killer feature" context hints.
func (r *WorkoutRepository) GetLastSetStats(ctx context.Context, userID, exerciseID int64) (*LastSetStats, error) {
	query := `
		SELECT ws.weight, ws.reps
		FROM workout_sets ws
		JOIN workout_sessions sess ON ws.workout_session_id = sess.id
		WHERE sess.user_id = $1 AND ws.exercise_id = $2
		ORDER BY ws.created_at DESC
		LIMIT 1
	`
	var stats LastSetStats

	err := r.db.GetContext(ctx, &stats, query, userID, exerciseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &stats, nil
}
