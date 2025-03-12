package app

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"messenger/internal/app/wsserver"
)

type IApp interface {
	Start() error
	Stop(ctx context.Context) error
}

type App struct {
	server wsserver.WSServer
	log    *slog.Logger
}

func New(log *slog.Logger, server wsserver.WSServer) IApp {
	return &App{
		server: server,
		log:    log,
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
	if err := a.server.Stop(ctx); err != nil {
		a.log.Error("failed to stop server", slog.String("error", err.Error()))
		return fmt.Errorf("failed to stop server: %w", err)
	}
	a.log.Info("application stopped")

	return nil
}
