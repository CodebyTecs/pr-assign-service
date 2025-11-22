package service

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/repository"
)

type CreatePRInput struct {
	ID     string
	Name   string
	Author string
}

type ReassignReviewerInput struct {
	PullRequestID string
	ReviewerID    string
}

type PRService interface {
	Create(ctx context.Context, input CreatePRInput) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	Reassign(ctx context.Context, input ReassignReviewerInput) (*domain.PullRequest, string, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error)
}

type prService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(prRepo repository.PRRepository, userRepo repository.UserRepository, teamRepo repository.TeamRepository) PRService {
	return &prService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *prService) Create(ctx context.Context, input CreatePRInput) (*domain.PullRequest, error) {
	_, err := s.prRepo.GetByID(ctx, input.ID)
	if err == nil {
		return nil, domain.ErrPRExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, input.Author)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	team, err := s.teamRepo.GetByName(ctx, author.TeamName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	users, err := s.userRepo.ListActiveByTeam(ctx, author.TeamName)
	if err != nil {
		return nil, err
	}

	candidates := make([]*domain.User, 0, len(team.Members))
	for _, u := range users {
		if u.ID == author.ID {
			continue
		}
		candidates = append(candidates, u)
	}

	selected := candidates
	if len(selected) > 2 {
		rand.Shuffle(len(selected), func(i, j int) {
			selected[i], selected[j] = selected[j], selected[i]
		})
		selected = selected[:2]
	}

	selectedReviewers := make([]string, 0, len(selected))
	for _, c := range selected {
		selectedReviewers = append(selectedReviewers, c.ID)
	}

	now := time.Now()
	pr := &domain.PullRequest{
		ID:        input.ID,
		Name:      input.Name,
		AuthorID:  input.Author,
		Status:    domain.PRStatusOpen,
		Reviewers: selectedReviewers,
		CreatedAt: &now,
	}

	if err := s.prRepo.Create(ctx, pr); err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prService) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	now := time.Now()
	pr.Status = domain.PRStatusMerged
	pr.MergedAt = &now

	err = s.prRepo.Update(ctx, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *prService) Reassign(ctx context.Context, input ReassignReviewerInput) (*domain.PullRequest, string, error) {
	pr, err := s.prRepo.GetByID(ctx, input.PullRequestID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", err
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.ErrPRMerged
	}

	foundIndex := -1
	for i, r := range pr.Reviewers {
		if r == input.ReviewerID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return nil, "", domain.ErrNotAssigned
	}

	oldReviewer, err := s.userRepo.GetByID(ctx, input.ReviewerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, "", domain.ErrNotFound
		}
		return nil, "", err
	}

	candidates, err := s.userRepo.ListActiveByTeam(ctx, oldReviewer.TeamName)
	if err != nil {
		return nil, "", err
	}

	filtered := make([]*domain.User, 0, len(candidates))
	for _, u := range candidates {
		if u.ID == input.ReviewerID {
			continue
		}
		if u.ID == pr.AuthorID {
			continue
		}
		alreadyReviewed := false
		for _, r := range pr.Reviewers {
			if r == u.ID {
				alreadyReviewed = true
				break
			}
		}
		if alreadyReviewed {
			continue
		}
		filtered = append(filtered, u)
	}

	if len(filtered) == 0 {
		return nil, "", domain.ErrNoCandidate
	}

	rand.Shuffle(len(filtered), func(i, j int) {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	})
	newReviewer := filtered[0]

	pr.Reviewers[foundIndex] = newReviewer.ID

	err = s.prRepo.Update(ctx, pr)
	if err != nil {
		return nil, "", err
	}

	return pr, newReviewer.ID, nil
}

func (s *prService) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	_, err := s.userRepo.GetByID(ctx, reviewerID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	prs, err := s.prRepo.ListByReviewer(ctx, reviewerID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
