package main

import (
	"github.com/syols/go-devops/internal"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func clientMock(settings internal.Settings) internal.Client {
	newClient := internal.NewHTTPClient(settings)
	newClient.CollectMetrics()
	newClient.SendMetrics()
	return newClient
}

func mockSettings() internal.Settings {
	settings := internal.Settings{
		Address: internal.Address{
			Host: "0.0.0.0",
			Port: 51791,
		},
		Agent: internal.Agent{},
	}
	return settings
}

func handlers(t *testing.T) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/update/gauge/Alloc/0", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("Alloc metric not updated")
	})
	r.HandleFunc("/update/counter/PollCount/0", func(w http.ResponseWriter, r *http.Request) {
		log.Fatal("Count metric not updated")
	})
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(r.URL.Path)
	})
	return r
}

func TestAgent(t *testing.T) {
	settings := mockSettings()
	listener, err := net.Listen("tcp", settings.GetAddress())
	if err != nil {
		log.Fatal(err)
	}

	server := httptest.NewUnstartedServer(handlers(t))
	err = server.Listener.Close()
	if err != nil {
		log.Fatal(err)
	}
	server.Listener = listener
	server.Start()
	defer server.Close()

	newClient := clientMock(settings)
	newClient.CollectMetrics()
	newClient.SendMetrics()
}
