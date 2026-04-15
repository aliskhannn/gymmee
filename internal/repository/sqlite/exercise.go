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

type dbExercise struct {
	ID          int64  `db:"id"`
	UserID      *int64 `db:"user_id"`
	Name        string `db:"name"`
	MuscleGroup string `db:"muscle_group"`
}

func toDBExercise(d *domain.Exercise) *dbExercise {
	return &dbExercise{
		ID:          d.ID,
		UserID:      d.UserID,
		Name:        d.Name,
		MuscleGroup: d.MuscleGroup,
	}
}

func toDomainExercise(dbE *dbExercise) domain.Exercise {
	return domain.Exercise{
		ID:          dbE.ID,
		UserID:      dbE.UserID,
		Name:        dbE.Name,
		MuscleGroup: dbE.MuscleGroup,
	}
}

// NewExerciseRepository creates a new instance of ExerciseRepository.
func NewExerciseRepository(db *sqlx.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

// GetSystemExercises retrieves all basic exercises provided by the system.
func (r *ExerciseRepository) GetSystemExercises(ctx context.Context) ([]domain.Exercise, error) {
	query := `SELECT * FROM exercises WHERE user_id IS NULL ORDER BY name ASC`
	var dbExercises []dbExercise

	if err := r.db.SelectContext(ctx, &dbExercises, query); err != nil {
		return nil, err
	}

	exercises := make([]domain.Exercise, 0, len(dbExercises))
	for _, e := range dbExercises {
		exercises = append(exercises, toDomainExercise(&e))
	}
	return exercises, nil
}

// GetUserExercises retrieves custom exercises created by a specific user.
func (r *ExerciseRepository) GetUserExercises(ctx context.Context, userID int64) ([]domain.Exercise, error) {
	query := `SELECT * FROM exercises WHERE user_id = $1 ORDER BY name ASC`
	var dbExercises []dbExercise

	if err := r.db.SelectContext(ctx, &dbExercises, query, userID); err != nil {
		return nil, err
	}

	exercises := make([]domain.Exercise, 0, len(dbExercises))
	for _, e := range dbExercises {
		exercises = append(exercises, toDomainExercise(&e))
	}
	return exercises, nil
}

// Create adds a new custom exercise for a user.
func (r *ExerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	dbe := toDBExercise(exercise)
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

	if err := stmt.QueryRowxContext(ctx, dbe).Scan(&exercise.ID); err != nil {
		return err
	}
	return nil
}
