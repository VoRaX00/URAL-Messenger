package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
}

func MustConfig[T any](path string) T {
	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

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
