package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
	"net/http"
	"time"
)

func Update(metrics store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricValue := chi.URLParam(r, "value")
		metricName := chi.URLParam(r, "name")

		payload, err := models.NewMetric(metricName, metricType, metricValue, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		oldPayload, isOk := metrics.Metrics[metricName]
		if isOk {
			if metricType != oldPayload.MetricType {
				http.Error(w, "wrong type name", http.StatusNotImplemented)
				return
			}
			if payload.MetricType == "counter" {
				*payload.CounterValue += *oldPayload.CounterValue
			}
		}

		metrics.Metrics[payload.Name] = payload
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

		err := payload.Check()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		calc := payload.CalculateHash(key)
		fmt.Println(calc)
		if payload.Hash != payload.CalculateHash(key) {
			http.Error(w, "wrong hash sum", http.StatusBadRequest)
		}

		oldPayload, isOk := metrics.Metrics[payload.Name]
		if isOk {
			if payload.MetricType != oldPayload.MetricType {
				http.Error(w, "wrong type name", http.StatusNotImplemented)
				return
			}
			if payload.MetricType == "counter" {
				*payload.CounterValue += *oldPayload.CounterValue
			}
		}

		metrics.Metrics[payload.Name] = payload
		if metrics.SaveInterval == 0 || metrics.Store.Type() == "database" {
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			err := metrics.Save(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
			err := payload.Check()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if payload.Hash != payload.CalculateHash(key) {
				http.Error(w, "wrong hash sum", http.StatusBadRequest)
			}

			oldPayload, isOk := metrics.Metrics[payload.Name]
			if isOk {
				if payload.MetricType != oldPayload.MetricType {
					http.Error(w, "wrong type name", http.StatusNotImplemented)
					return
				}
				if payload.MetricType == "counter" {
					*payload.CounterValue += *oldPayload.CounterValue
				}
			}
			metrics.Metrics[payload.Name] = payload
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(payloads[0]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
