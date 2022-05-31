package main

import (
	"github.com/syols/go-devops/internal/client"
	"github.com/syols/go-devops/internal/utils"
	"log"
	"os"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)
	var settings utils.Settings
	settings.LoadSettings()

	newClient := client.NewHttpClient(settings)
	pollTicker := time.NewTicker(settings.Agent.PollInterval * time.Millisecond)
	reportTicker := time.NewTicker(settings.Agent.ReportInterval * time.Millisecond)

	var runtimeMetrics = utils.NewGaugeMetricsValues(settings.Metrics.RuntimeMetrics)
	var count uint64

	for {
		select {
		case <-pollTicker.C:
			client.CollectMetrics(runtimeMetrics)
			count++
		case <-reportTicker.C:
			newClient.SendMetrics(runtimeMetrics, count)
		}
	}
}
