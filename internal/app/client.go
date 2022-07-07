package app

import (
	"bytes"
	"fmt"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
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

func (c *Client) SendMetrics() {
	for metricName, metricValue := range c.metrics {
		c.post(metricName, fmt.Sprintf("%f", metricValue), "gauge")
	}
	c.count++
	c.post("PollCount", strconv.Itoa(int(c.count)), "counter")
	c.post("RandomValue", fmt.Sprintf("%f", rand.Float64()), "gauge")
}

func CollectMetrics() map[string]models.Metric {
	metrics := make(map[string]models.Metric)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics["Alloc"] = models.GaugeMetric(stats.Alloc)
	metrics["BuckHashSys"] = models.GaugeMetric(stats.BuckHashSys)
	metrics["Frees"] = models.GaugeMetric(stats.Frees)
	metrics["GCCPUFraction"] = models.GaugeMetric(stats.GCCPUFraction)
	metrics["GCSys"] = models.GaugeMetric(stats.GCSys)
	metrics["HeapAlloc"] = models.GaugeMetric(stats.HeapAlloc)
	metrics["HeapIdle"] = models.GaugeMetric(stats.HeapIdle)
	metrics["HeapInuse"] = models.GaugeMetric(stats.HeapInuse)
	metrics["HeapObjects"] = models.GaugeMetric(stats.HeapObjects)
	metrics["HeapReleased"] = models.GaugeMetric(stats.HeapReleased)
	metrics["HeapSys"] = models.GaugeMetric(stats.HeapSys)
	metrics["LastGC"] = models.GaugeMetric(stats.LastGC)
	metrics["Lookups"] = models.GaugeMetric(stats.Lookups)
	metrics["MCacheInuse"] = models.GaugeMetric(stats.MCacheInuse)
	metrics["MCacheSys"] = models.GaugeMetric(stats.MCacheSys)
	metrics["MSpanInuse"] = models.GaugeMetric(stats.MSpanInuse)
	metrics["MSpanSys"] = models.GaugeMetric(stats.MSpanSys)
	metrics["Mallocs"] = models.GaugeMetric(stats.Mallocs)
	metrics["NextGC"] = models.GaugeMetric(stats.NextGC)
	metrics["NumForcedGC"] = models.GaugeMetric(stats.NumForcedGC)
	metrics["NumGC"] = models.GaugeMetric(stats.NumGC)
	metrics["OtherSys"] = models.GaugeMetric(stats.OtherSys)
	metrics["PauseTotalNs"] = models.GaugeMetric(stats.PauseTotalNs)
	metrics["StackInuse"] = models.GaugeMetric(stats.StackInuse)
	metrics["StackSys"] = models.GaugeMetric(stats.StackSys)
	metrics["Sys"] = models.GaugeMetric(stats.Sys)
	metrics["TotalAlloc"] = models.GaugeMetric(stats.TotalAlloc)
	return metrics
}

func (c *Client) post(metricName string, metricValue string, metricAlias string) {
	endpoint := c.endpoint(metricAlias, metricName, metricValue)
	payload := c.payload(metricAlias, metricName, metricValue)
	request, err := http.NewRequest(http.MethodPost, endpoint.String(), bytes.NewBufferString(payload))
	if err != nil {
		log.Print(err.Error())
	}

	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(payload)))

	response, err := c.Client.Do(request)
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = response.Body.Close()
	if err != nil {
		log.Print(err.Error())
	}

	defer c.Client.CloseIdleConnections()
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
