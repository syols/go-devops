package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/syols/go-devops/internal"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSettingsMock() internal.Settings {
	settings := internal.Settings{
		Server: internal.ServerSettings{
			Address: internal.Address{
				Host: "0.0.0.0",
				Port: 51791,
			},
		},
		Agent: internal.Agent{},
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
		log.Printf(r.URL.Path)
	})
	return r
}

func TestAgent(t *testing.T) {
	settings := newSettingsMock()
	listener, err := net.Listen("tcp", settings.GetAddress())
	assert.NoError(t, err)

	server := httptest.NewUnstartedServer(handlers(t))
	err = server.Listener.Close()
	assert.NoError(t, err)

	server.Listener = listener
	server.Start()
	defer server.Close()

	client := internal.NewHTTPClient(settings)
	metrics := internal.CollectMetrics()
	client.SetMetrics(metrics)
	client.SendMetrics()
}
