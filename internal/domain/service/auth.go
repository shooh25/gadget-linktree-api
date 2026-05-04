package service

import "gadget-linktree-api/internal/domain/model/user"

// アプリ内部用のtoken生成
type TokenGenerator interface {
	Generate(u *user.User) (string, error)
}

// tokenの検証
type TokenVerifier interface {
	Verify(token string) (string, error)
}

// Google IdToken検証
type ExternalUserInfo struct {
	GoogleId    string
	DisplayName string
	AvatarImage string
}

type IdentityProvider interface {
	VerifyToken(idToken string) (*ExternalUserInfo, error)
}
