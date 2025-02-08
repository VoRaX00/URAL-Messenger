package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"messenger/internal/app/wsserver"
	"messenger/internal/config"
	"os"
	"os/signal"
	"syscall"
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

	configPath := config.FetchConfigPath()
	cfg := config.MustConfig[config.Config](configPath)
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	log := setupLogger(cfg.Env)

	log.Info("starting application")
	wsSrv := wsserver.NewWsServer(fmt.Sprintf("localhost:%d", cfg.Server.Port), log)
	go wsSrv.MustStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Info("stopping application", slog.String("signal", sign.String()))

	wsSrv.Stop(context.Background())
	log.Info("application stopped")

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
