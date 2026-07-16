package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi/v5"
)

func TestRouterWithServer(t *testing.T) {

	r := chi.NewRouter()

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	ts := httptest.NewServer(r)

	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	
	if string(bodyBytes) != "pong" {
		t.Errorf("expected 'pong', got '%s'", string(bodyBytes))
	}
}