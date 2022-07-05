package metric

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

type Metric interface {
	TypeName() string

	String() string
	Payload(name string, key *string) Payload

	FromString(value string) (Metric, error)
	FromPayload(value Payload, key *string) (Metric, error)
}

type Payload struct {
	Name         string         `json:"id"`
	MetricType   string         `json:"type"`
	CounterValue *CounterMetric `json:"delta,omitempty"`
	GaugeValue   *GaugeMetric   `json:"value,omitempty"`
	Hash         string         `json:"hash,omitempty"`
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
	return fmt.Sprintf("%s:%s:%d", p.Name, p.MetricType, p.Metric())
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
		h.Write([]byte(payload.String()))
		payload.Hash = string(h.Sum(nil))
	}

	return payload
}

func NewMetric(typeName string) (Metric, error) {
	for _, v := range [...]Metric{GaugeMetric(0), CounterMetric(0)} {
		if typeName == v.TypeName() {
			return v, nil
		}
	}
	return nil, errors.New("wrong metric type")
}