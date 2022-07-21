package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"sync"
	"time"
)

type Client struct {
	Client  http.Client
	metrics map[string]float64
	mutex   sync.RWMutex
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
		Path:   "/updates/",
	}

	return Client{
		Client:  client,
		metrics: map[string]float64{},
		url:     uri.String(),
		key:     settings.Server.Key,
	}
}

func (c *Client) SetMetrics(metrics map[string]float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name, value := range metrics {
		v := value
		c.metrics[name] = v
	}
}

func (c *Client) SendMetrics() {
	var result []models.Metric
	c.mutex.RLock()
	for name, value := range c.metrics {
		v := value
		metric := models.Metric{
			Name:       name,
			MetricType: "gauge",
			GaugeValue: &v,
		}
		metric.Hash = metric.CalculateHash(c.key)
		result = append(result, metric)
	}

	c.count++
	pollCount := models.Metric{Name: "PollCount", CounterValue: &c.count, MetricType: "counter"}
	pollCount.Hash = pollCount.CalculateHash(c.key)
	result = append(result, pollCount)
	defer c.mutex.RUnlock()
	err := c.send(result)
	if err != nil {
		log.Print(err)
	}
}

func (c *Client) send(metric []models.Metric) error {
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
	return nil
}

func (c *Client) CollectMetrics() {
	v, _ := mem.VirtualMemory()
	log.Println("CollectMetrics")

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	result := map[string]float64{
		"Alloc":         float64(stats.Alloc),
		"BuckHashSys":   float64(stats.BuckHashSys),
		"TotalMemory":   float64(v.Total),
		"FreeMemory":    float64(v.Free),
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
	c.SetMetrics(result)
}

func (c *Client) CollectAdditionalMetrics() {
	log.Println("CollectAdditionalMetrics")
	result := make(map[string]float64)
	info, err := cpu.Percent(time.Second*10, true)
	if err != nil {
		log.Fatal(err)
	}
	for i, cp := range info {
		name := fmt.Sprintf("CPUutilization%d", i)
		result[name] = cp
	}
	c.SetMetrics(result)
}
