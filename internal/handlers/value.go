package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
	"net/http"
)

func Value(metrics store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		value, isOk := metrics.Metrics[metricName]
		if !isOk {
			http.Error(w, "value not found", http.StatusNotFound)
			return
		}

		if value.MetricType != metricType {
			http.Error(w, "wrong type name", http.StatusNotImplemented)
			return
		}

		if _, err := w.Write([]byte(value.Value())); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ValueJSON(metrics store.MetricsStorage) http.HandlerFunc {
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

		if err := payload.Check(); err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		value, isOk := metrics.Metrics[payload.Name]
		if isOk {
			if value.MetricType != payload.MetricType {
				http.Error(w, "wrong type name", http.StatusNotImplemented)
				return
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
