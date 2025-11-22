package service

import (
	"context"
	"errors"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/repository"
)

type CreateTeamInput struct {
	Name    string
	Members []CreateTeamMemberInput
}

type CreateTeamMemberInput struct {
	UserID   string
	Username string
	IsActive bool
}

type teamService struct {
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewTeamService(userRepo repository.UserRepository, teamRepo repository.TeamRepository) TeamService {
	return &teamService{
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

type TeamService interface {
	Create(ctx context.Context, input CreateTeamInput) (*domain.Team, error)
	Get(ctx context.Context, teamName string) (*domain.Team, error)
}

func (s *teamService) Create(ctx context.Context, input CreateTeamInput) (*domain.Team, error) {
	_, err := s.teamRepo.GetByName(ctx, input.Name)
	if err == nil {
		return nil, domain.ErrTeamExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	team := &domain.Team{
		Name:    input.Name,
		Members: make([]domain.TeamMember, 0, len(input.Members)),
	}

	err = s.teamRepo.Create(ctx, team)
	if err != nil {
		return nil, err
	}

	for _, member := range input.Members {
		user, err := s.userRepo.GetByID(ctx, member.UserID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}

		if user == nil {
			user = &domain.User{
				ID:       member.UserID,
				Username: member.Username,
				TeamName: input.Name,
				IsActive: member.IsActive,
			}
			err := s.userRepo.Create(ctx, user)
			if err != nil {
				return nil, err
			}
		} else {
			user.Username = member.Username
			user.TeamName = input.Name
			user.IsActive = member.IsActive
			err := s.userRepo.Update(ctx, user)
			if err != nil {
				return nil, err
			}
		}

		team.Members = append(team.Members, domain.TeamMember{
			UserID:   member.UserID,
			Username: member.Username,
			IsActive: member.IsActive,
		})
	}

	return team, nil
}

func (s *teamService) Get(ctx context.Context, teamName string) (*domain.Team, error) {
	team, err := s.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	users, err := s.userRepo.ListActiveByTeam(ctx, teamName) // надо реализовать
	if err != nil {
		return nil, err
	}

	team.Members = make([]domain.TeamMember, 0, len(users))
	for _, u := range users {
		team.Members = append(team.Members, domain.TeamMember{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	return team, nil
}
