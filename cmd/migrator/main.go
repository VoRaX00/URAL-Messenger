package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"messenger/internal/config"
	_ "messenger/migrations"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	path := config.FetchConfigPath()
	if path == "" {
		panic("path is empty")
	}
	cfg := config.MustConfig[Config](path)
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	db, err := sqlx.Open(cfg.DB.Driver, fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SSLMode))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if cfg.DB.IsDrop {
		if err = goose.DownTo(db.DB, cfg.MigrationsPath, 0); err != nil {
			panic(err)
		}
	}

	if err = goose.Up(db.DB, cfg.MigrationsPath); err != nil {
		panic(err)
	}
}
