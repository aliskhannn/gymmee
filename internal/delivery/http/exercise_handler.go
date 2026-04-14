package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aliskhannn/gym-log/internal/domain"
)

// ExerciseService defines methods for managing exercises in the HTTP layer.
type ExerciseService interface {
	GetAllAvailable(ctx context.Context, userID int64) ([]domain.Exercise, error)
}

// ExerciseHandler manages HTTP traffic for exercise catalog features.
type ExerciseHandler struct {
	exerciseService ExerciseService
	userService     UserService
}

// NewExerciseHandler creates a new instance of ExerciseHandler.
func NewExerciseHandler(es ExerciseService, us UserService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService: es, userService: us}
}

// GetAll returns a unified list of system and user-specific exercises.
func (h *ExerciseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tgUser, ok := r.Context().Value(UserContextKey).(*TelegramUser)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetOrCreateUser(r.Context(), tgUser.ID, &tgUser.Username)
	if err != nil {
		http.Error(w, "Failed to identify user", http.StatusInternalServerError)
		return
	}

	exercises, err := h.exerciseService.GetAllAvailable(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exercises)
}