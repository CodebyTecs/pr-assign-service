package handlers

import (
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/service"
)

type StatsHandler struct {
	statsService service.StatsService
}

func NewStatsHandler(statsService service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := h.statsService.Get(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, stats)
}
