package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Create(ctx context.Context, user *domain.User) error
	ListActiveByTeam(ctx context.Context, teamName string) ([]*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	const q = `
	SELECT user_id, username, team_name, is_active
	FROM users
	WHERE user_id = $1
	`

	var u domain.User

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&u.ID,
		&u.Username,
		&u.TeamName,
		&u.IsActive,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `
	UPDATE users
	SET username = $2, team_name = $3, is_active = $4
	WHERE user_id = $1
	`

	res, err := r.db.ExecContext(ctx, q, user.ID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
	INSERT INTO users (user_id, username, team_name, is_active)
	VALUES ($1, $2, $3, $4)
`
	_, err := r.db.ExecContext(ctx, q, user.ID, user.Username, user.TeamName, user.IsActive)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) ListActiveByTeam(ctx context.Context, teamName string) ([]*domain.User, error) {
	const q = `
	SELECT user_id, username, team_name, is_active
	FROM users
	WHERE team_name = $1
	`

	rows, err := r.db.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var result []*domain.User

	for rows.Next() {
		item := &domain.User{}
		if err := rows.Scan(&item.ID, &item.Username, &item.TeamName, &item.IsActive); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []*domain.User{}, nil
	}

	return result, nil
}
