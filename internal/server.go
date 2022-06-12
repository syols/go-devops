package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

type Server struct {
	server   http.Server
	metrics  Metrics
	settings Settings
}

func NewServer(settings Settings) Server {
	return Server{
		metrics:  Metrics{},
		settings: settings,
	}
}

func (s *Server) Run() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post("/update/{type}/{name}/{value}", s.postMetricHandler)
	router.Get("/value/{type}/{name}", s.getMetricHandler)

	server := http.Server{
		Addr:    s.settings.GetAddress(),
		Handler: router,
	}

	log.Printf("Server starts at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Listen error: %s", err)
	}
}

func (s *Server) postMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricValue := chi.URLParam(r, "value")
	metricName := chi.URLParam(r, "name")

	newMetric, err := s.metrics.NewMetric(metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	metric, isOk := s.metrics[metricName]
	if !isOk {
		metric = newMetric
	}

	updatedMetric, err := metric.Add(metricValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.metrics[metricName] = updatedMetric
}

func (s *Server) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	if _, err := s.metrics.NewMetric(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	metricName := chi.URLParam(r, "name")
	metric, err := s.metrics.getMetric(metricName, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(metric.String())); err != nil {
		log.Printf("write metric error: %s", err)
	}
}
