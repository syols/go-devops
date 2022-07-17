package handlers

import (
	"github.com/syols/go-devops/internal/store"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("OK")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Ping(metrics store.MetricsStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if err := metrics.Check(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
