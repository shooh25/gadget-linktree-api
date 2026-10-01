package handler

import (
	"context"
	"gadget-linktree-api/internal/application"
	"gadget-linktree-api/internal/domain/service"
	authv1 "gadget-linktree-api/gen/auth/v1"
	"log/slog"

	"connectrpc.com/connect"
)

type AuthHandler struct {
	userAuthService  *application.UserAuthService
	identityProvider service.IdentityProvider
}

func NewAuthHandler(authService *application.UserAuthService, identityProvider service.IdentityProvider) *AuthHandler {
	return &AuthHandler{
		userAuthService:  authService,
		identityProvider: identityProvider,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *connect.Request[authv1.LoginRequest]) (*connect.Response[authv1.LoginResponse], error) {
	// バリデーション
	if req.Msg.GetIdToken() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, nil)
	}

	// Google IdTokenの検証
	externalUser, err := h.identityProvider.VerifyToken(req.Msg.GetIdToken())
	if err != nil {
		slog.Error("failed to verify google token", "err", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	// アプリケーション層の認証処理
	token, err := h.userAuthService.Authenticate(*externalUser)
	if err != nil {
		slog.Error("failed to authenticate user", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// 結果の返却
	return connect.NewResponse(&authv1.LoginResponse{
		AccessToken: token,
	}), nil
}
