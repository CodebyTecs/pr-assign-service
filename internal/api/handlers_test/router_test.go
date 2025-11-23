package handlers_test

import (
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/api/handlers"
	"github.com/CodebyTecs/pr-assign-service/internal/app"
)

func newTestRouter() http.Handler {
	userSvc := &mockUserService{}
	teamSvc := &mockTeamService{}
	prSvc := &mockPRService{}
	statsSvc := &mockStatsService{}

	userHandler := handlers.NewUserHandler(userSvc, prSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	prHandler := handlers.NewPRHandler(prSvc)
	statsHandler := handlers.NewStatsHandler(statsSvc)

	r := app.NewRouter(userHandler, teamHandler, prHandler, statsHandler)

	return r.Handler()
}
