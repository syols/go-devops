package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/internal"
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

func mockSettings() internal.Settings {
	settings := internal.Settings{
		Server: internal.ServerSettings{
			Address: internal.Address{
				Host: "0.0.0.0",
				Port: 51792,
			},
		},
		Agent: internal.Agent{},
	}
	return settings
}

func TestServer(t *testing.T) {
	settings := mockSettings()
	log.SetOutput(os.Stdout)
	newServer := internal.NewServer(settings)
	go newServer.Run()
	time.Sleep(3 * time.Second)
	client := http.Client{}
	check(t, MockRoute{
		route:  "http://0.0.0.0:51792/update/counter/PollCount/1",
		value:  "",
		method: http.MethodPost,
	}, client)
	check(t, MockRoute{
		route:  "http://0.0.0.0:51792/value/counter/PollCount",
		value:  "1",
		method: http.MethodGet,
	}, client)
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