package internal

import (
	"errors"
	"fmt"
	"strconv"
)

const GaugeMetric = "gauge"
const CounterMetric = "counter"
const PollCountMetricName = "PollCount"
const RandomValueMetricName = "RandomValue"
const MetricNotFoundMessage = "metric not found"
const IncorrectMetricValue = "incorrect metric value"

type Metric interface {
	GetValue(name string) (string, error)
	SetValue(name string, value string) error
}

type GaugeMetricsValues map[string]float64
type CounterMetricsValues map[string]uint64

func (g *GaugeMetricsValues) GetValue(name string) (string, error) {
	if value, ok := (*g)[name]; !ok {
		return "", errors.New(MetricNotFoundMessage)
	} else {
		return fmt.Sprintf("%.3f", value), nil
	}
}

func (g *GaugeMetricsValues) SetValue(name string, value string) error {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return errors.New(IncorrectMetricValue)
	}
	(*g)[name] = val
	return nil
}

func (g *CounterMetricsValues) GetValue(name string) (string, error) {
	if value, ok := (*g)[name]; !ok {
		return "", errors.New(MetricNotFoundMessage)
	} else {
		return strconv.FormatUint(value, 10), nil
	}
}

func (g *CounterMetricsValues) SetValue(name string, value string) error {
	val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return errors.New(IncorrectMetricValue)
	}

	if _, ok := (*g)[name]; !ok {
		(*g)[name] = val
	} else {
		(*g)[name] += val
	}
	return nil
}
