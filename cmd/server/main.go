package main

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	if server, err := app.NewServer(config.NewConfig()); err == nil {
		server.Run()
	}
}
