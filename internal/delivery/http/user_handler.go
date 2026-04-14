package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aliskhannn/gym-log/internal/domain"
)

// UserService defines the methods required by the UserHandler.
type UserService interface {
	GetOrCreateUser(ctx context.Context, telegramID int64, username *string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
}

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	userService UserService
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe returns the current authenticated user's profile.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	tgUser, ok := r.Context().Value(UserContextKey).(*TelegramUser)
	if !ok {
		http.Error(w, "Internal Server Error: user context missing", http.StatusInternalServerError)
		return
	}

	var username *string
	if tgUser.Username != "" {
		username = &tgUser.Username
	}

	user, err := h.userService.GetOrCreateUser(r.Context(), tgUser.ID, username)
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type UpdateProfileRequest struct {
	Height          *float64 `json:"height"`
	Weight          *float64 `json:"weight"`
	TargetWeight    *float64 `json:"target_weight"`
	BarbellWeight   float64  `json:"barbell_weight"`
	AvailablePlates string   `json:"available_plates"`
}

// UpdateMe updates user's physical parameters and gym settings.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	tgUser, _ := r.Context().Value(UserContextKey).(*TelegramUser)
	user, _ := h.userService.GetOrCreateUser(r.Context(), tgUser.ID, &tgUser.Username)

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	user.Height = req.Height
	user.Weight = req.Weight
	user.TargetWeight = req.TargetWeight
	user.BarbellWeight = req.BarbellWeight
	user.AvailablePlates = req.AvailablePlates

	if err := h.userService.UpdateUser(r.Context(), user); err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
