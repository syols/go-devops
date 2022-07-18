package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
	"net/http"
)

func Update(metrics store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricValue := chi.URLParam(r, "value")
		metricName := chi.URLParam(r, "name")

		payload := models.NewMetric(metricName, metricType, metricValue, nil)
		if !update(w, payload, nil, metrics) {
			return
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateJSON(metrics store.MetricsStorage, key *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var payload models.Metric
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if !update(w, payload, key, metrics) {
			return
		}
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(payload); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdatesJSON(metrics store.MetricsStorage, key *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var payloads []models.Metric
		if err := decoder.Decode(&payloads); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		for _, payload := range payloads {
			if !update(w, payload, key, metrics) {
				return
			}
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(payloads[0]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func update(w http.ResponseWriter, payload models.Metric, key *string, metrics store.MetricsStorage) bool {
	err := payload.Check()
	if err, ok := err.(validator.ValidationErrors); ok {
		if err[0].Tag() == "metric" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return false
		}
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return false
	}

	hash := payload.CalculateHash(key)
	if payload.Hash != hash {
		http.Error(w, "wrong hash sum", http.StatusBadRequest)
		return false
	}

	oldPayload, isOk := metrics.Metrics[payload.Name]
	if isOk {
		if payload.MetricType != oldPayload.MetricType {
			http.Error(w, "wrong type name", http.StatusNotImplemented)
			return false
		}
		if payload.MetricType == "counter" {
			*payload.CounterValue += *oldPayload.CounterValue
		}
	}

	payload.Hash = payload.CalculateHash(key)
	metrics.Metrics[payload.Name] = payload
	return true
}
