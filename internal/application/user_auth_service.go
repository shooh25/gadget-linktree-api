package application

import (
	"gadget-linktree-api/internal/domain/model/user"
	"gadget-linktree-api/internal/domain/repository"
	"log/slog"
)

type UserAuthService struct {
	repo     repository.UserRepository
	tokenGen repository.TokenGenerator
}

func NewUserAuthService(repo repository.UserRepository, tokenGen repository.TokenGenerator) *UserAuthService {
	return &UserAuthService{
		repo:     repo,
		tokenGen: tokenGen,
	}
}

type AuthDto struct {
	GoogleId    string
	DisplayName string
	AvatarImage string
}

func (s *UserAuthService) Authenticate(dto AuthDto) (string, error) {
	// ユーザーをGoogleIDで検索
	foundUser, err := s.repo.FindByGoogleId(dto.GoogleId)
	if err != nil {
		slog.Error("failed to find user", "err", err, "googleId", dto.GoogleId)
		return "", err
	}

	var u *user.User
	if foundUser == nil {
		// 存在しなければ新規ユーザー情報構築
		slog.Info("user not found, creating new user", "googleId", dto.GoogleId)
		userId := user.GenerateUserId()
		avatarImageUrl := user.NewAvatarImageURL(dto.AvatarImage)
		profile := user.NewProfile(dto.DisplayName, avatarImageUrl, user.NewBio(""), []user.SocialLink{})
		u = user.NewUser(userId, dto.GoogleId, profile)

		// ユーザー保存
		if err := s.repo.Save(*u); err != nil {
			slog.Error("failed to save new user", "err", err, "googleId", dto.GoogleId)
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

