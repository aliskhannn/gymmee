package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/aliskhannn/gymmee/internal/domain"
)

var (
	// ErrNotFound is returned when a requested record is not found in the database.
	ErrNotFound = errors.New("record not found")
)

// UserRepository is the SQLite implementation of UserRepository.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (telegram_id, username, height, weight, target_weight)
		VALUES (:telegram_id, :username, :height, :weight, :target_weight)
		RETURNING id, created_at, updated_at
	`

	// NamedQueryContext позволяет использовать поля структуры напрямую
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, user).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user domain.User

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE telegram_id = $1`
	var user domain.User

	err := r.db.GetContext(ctx, &user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Возвращаем nil без ошибки, если юзер просто еще не зарегистрирован
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET username = :username, 
		    height = :height, 
		    weight = :weight, 
		    target_weight = :target_weight,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, user)
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
