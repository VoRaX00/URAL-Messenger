package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"messenger/internal/config"
	"messenger/internal/wsserver"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	cfg := config.MustConfig()
	log := setupLogger(cfg.Env)

	log.Info("starting server")
	wsSrv := wsserver.NewWsServer(fmt.Sprintf("localhost:%d", cfg.Server.Port), log)
	if err := wsSrv.Start(); err != nil {
		log.Error("error starting server")
		panic("Error starting server")
	}
	log.Info("server started")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return logger
}
