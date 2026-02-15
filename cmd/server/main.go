package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"ai-companion-be/internal/ai"
	"ai-companion-be/internal/config"
	"ai-companion-be/internal/database"
	"ai-companion-be/internal/handler"
	"ai-companion-be/internal/repository"
	"ai-companion-be/internal/router"
	"ai-companion-be/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load .env file if present (no error if missing, e.g. in production).
	_ = godotenv.Load()

	cfg := config.Load()

	// Database.
	ctx := context.Background()
	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	slog.Info("connected to database")

	// Run migrations.
	if err := database.RunMigrations(ctx, pool, "migrations"); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("migrations completed")

	// Repositories.
	userRepo := repository.NewUserRepository(pool)
	companionRepo := repository.NewCompanionRepository(pool)
	storyRepo := repository.NewStoryRepository(pool)
	messageRepo := repository.NewMessageRepository(pool)
	relationshipRepo := repository.NewRelationshipRepository(pool)
	memoryRepo := repository.NewMemoryRepository(pool)
	insightsRepo := repository.NewInsightsRepository(pool)

	// AI client.
	aiClient := ai.NewClient(cfg.OpenAI)

	// Services.
	authSvc := service.NewAuthService(userRepo, cfg.JWT)
	companionSvc := service.NewCompanionService(companionRepo)
	storySvc := service.NewStoryService(storyRepo, relationshipRepo, insightsRepo)
	messageSvc := service.NewMessageService(messageRepo, relationshipRepo, companionRepo, aiClient, insightsRepo)
	relationshipSvc := service.NewRelationshipService(relationshipRepo)
	memorySvc := service.NewMemoryService(memoryRepo)
	insightsSvc := service.NewInsightsService(insightsRepo, relationshipRepo)

	// Handlers.
	authH := handler.NewAuthHandler(authSvc)
	companionH := handler.NewCompanionHandler(companionSvc)
	storyH := handler.NewStoryHandler(storySvc)
	messageH := handler.NewMessageHandler(messageSvc)
	relationshipH := handler.NewRelationshipHandler(relationshipSvc)
	memoryH := handler.NewMemoryHandler(memorySvc)
	insightsH := handler.NewInsightsHandler(insightsSvc)

	// Router.
	r := router.New(cfg, authH, companionH, storyH, messageH, relationshipH, memoryH, insightsH)

	// Server.
	srv := &http.Server{
		Addr:         cfg.Server.Addr(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Graceful shutdown.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "addr", cfg.Server.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-done
	slog.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("server stopped")
}
