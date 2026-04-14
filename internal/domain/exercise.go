package domain

// Exercise represents a physical movement.
// If UserID is nil, it represents a global system-provided exercise.
type Exercise struct {
	ID          int64
	UserID      *int64
	Name        string
	MuscleGroup string
}
