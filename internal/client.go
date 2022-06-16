package internal

import (
	"bytes"
	"fmt"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

type Client struct {
	Client  http.Client
	scheme  string
	address string
	metrics map[string]metric.Metric
	count   uint64
}

func NewHTTPClient(sets settings.Settings) Client {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
	client := http.Client{Transport: transport}
	return Client{
		Client:  client,
		scheme:  "http",
		address: sets.GetAddress(),
		metrics: map[string]metric.Metric{},
	}
}

func (c *Client) SetMetrics(metrics map[string]metric.Metric) {
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

func CollectMetrics() map[string]metric.Metric {
	metrics := make(map[string]metric.Metric)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics["Alloc"] = metric.GaugeMetric(stats.Alloc)
	metrics["BuckHashSys"] = metric.GaugeMetric(stats.BuckHashSys)
	metrics["Frees"] = metric.GaugeMetric(stats.Frees)
	metrics["GCCPUFraction"] = metric.GaugeMetric(stats.GCCPUFraction)
	metrics["GCSys"] = metric.GaugeMetric(stats.GCSys)
	metrics["HeapAlloc"] = metric.GaugeMetric(stats.HeapAlloc)
	metrics["HeapIdle"] = metric.GaugeMetric(stats.HeapIdle)
	metrics["HeapInuse"] = metric.GaugeMetric(stats.HeapInuse)
	metrics["HeapObjects"] = metric.GaugeMetric(stats.HeapObjects)
	metrics["HeapReleased"] = metric.GaugeMetric(stats.HeapReleased)
	metrics["HeapSys"] = metric.GaugeMetric(stats.HeapSys)
	metrics["LastGC"] = metric.GaugeMetric(stats.LastGC)
	metrics["Lookups"] = metric.GaugeMetric(stats.Lookups)
	metrics["MCacheInuse"] = metric.GaugeMetric(stats.MCacheInuse)
	metrics["MCacheSys"] = metric.GaugeMetric(stats.MCacheSys)
	metrics["MSpanInuse"] = metric.GaugeMetric(stats.MSpanInuse)
	metrics["MSpanSys"] = metric.GaugeMetric(stats.MSpanSys)
	metrics["Mallocs"] = metric.GaugeMetric(stats.Mallocs)
	metrics["NextGC"] = metric.GaugeMetric(stats.NextGC)
	metrics["NumForcedGC"] = metric.GaugeMetric(stats.NumForcedGC)
	metrics["NumGC"] = metric.GaugeMetric(stats.NumGC)
	metrics["OtherSys"] = metric.GaugeMetric(stats.OtherSys)
	metrics["PauseTotalNs"] = metric.GaugeMetric(stats.PauseTotalNs)
	metrics["StackInuse"] = metric.GaugeMetric(stats.StackInuse)
	metrics["StackSys"] = metric.GaugeMetric(stats.StackSys)
	metrics["Sys"] = metric.GaugeMetric(stats.Sys)
	metrics["TotalAlloc"] = metric.GaugeMetric(stats.TotalAlloc)
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
	response, err := c.Client.Do(request)
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
