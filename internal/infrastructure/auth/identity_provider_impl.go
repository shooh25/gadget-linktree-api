package auth

import (
	"context"
	"gadget-linktree-api/internal/domain/service"
	"log/slog"

	"google.golang.org/api/idtoken"
)

type IdentityProviderImpl struct {
	clientID string
}

func NewIdentityProviderImpl(clientID string) *IdentityProviderImpl {
	return &IdentityProviderImpl{clientID: clientID}
}

func (g *IdentityProviderImpl) VerifyToken(idToken string) (*service.ExternalUserInfo, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, g.clientID)
	if err != nil {
		slog.Error("invalid id token", "err", err)
		return nil, err
	}

	return &service.ExternalUserInfo{
		GoogleId:    payload.Subject,
		DisplayName: payload.Claims["name"].(string),
		AvatarImage: payload.Claims["picture"].(string),
	}, nil
}