package service

import (
	"context"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/repository"
)

type StatsService interface {
	Get(ctx context.Context) (*domain.Stats, error)
}

type statsService struct {
	prRepo repository.PRRepository
}

func NewStatsService(prRepo repository.PRRepository) StatsService {
	return statsService{
		prRepo: prRepo,
	}
}

func (s statsService) Get(ctx context.Context) (*domain.Stats, error) {
	total, err := s.prRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}

	open, err := s.prRepo.CountByStatus(ctx, domain.PRStatusOpen)
	if err != nil {
		return nil, err
	}

	merged, err := s.prRepo.CountByStatus(ctx, domain.PRStatusMerged)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.prRepo.CountAssignmentsByReviewer(ctx)
	if err != nil {
		return nil, err
	}

	stats := &domain.Stats{
		TotalPR:        total,
		OpenPR:         open,
		MergedPR:       merged,
		ReviewsPerUser: reviewers,
	}

	return stats, nil
}
