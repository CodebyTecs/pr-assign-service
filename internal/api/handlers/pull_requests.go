package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/service"
)

type PRHandler struct {
	prService service.PRService
}

func NewPRHandler(prService service.PRService) *PRHandler {
	return &PRHandler{
		prService: prService,
	}
}

type prCreateRequest struct {
	PRId     string `json:"pull_request_id"`
	PRName   string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req prCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if req.PRId == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request id is required")
		return
	}
	if req.PRName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_name is required")
		return
	}
	if req.AuthorID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "author_id is required")
		return
	}

	input := service.CreatePRInput{
		ID:     req.PRId,
		Name:   req.PRName,
		Author: req.AuthorID,
	}

	pr, err := h.prService.Create(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "author or team not found")
			return

		case errors.Is(err, domain.ErrPRExists):
			writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
			return

		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
			return
		}
	}

	resp := struct {
		PR *domain.PullRequest `json:"pr"`
	}{
		PR: pr,
	}

	writeJSON(w, http.StatusCreated, resp)
}

type prMergeRequest struct {
	PRId string `json:"pull_request_id"`
}

func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req prMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if req.PRId == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}

	pr, err := h.prService.Merge(ctx, req.PRId)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "pullRequest not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}

	resp := struct {
		PR *domain.PullRequest `json:"pr"`
	}{
		PR: pr,
	}

	writeJSON(w, http.StatusOK, resp)
}

type prReassignRequest struct {
	PRId       string `json:"pull_request_id"`
	ReviewerID string `json:"old_user_id"`
}

func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req prReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if req.PRId == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}
	if req.ReviewerID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "old_user_id is required")
		return
	}

	input := service.ReassignReviewerInput{
		PullRequestID: req.PRId,
		ReviewerID:    req.ReviewerID,
	}

	pr, id, err := h.prService.Reassign(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			writeError(w, http.StatusNotFound, "NOT_FOUND", "PR or user not found")
			return
		case errors.Is(err, domain.ErrPRMerged):
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
			return
		case errors.Is(err, domain.ErrNotAssigned):
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
			return
		case errors.Is(err, domain.ErrNoCandidate):
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
			return
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
			return
		}
	}

	resp := struct {
		PR         *domain.PullRequest `json:"pr"`
		ReplacedBy string              `json:"replaced_by"`
	}{
		PR:         pr,
		ReplacedBy: id,
	}

	writeJSON(w, http.StatusOK, resp)
}
