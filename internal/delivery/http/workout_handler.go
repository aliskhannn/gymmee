package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aliskhannn/gym-log/internal/domain"
	"github.com/aliskhannn/gym-log/internal/service"
)

// WorkoutService defines the methods required by the WorkoutHandler.
type WorkoutService interface {
	StartSession(ctx context.Context, userID int64, planDayID *int64) (*domain.WorkoutSession, error)
	GetExerciseHint(ctx context.Context, user *domain.User, exerciseID int64) (*service.SetHintResult, error)
}

// WorkoutHandler handles HTTP requests related to workouts.
type WorkoutHandler struct {
	workoutService WorkoutService
	userService    UserService
}

// NewWorkoutHandler creates a new instance of WorkoutHandler.
func NewWorkoutHandler(ws WorkoutService, us UserService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: ws,
		userService:    us,
	}
}

// StartSessionRequest represents the payload for starting a workout.
type StartSessionRequest struct {
	PlanDayID *int64 `json:"plan_day_id"`
}

// StartSession handles POST requests to start or resume a workout session.
func (h *WorkoutHandler) StartSession(w http.ResponseWriter, r *http.Request) {
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

	var req StartSessionRequest
	if r.Body != nil && r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
	}

	session, err := h.workoutService.StartSession(r.Context(), user.ID, req.PlanDayID)
	if err != nil {
		http.Error(w, "Failed to start session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetHint handles GET requests to fetch the last set stats and plate calculations.
func (h *WorkoutHandler) GetHint(w http.ResponseWriter, r *http.Request) {
	tgUser, ok := r.Context().Value(UserContextKey).(*TelegramUser)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	exerciseIDStr := r.URL.Query().Get("exercise_id")
	if exerciseIDStr == "" {
		http.Error(w, "Missing exercise_id parameter", http.StatusBadRequest)
		return
	}

	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid exercise_id format", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetOrCreateUser(r.Context(), tgUser.ID, &tgUser.Username)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	hint, err := h.workoutService.GetExerciseHint(r.Context(), user, exerciseID)
	if err != nil {
		http.Error(w, "Failed to get exercise hint", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if hint == nil {
		w.Write([]byte(`{}`))
		return
	}

	json.NewEncoder(w).Encode(hint)
}
