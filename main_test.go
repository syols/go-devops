package main

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"log"
	"os"
	"testing"
)

func TestStartApplication(t *testing.T) {
	log.SetOutput(os.Stdout)
	settings := config.NewConfig()
	server := app.NewServer(settings)
	client := app.NewHTTPClient(settings)
	go server.Run()
	go client.SetMetrics(app.CollectMetrics())
	go client.SendMetrics()
}
