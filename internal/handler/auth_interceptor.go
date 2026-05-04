package handler

import (
	"context"
	"gadget-linktree-api/internal/domain/service"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const userIDKey contextKey = "user_id"

// 認証チェックを行うミドルウェア
func AuthInterceptor(verifier service.TokenVerifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// 認証をスキップする公開メソッドの判定
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// メタデータからトークンの取得
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is missing")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization header is missing")
		}

		// Bearer トークンの抽出
		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		if token == authHeaders[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
		}

		// トークンの検証
		userId, err := verifier.Verify(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		// ContextにUserIDをセットして後続の処理を実行
		newCtx := context.WithValue(ctx, userIDKey, userId)
		return handler(newCtx, req)
	}
}

func isPublicMethod(fullMethod string) bool {
	publicMethods := []string{
		"/auth.v1.AuthService/Login",
	}
	for _, m := range publicMethods {
		if m == fullMethod {
			return true
		}
	}
	return false
}

// GetUserID はContextからUserIDを取得するためのヘルパー関数です
func GetUserID(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIDKey).(string)
	return userId, ok
}
