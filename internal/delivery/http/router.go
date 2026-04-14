package http

import (
	"net/http"
)

// SetupRouter creates and configures the standard HTTP multiplexer.
func SetupRouter(botToken string, userHandler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// API Endpoints
	mux.HandleFunc("GET /api/me", AuthMiddleware(botToken, userHandler.GetMe))

	return mux
}
