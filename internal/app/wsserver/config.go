package wsserver

import "time"

type ServerConfig struct {
	Addr    string        `yaml:"addr"`
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}
