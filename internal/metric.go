package internal

import (
	"errors"
	"fmt"
	"strconv"
)

type Metrics map[string]Metric
type GaugeMetric float64
type CounterMetric uint64

type Metric interface {
	TypeName() string
	String() string
	Add(value string) (metric Metric, err error)
}

func (m Metrics) getMetric(metricName, metricType string) (Metric, error) {
	metric, isOk := m[metricName]
	if !isOk {
		return nil, errors.New("metric not found")
	}

	if metricType != metric.TypeName() {
		return nil, errors.New("metric not found")
	}
	return metric, nil
}

func (m Metrics) NewMetric(typeName string) (Metric, error) {
	for _, v := range [...]Metric{GaugeMetric(0), CounterMetric(0)} {
		if typeName == v.TypeName() {
			return v, nil
		}
	}
	return nil, errors.New("wrong metric type")
}

func (c CounterMetric) Add(value string) (metric Metric, err error) {
	val, _err := strconv.ParseUint(value, 10, 64)
	return CounterMetric(uint64(c) + val), _err
}

func (g GaugeMetric) Add(value string) (metric Metric, err error) {
	val, _err := strconv.ParseFloat(value, 64)
	return GaugeMetric(val), _err
}

func (g GaugeMetric) TypeName() string {
	return "gauge"
}

func (c CounterMetric) TypeName() string {
	return "counter"
}

func (g GaugeMetric) String() string {
	return fmt.Sprintf("%.3f", g)
}

func (c CounterMetric) String() string {
	return strconv.FormatUint(uint64(c), 10)
}
