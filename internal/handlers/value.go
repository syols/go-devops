package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/syols/go-devops/internal/model"
	"github.com/syols/go-devops/internal/store"
	"log"
	"net/http"
)

func Value(metrics store.MetricsStorage, _ *string, w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if _, err := model.NewMetric(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := metrics.GetMetric(metricName, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(currentMetric.String())); err != nil {
		log.Print(err.Error())
	}
}

func ValueJSON(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var metricPayload model.Payload
	if err := decoder.Decode(&metricPayload); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	_, err := model.NewMetric(metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := metrics.GetMetric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(currentMetric.Payload(metricPayload.Name, key)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
