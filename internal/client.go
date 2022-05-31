package client

import (
	"bytes"
	"fmt"
	"github.com/syols/go-devops/internal/utils"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const SchemeName = "http"
const ActionType = "update"
const RequestErrorMessage = "Request error: %s"

type Client struct {
	client  http.Client
	scheme  string
	address string
}

func NewHttpClient(settings utils.Settings) Client {
	client := http.Client{}
	return Client{
		client:  client,
		scheme:  SchemeName,
		address: settings.GetAddress(),
	}
}

func (h *Client) SendMetrics(values utils.GaugeMetricsValues, counter uint64) {
	for metricName, metricValue := range values {
		h.post(metricName, fmt.Sprintf("%f", metricValue), utils.GaugeMetric)
	}
	h.post(utils.PollCountMetricName, strconv.Itoa(int(counter)), utils.CounterMetric)
	h.post(utils.RandomValueMetricName, fmt.Sprintf("%f", rand.Float64()), utils.GaugeMetric)
}

func (h *Client) post(metricName string, metricValue string, metricAlias string) {
	endpoint := h.endpoint(metricAlias, metricName, metricValue)
	payload := h.payload(metricAlias, metricName, metricValue)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBufferString(payload))
	if err != nil {
		log.Fatalf(RequestErrorMessage, err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	if response, respErr := h.client.Do(request); respErr != nil {
		log.Printf(RequestErrorMessage, respErr)
	} else {
		log.Printf(metricName, response.Status)
	}
}

func (h *Client) endpoint(alias string, metricName string, metricValue string) url.URL {
	return url.URL{
		Scheme: h.scheme,
		Host:   h.address,
		Path:   strings.Join([]string{ActionType, alias, metricName, metricValue}, "/"),
	}
}

func (h *Client) payload(alias string, metricName string, metricValue string) string {
	return strings.Join([]string{ActionType, alias, metricName, metricValue}, "::")
}
