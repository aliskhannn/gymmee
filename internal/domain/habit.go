package domain

import "time"

// Habit represents a daily goal or routine tracked by a User.
type Habit struct {
	ID        int64
	UserID    int64
	Name      string
	CreatedAt time.Time
}

// HabitLog represents the completion status of a Habit on a specific Date.
type HabitLog struct {
	ID        int64
	HabitID   int64
	Date      time.Time
	Completed bool
}
