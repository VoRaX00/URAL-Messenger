package app

import (
	"context"
	_ "github.com/lib/pq"
	"log/slog"
	"messenger/internal/app/wsserver"
	"messenger/internal/storage"
	"messenger/internal/storage/postgres"
	"messenger/internal/storage/redis"
)

type IApp interface {
	Start()
	Stop(ctx context.Context)
}

type App struct {
	server  wsserver.WSServer
	storage []storage.Storage
}

func New(log *slog.Logger, cfgServer wsserver.Config, cfgPg postgres.Config, cfgRedis redis.Config) IApp {
	pgClient := postgres.New(cfgPg)
	pgClient.MustConnect()

	redisClient := redis.New(cfgRedis)
	redisClient.MustConnect()

	wsServer := wsserver.New(cfgServer.Addr, log)
	return &App{
		server: wsServer,
		storage: []storage.Storage{
			pgClient,
			redisClient,
		},
	}
}

func (a *App) Start() {
	go a.server.MustStart()
}

func (a *App) Stop(ctx context.Context) {
	for i := range a.storage {
		a.storage[i].MustClose()
	}
	a.server.MustStop(ctx)
}
