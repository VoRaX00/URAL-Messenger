package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env    string       `yaml:"env" env-default:"local"`
	Server ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type RedisConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
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
