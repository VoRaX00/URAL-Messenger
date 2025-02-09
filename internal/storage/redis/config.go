package redis

import "time"

type Config struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"-"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}
