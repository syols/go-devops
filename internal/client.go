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

type Client struct {
	client  http.Client
	scheme  string
	address string
	metrics Metrics
	count   uint64
}

func NewHTTPClient(settings Settings) Client {
	client := http.Client{}
	return Client{
		client:  client,
		scheme:  "http",
		address: settings.GetAddress(),
		metrics: Metrics{},
	}
}

func (c *Client) SetMetrics(metrics Metrics) {
	c.metrics = metrics
}

func (c *Client) SendMetrics() {
	for metricName, metricValue := range c.metrics {
		c.post(metricName, fmt.Sprintf("%f", metricValue), "gauge")
	}
	c.count++
	c.post("PollCount", strconv.Itoa(int(c.count)), "counter")
	c.post("RandomValue", fmt.Sprintf("%f", rand.Float64()), "gauge")
}

func CollectMetrics() Metrics {
	metrics := make(Metrics)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics["Alloc"] = GaugeMetric(stats.Alloc)
	metrics["BuckHashSys"] = GaugeMetric(stats.BuckHashSys)
	metrics["Frees"] = GaugeMetric(stats.Frees)
	metrics["GCCPUFraction"] = GaugeMetric(stats.GCCPUFraction)
	metrics["GCSys"] = GaugeMetric(stats.GCSys)
	metrics["HeapAlloc"] = GaugeMetric(stats.HeapAlloc)
	metrics["HeapIdle"] = GaugeMetric(stats.HeapIdle)
	metrics["HeapInuse"] = GaugeMetric(stats.HeapInuse)
	metrics["HeapObjects"] = GaugeMetric(stats.HeapObjects)
	metrics["HeapReleased"] = GaugeMetric(stats.HeapReleased)
	metrics["HeapSys"] = GaugeMetric(stats.HeapSys)
	metrics["LastGC"] = GaugeMetric(stats.LastGC)
	metrics["Lookups"] = GaugeMetric(stats.Lookups)
	metrics["MCacheInuse"] = GaugeMetric(stats.MCacheInuse)
	metrics["MCacheSys"] = GaugeMetric(stats.MCacheSys)
	metrics["MSpanInuse"] = GaugeMetric(stats.MSpanInuse)
	metrics["MSpanSys"] = GaugeMetric(stats.MSpanSys)
	metrics["Mallocs"] = GaugeMetric(stats.Mallocs)
	metrics["NextGC"] = GaugeMetric(stats.NextGC)
	metrics["NumForcedGC"] = GaugeMetric(stats.NumForcedGC)
	metrics["NumGC"] = GaugeMetric(stats.NumGC)
	metrics["OtherSys"] = GaugeMetric(stats.OtherSys)
	metrics["PauseTotalNs"] = GaugeMetric(stats.PauseTotalNs)
	metrics["StackInuse"] = GaugeMetric(stats.StackInuse)
	metrics["StackSys"] = GaugeMetric(stats.StackSys)
	metrics["Sys"] = GaugeMetric(stats.Sys)
	metrics["TotalAlloc"] = GaugeMetric(stats.TotalAlloc)
	return metrics
}

func (c *Client) post(metricName string, metricValue string, metricAlias string) {
	endpoint := c.endpoint(metricAlias, metricName, metricValue)
	payload := c.payload(metricAlias, metricName, metricValue)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBufferString(payload))
	if err != nil {
		log.Fatalf("Request error: %s", err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	response, err := c.client.Do(request)
	if err != nil {
		log.Printf("Request error: %s", err)
		return
	}
	defer response.Body.Close()
	log.Printf(metricName, response.Status)
}

func (c *Client) endpoint(alias string, metricName string, metricValue string) url.URL {
	return url.URL{
		Scheme: c.scheme,
		Host:   c.address,
		Path:   strings.Join([]string{"update", alias, metricName, metricValue}, "/"),
	}
}

func (c *Client) payload(alias string, metricName string, metricValue string) string {
	return strings.Join([]string{"update", alias, metricName, metricValue}, "::")
}
