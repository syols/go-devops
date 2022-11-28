package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
)

// Value godoc
// @Tags Value
// @Summary Request metric
// @Produce json
// @Param type path int64 true "Metric type"
// @Param name path int64 true "Metric name"
// @Success 200 {object} Metric
// @Failure 400 {string} string "StatusNotFound"
// @Failure 403 {string} string "StatusNotImplemented"
// @Router /value/{type}/{name} [get]
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

// ValueJSON godoc
// @Tags ValueJSON
// @Summary Get metric
// @Description Get metric request
// @Produce json
// @Param Metric body Metric true "Metric"
// @Success 200 {object} Metric
// @Failure 415 {string} string "StatusUnsupportedMediaType"
// @Failure 422 {string} string "StatusUnprocessableEntity"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /value/ [post]
func ValueJSON(metrics store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != ContentType {
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
		if err, ok := err.(validator.ValidationErrors); ok {
			if err[0].Tag() == "metricType" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		value, isOk := metrics.Metrics[payload.Name]
		if !isOk {
			http.Error(w, "value not found", http.StatusNotFound)
			return
		}

		if value.MetricType != payload.MetricType {
			http.Error(w, "wrong type name", http.StatusNotImplemented)
			return
		}

		value.Hash = value.CalculateHash(metrics.Key)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
