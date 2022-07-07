package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
	"net/http"
)

func Update(metrics store.MetricsStorage, _ *string, w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricValue := chi.URLParam(r, "value")
	metricName := chi.URLParam(r, "name")

	createdMetric, err := models.NewMetric(metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := metrics.Metric(metricName, metricType)
	if err != nil {
		currentMetric = createdMetric
	}

	if createdMetric.TypeName() != currentMetric.TypeName() {
		http.Error(w, "wrong createdMetric type", http.StatusBadRequest)
		return
	}

	updatedMetric, err := currentMetric.FromString(metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metrics.UpdateMetric(metricName, updatedMetric)
}

func UpdateJSON(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var metricPayload models.Payload
	if err := decoder.Decode(&metricPayload); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdMetric, err := models.NewMetric(metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := metrics.Metric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		currentMetric = createdMetric
	}

	if createdMetric.TypeName() != currentMetric.TypeName() {
		http.Error(w, "wrong metric type", http.StatusBadRequest)
		return
	}

	payload, err := currentMetric.FromPayload(metricPayload, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	metrics.UpdateMetric(metricPayload.Name, payload)

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(currentMetric.Payload(metricPayload.Name, key)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdatesJSON(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var metricPayloads []models.Payload
	if err := decoder.Decode(&metricPayloads); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	for _, metricPayload := range metricPayloads {
		createdMetric, err := models.NewMetric(metricPayload.MetricType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		currentMetric, err := metrics.Metric(metricPayload.Name, metricPayload.MetricType)
		if err != nil {
			currentMetric = createdMetric
		}

		if createdMetric.TypeName() != currentMetric.TypeName() {
			http.Error(w, "wrong metric type", http.StatusBadRequest)
			return
		}

		payload, err := currentMetric.FromPayload(metricPayload, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metrics.UpdateMetric(metricPayload.Name, payload)
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(metricPayloads[0]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
