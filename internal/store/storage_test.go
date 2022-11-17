package store

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/models"
	mock_store "github.com/syols/go-devops/mock"
)

func newStore(t *testing.T) ([]models.Metric, context.Context, *mock_store.MockStore) {
	value := 1.1
	metric := models.Metric{
		Name:       "testGauge",
		MetricType: models.GaugeName,
		GaugeValue: &value,
	}
	metrics := []models.Metric{metric}
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_store.NewMockStore(ctrl)
	return metrics, ctx, store
}

func TestNewStorage(t *testing.T) {
	cfg := config.NewConfig()
	ctx := context.Background()
	storage, err := NewMetricsStorage(cfg)
	model := models.Metric{
		Name: "41",
	}
	storage.Metrics["some"] = model
	storage.Load(ctx)
	assert.NoError(t, storage.Save(ctx))
	assert.NoError(t, err)
}

func TestLoadStorage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock_store.NewMockStore(ctrl)
	store.EXPECT().Load(ctx).AnyTimes()
	metrics := MetricsStorage{
		Metrics: make(Metrics),
		Store:   store,
	}
	metrics.Load(ctx)
}

func TestSaveStorage(t *testing.T) {
	metrics, ctx, store := newStore(t)
	store.EXPECT().Save(ctx, metrics).AnyTimes()
	err := store.Save(ctx, metrics)
	assert.NoError(t, err)
}
