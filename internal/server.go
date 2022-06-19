package internal

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/syols/go-devops/internal/metric"
	"github.com/syols/go-devops/internal/settings"
	"github.com/syols/go-devops/internal/store"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Server struct {
	server  http.Server
	metrics store.MetricsStorage
	sets    settings.Settings
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func NewServer(sets settings.Settings) Server {
	return Server{
		metrics: store.NewMetricsStorage(sets),
		sets:    sets,
	}
}

func (s *Server) Run() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	s.shutdown(sign)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(compressMiddleware)

	router.Get("/", s.healthcheckHandler)
	router.Get("/value/{type}/{name}", s.valueMetricHandler)
	router.Post("/update/{type}/{name}/{value}", s.updateMetricHandler)
	router.Post("/update/", s.updateJSONMetricHandler)
	router.Post("/value/", s.valueJSONMetricHandler)

	server := http.Server{
		Addr:    s.sets.GetAddress(),
		Handler: router,
	}

	log.Printf("Server starts at %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Listen error: %s", err)
	}
}

func (s *Server) updateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricValue := chi.URLParam(r, "value")
	metricName := chi.URLParam(r, "name")

	createdMetric, err := metric.NewMetric(metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := s.metrics.GetMetric(metricName, metricType)
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
	s.metrics.SetMetric(metricName, updatedMetric)
}

func (s *Server) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("write error: %s", err)
	}
}

func (s *Server) valueMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "type")
	metricName := chi.URLParam(r, "name")

	if _, err := metric.NewMetric(metricType); err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := s.metrics.GetMetric(metricName, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(currentMetric.String())); err != nil {
		log.Printf("write currentMetric error: %s", err)
	}
}

func (s *Server) updateJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
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

	createdMetric, err := metric.NewMetric(metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}

	currentMetric, err := s.metrics.GetMetric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		currentMetric = createdMetric
	}

	if createdMetric.TypeName() != currentMetric.TypeName() {
		http.Error(w, "wrong metric type", http.StatusBadRequest)
		return
	}
	s.metrics.SetMetric(metricPayload.Name, currentMetric.FromPayload(metricPayload))

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(currentMetric.Payload(metricPayload.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) valueJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
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

	currentMetric, err := s.metrics.GetMetric(metricPayload.Name, metricPayload.MetricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(currentMetric.Payload(metricPayload.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func compressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Print(err.Error())
			return
		}
		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Print(err.Error())
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func (s *Server) shutdown(sign chan os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		<-sign
		log.Println("Exiting")
		if err := s.server.Shutdown(ctx); err == nil {
			log.Println("Server shutdown")
		}
		s.metrics.Save()
		os.Exit(0)
	}()
}
