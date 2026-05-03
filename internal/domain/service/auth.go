package service

import "gadget-linktree-api/internal/domain/model/user"

// アプリ内部用のtoken生成
type TokenGenerator interface {
	Generate(u *user.User) (string, error)
}

// Google Idトークン検証
type ExternalUserInfo struct {
	GoogleId    string
	DisplayName string
	AvatarImage string
}

type IdentityProvider interface {
	VerifyToken(idToken string) (*ExternalUserInfo, error)
}
