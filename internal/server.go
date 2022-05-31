package server

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/syols/go-devops/internal/utils"
	"log"
	"net/http"
)

const ListenError = "Write error: %s"
const StartListening = "Server starts at %s"
const MetricNotFound = "metric type not found"

type Server struct {
	address        string
	gaugeMetrics   utils.GaugeMetricsValues
	counterMetrics utils.CounterMetricsValues
	isCheck        bool
}

func NewServer(settings utils.Settings) Server {
	return Server{
		address:        settings.GetAddress(),
		gaugeMetrics:   utils.NewGaugeMetricsValues(settings.Metrics.RuntimeMetrics),
		counterMetrics: utils.NewCounterMetricsValues(),
		isCheck:        false,
	}
}

func (s *Server) Run() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	s.routing(router)
	log.Printf(StartListening, s.address)
	if err := http.ListenAndServe(s.address, router); err != nil {
		log.Fatalf(ListenError, err)
	}
}

func (s *Server) getMetric(metricType string) (utils.Metric, error) {
	switch metricType {
	case utils.GaugeMetric:
		return &s.gaugeMetrics, nil
	case utils.CounterMetric:
		return &s.counterMetrics, nil
	default:
		return nil, errors.New(MetricNotFound)
	}
}
