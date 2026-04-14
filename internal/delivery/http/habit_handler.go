package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aliskhannn/gym-log/internal/domain"
	"github.com/aliskhannn/gym-log/internal/service"
)

// HabitService defines methods for managing habits in the HTTP layer.
type HabitService interface {
	GetDailyHabits(ctx context.Context, userID int64, date time.Time) ([]service.HabitWithStatus, error)
	ToggleHabit(ctx context.Context, habitID int64, completed bool) error
	CreateHabit(ctx context.Context, userID int64, name string) (*domain.Habit, error)
}

// HabitHandler manages HTTP traffic for habit-related features.
type HabitHandler struct {
	habitService HabitService
	userService  UserService
}

// NewHabitHandler creates a new instance of HabitHandler.
func NewHabitHandler(hs HabitService, us UserService) *HabitHandler {
	return &HabitHandler{habitService: hs, userService: us}
}

// GetDaily returns habits and their status for the current day.
func (h *HabitHandler) GetDaily(w http.ResponseWriter, r *http.Request) {
	tgUser, _ := r.Context().Value(UserContextKey).(*TelegramUser)
	user, err := h.userService.GetOrCreateUser(r.Context(), tgUser.ID, &tgUser.Username)
	if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	habits, err := h.habitService.GetDailyHabits(r.Context(), user.ID, time.Now())
	if err != nil {
		http.Error(w, "Failed to fetch daily habits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(habits)
}

type ToggleRequest struct {
	HabitID   int64 `json:"habit_id"`
	Completed bool  `json:"completed"`
}

// Toggle handles checking/unchecking a habit.
func (h *HabitHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	var req ToggleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.habitService.ToggleHabit(r.Context(), req.HabitID, req.Completed); err != nil {
		http.Error(w, "Failed to toggle habit", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
