package main

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// gzipHandler wraps an http.Handler to compress responses with gzip
func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Create a gzip writer
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the ResponseWriter
		w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter wraps the standard ResponseWriter to write gzipped data
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
