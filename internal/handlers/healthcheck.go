package handlers

import (
	"github.com/syols/go-devops/internal/store"
	"log"
	"net/http"
)

func Healthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Print(err.Error())
	}
}

func Ping(metrics store.MetricsStorage, _ *string, w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	err := metrics.Check()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
