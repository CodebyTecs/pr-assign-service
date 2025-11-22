package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
)

type PRRepository interface {
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	Create(ctx context.Context, pr *domain.PullRequest) error
	Update(ctx context.Context, pr *domain.PullRequest) error
	ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)

	CountAll(ctx context.Context) (int, error)
	CountByStatus(ctx context.Context, status domain.PRStatus) (int, error)
	CountAssignmentsByReviewer(ctx context.Context) ([]domain.UserReviewStat, error)
}

type prRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) PRRepository {
	return &prRepository{db: db}
}

func (r *prRepository) GetByID(ctx context.Context, id string) (*domain.PullRequest, error) {
	const q = `
	SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at
	FROM pull_requests
	WHERE pull_request_id = $1
	`

	var pr domain.PullRequest

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		pq.Array(&pr.Reviewers),
		&pr.CreatedAt,
		&pr.MergedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &pr, nil
}

func (r *prRepository) Create(ctx context.Context, pr *domain.PullRequest) error {
	const q = `
	INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, q,
		pr.ID,
		pr.Name,
		pr.AuthorID,
		pr.Status,
		pq.Array(pr.Reviewers),
		pr.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *prRepository) Update(ctx context.Context, pr *domain.PullRequest) error {
	const q = `
	UPDATE pull_requests
	SET status = $2, assigned_reviewers = $3, merged_at = $4
	WHERE pull_request_id = $1
	`

	res, err := r.db.ExecContext(ctx, q,
		pr.ID,
		pr.Status,
		pq.Array(pr.Reviewers),
		pr.MergedAt,
	)
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

func (r *prRepository) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	const q = `
	SELECT pull_request_id, pull_request_name, author_id, status
	FROM pull_requests
	WHERE $1 = ANY(assigned_reviewers)
	`

	rows, err := r.db.QueryContext(ctx, q, reviewerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var result []domain.PullRequestShort

	for rows.Next() {
		item := domain.PullRequestShort{}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.AuthorID,
			&item.Status,
		); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []domain.PullRequestShort{}, nil
	}

	return result, nil
}

func (r *prRepository) CountAll(ctx context.Context) (int, error) {
	const q = `
	SELECT COUNT(*)
	FROM pull_requests
	`

	var count int
	err := r.db.QueryRowContext(ctx, q).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *prRepository) CountByStatus(ctx context.Context, status domain.PRStatus) (int, error) {
	const q = `
	SELECT COUNT(*)
	FROM pull_requests
	WHERE status = $1
	`

	var count int
	err := r.db.QueryRowContext(ctx, q, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *prRepository) CountAssignmentsByReviewer(ctx context.Context) ([]domain.UserReviewStat, error) {
	const q = `
	SELECT reviewer_id, COUNT(*)
	FROM pull_requests
	CROSS JOIN LATERAL unnest(assigned_reviewers) AS reviewer_id
	GROUP BY reviewer_id
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("failed to close rows:", err)
		}
	}()

	var result []domain.UserReviewStat

	for rows.Next() {
		var s domain.UserReviewStat
		if err := rows.Scan(&s.UserID, &s.ReviewsCount); err != nil {
			return nil, err
		}
		result = append(result, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []domain.UserReviewStat{}, nil
	}

	return result, nil
}
