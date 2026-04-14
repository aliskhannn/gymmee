package http

import (
	"net/http"
)

// SetupRouter creates and configures the standard HTTP multiplexer.
func SetupRouter(
	botToken string,
	userHandler *UserHandler,
	workoutHandler *WorkoutHandler,
	habitHandler *HabitHandler,
	exerciseHandler *ExerciseHandler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoints - Users
	mux.HandleFunc("GET /api/me", AuthMiddleware(botToken, userHandler.GetMe))

	// API Endpoints - Workouts
	mux.HandleFunc("POST /api/workouts/start", AuthMiddleware(botToken, workoutHandler.StartSession))
	mux.HandleFunc("GET /api/workouts/hints", AuthMiddleware(botToken, workoutHandler.GetHint))
	mux.HandleFunc("POST /api/workouts/sets", AuthMiddleware(botToken, workoutHandler.AddSet))
	mux.HandleFunc("POST /api/workouts/finish", AuthMiddleware(botToken, workoutHandler.FinishSession))

	// API Endpoints - Habits
	mux.HandleFunc("GET /api/habits/daily", AuthMiddleware(botToken, habitHandler.GetDaily))
	mux.HandleFunc("POST /api/habits/toggle", AuthMiddleware(botToken, habitHandler.Toggle))
	mux.HandleFunc("POST /api/habits", AuthMiddleware(botToken, habitHandler.Create))

	// API Endpoints - Exercises
	mux.HandleFunc("GET /api/exercises", AuthMiddleware(botToken, exerciseHandler.GetAll))

	return mux
}
