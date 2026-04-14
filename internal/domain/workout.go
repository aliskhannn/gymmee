package domain

import "time"

// WorkoutSession represents a specific training session performed by a User.
// It can optionally be linked to a structured PlanDay.
type WorkoutSession struct {
	ID        int64
	UserID    int64
	PlanDayID *int64
	StartedAt time.Time
	EndedAt   *time.Time
}

// WorkoutSet represents a single set of an Exercise performed during a WorkoutSession.
type WorkoutSet struct {
	ID               int64
	WorkoutSessionID int64
	ExerciseID       int64
	Weight           float64
	Reps             int
	CreatedAt        time.Time
}
