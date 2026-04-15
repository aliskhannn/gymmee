package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/aliskhannn/gymmee/internal/domain"
)

// WorkoutRepository is the SQLite implementation for workout sessions and sets.
type WorkoutRepository struct {
	db *sqlx.DB
}

type dbWorkoutSession struct {
	ID        int64      `db:"id"`
	UserID    int64      `db:"user_id"`
	PlanDayID *int64     `db:"plan_day_id"`
	StartedAt time.Time  `db:"started_at"`
	EndedAt   *time.Time `db:"ended_at"`
}

func toDBWorkoutSession(d *domain.WorkoutSession) *dbWorkoutSession {
	return &dbWorkoutSession{
		ID:        d.ID,
		UserID:    d.UserID,
		PlanDayID: d.PlanDayID,
		StartedAt: d.StartedAt,
		EndedAt:   d.EndedAt,
	}
}

func toDomainWorkoutSession(dbW *dbWorkoutSession) *domain.WorkoutSession {
	return &domain.WorkoutSession{
		ID:        dbW.ID,
		UserID:    dbW.UserID,
		PlanDayID: dbW.PlanDayID,
		StartedAt: dbW.StartedAt,
		EndedAt:   dbW.EndedAt,
	}
}

type dbWorkoutSet struct {
	ID               int64     `db:"id"`
	WorkoutSessionID int64     `db:"workout_session_id"`
	ExerciseID       int64     `db:"exercise_id"`
	Weight           float64   `db:"weight"`
	Reps             int       `db:"reps"`
	CreatedAt        time.Time `db:"created_at"`
}

func toDBWorkoutSet(d *domain.WorkoutSet) *dbWorkoutSet {
	return &dbWorkoutSet{
		ID:               d.ID,
		WorkoutSessionID: d.WorkoutSessionID,
		ExerciseID:       d.ExerciseID,
		Weight:           d.Weight,
		Reps:             d.Reps,
		CreatedAt:        d.CreatedAt,
	}
}

// LastSetStats is a DTO used internally, so we leave it as is.
type LastSetStats struct {
	Weight float64 `db:"weight"`
	Reps   int     `db:"reps"`
}

// NewWorkoutRepository creates a new instance of WorkoutRepository.
func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

// CreateSession starts a new workout session for a user.
func (r *WorkoutRepository) CreateSession(ctx context.Context, session *domain.WorkoutSession) error {
	dbws := toDBWorkoutSession(session)
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

	if err := stmt.QueryRowxContext(ctx, dbws).Scan(&session.ID, &session.StartedAt); err != nil {
		return err
	}
	return nil
}

// GetActiveSession returns an ongoing workout session (where ended_at IS NULL).
func (r *WorkoutRepository) GetActiveSession(ctx context.Context, userID int64) (*domain.WorkoutSession, error) {
	query := `SELECT * FROM workout_sessions WHERE user_id = $1 AND ended_at IS NULL LIMIT 1`
	var dbws dbWorkoutSession

	err := r.db.GetContext(ctx, &dbws, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomainWorkoutSession(&dbws), nil
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
	dbs := toDBWorkoutSet(set)
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

	if err := stmt.QueryRowxContext(ctx, dbs).Scan(&set.ID, &set.CreatedAt); err != nil {
		return err
	}
	return nil
}

// GetLastSetStats retrieves the weight and reps from the user's most recent interaction with an exercise.
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

// GetHistory retrieves all completed workout sessions for a user.
func (r *WorkoutRepository) GetHistory(ctx context.Context, userID int64) ([]domain.WorkoutSession, error) {
	query := `SELECT * FROM workout_sessions WHERE user_id = $1 AND ended_at IS NOT NULL ORDER BY started_at DESC`
	var dbSessions []dbWorkoutSession
	if err := r.db.SelectContext(ctx, &dbSessions, query, userID); err != nil {
		return nil, err
	}

	sessions := make([]domain.WorkoutSession, 0, len(dbSessions))
	for _, s := range dbSessions {
		sessions = append(sessions, *toDomainWorkoutSession(&s))
	}
	return sessions, nil
}
