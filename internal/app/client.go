package app

import (
	"bytes"
	"context"
	cryprorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	pb "github.com/syols/go-devops/internal/rpc/proto"
)

type Transport interface {
	send(metrics []models.Metric) error
}

// Client struct
type Client struct {
	metrics        map[string]float64
	key            *string
	transport      Transport
	url            string
	count          uint64
	pollInterval   time.Duration
	reportInterval time.Duration
	mutex          sync.RWMutex
}

// HTTPTransport struct
type HTTPTransport struct {
	url       string
	http      http.Client
	publicKey *rsa.PublicKey
}

// GrpcTransport struct
type GrpcTransport struct {
	grpc pb.GoDevopsClient
}

func NewHTTPTransport(settings config.Config) HTTPTransport {
	transport := &http.Transport{
		MaxIdleConns:        40,
		MaxIdleConnsPerHost: 40,
	}

	httpClient := http.Client{Transport: transport}
	uri := url.URL{
		Scheme: "http",
		Host:   settings.Server.Address.String(),
		Path:   "/updates/",
	}

	var publicKey *rsa.PublicKey
	if settings.Store.CryptoKeyFilePath != nil {
		byteArr, err := os.ReadFile(*settings.Store.CryptoKeyFilePath) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		if blocks, err := pem.Decode(byteArr); err == nil {
			pubInterface, _ := x509.ParsePKIXPublicKey(blocks.Bytes)
			publicKey = pubInterface.(*rsa.PublicKey)
		}
	}

	return HTTPTransport{
		url:       uri.String(),
		http:      httpClient,
		publicKey: publicKey,
	}
}

func NewGrpcTransport(settings config.Config) GrpcTransport {
	port := fmt.Sprintf(":%d", settings.Grpc.Address.Port)
	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return GrpcTransport{
		grpc: pb.NewGoDevopsClient(conn),
	}
}

func NewTransport(settings config.Config) Transport {
	if settings.Grpc != nil {
		return NewGrpcTransport(settings)
	}

	return NewHTTPTransport(settings)
}

// NewClient creates new HTTP client struct
func NewClient(settings config.Config) Client {
	return Client{
		transport:      NewTransport(settings),
		metrics:        map[string]float64{},
		key:            settings.Server.Key,
		pollInterval:   settings.Agent.PollInterval,
		reportInterval: settings.Agent.ReportInterval,
	}
}

// SetMetrics set metric to store
func (c *Client) SetMetrics(metrics map[string]float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for name := range metrics {
		c.metrics[name] = metrics[name]
	}
}

// SendMetrics send metric
func (c *Client) SendMetrics(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	reportInterval := time.NewTicker(c.reportInterval)
	for {
		select {
		case <-reportInterval.C:
			func() {
				log.Println("SendMetrics")
				var result []models.Metric
				c.mutex.RLock()
				defer c.mutex.RUnlock()
				for name, value := range c.metrics {
					v := value
					metric := models.Metric{
						Name:       name,
						MetricType: models.GaugeName,
						GaugeValue: &v,
					}
					metric.Hash = metric.CalculateHash(c.key)
					result = append(result, metric)
				}

				c.count++
				pollCount := models.Metric{Name: "PollCount", CounterValue: &c.count, MetricType: models.CounterName}
				pollCount.Hash = pollCount.CalculateHash(c.key)
				result = append(result, pollCount)
				if err := c.transport.send(result); err != nil {
					log.Print(err)
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}

func (h HTTPTransport) send(metric []models.Metric) error {
	requestBytes, err := json.Marshal(metric)
	encryptedBytes := TryEncrypt(requestBytes, h.publicKey)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", h.url, bytes.NewBuffer(encryptedBytes))
	if err != nil {
		return err
	}

	realIP, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", realIP)

	resp, err := h.http.Do(req)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func (g GrpcTransport) send(metrics []models.Metric) error {
	for _, metric := range metrics {
		message := pb.MetricMessage{
			Name: metric.Name,
		}

		if metric.MetricType == models.CounterName {
			counter := pb.MetricMessage_Counter{
				Counter: *metric.CounterValue,
			}
			message.Value = &counter
		}

		if metric.MetricType == models.GaugeName {
			gauge := pb.MetricMessage_Gauge{
				Gauge: *metric.GaugeValue,
			}
			message.Value = &gauge
		}

		if _, err := g.grpc.UpdateMetric(context.Background(), &message); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) CollectMetrics(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	pollInterval := time.NewTicker(c.pollInterval)
	for {
		select {
		case <-pollInterval.C:
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
		case <-ctx.Done():
			return
		}
	}
}

// CollectAdditionalMetrics collect additional metrics
func (c *Client) CollectAdditionalMetrics(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	pollInterval := time.NewTicker(c.pollInterval)
	for {
		select {
		case <-pollInterval.C:
			log.Println("CollectAdditionalMetrics")
			result := make(map[string]float64)
			info, err := cpu.Percent(time.Second*10, true)
			if err != nil {
				wg.Done()
				log.Fatal(err)
			}
			for i, cp := range info {
				name := fmt.Sprintf("CPUutilization%d", i)
				result[name] = cp
			}
			c.SetMetrics(result)
		case <-ctx.Done():
			return
		}
	}
}

func TryEncrypt(msg []byte, key *rsa.PublicKey) []byte {
	if key == nil {
		return msg
	}

	hash := sha512.New()
	result, err := rsa.EncryptOAEP(hash, cryprorand.Reader, key, msg, nil)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
