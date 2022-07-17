package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/app"
	"github.com/syols/go-devops/internal/models"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type Mock struct {
	route    string
	request  *models.Metric
	response models.Metric
	method   string
}

type MetricSuite struct {
	suite.Suite
	settings config.Config
	client   http.Client
	server   app.Server
}

func (suite *MetricSuite) SetupTest() {
	list, err := net.Listen("tcp", ":0")
	suite.NoError(err)
	port := list.Addr().(*net.TCPAddr).Port
	err = list.Close()
	suite.NoError(err)
	suite.settings = config.Config{
		Server: config.ServerConfig{
			Address: config.Address{
				Host: "0.0.0.0",
				Port: uint16(port),
			},
		},
		Store: config.StoreConfig{
			StoreInterval: time.Second * 10,
		},
	}
	suite.client = http.Client{Transport: &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}}
}

func (suite *MetricSuite) BeforeTest(suiteName, testName string) {
	server, err := app.NewServer(suite.settings)
	suite.NoError(err)
	go server.Run()
	time.Sleep(time.Second)
}

func (suite *MetricSuite) TestUpdateGauge() {
	value := 1.1
	metric := models.Metric{
		Name:       "testGauge",
		MetricType: "gauge",
		GaugeValue: &value,
	}
	mocks := []Mock{
		{
			route:    "/update/gauge/testGauge/1.1",
			request:  nil,
			response: metric,
			method:   "POST",
		},
		{
			route:    "/value/gauge/testGauge",
			request:  nil,
			response: metric,
			method:   "GET",
		},
	}
	for _, mock := range mocks {
		suite.check(mock)
	}
}

func (suite *MetricSuite) TestUpdateGaugeJSON() {
	value := 1.1
	metric := models.Metric{
		Name:       "testGauge",
		MetricType: "gauge",
		GaugeValue: &value,
	}
	mock := Mock{
		route:    "/update/",
		request:  &metric,
		response: metric,
		method:   "POST",
	}
	suite.check(mock)
}

func (suite *MetricSuite) TestUpdateCounter() {
	value := uint64(1)
	metric := models.Metric{
		Name:         "test",
		MetricType:   "counter",
		CounterValue: &value,
	}
	updatedValue := uint64(3)
	updatedMetric := models.Metric{
		Name:         "test",
		MetricType:   "counter",
		CounterValue: &updatedValue,
	}
	mocks := []Mock{
		{
			route:    "/update/counter/test/1",
			request:  nil,
			response: metric,
			method:   "POST",
		}, {
			route:    "/update/counter/test/2",
			request:  nil,
			response: updatedMetric,
			method:   "POST",
		}, {
			route:    "/value/counter/test",
			request:  nil,
			response: updatedMetric,
			method:   "GET",
		},
	}
	for _, mock := range mocks {
		suite.check(mock)
	}
}

func (suite *MetricSuite) TestUpdateCounterJSON() {
	value := uint64(1)
	metric := models.Metric{
		Name:         "test",
		MetricType:   "counter",
		CounterValue: &value,
	}
	updatedValue := uint64(2)
	updatedMetric := models.Metric{
		Name:         "test",
		MetricType:   "counter",
		CounterValue: &updatedValue,
	}
	mocks := []Mock{
		{
			route:    "/update/",
			request:  &metric,
			response: metric,
			method:   "POST",
		}, {
			route:    "/update/",
			request:  &metric,
			response: updatedMetric,
			method:   "POST",
		},
	}
	for _, mock := range mocks {
		suite.check(mock)
	}
}

func (suite *MetricSuite) check(mock Mock) {
	uri := url.URL{
		Scheme: "http",
		Host:   suite.settings.Address(),
		Path:   mock.route,
	}

	requestBytes, err := json.Marshal(mock.request)
	suite.NoError(err)

	request, err := http.NewRequest(mock.method, uri.String(), bytes.NewReader(requestBytes))
	request.Header.Set("Content-Type", "application/json")
	response, err := suite.client.Do(request)
	suite.NoError(err)

	var responsePayload models.Metric
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&responsePayload)

	suite.NoError(err)
	suite.Equal(responsePayload, mock.response)
	suite.NoError(response.Body.Close())
}

func TestMetricSuite(t *testing.T) {
	suite.Run(t, new(MetricSuite))
}
