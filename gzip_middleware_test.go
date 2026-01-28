package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkGzipMiddleware(b *testing.B) {
	handler := gzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write enough data to make gzip actually do some work
		data := make([]byte, 1024)
		for i := range data {
			data[i] = 'a'
		}
		w.Write(data)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

func TestGzipMiddlewareCorrectness(t *testing.T) {
	handler := gzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("Expected Content-Encoding: gzip, got %v", resp.Header.Get("Content-Encoding"))
	}

	// Verify body is gzipped
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read compressed body: %v", err)
	}

	if string(body) != "Hello World" {
		t.Errorf("Expected body 'Hello World', got '%s'", string(body))
	}
}
