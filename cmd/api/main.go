package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gadget-linktree-api/gen/auth/v1/authv1connect"
	"gadget-linktree-api/internal/application"
	"gadget-linktree-api/internal/handler"
	"gadget-linktree-api/internal/infrastructure/auth"
	"gadget-linktree-api/internal/infrastructure/persistence"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// 設定の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	appSecret := os.Getenv("APP_SECRET")

	// 依存関係の初期化
	userRepo := persistence.NewUserRepositoryImpl()
	tokenGen := auth.NewTokenGeneratorImpl(appSecret)
	identityProvider := auth.NewIdentityProviderImpl(googleClientID)
	
	userAuthService := application.NewUserAuthService(userRepo, tokenGen)
	
	authHandler := handler.NewAuthHandler(userAuthService, identityProvider)

	// ルーターの設定
	r := chi.NewRouter()

	// CORS設定
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // フロントエンドのURL
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Accept-Encoding", "Content-Type", "Authorization", "Connect-Protocol-Version", "Connect-Timeout-Ms"},
		AllowCredentials: true,
	})
	r.Use(c.Handler)

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Connect RPC ハンドラの登録
	path, connectHandler := authv1connect.NewAuthServiceHandler(authHandler)
	r.Handle(path+"*", connectHandler)

	// 動作確認用エンドポイント
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// サーバーの起動
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Starting server on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
