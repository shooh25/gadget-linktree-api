package application

import (
	"gadget-linktree-api/internal/domain/model/user"
	"gadget-linktree-api/internal/domain/repository"
	"log/slog"
)


type UserGetService struct {
	repo repository.UserRepository
}

func NewUserGetService(repo repository.UserRepository) *UserGetService {
	return &UserGetService{
		repo: repo,
	}
}

func (s *UserGetService) GetUser(userId user.UserId) (*user.User, error) {
	u, err := s.repo.FindByUserId(userId)
	if err != nil {
		slog.Error("failed to get user", "err", err, "userId", userId.Value())
		return nil, err
	}
	return u, nil
}