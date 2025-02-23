package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log/slog"
	"messenger/internal/app"
	"messenger/internal/app/wsserver"
	"messenger/internal/config"
	"messenger/internal/handler"
	"messenger/internal/services/chat"
	"messenger/internal/services/message"
	"messenger/internal/storages/postgres"
	redisrepo "messenger/internal/storages/redis"
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

	pgClient, redisClient, log, server := setupDependencies(cfg)
	defer pgClient.Close()
	defer redisClient.Close()

	application := app.New(log, server, pgClient, redisClient)
	if err := application.Start(); err != nil {
		log.Error("error starting application", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	log.Info("stopping application", slog.String("signal", sign.String()))
	if err := application.Stop(context.Background()); err != nil {
		log.Error("error stopping application", err)
		os.Exit(1)
	}
}

func setupDependencies(cfg config.Config) (*sqlx.DB, *redis.Client, *slog.Logger, wsserver.WSServer) {
	log := setupLogger(cfg.Env)

	pgClient, err := setupPostgres("./config/postgres.yaml")
	if err != nil {
		log.Error("failed to setup", err)
	}
	redisClient := setupRedis("./config/redis.yaml")

	messageRepo := postgres.NewMessageRepository(pgClient)
	messageCacheRepo := redisrepo.NewMessageRepository(redisClient)
	chatRepo := postgres.NewChatRepository(pgClient)
	chatCacheRepo := redisrepo.NewChatRepository(redisClient)

	server := setupServer(log, messageRepo, messageCacheRepo,
		chatRepo, chatCacheRepo, "./config/wsserver.yaml")

	return pgClient, redisClient, log, server
}

func setupServer(log *slog.Logger,
	messageRepository message.Repository, messageCacheRepository message.CacheRepository,
	chatRepository chat.Repository, chatCacheRepository chat.CacheRepository,
	configPath string) wsserver.WSServer {
	serverConfig := config.MustConfig[wsserver.Config](configPath)

	messageService := message.NewMessageService(log, messageCacheRepository, messageRepository)
	chatService := chat.NewChat(log, chatRepository, chatCacheRepository)
	messengerHandler := handler.NewHandler(log, messageService, chatService)

	server := wsserver.New(log, messengerHandler, serverConfig)
	return server
}

func setupRedis(configPath string) *redis.Client {
	redisCfg := config.MustConfig[redisrepo.Config](configPath)
	redisCfg.Password = os.Getenv("REDIS_PASSWORD")

	db := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		MaxRetries:   redisCfg.MaxRetries,
		DialTimeout:  redisCfg.DialTimeout,
		ReadTimeout:  redisCfg.Timeout,
		WriteTimeout: redisCfg.Timeout,
	})
	return db
}

func setupPostgres(configPath string) (*sqlx.DB, error) {
	pgCfg := config.MustConfig[postgres.Config](configPath)
	pgCfg.Password = os.Getenv("DB_PASSWORD")

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
			pgCfg.Host, pgCfg.Port, pgCfg.DBName, pgCfg.User, pgCfg.Password, pgCfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres: %w", err)
	}
	return db, nil
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
