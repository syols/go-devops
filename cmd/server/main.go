package main

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	server := app.NewServer(config.NewConfig())
	server.Run()
}
