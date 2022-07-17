package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"strconv"
)

var validate *validator.Validate

type Metric struct {
	Name         string   `json:"id" db:"id" validate:"required"`
	MetricType   string   `json:"type" db:"metric_type" validate:"required_with='counter' 'gauge',metricValidation"`
	CounterValue *uint64  `json:"delta,omitempty" db:"counter_value"`
	GaugeValue   *float64 `json:"value,omitempty" db:"gauge_value"`
	Hash         string   `json:"hash,omitempty" db:"hash"`
}

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("metricValidation", metricValidation); err != nil {
		log.Println(err)
	}
}

func metricValidation(fl validator.FieldLevel) bool {
	metric, ok := fl.Parent().Interface().(Metric)
	if !ok {
		return false
	}

	return (metric.CounterValue != nil && metric.MetricType == "counter") ||
		(metric.GaugeValue != nil && metric.MetricType == "gauge")
}

func (p *Metric) String() string {
	if p.MetricType == "gauge" {
		return fmt.Sprintf("%s:gauge:%f", p.Name, *p.GaugeValue)
	}
	return fmt.Sprintf("%s:counter:%d", p.Name, *p.CounterValue)
}

func (p *Metric) Check() error {
	return validate.Struct(p)
}

func NewMetric(name, typeName, value string, key *string) (Metric, error) {
	payload := Metric{
		Name:       name,
		MetricType: typeName,
	}

	if payload.MetricType == "gauge" {
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			payload.GaugeValue = &v
		}
	}

	if payload.MetricType == "counter" {
		v, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			payload.CounterValue = &v
		}
	}

	payload.Hash = payload.CalculateHash(key)
	return payload, payload.Check()
}

func (p *Metric) CalculateHash(key *string) (result string) {
	if key != nil {
		h := hmac.New(sha256.New, []byte(*key))
		hashString := p.String()
		h.Write([]byte(hashString))
		result = fmt.Sprintf("%x", h.Sum(nil))
	}
	return
}
