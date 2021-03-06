package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/syols/go-devops/internal"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

type MockRoute struct { // добавился слайс тестов
	route    string
	response *string
	request  string
	method   string
}

func mockSettings(t *testing.T) settings.Settings {
	list, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	port := list.Addr().(*net.TCPAddr).Port
	err = list.Close()
	require.NoError(t, err)

	sets := settings.Settings{
		Server: settings.ServerSettings{
			Address: settings.Address{
				Host: "0.0.0.0",
				Port: uint16(port),
			},
		},
		Agent: settings.AgentSettings{},
		Store: settings.StoreSettings{
			StoreInterval: time.Second * 10,
		},
	}
	return sets
}

func mockClientServer(sets settings.Settings) http.Client {
	log.SetOutput(os.Stdout)
	server := internal.NewServer(sets)
	go server.Run()
	tr := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
	client := http.Client{Transport: tr}
	time.Sleep(time.Millisecond * 500)
	return client
}

func TestPlainServer(t *testing.T) {
	sets := mockSettings(t)
	client := mockClientServer(sets)
	defer client.CloseIdleConnections()

	//Test update - plain/text
	uri := url.URL{
		Scheme: "http",
		Host:   sets.GetAddress(),
		Path:   "/update/counter/PollCount/1",
	}
	checkPlainText(t, MockRoute{
		route:  uri.String(),
		method: http.MethodPost,
	}, client)

	responseString := "1"
	uri = url.URL{
		Scheme: "http",
		Host:   sets.GetAddress(),
		Path:   "/value/counter/PollCount",
	}
	checkPlainText(t, MockRoute{
		route:    uri.String(),
		response: &responseString,
		method:   http.MethodGet,
	}, client)
}

func TestJsonServer(t *testing.T) {
	log.SetOutput(os.Stdout)
	sets := mockSettings(t)
	client := mockClientServer(sets)
	defer client.CloseIdleConnections()

	uri := url.URL{
		Scheme: "http",
		Host:   sets.GetAddress(),
		Path:   "/update/",
	}
	payloadString := string(`{"id":"testGauge","type":"gauge","value":100}`)
	checkApplicationJSON(t, MockRoute{
		route:   uri.String(),
		request: payloadString,
		method:  http.MethodPost,
	}, client)

	uri = url.URL{
		Scheme: "http",
		Host:   sets.GetAddress(),
		Path:   "/value/",
	}
	checkApplicationJSON(t, MockRoute{
		route:    uri.String(),
		request:  payloadString,
		response: &payloadString,
		method:   http.MethodPost,
	}, client)
}

func checkPlainText(t *testing.T, test MockRoute, client http.Client) {
	request, err := http.NewRequest(test.method, test.route, bytes.NewBufferString(test.request))
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)

	if test.response != nil {
		if body, err := io.ReadAll(response.Body); err == nil {
			value := string(body)
			assert.Equal(t, value, *test.response, "Test failed")
		}
	}
	require.NoError(t, response.Body.Close())
}

func checkApplicationJSON(t *testing.T, test MockRoute, client http.Client) {
	request, err := http.NewRequest(test.method, test.route, bytes.NewBufferString(test.request))
	request.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	response, err := client.Do(request)
	require.NoError(t, err)

	if test.response != nil {
		var responsePayload, requestPayload metric.Payload
		log.Print("check")
		decoder := json.NewDecoder(response.Body)
		assert.NoError(t, decoder.Decode(&responsePayload))
		assert.NoError(t, json.Unmarshal([]byte(*test.response), &requestPayload))
		assert.Equal(t, responsePayload, requestPayload)
	}

	require.NoError(t, response.Body.Close())
}
