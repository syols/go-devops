package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/syols/go-devops/internal/models"
)

func Example() {
	value := uint64(1)
	metric := models.Metric{
		Name:         "test",
		MetricType:   models.CounterName,
		CounterValue: &value,
	}

	requestBytes, err := json.Marshal(&metric)

	resp, err := http.Post("/update/counter/test/2", "application/json", bytes.NewBuffer(requestBytes))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	metric = models.Metric{
		Name:       "test",
		MetricType: models.CounterName,
	}
	resp, err = http.Post("/value/", "application/json", bytes.NewBuffer(requestBytes))
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
}
