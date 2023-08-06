package main

import (
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
	"libraryService/internal/config"
	"libraryService/internal/http-server/handlers/receive"
	"libraryService/internal/http-server/handlers/save"
	"libraryService/internal/storage/sqlite"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Debug("Debug messages are enabled")

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", err)
		os.Exit(1)
	}

	log.Debug("Storage is initialized")

	router := chi.NewRouter()

	router.Get("/books", receive.New(log, db))
	router.Post("/add", save.New(log, db))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Error("Server stopped")
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
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
