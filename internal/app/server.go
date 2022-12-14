package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/handlers"
	"github.com/syols/go-devops/internal/models"
	pb "github.com/syols/go-devops/internal/rpc/proto"
	"github.com/syols/go-devops/internal/store"
)

// Server struct
type Server struct {
	pb.UnimplementedGoDevopsServer

	server     http.Server
	metrics    *store.MetricsStorage
	settings   config.Config
	privateKey *rsa.PrivateKey
}

// NewServer creates server struct
func NewServer(settings config.Config) (Server, error) {
	metrics, err := store.NewMetricsStorage(settings)
	if err != nil {
		return Server{}, err
	}

	var privateKey *rsa.PrivateKey
	if settings.Store.CryptoKeyFilePath != nil {
		byteArr, err := os.ReadFile(*settings.Store.CryptoKeyFilePath) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		blocks, _ := pem.Decode(byteArr)
		if blocks != nil {
			log.Fatal("Error at decode")
		}

		privateKey, err = x509.ParsePKCS1PrivateKey(blocks.Bytes)
		if err != nil {
			fmt.Print(err)
		}
	}

	return Server{
		metrics:    &metrics,
		settings:   settings,
		privateKey: privateKey,
	}, nil
}

// Run server
func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()
	s.shutdown(ctx)

	server := http.Server{
		Addr:    s.settings.Server.Address.String(),
		Handler: s.router(),
	}

	go s.runRpc()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func (s *Server) runRpc() {
	if s.settings.Grpc == nil {
		return
	}

	port := fmt.Sprintf(":%d", s.settings.Grpc.Address.Port)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGoDevopsServer(grpcServer, s)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(handlers.Compress)
	if s.settings.Server.TrustedSubnet != nil {
		router.Use(handlers.CheckSubnet(s.settings.Server.TrustedSubnet))
	}
	router.Use(handlers.Save(s.metrics))

	router.Get("/", handlers.Healthcheck)
	router.Get("/ping", handlers.Ping(s.metrics))
	router.Get("/value/{type}/{name}", handlers.Value(s.metrics))
	router.Post("/update/{type}/{name}/{value}", handlers.Update(s.metrics))
	router.Post("/update/", handlers.UpdateJSON(s.metrics, s.settings.Server.Key))
	router.Post("/updates/", handlers.UpdatesJSON(s.metrics, s.settings.Server.Key, s.privateKey))
	router.Post("/value/", handlers.ValueJSON(s.metrics))

	// pprof
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	router.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return router
}

func (s *Server) shutdown(ctx context.Context) {
	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		if s.settings.Store.DatabaseConnectionString == nil {
			err := s.metrics.Save(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}
		os.Exit(0)
	}()
}

func (s *Server) UpdateMetric(ctx context.Context, in *pb.MetricMessage) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	metric := models.Metric{
		Name:         in.Name,
		MetricType:   in.Type,
		CounterValue: in.Counter,
		GaugeValue:   in.Gauge,
	}
	err := metric.Check()
	if err, ok := err.(validator.ValidationErrors); ok {
		response.Error = err.Error()
		return &response, nil
	}

	hash := metric.CalculateHash(s.settings.Server.Key)
	if metric.Hash != hash {
		response.Error = "wrong hash sum"
		return &response, nil
	}
	value, isOk := s.metrics.Metrics.Load(metric.Name)
	oldPayload := value.(models.Metric)
	if isOk {
		if metric.MetricType != oldPayload.MetricType {
			response.Error = "wrong type name"
			return &response, nil
		}
		if metric.MetricType == models.CounterName {
			*metric.CounterValue += *oldPayload.CounterValue
		}
	}

	metric.Hash = metric.CalculateHash(s.settings.Server.Key)
	s.metrics.Metrics.Store(metric.Name, metric)
	response.Metric = &pb.MetricMessage{
		Name:    metric.Name,
		Type:    metric.MetricType,
		Counter: metric.CounterValue,
		Gauge:   metric.GaugeValue,
	}
	return &response, nil
}
