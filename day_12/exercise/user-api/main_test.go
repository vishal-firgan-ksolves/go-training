package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestMiddlewaresAndValidation(t *testing.T) {
	
	r := chi.NewRouter()
	r.Use(requestIdMiddleware)
	r.Use(requestLogger)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", getAllUsers)
		r.Post("/", createUser)
		r.Get("/{id}", getUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})

	// Test case 1: Retrieve all users & verify X-Request-ID header
	t.Run("GET users - verify Request ID header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status OK, got %d", rec.Code)
		}

		reqID := rec.Header().Get("X-Request-ID")
		if reqID == "" {
			t.Error("expected X-Request-ID header to be set, but it was empty")
		}
	})

	// Test case 2: POST users validation failure & consistent error envelope
	t.Run("POST users - invalid payload", func(t *testing.T) {
		body := `{"name": "", "email": "invalid-email"}`
		req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status BadRequest, got %d", rec.Code)
		}

		var env ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if env.Success {
			t.Error("expected success to be false")
		}
		if env.Error != "Validation Error" {
			t.Errorf("expected error field to be 'Validation Error', got '%s'", env.Error)
		}
		if !strings.Contains(env.Message, "Name") || !strings.Contains(env.Message, "Email") {
			t.Errorf("expected validation message to detail failures, got '%s'", env.Message)
		}
		
		reqIDHeader := rec.Header().Get("X-Request-ID")
		if env.RequestID == "" || env.RequestID != reqIDHeader {
			t.Errorf("expected RequestID in envelope (%s) to match header (%s)", env.RequestID, reqIDHeader)
		}
	})

	// Test case 3: POST users success
	t.Run("POST users - success", func(t *testing.T) {
		body := `{"name": "John Doe", "email": "john@example.com"}`
		req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status Created, got %d", rec.Code)
		}

		var u User
		if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
			t.Fatalf("failed to unmarshal user response: %v", err)
		}

		if u.Name != "John Doe" || u.Email != "john@example.com" {
			t.Errorf("unexpected user values: %+v", u)
		}
	})

	// Test case 4: GET non-existent user - verify 404 and envelope
	t.Run("GET user not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users/999", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status NotFound, got %d", rec.Code)
		}

		var env ErrorEnvelope
		if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if env.Success {
			t.Error("expected success to be false")
		}
		if env.Error != "Not Found" {
			t.Errorf("expected error field to be 'Not Found', got '%s'", env.Error)
		}
		if env.Message != "User not found" {
			t.Errorf("expected message to be 'User not found', got '%s'", env.Message)
		}
		
		reqIDHeader := rec.Header().Get("X-Request-ID")
		if env.RequestID == "" || env.RequestID != reqIDHeader {
			t.Errorf("expected RequestID in envelope (%s) to match header (%s)", env.RequestID, reqIDHeader)
		}
	})
}
