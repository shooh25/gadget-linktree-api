package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type TokenVerifierImpl struct {
	secret []byte
}

func NewTokenVerifierImpl(secret string) *TokenVerifierImpl {
	return &TokenVerifierImpl{
		secret: []byte(secret),
	}
}

func (v *TokenVerifierImpl) Verify(tokenString string) (string, error) {
	// トークンのパースと検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名方式のチェック
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return v.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", errors.New("token expired")
		}
		return "", errors.New("invalid token")
	}

	// 有効なトークンからUserID（sub）を取り出す
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return "", errors.New("invalid token payload")
		}
		return sub, nil
	}

	return "", errors.New("invalid token")
}
