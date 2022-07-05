package main

import (
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
	"testing"
)

func TestStartApplication(t *testing.T) {
	log.SetOutput(os.Stdout)
	sets := settings.NewSettings()
	server := internal.NewServer(sets)
	client := internal.NewHTTPClient(sets)
	go server.Run()
	go client.SetMetrics(internal.CollectMetrics())
	go client.SendMetrics()
}
