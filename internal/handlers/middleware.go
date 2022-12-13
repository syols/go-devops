package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/syols/go-devops/internal/store"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

const ContentType = "application/json"

// Compress middleware. Compress payload to gzip
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.NoCompression)
		if err != nil {
			log.Fatal(err)
		}

		defer func(gz *gzip.Writer) {
			err := gz.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(gz)

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// Save saves metrics to store
func Save(metrics store.MetricsStorage) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(ctx))
			if metrics.SaveInterval == 0 || metrics.Store.Type() == "database" {
				ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
				defer cancel()
				err := metrics.Save(ctx)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
	}
}

// CheckSubnet check trusted subnet
func CheckSubnet(trustedSubnet *string) func(http.Handler) http.Handler {
	_, trustedNetwork, err := net.ParseCIDR(*trustedSubnet)

	if err != nil {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				http.Error(w, "Trusted subnet contains incorrect data", http.StatusInternalServerError)
				return
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			agentIP := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(agentIP)
			if !trustedNetwork.Contains(ip) {
				http.Error(w, "No trusted subnet is specified", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
