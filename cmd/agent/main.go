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

func main() {
	log.SetOutput(os.Stdout)
	settings := config.NewConfig()
	var wg sync.WaitGroup
	client := app.NewHTTPClient(settings)
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	wg.Add(3)
	go client.CollectAdditionalMetrics(ctx, &wg)
	go client.CollectMetrics(ctx, &wg)
	go client.SendMetrics(ctx, &wg)
	wg.Wait()
	log.Print("Done!")
}
