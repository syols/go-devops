package internal

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

const ListenError = "Listen error: %s"
const StartListening = "Server starts at %s"
const MetricNotFound = "metric type not found"

const ShutdownError = "Error at server shutdown %v"
const WriteError = "Write error: %s"
const UpdateMetricRoute = "/update/{type}/{name}/{value}"
const ValueMetricRoute = "/value/{type}/{name}"
const TypeParam = "type"
const NameParam = "name"
const ValueParam = "value"

type Server struct {
	server         http.Server
	address        string
	gaugeMetrics   GaugeMetricsValues
	counterMetrics CounterMetricsValues
	isCheck        bool
}

func NewServer(settings Settings) Server {
	return Server{
		address:        settings.GetAddress(),
		gaugeMetrics:   GaugeMetricsValues{},
		counterMetrics: CounterMetricsValues{},
		isCheck:        false,
	}
}

func (s *Server) Run() context.CancelFunc {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Post(UpdateMetricRoute, s.postMetricHandler)
	router.Get(ValueMetricRoute, s.getMetricHandler)

	server := http.Server{
		Addr:    s.address,
		Handler: router,
	}

	ctx, cancel := context.WithCancel(context.Background())
	log.Printf(StartListening, s.address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf(ListenError, err)
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf(ShutdownError, err)
	}
	return cancel
}

func (s *Server) getMetric(metricType string) (Metric, error) {
	switch metricType {
	case GaugeMetric:
		return &s.gaugeMetrics, nil
	case CounterMetric:
		return &s.counterMetrics, nil
	default:
		return nil, errors.New(MetricNotFound)
	}
}

func (s *Server) postMetricHandler(w http.ResponseWriter, r *http.Request) {
	if metric, getMetricError := s.getMetric(chi.URLParam(r, TypeParam)); getMetricError != nil {
		http.Error(w, getMetricError.Error(), http.StatusNotImplemented)
	} else {
		name := chi.URLParam(r, NameParam)
		value := chi.URLParam(r, ValueParam)
		if setError := metric.SetValue(name, value); setError != nil {
			http.Error(w, setError.Error(), http.StatusBadRequest)
		}
		if _, err := w.Write([]byte(value)); err != nil {
			log.Printf(WriteError, err)
		}
	}
}

func (s *Server) getMetricHandler(w http.ResponseWriter, r *http.Request) {
	if metric, getMetricError := s.getMetric(chi.URLParam(r, TypeParam)); getMetricError == nil {
		if value, getError := metric.GetValue(chi.URLParam(r, NameParam)); getError != nil {
			http.Error(w, getError.Error(), http.StatusNotFound)
		} else if _, err := w.Write([]byte(value)); err != nil {
			log.Printf(WriteError, err)
		}
	}
}
