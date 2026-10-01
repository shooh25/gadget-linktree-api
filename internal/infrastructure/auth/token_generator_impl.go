package auth

import (
	"gadget-linktree-api/internal/domain/model/user"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type TokenGeneratorImpl struct {
	secret []byte
}

func NewTokenGeneratorImpl(secret string) *TokenGeneratorImpl {
	return &TokenGeneratorImpl{
		secret: []byte(secret),
	}
}

func (g *TokenGeneratorImpl) Generate(u *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.UserId().Value(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
		"iat": time.Now().Unix(),
	})

	return token.SignedString(g.secret)
}
