package application

import (
	"errors"
	"gadget-linktree-api/internal/domain/model/user"
	"gadget-linktree-api/internal/domain/repository"
	"log/slog"
)

type UserRegisterService struct {
	repo repository.UserRepository
}

type UserRegisterDto struct {
	GoogleId    string
	DisplayName string
	AvatarImage string
}

func NewUserRegisterService(repo repository.UserRepository) *UserRegisterService {
	return &UserRegisterService{
		repo: repo,
	}
}

func (s *UserRegisterService) UserRegister(dto UserRegisterDto) error {
	// ユーザーが既に存在するかチェック
	exists, err := s.repo.ExistsByGoogleId(dto.GoogleId)
	if err != nil {
		slog.Error("failed to check if user exists", "err", err, "googleId", dto.GoogleId)
		return err
	}
	if exists {
		slog.Error("user already exists", "googleId", dto.GoogleId)
		return errors.New("user already exists")
	}

	// ユーザー情報構築
	userId := user.GenerateUserId()
	avatarImageUrl := user.NewAvatarImageURL(dto.AvatarImage)
	profile := user.NewProfile(dto.DisplayName, avatarImageUrl, user.NewBio(""), []user.SocialLink{})
	newUser := user.NewUser(userId, dto.GoogleId, profile)

	// ユーザー保存
	err = s.repo.Save(*newUser)

	if err != nil {
		slog.Error("failed to save user", "err", err, "googleId", dto.GoogleId)
		return err
	}

	slog.Info("user saved successfully", "googleId", dto.GoogleId)
	return nil
}
