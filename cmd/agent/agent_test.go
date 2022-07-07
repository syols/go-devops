package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockSettings(t *testing.T) settings.Config {
	list, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	port := list.Addr().(*net.TCPAddr).Port
	err = list.Close()
	require.NoError(t, err)

	sets := settings.Config{
		Server: settings.ServerConfig{
			Address: settings.Address{
				Host: "0.0.0.0",
				Port: uint16(port),
			},
		},
		Agent: settings.AgentConfig{},
	}
	return sets
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
		log.Printf(r.URL.Path)
	})
	return r
}

func TestAgent(t *testing.T) {
	sets := mockSettings(t)
	listener, err := net.Listen("tcp", sets.GetAddress())
	assert.NoError(t, err)

	server := httptest.NewUnstartedServer(handlers(t))
	err = server.Listener.Close()
	assert.NoError(t, err)

	server.Listener = listener
	server.Start()
	defer server.Close()

	client := internal.NewHTTPClient(sets)
	metrics := internal.CollectMetrics()
	client.SetMetrics(metrics)
	client.SendMetrics()
	client.Client.CloseIdleConnections()
}
