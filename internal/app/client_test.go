package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/syols/go-devops/config"
)

func TestSetMetrics(t *testing.T) {
	duration := time.Duration(1)
	t.Setenv("ADDRESS", "0.0.0.0:8080")
	t.Setenv("REPORT_INTERVAL", duration.String())
	t.Setenv("POLL_INTERVAL", duration.String())
	t.Setenv("CLIENT_TIMEOUT", duration.String())
	t.Setenv("STORE_INTERVAL", duration.String())
	t.Setenv("RESTORE", "true")
	t.Setenv("KEY", "some_key")
	t.Setenv("STORE_FILE", "/tmp/devops-metrics-db.json")
	t.Setenv("DATABASE_DSN", "postgres://postgres:postgres@localhost/postgres?sslmode=disable")

	cfg := config.Config{}
	cfg.LoadFromEnvironment()

	client := NewClient(cfg)
	client.SetMetrics(map[string]float64{"some": 41})
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*100))
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	client.SendMetrics(ctx, &wg)
}
