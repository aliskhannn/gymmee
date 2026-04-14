package domain

import "time"

// User represents a core system user, typically linked via Telegram.
type User struct {
	ID              int64
	TelegramID      int64
	Username        *string
	Height          *float64
	Weight          *float64
	TargetWeight    *float64
	BarbellWeight   float64
	AvailablePlates string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
