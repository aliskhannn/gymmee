package domain

import "time"

// Plan represents a user's training program (e.g., "Full Body", "Upper/Lower Split").
type Plan struct {
	ID        int64
	UserID    int64
	Name      string
	IsActive  bool
	CreatedAt time.Time
}

// PlanDay maps a specific day of the week to a Plan and targeted muscle groups.
type PlanDay struct {
	ID          int64
	PlanID      int64
	DayOfWeek   int // 1 for Monday, 7 for Sunday
	MuscleGroup string
}
