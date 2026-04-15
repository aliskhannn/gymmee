package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

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

type dbUser struct {
	ID              int64     `db:"id"`
	TelegramID      int64     `db:"telegram_id"`
	Username        *string   `db:"username"`
	Height          *float64  `db:"height"`
	Weight          *float64  `db:"weight"`
	TargetWeight    *float64  `db:"target_weight"`
	BarbellWeight   float64   `db:"barbell_weight"`
	AvailablePlates string    `db:"available_plates"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func toDBUser(d *domain.User) *dbUser {
	return &dbUser{
		ID:              d.ID,
		TelegramID:      d.TelegramID,
		Username:        d.Username,
		Height:          d.Height,
		Weight:          d.Weight,
		TargetWeight:    d.TargetWeight,
		BarbellWeight:   d.BarbellWeight,
		AvailablePlates: d.AvailablePlates,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
	}
}

func toDomainUser(dbU *dbUser) *domain.User {
	return &domain.User{
		ID:              dbU.ID,
		TelegramID:      dbU.TelegramID,
		Username:        dbU.Username,
		Height:          dbU.Height,
		Weight:          dbU.Weight,
		TargetWeight:    dbU.TargetWeight,
		BarbellWeight:   dbU.BarbellWeight,
		AvailablePlates: dbU.AvailablePlates,
		CreatedAt:       dbU.CreatedAt,
		UpdatedAt:       dbU.UpdatedAt,
	}
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	dbu := toDBUser(user)

	query := `
       INSERT INTO users (telegram_id, username, height, weight, target_weight)
       VALUES (:telegram_id, :username, :height, :weight, :target_weight)
       RETURNING id, created_at, updated_at
    `

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, dbu).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var dbu dbUser

	err := r.db.GetContext(ctx, &dbu, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return toDomainUser(&dbu), nil
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	query := `SELECT * FROM users WHERE telegram_id = $1`
	var dbu dbUser

	err := r.db.GetContext(ctx, &dbu, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomainUser(&dbu), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	dbu := toDBUser(user)

	query := `
       UPDATE users 
       SET username = :username, 
           height = :height, 
           weight = :weight, 
           target_weight = :target_weight,
           updated_at = CURRENT_TIMESTAMP
       WHERE id = :id
    `

	result, err := r.db.NamedExecContext(ctx, query, dbu)
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
