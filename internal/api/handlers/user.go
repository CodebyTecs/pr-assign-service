package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/service"
)

type UserHandler struct {
	userService service.UserService
	prService   service.PRService
}

func NewUserHandler(userService service.UserService, prService service.PRService) *UserHandler {
	return &UserHandler{
		userService: userService,
		prService:   prService,
	}
}

type setIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type userResponse struct {
	User *domain.User `json:"user"`
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req setIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id is required")
		return
	}

	user, err := h.userService.UpdateActivity(ctx, req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, userResponse{User: user})
}

type userReviewResponse struct {
	UserID      string                    `json:"user_id"`
	PullRequest []domain.PullRequestShort `json:"pull_requests"`
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id is required")
		return
	}

	prs, err := h.prService.ListByReviewer(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "user not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	resp := userReviewResponse{
		UserID:      userID,
		PullRequest: prs,
	}

	writeJSON(w, http.StatusOK, resp)
}
