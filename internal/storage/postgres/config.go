package postgres

type Config struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	User     string `yaml:"username" env-required:"true"`
	Password string `yaml:"-"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}
