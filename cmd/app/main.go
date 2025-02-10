package main

import (
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"messenger/internal/app"
	"messenger/internal/app/wsserver"
	"messenger/internal/config"
	"messenger/internal/storage"
	"messenger/internal/storage/postgres"
	"messenger/internal/storage/redis"
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

	pgClient := setupPostgres("./config/postgres.yaml")
	redisClient := setupRedis("./config/redis.yaml")
	log := setupLogger(cfg.Env)

	server := setupServer(log, "./config/wsserver.yaml")

	log.Info("starting application")

	application := app.New(log, server, pgClient, redisClient)

	application.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Info("stopping application", slog.String("signal", sign.String()))
	application.Stop(context.Background())
	log.Info("application stopped")
}

func setupServer(log *slog.Logger, configPath string) wsserver.WSServer {
	serverConfig := config.MustConfig[wsserver.Config](configPath)
	server := wsserver.New(log, serverConfig)
	return server
}

func setupRedis(configPath string) storage.Storage {
	redisCfg := config.MustConfig[redis.Config](configPath)
	redisCfg.Password = os.Getenv("REDIS_PASSWORD")

	client := redis.New(redisCfg)
	return client
}

func setupPostgres(configPath string) storage.Storage {
	pgCfg := config.MustConfig[postgres.Config](configPath)
	pgCfg.Password = os.Getenv("DB_PASSWORD")

	pg := postgres.New(pgCfg)
	return pg
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
