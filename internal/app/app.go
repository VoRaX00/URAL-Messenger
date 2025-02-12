package app

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"messenger/internal/app/wsserver"
)

type IApp interface {
	Start() error
	Stop(ctx context.Context) error
}

type App struct {
	server      wsserver.WSServer
	pgClient    *sqlx.DB
	redisClient *redis.Client
	log         *slog.Logger
}

func New(log *slog.Logger, server wsserver.WSServer, pgClient *sqlx.DB, redisClient *redis.Client) IApp {
	return &App{
		server:      server,
		pgClient:    pgClient,
		redisClient: redisClient,
		log:         log,
	}
}

func (a *App) Start() error {
	a.log.Info("starting application")
	if err := a.server.Start(); err != nil {
		a.log.Error("failed to start server", slog.String("error", err.Error()))
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("stopping application")

	// Остановка сервера
	if err := a.server.Stop(ctx); err != nil {
		a.log.Error("failed to stop server", slog.String("error", err.Error()))
		return fmt.Errorf("failed to stop server: %w", err)
	}

	// Закрытие соединений с базой данных и Redis
	if err := a.pgClient.Close(); err != nil {
		a.log.Error("failed to close PostgreSQL connection", slog.String("error", err.Error()))
	}

	if err := a.redisClient.Close(); err != nil {
		a.log.Error("failed to close Redis connection", slog.String("error", err.Error()))
	}

	a.log.Info("application stopped")
	return nil
}
