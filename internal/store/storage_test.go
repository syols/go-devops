package store

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/syols/go-devops/internal/models"
	mock_store "github.com/syols/go-devops/mock"
)

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
	store.EXPECT().Save(ctx, metrics).AnyTimes()
	store.Save(ctx, metrics)
}
