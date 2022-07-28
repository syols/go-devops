package handlers

import (
	"compress/gzip"
	"context"
	"io"
	"log"
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

func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
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

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		next.ServeHTTP(w, r)
		log.Printf("%s::%s", uri, method)
	})
}

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
