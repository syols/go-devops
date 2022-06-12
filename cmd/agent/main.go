package main

import (
	"github.com/syols/go-devops/internal"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	settings := internal.NewSettings()

	client := internal.NewHTTPClient(settings)
	pollTicker := time.NewTicker(settings.Agent.PollInterval)
	reportTicker := time.NewTicker(settings.Agent.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			go client.SetMetrics(internal.CollectMetrics())
		case <-reportTicker.C:
			go client.SendMetrics()
		}
	}
}
