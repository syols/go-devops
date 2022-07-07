package internal

import (
	"bytes"
	"fmt"
	"github.com/syols/go-devops/internal/model"
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
	metrics map[string]model.Metric
	count   uint64
}

func NewHTTPClient(sets settings.Config) Client {
	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}
	client := http.Client{Transport: transport}
	return Client{
		Client:  client,
		scheme:  "http",
		address: sets.GetAddress(),
		metrics: map[string]model.Metric{},
	}
}

func (c *Client) SetMetrics(metrics map[string]model.Metric) {
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

func CollectMetrics() map[string]model.Metric {
	metrics := make(map[string]model.Metric)
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics["Alloc"] = model.GaugeMetric(stats.Alloc)
	metrics["BuckHashSys"] = model.GaugeMetric(stats.BuckHashSys)
	metrics["Frees"] = model.GaugeMetric(stats.Frees)
	metrics["GCCPUFraction"] = model.GaugeMetric(stats.GCCPUFraction)
	metrics["GCSys"] = model.GaugeMetric(stats.GCSys)
	metrics["HeapAlloc"] = model.GaugeMetric(stats.HeapAlloc)
	metrics["HeapIdle"] = model.GaugeMetric(stats.HeapIdle)
	metrics["HeapInuse"] = model.GaugeMetric(stats.HeapInuse)
	metrics["HeapObjects"] = model.GaugeMetric(stats.HeapObjects)
	metrics["HeapReleased"] = model.GaugeMetric(stats.HeapReleased)
	metrics["HeapSys"] = model.GaugeMetric(stats.HeapSys)
	metrics["LastGC"] = model.GaugeMetric(stats.LastGC)
	metrics["Lookups"] = model.GaugeMetric(stats.Lookups)
	metrics["MCacheInuse"] = model.GaugeMetric(stats.MCacheInuse)
	metrics["MCacheSys"] = model.GaugeMetric(stats.MCacheSys)
	metrics["MSpanInuse"] = model.GaugeMetric(stats.MSpanInuse)
	metrics["MSpanSys"] = model.GaugeMetric(stats.MSpanSys)
	metrics["Mallocs"] = model.GaugeMetric(stats.Mallocs)
	metrics["NextGC"] = model.GaugeMetric(stats.NextGC)
	metrics["NumForcedGC"] = model.GaugeMetric(stats.NumForcedGC)
	metrics["NumGC"] = model.GaugeMetric(stats.NumGC)
	metrics["OtherSys"] = model.GaugeMetric(stats.OtherSys)
	metrics["PauseTotalNs"] = model.GaugeMetric(stats.PauseTotalNs)
	metrics["StackInuse"] = model.GaugeMetric(stats.StackInuse)
	metrics["StackSys"] = model.GaugeMetric(stats.StackSys)
	metrics["Sys"] = model.GaugeMetric(stats.Sys)
	metrics["TotalAlloc"] = model.GaugeMetric(stats.TotalAlloc)
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
