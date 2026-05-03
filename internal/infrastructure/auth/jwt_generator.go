package auth

import (
	"gadget-linktree-api/internal/domain/model/user"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type JWTokenGenerator struct {
	secret []byte
}

func NewJWTokenGenerator(secret string) *JWTokenGenerator {
	return &JWTokenGenerator{
		secret: []byte(secret),
	}
}

func (g *JWTokenGenerator) Generate(u *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.UserId().Value(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
		"iat": time.Now().Unix(),
	})

	return token.SignedString(g.secret)
}
