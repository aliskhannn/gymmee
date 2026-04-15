// Package sqlite provides SQLite implementations for data access.
package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/aliskhannn/gymmee/internal/domain"
)

// ExerciseRepository is the SQLite implementation for exercise data operations.
type ExerciseRepository struct {
	db *sqlx.DB
}

// NewExerciseRepository creates a new instance of ExerciseRepository.
func NewExerciseRepository(db *sqlx.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// GetSystemExercises retrieves all basic exercises provided by the system.
func (r *ExerciseRepository) GetSystemExercises(ctx context.Context) ([]domain.Exercise, error) {
	query := `SELECT * FROM exercises WHERE user_id IS NULL ORDER BY name ASC`
	var exercises []domain.Exercise

	if err := r.db.SelectContext(ctx, &exercises, query); err != nil {
		return nil, err
	}
	return exercises, nil
}

// GetUserExercises retrieves custom exercises created by a specific user.
func (r *ExerciseRepository) GetUserExercises(ctx context.Context, userID int64) ([]domain.Exercise, error) {
	query := `SELECT * FROM exercises WHERE user_id = $1 ORDER BY name ASC`
	var exercises []domain.Exercise

	if err := r.db.SelectContext(ctx, &exercises, query, userID); err != nil {
		return nil, err
	}
	return exercises, nil
}

// Create adds a new custom exercise for a user.
func (r *ExerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	query := `
		INSERT INTO exercises (user_id, name, muscle_group)
		VALUES (:user_id, :name, :muscle_group)
		RETURNING id
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRowxContext(ctx, exercise).Scan(&exercise.ID)
}
