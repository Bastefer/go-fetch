package main

import (
	"context"
	"log/slog"
	"os/signal"
	"service-parser/internal/app/external"
	"service-parser/internal/app/handler"
	"service-parser/internal/app/repository"
	"service-parser/internal/app/service"
	"service-parser/internal/config"
	"service-parser/internal/db/connection"
	"service-parser/internal/logger/sl"
	"syscall"
	"time"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.AppEnv)
	log.Info("starting service", slog.String("env", cfg.AppEnv))
	log.Info("debug are enabled")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	storage, err := connection.NewStorage(ctx, cfg.PostgresDSN())
	if err != nil {
		log.Error("failed to init storage", slog.Any(sl.Error, err))
		os.Exit(1)
	}
	defer storage.Pool().Close()

	urls := []string{cfg.Source1, cfg.Source2, cfg.Source3}
	httpClient := external.NewHTTPClient(&http.Client{})
	productRepository := repository.NewProductRepository()
	clientRepository := repository.NewClientRepository()
	brandRepository := repository.NewBrandRepository()
	categoryRepository := repository.NewCategoryRepository()
	clientProductRepository := repository.NewClientProductRepository()
	taskRepository := repository.NewTaskRepository()

	fetchService := service.NewFetchService(log, storage.Pool(), httpClient, productRepository, clientRepository, brandRepository, categoryRepository, clientProductRepository, taskRepository, urls, cfg.ClientsSource)
	statisticService := service.NewStatisticService(log, storage.Pool(), productRepository, clientRepository, brandRepository, categoryRepository)
	handler := handler.New(fetchService, statisticService)

	router := gin.Default()

	router.POST("/download", handler.Download)
	router.GET("/stats", handler.Stats)
	srv := &http.Server{
		Addr:    cfg.Address(),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.Any(sl.Error, err))
		}
	}()

	log.Info("starting server", slog.String("address", cfg.Address()))

	<-ctx.Done()

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown server", slog.Any(sl.Error, err))
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
