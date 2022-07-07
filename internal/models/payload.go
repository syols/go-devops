package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/syols/go-devops/internal/errors"
)

type Payload struct {
	Name         string         `json:"id" db:"id"`
	MetricType   string         `json:"type" db:"metric_type"`
	CounterValue *CounterMetric `json:"delta,omitempty" db:"counter_value"`
	GaugeValue   *GaugeMetric   `json:"value,omitempty" db:"gauge_value"`
	Hash         string         `json:"hash,omitempty" db:"hash"`
}

func (p *Payload) Metric() (metric Metric) {
	if p.CounterValue != nil {
		metric = p.CounterValue
	}
	if p.GaugeValue != nil {
		metric = p.GaugeValue
	}
	return
}

func (p *Payload) String() string {
	if p.MetricType == GaugeMetric(0).TypeName() {
		return fmt.Sprintf("%s:gauge:%f", p.Name, *p.GaugeValue)
	}
	return fmt.Sprintf("%s:counter:%d", p.Name, *p.CounterValue)
}

func NewPayload(name string, key *string, value Metric) Payload {
	payload := Payload{
		Name:       name,
		MetricType: value.TypeName(),
	}

	if payload.MetricType == GaugeMetric(0).TypeName() {
		gauge := (value).(GaugeMetric)
		payload.GaugeValue = &gauge
	}

	if payload.MetricType == CounterMetric(0).TypeName() {
		counter := (value).(CounterMetric)
		payload.CounterValue = &counter
	}

	if key != nil {
		h := hmac.New(sha256.New, []byte(*key))
		hashString := payload.String()
		h.Write([]byte(hashString))
		payload.Hash = fmt.Sprintf("%x", h.Sum(nil))
	}

	return payload
}

func NewMetric(typeName string) (Metric, error) {
	for _, v := range [...]Metric{GaugeMetric(0), CounterMetric(0)} {
		if typeName == v.TypeName() {
			return v, nil
		}
	}
	return nil, errors.NewTypeNameError(typeName)
}