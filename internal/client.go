package internal

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

const SchemeName = "http"
const ActionType = "update"
const SkipMessage = "Skip field: %s"
const RequestErrorMessage = "Request error: %s"

type Client struct {
	client  http.Client
	scheme  string
	address string
	metric  GaugeMetricsValues
	count   uint64
}

func NewHTTPClient(settings Settings) Client {
	client := http.Client{}
	return Client{
		client:  client,
		scheme:  SchemeName,
		address: settings.GetAddress(),
		metric:  GaugeMetricsValues{},
		count:   0,
	}
}

func (c *Client) SendMetrics() {
	for metricName, metricValue := range c.metric {
		c.post(metricName, fmt.Sprintf("%f", metricValue), GaugeMetric)
	}
	c.post(PollCountMetricName, strconv.Itoa(int(c.count)), CounterMetric)
	c.post(RandomValueMetricName, fmt.Sprintf("%f", rand.Float64()), GaugeMetric)
}

func (c *Client) CollectMetrics() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	(*c).metric["Alloc"] = float64(stats.Alloc)
	(*c).metric["BuckHashSys"] = float64(stats.BuckHashSys)
	(*c).metric["Frees"] = float64(stats.Frees)
	(*c).metric["GCCPUFraction"] = stats.GCCPUFraction
	(*c).metric["GCSys"] = float64(stats.GCSys)
	(*c).metric["HeapAlloc"] = float64(stats.HeapAlloc)
	(*c).metric["HeapIdle"] = float64(stats.HeapIdle)
	(*c).metric["HeapInuse"] = float64(stats.HeapInuse)
	(*c).metric["HeapObjects"] = float64(stats.HeapObjects)
	(*c).metric["HeapReleased"] = float64(stats.HeapReleased)
	(*c).metric["HeapSys"] = float64(stats.HeapSys)
	(*c).metric["LastGC"] = float64(stats.LastGC)
	(*c).metric["Lookups"] = float64(stats.Lookups)
	(*c).metric["MCacheInuse"] = float64(stats.MCacheInuse)
	(*c).metric["MCacheSys"] = float64(stats.MCacheSys)
	(*c).metric["MSpanInuse"] = float64(stats.MSpanInuse)
	(*c).metric["MSpanSys"] = float64(stats.MSpanSys)
	(*c).metric["Mallocs"] = float64(stats.Mallocs)
	(*c).metric["NextGC"] = float64(stats.NextGC)
	(*c).metric["NumForcedGC"] = float64(stats.NumForcedGC)
	(*c).metric["NumGC"] = float64(stats.NumGC)
	(*c).metric["OtherSys"] = float64(stats.OtherSys)
	(*c).metric["PauseTotalNs"] = float64(stats.PauseTotalNs)
	(*c).metric["StackInuse"] = float64(stats.StackInuse)
	(*c).metric["StackSys"] = float64(stats.StackSys)
	(*c).metric["Sys"] = float64(stats.Sys)
	(*c).metric["TotalAlloc"] = float64(stats.TotalAlloc)
	c.count++
}

func (c *Client) post(metricName string, metricValue string, metricAlias string) {
	endpoint := c.endpoint(metricAlias, metricName, metricValue)
	payload := c.payload(metricAlias, metricName, metricValue)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBufferString(payload))
	if err != nil {
		log.Fatalf(RequestErrorMessage, err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	response, respErr := c.client.Do(request)
	if respErr != nil {
		log.Printf(RequestErrorMessage, respErr)
	} else {
		defer response.Body.Close()
		log.Printf(metricName, response.Status)
	}
}

func (c *Client) endpoint(alias string, metricName string, metricValue string) url.URL {
	return url.URL{
		Scheme: c.scheme,
		Host:   c.address,
		Path:   strings.Join([]string{ActionType, alias, metricName, metricValue}, "/"),
	}
}

func (c *Client) payload(alias string, metricName string, metricValue string) string {
	return strings.Join([]string{ActionType, alias, metricName, metricValue}, "::")
}
