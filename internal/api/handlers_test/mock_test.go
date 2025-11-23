package handlers_test

import (
	"context"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/service"
)

type mockUserService struct{}
type mockTeamService struct{}
type mockPRService struct{}
type mockStatsService struct{}

func (m *mockUserService) UpdateActivity(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	return &domain.User{
		ID:       userID,
		Username: "Test",
		TeamName: "Test",
		IsActive: isActive,
	}, nil
}

func (m *mockUserService) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	return []domain.PullRequestShort{}, nil
}

func (m *mockTeamService) Create(ctx context.Context, input service.CreateTeamInput) (*domain.Team, error) {
	return &domain.Team{
		Name:    input.Name,
		Members: []domain.TeamMember{},
	}, nil
}

func (m *mockTeamService) Get(ctx context.Context, teamName string) (*domain.Team, error) {
	return &domain.Team{
		Name:    teamName,
		Members: []domain.TeamMember{},
	}, nil
}

func (m *mockPRService) Create(ctx context.Context, input service.CreatePRInput) (*domain.PullRequest, error) {
	return &domain.PullRequest{
		ID:       input.ID,
		Name:     input.Name,
		AuthorID: input.Author,
		Status:   domain.PRStatusOpen,
	}, nil
}

func (m *mockPRService) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	return &domain.PullRequest{
		ID:     id,
		Status: domain.PRStatusMerged,
	}, nil
}

func (m *mockPRService) Reassign(ctx context.Context, input service.ReassignReviewerInput) (*domain.PullRequest, string, error) {
	return &domain.PullRequest{
		ID:        input.PullRequestID,
		Status:    domain.PRStatusOpen,
		Reviewers: []string{"id-new"},
	}, "id-new", nil
}

func (m *mockPRService) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequestShort, error) {
	return []domain.PullRequestShort{}, nil
}

func (m *mockStatsService) Get(ctx context.Context) (*domain.Stats, error) {
	return &domain.Stats{
		TotalPR:        0,
		OpenPR:         0,
		MergedPR:       0,
		ReviewsPerUser: []domain.UserReviewStat{},
	}, nil
}
