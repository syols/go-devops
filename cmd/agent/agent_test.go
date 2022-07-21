package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
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
	listener, err := net.Listen("tcp", settings.Address())
	assert.NoError(t, err)

	server := httptest.NewUnstartedServer(handlers(t))
	err = server.Listener.Close()
	assert.NoError(t, err)

	server.Listener = listener
	server.Start()
	defer server.Close()

	client := app.NewHTTPClient(settings)
	client.CollectMetrics()
	client.SendMetrics()
	assert.NoError(t, err)
}
