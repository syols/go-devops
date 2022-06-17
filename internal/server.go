package internal

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"github.com/syols/go-devops/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server  http.Server
	metrics store.MetricsStorage
	sets    settings.Settings
}

func NewServer(sets settings.Settings) Server {
	return Server{
		metrics: store.NewMetricsStorage(sets),
		sets:    sets,
	}
}

func (s *Server) Run() {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	s.shutdown(sign)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/value/{type}/{name}", s.valueMetricHandler)
	router.Post("/update/{type}/{name}/{value}", s.updateMetricHandler)
	router.Post("/update/", s.updateJsonMetricHandler)
	router.Post("/value/", s.valueJsonMetricHandler)

	server := http.Server{
		Addr:    s.sets.GetAddress(),
		Handler: router,
	}

	log.Printf("Server starts at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Listen error: %s", err)
	}
}

func (s *Server) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricValue := chi.URLParam(r, "value")
	metricName := chi.URLParam(r, "name")

	result, err := metric.NewMetric(metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	oldMetric, err := s.metrics.GetMetric(metricName, metricType)
	if err != nil {
		oldMetric = result
	}

	if result.TypeName() != oldMetric.TypeName() {
		http.Error(w, "wrong result type", http.StatusBadRequest)
		return
	}

	updatedMetric, err := oldMetric.FromString(metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metrics.SetMetric(metricName, updatedMetric)
}

func (s *Server) valueMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if _, err := metric.NewMetric(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	oldMetric, err := s.metrics.GetMetric(metricName, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(oldMetric.String())); err != nil {
		log.Printf("write oldMetric error: %s", err)
	}
}

func (s *Server) updateJsonMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var metricPayload metric.Payload
	if err := decoder.Decode(&metricPayload); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	newMetric, err := metric.NewMetric(metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	oldMetric, err := s.metrics.GetMetric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		oldMetric = newMetric
	}

	if newMetric.TypeName() != oldMetric.TypeName() {
		http.Error(w, "wrong metric type", http.StatusBadRequest)
		return
	}

	s.metrics.SetMetric(metricPayload.Name, oldMetric.FromPayload(metricPayload))
}

func (s *Server) valueJsonMetricHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var metricPayload metric.Payload
	if err := decoder.Decode(&metricPayload); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	_, err := metric.NewMetric(metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	oldMetric, err := s.metrics.GetMetric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(oldMetric.Payload(metricPayload.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) shutdown(sign chan os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		<-sign
		log.Print("\n\nExiting")
		if err := s.server.Shutdown(ctx); err == nil {
			log.Println("Server shutdown")
		}

		s.metrics.Save()

		os.Exit(0)
	}()
}
