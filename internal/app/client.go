package app

import (
	"bytes"
	"encoding/json"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
)

type Client struct {
	Client  http.Client
	metrics []models.Metric
	count   uint64
	url     string
	key     *string
}

func NewHTTPClient(settings config.Config) Client {
	transport := &http.Transport{
		MaxIdleConns:        40,
		MaxIdleConnsPerHost: 40,
	}
	client := http.Client{Transport: transport}

	uri := url.URL{
		Scheme: "http",
		Host:   settings.Address(),
		Path:   "/update/",
	}

	return Client{
		Client:  client,
		metrics: []models.Metric{},
		url:     uri.String(),
		key:     settings.Server.Key,
	}
}

func (c *Client) SetMetrics(metrics []models.Metric) {
	c.metrics = metrics
}

func (c *Client) SendMetrics() error {
	c.count++
	pollCount := models.Metric{Name: "PollCount", CounterValue: &c.count, MetricType: "counter"}
	pollCount.Hash = pollCount.CalculateHash(c.key)

	for _, metric := range append([]models.Metric{pollCount}, c.metrics...) {
		requestBytes, err := json.Marshal(metric)
		if err != nil {
			return err
		}
		resp, err := http.Post(c.url, "application/json", bytes.NewBuffer(requestBytes))
		if err != nil {
			return err
		}

		err = resp.Body.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func CollectMetrics(key *string) []models.Metric {
	metrics := make([]models.Metric, 0)

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
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
		"RandomValue":   rand.Float64(),
	}

	for name, value := range val {
		val := value
		metric := models.Metric{
			Name:       name,
			MetricType: "gauge",
			GaugeValue: &val,
		}
		metric.Hash = metric.CalculateHash(key)
		metrics = append(metrics, metric)
	}
	return metrics
}
