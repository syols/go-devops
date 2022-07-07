package store

import (
	"errors"
	"github.com/syols/go-devops/internal/model"
	"github.com/syols/go-devops/internal/settings"
	"log"
	"time"
)

type Store interface {
	Save(value []model.Payload) error
	Load() ([]model.Payload, error)
	Type() string
	Check() error
}

type MetricsStorage struct {
	metrics      map[string]model.Metric
	store        Store
	saveInterval time.Duration
	key          *string
}

func NewStore(sets settings.Config) Store {
	if sets.Store.DatabaseConnectionString != nil {
		return NewDatabaseStore(*sets.Store.DatabaseConnectionString)
	}
	if sets.Store.StoreFile != nil {
		return NewFileStore(*sets.Store.StoreFile)
	}
	return NewFileStore("tmp.json")
}

func NewMetricsStorage(sets settings.Config) MetricsStorage {
	metrics := MetricsStorage{
		metrics:      map[string]model.Metric{},
		store:        NewStore(sets),
		saveInterval: sets.Store.StoreInterval,
		key:          sets.Server.Key,
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

func (m MetricsStorage) SetMetric(metricName string, value model.Metric) {
	m.metrics[metricName] = value
	if m.saveInterval == 0 || m.store.Type() == "database" {
		m.Save()
	}
}

func (m MetricsStorage) GetMetric(metricName, metricType string) (model.Metric, error) {
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
		value, err := model.NewMetric(payload.MetricType)
		if err != nil {
			log.Print(err.Error())
		}
		m.metrics[payload.Name], err = value.FromPayload(payload, m.key)
		if err != nil {
			log.Print(err.Error())
		}
	}
}

func (m MetricsStorage) Save() {
	var payload []model.Payload
	for k, v := range m.metrics {
		payload = append(payload, v.Payload(k, m.key))
	}

	if len(payload) > 0 {
		if err := m.store.Save(payload); err != nil {
			log.Print(err.Error())
		}
	}
}

func (m MetricsStorage) Check() error {
	return m.store.Check()
}
