package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/CodebyTecs/pr-assign-service/internal/api/handlers"
	"github.com/CodebyTecs/pr-assign-service/internal/config"
	"github.com/CodebyTecs/pr-assign-service/internal/repository"
	"github.com/CodebyTecs/pr-assign-service/internal/service"
	"github.com/joho/godotenv"
)

type Env struct {
	Config   *config.Config
	Postgres *sql.DB
}

func New() (*Env, error) {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	postgre, err := provideDB(
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		),
	)
	if err != nil {
		return nil, err
	}

	return &Env{
		Config:   cfg,
		Postgres: postgre,
	}, nil
}

func (e *Env) Run() error {
	userRepo := repository.NewUserRepository(e.Postgres)
	teamRepo := repository.NewTeamRepository(e.Postgres)
	prRepo := repository.NewPRRepository(e.Postgres)

	userSvc := service.NewUserService(userRepo)
	teamSvc := service.NewTeamService(userRepo, teamRepo)
	prSvc := service.NewPRService(prRepo, userRepo, teamRepo)
	statsSvc := service.NewStatsService(prRepo)

	userHandler := handlers.NewUserHandler(userSvc, prSvc)
	teamHandler := handlers.NewTeamHandler(teamSvc)
	prHandler := handlers.NewPRHandler(prSvc)
	statsHandler := handlers.NewStatsHandler(statsSvc)

	router := http.NewServeMux()
	router.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	router.HandleFunc("/users/getReview", userHandler.GetReview)
	router.HandleFunc("/team/add", teamHandler.Add)
	router.HandleFunc("/team/get", teamHandler.Get)
	router.HandleFunc("/pullRequest/create", prHandler.Create)
	router.HandleFunc("/pullRequest/reassign", prHandler.Reassign)
	router.HandleFunc("/pullRequest/merge", prHandler.Merge)
	router.HandleFunc("/stats", statsHandler.Get)

	addr := fmt.Sprintf("%s:%s", e.Config.HTTPServer.Address, e.Config.HTTPServer.Port)

	fmt.Println("Server listening on", addr)
	return http.ListenAndServe(addr, router)
}

func provideDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open connection: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("can't ping database: %w", err)
	}
	return db, nil
}

func (e *Env) Close() error {
	var firstErr error
	if e.Postgres != nil {
		if err := e.Postgres.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
