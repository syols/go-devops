package metric

import (
	"errors"
)

type Metric interface {
	TypeName() string

	String() string
	Payload(name string) Payload

	FromString(value string) (Metric, error)
	FromPayload(value Payload) Metric
}

type Payload struct {
	Name         string   `json:"id"`
	MetricType   string   `json:"type"`
	CounterValue *uint64  `json:"delta,omitempty"`
	GaugeValue   *float64 `json:"value,omitempty"`
}

func NewMetric(typeName string) (Metric, error) {
	for _, v := range [...]Metric{GaugeMetric(0), CounterMetric(0)} {
		if typeName == v.TypeName() {
			return v, nil
		}
	}
	return nil, errors.New("wrong metric type")
}
