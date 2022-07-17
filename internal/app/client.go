package app

import (
	"bytes"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
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
	metrics map[string]models.Metric
	count   uint64
}

func NewHTTPClient(settings config.Config) Client {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
	client := http.Client{Transport: transport}
	return Client{
		Client:  client,
		scheme:  "http",
		address: settings.Address(),
		metrics: map[string]models.Metric{},
	}
}

func (c *Client) SetMetrics(metrics map[string]models.Metric) {
	c.metrics = metrics
}

func (c *Client) SendMetrics() error {
	for _, metricValue := range c.metrics {
		err := c.post(metricValue)
		if err != nil {
			return err
		}
	}
	c.count++
	pollCount := models.Metric{Name: "PollCount", CounterValue: &c.count, MetricType: "counter"}
	err := c.post(pollCount)
	if err != nil {
		return err
	}

	value := rand.Float64()
	randomValue := models.Metric{Name: "RandomValue", GaugeValue: &value, MetricType: "gauge"}
	err = c.post(randomValue)
	if err != nil {
		return err
	}
	return nil
}

func CollectMetrics() store.Metrics {
	metrics := make(map[string]models.Metric)

	var stats runtime.MemStats
	val := map[string]float64{
		"Alloc":         float64(stats.Alloc),
		"BuckHashSys":   float64(stats.BuckHashSys),
		"Frees":         float64(stats.Frees),
		"GCCPUFraction": stats.GCCPUFraction,
		"GCSys":         float64(stats.GCSys),
		"HeapAlloc":     float64(stats.HeapAlloc),
		"HeapIdle":      float64(stats.HeapIdle),
		"HeapInuse":     float64(stats.HeapInuse),
		"HeapObjects":   float64(stats.HeapObjects),
		"HeapReleased":  float64(stats.HeapReleased),
		"HeapSys":       float64(stats.HeapSys),
		"LastGC":        float64(stats.LastGC),
		"Lookups":       float64(stats.Lookups),
		"MCacheInuse":   float64(stats.MCacheInuse),
		"MCacheSys":     float64(stats.MCacheSys),
		"MSpanInuse":    float64(stats.MSpanInuse),
		"MSpanSys":      float64(stats.MSpanSys),
		"Mallocs":       float64(stats.Mallocs),
		"NextGC":        float64(stats.NextGC),
		"NumForcedGC":   float64(stats.NumForcedGC),
		"NumGC":         float64(stats.NumGC),
		"OtherSys":      float64(stats.OtherSys),
		"PauseTotalNs":  float64(stats.PauseTotalNs),
		"StackInuse":    float64(stats.StackInuse),
		"StackSys":      float64(stats.StackSys),
		"Sys":           float64(stats.Sys),
		"TotalAlloc":    float64(stats.TotalAlloc),
	}

	for name, value := range val {
		metrics[name] = models.Metric{
			Name:       name,
			MetricType: "gauge",
			GaugeValue: &value,
		}
	}
	return metrics
}

func (c *Client) post(metric models.Metric) error {
	payload := c.payload(metric)
	request, err := http.NewRequest(http.MethodPost, c.endpoint(metric), bytes.NewBufferString(payload))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))
	response, err := c.Client.Do(request)
	if err != nil {
		return err
	}

	err = response.Body.Close()
	if err != nil {
		return err
	}

	defer c.Client.CloseIdleConnections()
	return nil
}

func (c *Client) endpoint(metric models.Metric) string {
	result := url.URL{
		Scheme: c.scheme,
		Host:   c.address,
		Path:   strings.Join([]string{"update", metric.MetricType, metric.Name, metric.String()}, "/"),
	}
	return result.String()
}

func (c *Client) payload(metric models.Metric) string {
	return strings.Join([]string{"update", metric.MetricType, metric.Name, metric.String()}, "::")
}
