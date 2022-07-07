package main

import (
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	sets := settings.NewConfig()
	client := internal.NewHTTPClient(sets)
	pollTicker := time.NewTicker(sets.Agent.PollInterval)
	reportTicker := time.NewTicker(sets.Agent.ReportInterval)

	for {
		select {
		case <-pollTicker.C:
			go client.SetMetrics(internal.CollectMetrics())
		case <-reportTicker.C:
			go client.SendMetrics()
		}
	}
}
