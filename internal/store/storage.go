package store

import (
	"errors"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"log"
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
		metrics.LoadMetrics()
	}

	if sets.Store.StoreInterval > 0 {
		ticker := time.NewTicker(metrics.saveInterval)
		go func() {
			for {
				<-ticker.C
				metrics.Save()
			}
		}()
	}
	return metrics
}

func (m MetricsStorage) SetMetric(metricName string, value metric.Metric) {
	m.metrics[metricName] = value
	if m.saveInterval == 0 {
		m.Save()
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

func (m MetricsStorage) LoadMetrics() {
	metricsPayload, err := m.store.Load()
	if err != nil {
		log.Print(err.Error())
		return
	}

	for _, payload := range metricsPayload {
		value, err := metric.NewMetric(payload.MetricType)
		if err != nil {
			log.Print(err.Error())
		}
		m.metrics[payload.Name] = value.FromPayload(payload)
	}
}

func (m MetricsStorage) Save() {
	var payload []metric.Payload
	for k, v := range m.metrics {
		payload = append(payload, v.Payload(k))
	}
	if err := m.store.Save(payload); err != nil {
		log.Print(err.Error())
	}
}
