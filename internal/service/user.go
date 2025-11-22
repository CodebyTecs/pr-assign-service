package service

import (
	"context"
	"errors"

	"github.com/CodebyTecs/pr-assign-service/internal/domain"
	"github.com/CodebyTecs/pr-assign-service/internal/repository"
)

type UserService interface {
	UpdateActivity(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		repo: repository,
	}
}

func (s *userService) UpdateActivity(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	user.IsActive = isActive

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
