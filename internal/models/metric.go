package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
)

const GaugeName = "gauge"
const CounterName = "counter"

var validate *validator.Validate

type Metric struct {
	Name         string   `json:"id" db:"id" validate:"required"`
	MetricType   string   `json:"type" db:"metric_type" validate:"metricType,metric"`
	CounterValue *uint64  `json:"delta,omitempty" db:"counter_value" validate:"omitempty,metricCounter"`
	GaugeValue   *float64 `json:"value,omitempty" db:"gauge_value" validate:"omitempty,metricGauge"`
	Hash         string   `json:"hash,omitempty" db:"hash"`
}

func init() {
	validate = validator.New()
	if err := validate.RegisterValidation("metric", metricValidation); err != nil {
		log.Fatal(err)
	}
	if err := validate.RegisterValidation("metricType", metricTypeValidation); err != nil {
		log.Fatal(err)
	}
	if err := validate.RegisterValidation("metricGauge", metricGaugeValidation); err != nil {
		log.Fatal(err)
	}
	if err := validate.RegisterValidation("metricCounter", metricCounterValidation); err != nil {
		log.Fatal(err)
	}
}

func metricGaugeValidation(fl validator.FieldLevel) bool {
	if metric, ok := fl.Parent().Interface().(Metric); ok {
		return metric.MetricType == GaugeName && metric.GaugeValue != nil
	}
	return false
}

func metricCounterValidation(fl validator.FieldLevel) bool {
	if metric, ok := fl.Parent().Interface().(Metric); ok {
		return metric.MetricType == CounterName && metric.CounterValue != nil
	}
	return false
}

func metricTypeValidation(fl validator.FieldLevel) bool {
	if metric, ok := fl.Parent().Interface().(Metric); ok {
		return metric.MetricType == CounterName || metric.MetricType == GaugeName
	}
	return false
}

func metricValidation(fl validator.FieldLevel) bool {
	if metric, ok := fl.Parent().Interface().(Metric); ok {
		return metric.CounterValue != nil || metric.GaugeValue != nil
	}
	return false
}

func (p *Metric) String() string {
	if p.MetricType == GaugeName {
		return fmt.Sprintf("%s:gauge:%f", p.Name, *p.GaugeValue)
	}
	return fmt.Sprintf("%s:counter:%d", p.Name, *p.CounterValue)
}

func (p *Metric) Value() string {
	if p.MetricType == GaugeName {
		return fmt.Sprintf("%.3f", *p.GaugeValue)
	}
	return strconv.FormatUint(*p.CounterValue, 10)
}

func (p *Metric) Check() error {
	return validate.Struct(p)
}

func NewMetric(name, typeName, value string, key *string) Metric {
	payload := Metric{
		Name:       name,
		MetricType: typeName,
	}

	if payload.MetricType == GaugeName {
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			payload.GaugeValue = &v
		}
	}

	if payload.MetricType == CounterName {
		v, err := strconv.ParseUint(value, 10, 64)
		if err == nil {
			payload.CounterValue = &v
		}
	}

	payload.Hash = payload.CalculateHash(key)
	return payload
}

func (p *Metric) CalculateHash(key *string) (result string) {
	if key != nil {
		h := hmac.New(sha256.New, []byte(*key))
		h.Write([]byte(p.String()))
		result = fmt.Sprintf("%x", h.Sum(nil))
	}
	return
}
