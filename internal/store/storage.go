package store

import (
	"errors"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Store interface {
	Save(value []metric.Payload) error
	Load() ([]metric.Payload, error)
}

type MetricsStorage struct {
	metrics      map[string]metric.Metric
	store        Store
	saveInterval time.Duration
}

func NewStore(sets settings.Settings) Store {
	if sets.Store.StoreFile != nil {
		return NewFileStore(*sets.Store.StoreFile)
	}
	return NewFileStore("tmp.json")
}

func NewMetricsStorage(sets settings.Settings) MetricsStorage {
	metrics := MetricsStorage{
		metrics:      map[string]metric.Metric{},
		store:        NewStore(sets),
		saveInterval: sets.Store.StoreInterval,
	}

	if sets.Store.Restore {
		if err := metrics.LoadMetrics(); err != nil {
			log.Printf(err.Error())
		}
	}

	if sets.Store.StoreInterval > 0 {
		ticker := time.NewTicker(metrics.saveInterval)
		go func() {
			for {
				<-ticker.C
				err := metrics.SaveMetrics()
				if err != nil {
					log.Printf(err.Error())
				}
			}
		}()
	}
	metrics.onExit()
	return metrics
}

func (m MetricsStorage) SetMetric(metricName string, value metric.Metric) {
	m.metrics[metricName] = value

	if m.saveInterval == 0 {
		if err := m.SaveMetrics(); err != nil {
			log.Printf(err.Error())
		}
	}
}

func (m MetricsStorage) GetMetric(metricName, metricType string) (metric.Metric, error) {
	value, isOk := m.metrics[metricName]
	if !isOk {
		return nil, errors.New("value not found, wrong metric name")
	}

	if metricType != value.TypeName() {
		return nil, errors.New("value not found, wrong metric type")
	}
	return value, nil
}

func (m MetricsStorage) LoadMetrics() error {
	metricsPayload, err := m.store.Load()
	if err != nil {
		log.Fatalf(err.Error())
	}

	for _, payload := range metricsPayload {
		value, err := metric.NewMetric(payload.MetricType)
		if err != nil {
			log.Fatalf(err.Error())
		}
		m.metrics[payload.Name] = value.FromPayload(payload)
	}
	return err
}

func (m MetricsStorage) SaveMetrics() error {
	var payload []metric.Payload
	for k, v := range m.metrics {
		payload = append(payload, v.Payload(k))
	}

	return m.store.Save(payload)
}

func (m MetricsStorage) onExit() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sign
		log.Print("Exiting...")
		if err := m.SaveMetrics(); err != nil {
			log.Printf(err.Error())
		}
		os.Exit(0)
	}()
}
