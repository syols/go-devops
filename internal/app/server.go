package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/syols/go-devops/config"
	"github.com/syols/go-devops/internal/handlers"
	"github.com/syols/go-devops/internal/store"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	server   http.Server
	metrics  store.MetricsStorage
	settings config.Config
}

type Handler func(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request)

func NewServer(settings config.Config) (Server, error) {
	metrics, err := store.NewMetricsStorage(settings)
	if err != nil {
		return Server{}, err
	}
	return Server{
		metrics:  metrics,
		settings: settings,
	}, nil
}

func (s *Server) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	s.shutdown(ctx)

	server := http.Server{
		Addr:    s.settings.Address(),
		Handler: s.router(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(handlers.Compress)
	router.Use(handlers.Logging)
	router.Use(handlers.Save(s.metrics))

	router.Get("/", handlers.Healthcheck)
	router.Get("/ping", handlers.Ping(s.metrics))
	router.Get("/value/{type}/{name}", handlers.Value(s.metrics))
	router.Post("/update/{type}/{name}/{value}", handlers.Update(s.metrics))
	router.Post("/update/", handlers.UpdateJSON(s.metrics, s.settings.Server.Key))
	router.Post("/updates/", handlers.UpdatesJSON(s.metrics, s.settings.Server.Key))
	router.Post("/value/", handlers.ValueJSON(s.metrics))
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
