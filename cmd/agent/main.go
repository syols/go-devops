package main

import (
	"github.com/syols/go-devops/internal"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	var settings internal.Settings
	settings.LoadSettings(internal.ConfigPath)

	newClient := internal.NewHTTPClient(settings)
	pollTicker := time.NewTicker(settings.Agent.PollInterval * time.Millisecond)
	reportTicker := time.NewTicker(settings.Agent.ReportInterval * time.Millisecond)

	for {
		select {
		case <-pollTicker.C:
			newClient.CollectMetrics()
		case <-reportTicker.C:
			newClient.SendMetrics()
		}
	}
}
