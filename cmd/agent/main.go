package main

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	settings := config.NewConfig()
	client := app.NewHTTPClient(settings)
	pollTicker := time.NewTicker(settings.Agent.PollInterval)
	reportTicker := time.NewTicker(settings.Agent.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			go client.CollectMetrics()
			go client.CollectAdditionalMetrics()
		case <-reportTicker.C:
			go client.SendMetrics()
		}
	}
}
