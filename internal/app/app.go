package app

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"messenger/internal/app/wsserver"
	"messenger/internal/storage"
	"messenger/internal/storage/postgres"
)

type IApp interface {
	Start()
	Stop(ctx context.Context)
}

type App struct {
	server  wsserver.WSServer
	storage storage.Storage
}

func New(log *slog.Logger, cfgServer wsserver.ServerConfig, cfgStorage postgres.ConfigPostgres) IApp {
	getStorage := postgres.NewRepository(cfgStorage)
	err := getStorage.Connect()
	if err != nil {
		panic(err)
	}

	wsServer := wsserver.NewWsServer(fmt.Sprintf("%s:%d", cfgServer.Addr, cfgServer.Port), log)
	return &App{
		server:  wsServer,
		storage: getStorage,
	}
}

func (a *App) Start() {
	go a.server.MustStart()
}

func (a *App) Stop(ctx context.Context) {
	a.storage.MustClose()
	a.server.MustStop(ctx)
}
