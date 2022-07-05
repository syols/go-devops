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
	if g.Payload(value.Name, key).Hash != value.Hash {
		return nil, errors.New("wrong hash sum")
	}
	return value.GaugeValue, nil
}

func (g GaugeMetric) Payload(name string, key *string) Payload {
	return NewPayload(name, key, g)
}
