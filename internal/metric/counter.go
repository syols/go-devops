package metric

import (
	"strconv"
)

type CounterMetric uint64

func (c CounterMetric) TypeName() string {
	return "counter"
}

func (c CounterMetric) String() string {
	return strconv.FormatUint(uint64(c), 10)
}

func (c CounterMetric) FromString(value string) (Metric, error) {
	val, _err := strconv.ParseUint(value, 10, 64)
	return CounterMetric(uint64(c) + val), _err
}

func (c CounterMetric) FromPayload(value Payload) Metric {
	return CounterMetric(uint64(c) + *value.CounterValue)
}

func (c CounterMetric) Payload(name string) Payload {
	value := uint64(c)
	return Payload{
		Name:         name,
		MetricType:   c.TypeName(),
		CounterValue: &value,
	}
}
