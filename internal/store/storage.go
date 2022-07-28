package store

import (
	"context"
	"log"
	"time"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
)

type Store interface {
	Save(ctx context.Context, value []models.Metric) error
	Load(ctx context.Context) ([]models.Metric, error)
	Type() string
	Check() error
}

type Metrics map[string]models.Metric
type MetricsStorage struct {
	Metrics
	Store

	SaveInterval time.Duration
	Key          *string
}

func NewStore(settings config.Config) (Store, error) {
	if settings.Store.DatabaseConnectionString != nil {
		return NewDatabaseStore(*settings.Store.DatabaseConnectionString)
	}

	if settings.Store.StoreFile != nil {
		return NewFileStore(*settings.Store.StoreFile), nil
	}

	return NewFileStore("tmp.json"), nil
}

func NewMetricsStorage(settings config.Config) (MetricsStorage, error) {
	store, err := NewStore(settings)
	if err != nil {
		return MetricsStorage{}, err
	}

	metrics := MetricsStorage{
		Metrics:      make(Metrics),
		Store:        store,
		SaveInterval: settings.Store.StoreInterval,
		Key:          settings.Server.Key,
	}

	if settings.Store.Restore {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		metrics.Load(ctx)
		defer cancel()
	}

	if settings.Store.StoreInterval > 0 {
		ticker := time.NewTicker(metrics.SaveInterval)
		go func() {
			for {
				<-ticker.C
				if err := metrics.Save(context.Background()); err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
	return metrics, nil
}

func (m MetricsStorage) Load(ctx context.Context) {
	if metricsPayload, err := m.Store.Load(ctx); err == nil {
		for _, payload := range metricsPayload {
			m.Metrics[payload.Name] = payload
		}
	}
}

func (m MetricsStorage) Save(ctx context.Context) error {
	length := len(m.Metrics)
	if length == 0 {
		return nil
	}

	var result []models.Metric
	for _, v := range m.Metrics {
		result = append(result, v)
	}

	return m.Store.Save(ctx, result)
}

func (m MetricsStorage) Check() error {
	return m.Store.Check()
}
