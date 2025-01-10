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
	DB     DBConfig     `yaml:"db"`
}

type ServerConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	User     string `yaml:"username" env-required:"true"`
	Password string `yaml:"-"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

func MustConfig() Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config file path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found")
	}

	var config Config
	err := cleanenv.ReadConfig(path, &config)
	if err != nil {
		panic(err)
	}

	config.DB.Password = os.Getenv("DB_PASSWORD")
	return config
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "", "path to config file")
	flag.Parse()
	return path
}
