package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/settings"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

type MockRoute struct { // добавился слайс тестов
	route  string
	value  string
	method string
}

func mockSettings() settings.Settings {
	sets := settings.Settings{
		Server: settings.ServerSettings{
			Address: settings.Address{
				Host: "0.0.0.0",
				Port: 8081,
			},
		},
		Agent: settings.AgentSettings{},
	}
	return sets
}

func TestServer(t *testing.T) {
	sets := mockSettings()
	log.SetOutput(os.Stdout)
	newServer := internal.NewServer(sets)
	go newServer.Run()
	tr := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
	client := http.Client{Transport: tr}
	time.Sleep(3 * time.Second)
	check(t, MockRoute{
		route:  "http://0.0.0.0:8081/update/counter/PollCount/1",
		value:  "",
		method: http.MethodPost,
	}, client)
	check(t, MockRoute{
		route:  "http://0.0.0.0:8081/value/counter/PollCount",
		value:  "1",
		method: http.MethodGet,
	}, client)
	client.CloseIdleConnections()
}

func check(t *testing.T, test MockRoute, client http.Client) {
	request, err := http.NewRequest(test.method, test.route, bytes.NewBufferString(test.value))
	if err != nil {
		t.Errorf(err.Error())
	}

	response, err := client.Do(request)
	require.NoError(t, err)

	if body, err := io.ReadAll(response.Body); err == nil {
		value := string(body)
		assert.Equal(t, value, test.value, "Test failed")
	}
	defer response.Body.Close()
}
