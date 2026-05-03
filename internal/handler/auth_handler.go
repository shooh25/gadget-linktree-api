package handler

import (
	"encoding/json"
	"gadget-linktree-api/internal/application"
	"gadget-linktree-api/internal/domain/service"
	"log/slog"
	"net/http"
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

type LoginRequest struct {	
	IdToken string `json:"idToken"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("failed to decode login request", "err", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IdToken == "" {
		http.Error(w, "idToken is required", http.StatusBadRequest)
		return
	}

	// Google Idトークンの検証
	externalUser, err := h.identityProvider.VerifyToken(req.IdToken)
	if err != nil {
		slog.Error("failed to verify google token", "err", err)
		http.Error(w, "Invalid Id token", http.StatusUnauthorized)
		return
	}

	// アプリケーション層の認証処理
	token, err := h.userAuthService.Authenticate(*externalUser)
	if err != nil {
		slog.Error("failed to authenticate user", "err", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	// 結果の返却
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(LoginResponse{Token: token}); err != nil {
		slog.Error("failed to encode login response", "err", err)
	}
}
