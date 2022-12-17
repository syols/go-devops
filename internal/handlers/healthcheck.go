package handlers

import (
	"net/http"

	"github.com/syols/go-devops/internal/store"
)

// Healthcheck godoc
// @Summary Healthcheck
// @Success 200 {object} OK
// @Router /healthcheck [get]
func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Ping godoc
// @Summary Ping
// @Success 200 {object} OK
// @Router /ping [get]
func Ping(metrics *store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := metrics.Check(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
