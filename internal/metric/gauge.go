package metric

import (
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

func (g GaugeMetric) FromPayload(value Payload) Metric {
	return GaugeMetric(*value.GaugeValue)
}

func (g GaugeMetric) Payload(name string) Payload {
	value := float64(g)
	return Payload{
		Name:       name,
		MetricType: g.TypeName(),
		GaugeValue: &value,
	}
}
