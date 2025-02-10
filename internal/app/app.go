package app

import (
	"context"
	_ "github.com/lib/pq"
	"log/slog"
	"messenger/internal/app/wsserver"
	"messenger/internal/storage"
)

type IApp interface {
	Start()
	Stop(ctx context.Context)
}

type App struct {
	server   wsserver.WSServer
	storages []storage.Storage
}

func New(log *slog.Logger, server wsserver.WSServer, storages ...storage.Storage) IApp {
	return &App{
		server:   server,
		storages: storages,
	}
}

func (a *App) Start() {
	for i := range a.storages {
		a.storages[i].MustConnect()
	}
	go a.server.MustStart()
}

func (a *App) Stop(ctx context.Context) {
	a.server.MustStop(ctx)
	for i := range a.storages {
		a.storages[i].MustClose()
	}
}
