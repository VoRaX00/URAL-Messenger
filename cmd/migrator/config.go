package main

type Config struct {
	DB             DBConfig `yaml:"db"`
	MigrationsPath string   `yaml:"migrations_path"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"-"`
	Username string `yaml:"username"`
	SSLMode  string `yaml:"ssl_mode"`
	IsDrop   bool   `yaml:"is_drop"`
}
