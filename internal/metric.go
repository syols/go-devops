package utils

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
)

const GaugeMetric = "gauge"
const CounterMetric = "counter"
const PollCountMetricName = "PollCount"
const RandomValueMetricName = "RandomValue"
const MetricNotFoundMessage = "metric not found"
const SkipMessage = "Skip field: %s"
const IncorrectMetricValue = "incorrect metric value"

type Metric interface {
	GetValue(name string) (string, error)
	SetValue(name string, value string) error
}

type GaugeMetricsValues map[string]float64
type CounterMetricsValues map[string]uint64

func NewGaugeMetricsValues(metrics []string) GaugeMetricsValues {
	var result = make(GaugeMetricsValues)
	for _, m := range metrics {
		result[m] = 0
	}
	result[RandomValueMetricName] = 0
	return result
}

func NewCounterMetricsValues() CounterMetricsValues {
	return CounterMetricsValues{PollCountMetricName: 0}
}

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

func (g *GaugeMetricsValues) collect() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	v := reflect.ValueOf(stats)
	vsType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := vsType.Field(i).Name
		if _, ok := (*g)[field]; ok {
			value := v.Field(i).Interface()
			switch value.(type) {
			case uint64:
				(*g)[field] = float64(value.(uint64))
			case uint32:
				(*g)[field] = float64(value.(uint32))
			case float64:
				(*g)[field] = value.(float64)
			default:
				fmt.Printf(SkipMessage, field)
			}
		}
	}
}
