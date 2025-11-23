package app

import (
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/api/handlers"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(userHandler *handlers.UserHandler, teamHandler *handlers.TeamHandler, prHandler *handlers.PRHandler, statsHandler *handlers.StatsHandler) *Router {
	mux := http.NewServeMux()

	mux.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	mux.HandleFunc("/users/getReview", userHandler.GetReview)

	mux.HandleFunc("/team/add", teamHandler.Add)
	mux.HandleFunc("/team/get", teamHandler.Get)

	mux.HandleFunc("/pullRequest/create", prHandler.Create)
	mux.HandleFunc("/pullRequest/merge", prHandler.Merge)
	mux.HandleFunc("/pullRequest/reassign", prHandler.Reassign)

	mux.HandleFunc("/stats", statsHandler.Get)

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
