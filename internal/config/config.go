package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"messenger/internal/app/wsserver"
	"messenger/internal/storage/postgres"
	"messenger/internal/storage/redis"
	"os"
)

type Config struct {
	Env         string          `yaml:"env" env-default:"local"`
	Server      wsserver.Config `yaml:"server"`
	PGConfig    postgres.Config `yaml:"postgres"`
	RedisConfig redis.Config    `yaml:"redis"`
}

func MustConfig[T any](path string) T {
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var config T
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func FetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	return path
}
