package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
)

// @Title Agent API
// @Description Агент для получения метрик
// @Version 0.1.0

// @Contact.email some@mail.com

//Example: go run -ldflags "-X main.buildVersion=%%" main.go
var buildVersion string
var buildDate string
var buildCommit string

func main() {
	log.SetOutput(os.Stdout)
	log.Printf("Build version: %s", config.ReplaceNoneValue(buildVersion))
	log.Printf("Build date: %s", config.ReplaceNoneValue(buildDate))
	log.Printf("Build commit: %s", config.ReplaceNoneValue(buildCommit))

	settings := config.NewConfig()
	var wg sync.WaitGroup
	client := app.NewClient(settings)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	wg.Add(3)
	go client.CollectAdditionalMetrics(ctx, &wg)
	go client.CollectMetrics(ctx, &wg)
	go client.SendMetrics(ctx, &wg)
	wg.Wait()
	log.Print("Done!")
}
