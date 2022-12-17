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
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/handlers"
	"github.com/syols/go-devops/internal/models"
	pb "github.com/syols/go-devops/internal/rpc/proto"
	"github.com/syols/go-devops/internal/store"
)

// Server struct
type Server struct {
	pb.UnimplementedGoDevopsServer

	httpServer http.Server
	grpcServer *grpc.Server
	metrics    *store.MetricsStorage
	settings   config.Config
	privateKey *rsa.PrivateKey
}

// NewServer creates httpServer struct
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
		metrics:    metrics,
		settings:   settings,
		privateKey: privateKey,
	}, nil
}

// Run httpServer
func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s.shutdown(ctx)

	s.httpServer = http.Server{
		Addr:    s.settings.Server.Address.String(),
		Handler: s.router(),
	}

	var wg sync.WaitGroup
	if s.settings.Grpc != nil {
		var opts []grpc.ServerOption
		s.grpcServer = grpc.NewServer(opts...)
		pb.RegisterGoDevopsServer(s.grpcServer, s)
		wg.Add(1)
		go s.runGrpc()
	}

	wg.Add(1)
	go s.runHTTPServer()
	wg.Wait()
}

func (s *Server) runHTTPServer() {
	log.Println("runHTTPServer")
	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) runGrpc() {
	log.Println("runGrpc")
	port := fmt.Sprintf(":%d", s.settings.Grpc.Address.Port)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.grpcServer.Serve(listen); err != nil {
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
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}

		if s.grpcServer != nil {
			s.grpcServer.GracefulStop()
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
	metric := models.Metric{
		Name: in.Name,
	}
	value := in.GetValue()
	switch value.(type) {
	case *pb.MetricMessage_Gauge:
		gauge := in.GetGauge()
		metric.MetricType = models.GaugeName
		metric.GaugeValue = &gauge
	case *pb.MetricMessage_Counter:
		counter := in.GetCounter()
		metric.MetricType = models.CounterName
		metric.CounterValue = &counter
	}

	var response pb.UpdateMetricResponse

	err := metric.Check()
	if err, ok := err.(validator.ValidationErrors); ok {
		msg := fmt.Sprintf("validation error: %s", err.Error())
		return nil, status.Error(codes.InvalidArgument, msg)
	}

	hash := metric.CalculateHash(s.settings.Server.Key)
	if metric.Hash != hash {
		return nil, status.Error(codes.InvalidArgument, "wrong hash sum")
	}

	loadedValue, isOk := s.metrics.Metrics.Load(metric.Name)
	if isOk {
		oldPayload := loadedValue.(models.Metric)
		if metric.MetricType != oldPayload.MetricType {
			return &response, nil
		}
		if metric.MetricType == models.CounterName {
			*metric.CounterValue += *oldPayload.CounterValue
		}
	}

	metric.Hash = metric.CalculateHash(s.settings.Server.Key)
	s.metrics.Metrics.Store(metric.Name, metric)
	response.Metric = in
	return &response, nil
}
