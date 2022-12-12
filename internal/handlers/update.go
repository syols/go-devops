package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/syols/go-devops/internal/models"
	"github.com/syols/go-devops/internal/store"
)

// Update godoc
// @Tags Update
// @Summary Update metric
// @Description Update metric request
// @Produce json
// @Param type path int64 true "Metric type"
// @Param value path float64 true "Metric value"
// @Param name path int64 true "Metric name"
// @Success 200 {object} Metric
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /update/{type}/{name}/{value} [post]
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

// UpdateJSON godoc
// @Tags UpdateJSON
// @Summary Update metric
// @Description Update metric request
// @Produce json
// @Param Metric body Metric true "Metric"
// @Success 200 {object} Metric
// @Failure 415 {string} string "StatusUnsupportedMediaType"
// @Failure 422 {string} string "StatusUnprocessableEntity"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /update/ [post]
func UpdateJSON(metrics store.MetricsStorage, key *string) http.HandlerFunc {
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

// UpdatesJSON godoc
// @Tags UpdatesJSON
// @Summary Update metrics
// @Description Update metrics request
// @Produce json
// @Param []Metric body []Metric true "Metric list"
// @Success 200 {object} Metric
// @Failure 415 {string} string "StatusUnsupportedMediaType"
// @Failure 422 {string} string "StatusUnprocessableEntity"
// @Failure 500 {string} string "StatusInternalServerError"
// @Router /updates/ [post]
func UpdatesJSON(metrics store.MetricsStorage, key *string, privateKey *rsa.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != ContentType {
			http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
			return
		}

		buf, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		buf, err = tryDecrypt(buf, privateKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		var payloads []models.Metric
		if err = json.Unmarshal(buf, &payloads); err != nil {
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
		if payload.MetricType == models.CounterName {
			*payload.CounterValue += *oldPayload.CounterValue
		}
	}

	payload.Hash = payload.CalculateHash(key)
	metrics.Metrics[payload.Name] = payload
	return true
}

func tryDecrypt(msg []byte, key *rsa.PrivateKey) ([]byte, error) {
	if key == nil {
		return msg, nil
	}

	hash := sha512.New()
	result, err := rsa.DecryptOAEP(hash, rand.Reader, key, msg, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
