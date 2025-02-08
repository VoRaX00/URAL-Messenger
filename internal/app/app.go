package app

import (
	"log/slog"
	"messenger/internal/app/wsserver"
	"messenger/internal/storage"
	"messenger/internal/storage/postgres"
)

type App struct {
	Server  *wsserver.WSServer
	Storage storage.Storage
}

func New(log *slog.Logger, cfg postgres.ConfigPostgres) *App {
	strg := postgres.NewRepository(cfg)
	err := strg.Connect()
	if err != nil {
		panic(err)
	}

	return &App{
		Server:  nil,
		Storage: strg,
	}
}
