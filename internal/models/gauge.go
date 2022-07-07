package models

import (
	"fmt"
	"github.com/syols/go-devops/internal/errors"
	"strconv"
)

type GaugeMetric float64

func (g GaugeMetric) TypeName() string {
	return "gauge"
}

func (g GaugeMetric) String() string {
	return fmt.Sprintf("%.3f", g)
}

func (g GaugeMetric) FromString(value string) (Metric, error) {
	val, _err := strconv.ParseFloat(value, 64)
	return GaugeMetric(val), _err
}

func (g GaugeMetric) FromPayload(value Payload, key *string) (Metric, error) {
	if value.GaugeValue.TypeName() != value.MetricType {
		return value.GaugeValue, errors.NewTypeNameError(value.MetricType)
	}

	payload := value.GaugeValue.Payload(value.Name, key)
	if payload.Hash != value.Hash {
		return value.GaugeValue, errors.NewHashSumError(value.MetricType)
	}
	return value.GaugeValue, nil
}

func (g GaugeMetric) Payload(name string, key *string) Payload {
	return NewPayload(name, key, g)
}
