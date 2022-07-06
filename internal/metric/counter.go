package metric

import (
	"errors"
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

func (c CounterMetric) FromPayload(value Payload, key *string) (Metric, error) {
	if value.CounterValue.TypeName() != value.MetricType {
		return value.CounterValue, errors.New("wrong type name")
	}

	result := CounterMetric(uint64(c) + uint64(*value.CounterValue))

	payload := value.CounterValue.Payload(value.Name, key)
	if payload.Hash != value.Hash {
		return result, errors.New("wrong hash sum")
	}
	return result, nil
}

func (c CounterMetric) Payload(name string, key *string) Payload {
	return NewPayload(name, key, c)
}
