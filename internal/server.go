package internal

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/syols/go-devops/internal/handlers"
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
	sets    settings.Config
}

type Handler func(metrics store.MetricsStorage, key *string, w http.ResponseWriter, r *http.Request)

func NewServer(sets settings.Config) Server {
	return Server{
		metrics: store.NewMetricsStorage(sets),
		sets:    sets,
	}
}

func (s *Server) handler(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(s.metrics, s.sets.Server.Key, w, r)
	}
}

func (s *Server) Run() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	s.shutdown(sign)

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

	server := http.Server{
		Addr:    s.sets.GetAddress(),
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Print(err.Error())
	}
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
