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
) *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoints - Users
	mux.HandleFunc("GET /api/me", AuthMiddleware(botToken, userHandler.GetMe))

	// API Endpoints - Workouts
	mux.HandleFunc("POST /api/workouts/start", AuthMiddleware(botToken, workoutHandler.StartSession))
	mux.HandleFunc("GET /api/workouts/hints", AuthMiddleware(botToken, workoutHandler.GetHint))

	// API Endpoints - Habits
	mux.HandleFunc("GET /api/habits/daily", AuthMiddleware(botToken, habitHandler.GetDaily))
	mux.HandleFunc("POST /api/habits/toggle", AuthMiddleware(botToken, habitHandler.Toggle))

	return mux
}
