package store

import (
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/errors"
	"github.com/syols/go-devops/internal/models"
	"log"
	"time"
)

type Store interface {
	Save(value []models.Payload) error
	Load() ([]models.Payload, error)
	Type() string
	Check() error
}

type MetricsStorage struct {
	metrics      map[string]models.Metric
	store        Store
	saveInterval time.Duration
	key          *string
}

func NewStore(settings config.Config) Store {
	if settings.Store.DatabaseConnectionString != nil {
		return NewDatabaseStore(*settings.Store.DatabaseConnectionString)
	}
	if settings.Store.StoreFile != nil {
		return NewFileStore(*settings.Store.StoreFile)
	}
	return NewFileStore("tmp.json")
}

func NewMetricsStorage(settings config.Config) MetricsStorage {
	metrics := MetricsStorage{
		metrics:      map[string]models.Metric{},
		store:        NewStore(settings),
		saveInterval: settings.Store.StoreInterval,
		key:          settings.Server.Key,
	}

	if settings.Store.Restore {
		metrics.Load()
	}

	if settings.Store.StoreInterval > 0 {
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

func (m MetricsStorage) UpdateMetric(metricName string, value models.Metric) {
	m.metrics[metricName] = value
	if m.saveInterval == 0 || m.store.Type() == "database" {
		m.Save()
	}
}

func (m MetricsStorage) Metric(metricName, metricType string) (models.Metric, error) {
	value, isOk := m.metrics[metricName]
	if !isOk {
		return nil, errors.NewValueNotFound(metricName)
	}

	if metricType != value.TypeName() {
		return nil, errors.NewValueNotFound(metricType)
	}
	return value, nil
}

func (m MetricsStorage) Load() {
	metricsPayload, err := m.store.Load()
	if err != nil {
		log.Print(err.Error())
		return
	}

	for _, payload := range metricsPayload {
		value, err := models.NewMetric(payload.MetricType)
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
	var payload []models.Payload
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
