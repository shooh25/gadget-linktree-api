package application

import (
	"gadget-linktree-api/internal/domain/model/user"
	"gadget-linktree-api/internal/domain/repository"
	"gadget-linktree-api/internal/domain/service"
	"log/slog"
)

type UserAuthService struct {
	repo     repository.UserRepository
	tokenGen service.TokenGenerator
}

func NewUserAuthService(repo repository.UserRepository, tokenGen service.TokenGenerator) *UserAuthService {
	return &UserAuthService{
		repo:     repo,
		tokenGen: tokenGen,
	}
}

func (s *UserAuthService) Authenticate(info service.ExternalUserInfo) (string, error) {
	// ユーザーをGoogleIDで検索
	foundUser, err := s.repo.FindByGoogleId(info.GoogleId)
	if err != nil {
		slog.Error("failed to find user", "err", err, "googleId", info.GoogleId)
		return "", err
	}

	var u *user.User
	if foundUser == nil {
		// 存在しなければ新規ユーザー情報構築
		slog.Info("user not found, creating new user", "googleId", info.GoogleId)
		userId := user.GenerateUserId()
		avatarImageUrl := user.NewAvatarImageURL(info.AvatarImage)
		profile := user.NewProfile(info.DisplayName, avatarImageUrl, user.NewBio(""), []user.SocialLink{})
		u = user.NewUser(userId, info.GoogleId, profile)

		// ユーザー保存
		if err := s.repo.Save(*u); err != nil {
			slog.Error("failed to save new user", "err", err, "googleId", info.GoogleId)
			return "", err
		}
	} else {
		u = foundUser
	}

	// トークン発行
	token, err := s.tokenGen.Generate(u)
	if err != nil {
		slog.Error("failed to generate token", "err", err, "userId", u.UserId().Value())
		return "", err
	}

	return token, nil
}

