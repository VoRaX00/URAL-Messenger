package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"messenger/internal/config"
	"messenger/internal/wsserver"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	conf := config.MustConfig()
	wsSrv := wsserver.NewWsServer(fmt.Sprintf("localhost:%d", conf.Server.Port))
	if err := wsSrv.Start(); err != nil {
		panic("Error starting server")
	}
}
