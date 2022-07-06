package metric

import (
	"errors"
	"fmt"
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
		return value.GaugeValue, errors.New("wrong type name")
	}

	payload := g.Payload(value.Name, key)
	if payload.Hash != value.Hash {
		return value.GaugeValue, errors.New("wrong hash sum")
	}
	return value.GaugeValue, nil
}

func (g GaugeMetric) Payload(name string, key *string) Payload {
	return NewPayload(name, key, g)
}
