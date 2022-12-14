package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
)

func mockSettings(t *testing.T) config.Config {
	list, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	port := list.Addr().(*net.TCPAddr).Port
	err = list.Close()
	require.NoError(t, err)

	settings := config.Config{
		Server: config.ServerConfig{
			Address: config.Address{
				Host: "0.0.0.0",
				Port: uint16(port),
			},
		},
		Store: config.StoreConfig{
			StoreInterval: time.Second,
		},
		Agent: config.AgentConfig{
			PollInterval:   time.Second,
			ReportInterval: time.Second,
		},
	}
	return settings
}

func handlers(t *testing.T) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/update/gauge/Alloc/0", func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "Alloc metric not updated")
	})
	r.HandleFunc("/update/counter/PollCount/0", func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "Count metric not updated")
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})
	return r
}

func TestAgent(t *testing.T) {
	settings := mockSettings(t)
	listener, err := net.Listen("tcp", settings.Server.Address.String())
	assert.NoError(t, err)

	server := httptest.NewUnstartedServer(handlers(t))
	err = server.Listener.Close()
	assert.NoError(t, err)

	server.Listener = listener
	server.Start()
	defer server.Close()

	var wg sync.WaitGroup
	client := app.NewClient(settings)
	wg.Add(3)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go client.CollectAdditionalMetrics(ctx, &wg)
	go client.CollectMetrics(ctx, &wg)
	go client.SendMetrics(ctx, &wg)
	cancel()
	wg.Wait()
	assert.NoError(t, err)
}
