package http

import (
	"net/http"
)

// SetupRouter creates and configures the standard HTTP multiplexer.
func SetupRouter(botToken string, userHandler *UserHandler, workoutHandler *WorkoutHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("GET /api/me", AuthMiddleware(botToken, userHandler.GetMe))

	// API Endpoints - Workouts
	mux.HandleFunc("POST /api/workouts/start", AuthMiddleware(botToken, workoutHandler.StartSession))
	mux.HandleFunc("GET /api/workouts/hints", AuthMiddleware(botToken, workoutHandler.GetHint))

	return mux
}
