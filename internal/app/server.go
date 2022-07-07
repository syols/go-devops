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
	"time"
)

type Server struct {
	server   http.Server
	metrics  store.MetricsStorage
	settings config.Config
}

type Handler func(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request)

func NewServer(settings config.Config) Server {
	return Server{
		metrics:  store.NewMetricsStorage(settings),
		settings: settings,
	}
}

func (s *Server) handler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(s.metrics, s.settings.Server.Key, w, r)
	}
}

func (s *Server) Run() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	s.shutdown(sign)

	server := http.Server{
		Addr:    s.settings.Address(),
		Handler: s.router(),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Print(err.Error())
	}
}

func (s *Server) router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(handlers.Compress)

	router.Get("/", handlers.Healthcheck)
	router.Get("/ping", s.handler(handlers.Ping))
	router.Get("/value/{type}/{name}", s.handler(handlers.Value))
	router.Post("/update/{type}/{name}/{value}", s.handler(handlers.Update))
	router.Post("/update/", s.handler(handlers.UpdateJSON))
	router.Post("/updates/", s.handler(handlers.UpdatesJSON))
	router.Post("/value/", s.handler(handlers.ValueJSON))
	return router
}

func (s *Server) shutdown(sign chan os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		<-sign
		if err := s.server.Shutdown(ctx); err != nil {
			log.Println(err.Error())
		}
		s.metrics.Save()
		os.Exit(0)
	}()
}
