package handler

import (
	"context"
	"gadget-linktree-api/internal/application"
	"gadget-linktree-api/internal/domain/service"
	authv1 "gadget-linktree-api/proto/auth/v1"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	userAuthService  *application.UserAuthService
	identityProvider service.IdentityProvider
}

func NewAuthHandler(authService *application.UserAuthService, identityProvider service.IdentityProvider) *AuthHandler {
	return &AuthHandler{
		userAuthService:  authService,
		identityProvider: identityProvider,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	// バリデーション
	if req.GetIdToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "id_token is required")
	}

	// Google IdTokenの検証
	externalUser, err := h.identityProvider.VerifyToken(req.GetIdToken())
	if err != nil {
		slog.Error("failed to verify google token", "err", err)
		return nil, status.Error(codes.Unauthenticated, "invalid ID token")
	}

	// アプリケーション層の認証処理
	token, err := h.userAuthService.Authenticate(*externalUser)
	if err != nil {
		slog.Error("failed to authenticate user", "err", err)
		return nil, status.Error(codes.Internal, "authentication failed")
	}

	// 結果の返却
	return &authv1.LoginResponse{
		AccessToken: token,
	}, nil
}
