package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
)

type TeamRepository interface {
	GetByName(ctx context.Context, name string) (*domain.Team, error)
	Create(ctx context.Context, team *domain.Team) error
}

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	const q = `
	SELECT team_name
	FROM teams
	WHERE team_name = $1
	`

	var team domain.Team

	err := r.db.QueryRowContext(ctx, q, name).Scan(&team.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &team, nil
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {
	const q = `
	INSERT INTO teams (team_name)
	VALUES ($1)
	`

	_, err := r.db.ExecContext(ctx, q, team.Name)
	if err != nil {
		return err
	}

	return nil
}
