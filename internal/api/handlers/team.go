package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/service"
)

type TeamHandler struct {
	teamService service.TeamService
}

func NewTeamHandler(teamService service.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	team, err := h.teamService.Get(ctx, teamName)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "team not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	input := service.CreateTeamInput{
		Name:    req.Name,
		Members: make([]service.CreateTeamMemberInput, 0, len(req.Members)),
	}

	for _, member := range req.Members {
		input.Members = append(input.Members, service.CreateTeamMemberInput{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	team, err := h.teamService.Create(ctx, input)
	if err != nil {
		if errors.Is(err, domain.ErrTeamExists) {
			writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	resp := struct {
		Team *domain.Team `json:"team"`
	}{
		Team: team,
	}

	writeJSON(w, http.StatusCreated, resp)
}
