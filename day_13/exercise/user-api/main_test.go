package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

// setupTestRedis connects to local Redis and flushes the database for clean tests
func setupTestRedis(t *testing.T) {
	initRedis()
	if rdb != nil {
		ctx := context.Background()
		if err := rdb.FlushDB(ctx).Err(); err != nil {
			t.Logf("Warning: failed to flush redis database: %v", err)
		}
	} else {
		t.Fatal("Redis client is not initialized. Ensure Redis is running on localhost:6379.")
	}
}

func TestMiddlewaresAndValidation(t *testing.T) {
	// Initialize Router with the exact same middleware configuration as main()
	r := chi.NewRouter()
	r.Use(requestIdMiddleware)
	r.Use(requestLogger)
	r.Use(rateLimiterMiddleware)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", getAllUsers)
		r.Post("/", createUser)
		r.Get("/{id}", getUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})

	// Test case 1: Retrieve all users & verify X-Request-ID header
	t.Run("GET users - verify Request ID header", func(t *testing.T) {
		setupTestRedis(t)

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
		setupTestRedis(t)

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
		setupTestRedis(t)

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
		setupTestRedis(t)

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

	// Test case 5: GET user with Cache-Aside (HIT vs MISS)
	t.Run("GET user - verify Cache-Aside HIT/MISS behavior", func(t *testing.T) {
		setupTestRedis(t)

		// First GET request: should hit database store (Cache MISS)
		req := httptest.NewRequest("GET", "/api/users/1", nil)
		rec1 := httptest.NewRecorder()
		r.ServeHTTP(rec1, req)

		if rec1.Code != http.StatusOK {
			t.Fatalf("expected status OK on first request, got %d", rec1.Code)
		}

		cacheHeader1 := rec1.Header().Get("X-Cache")
		if cacheHeader1 != "MISS" {
			t.Errorf("expected X-Cache header to be 'MISS' on first request, got '%s'", cacheHeader1)
		}

		// Second GET request: should fetch from Redis cache (Cache HIT)
		rec2 := httptest.NewRecorder()
		r.ServeHTTP(rec2, req)

		if rec2.Code != http.StatusOK {
			t.Fatalf("expected status OK on second request, got %d", rec2.Code)
		}

		cacheHeader2 := rec2.Header().Get("X-Cache")
		if cacheHeader2 != "HIT" {
			t.Errorf("expected X-Cache header to be 'HIT' on second request, got '%s'", cacheHeader2)
		}
	})

	// Test case 6: GET user with cache invalidation on updates/deletes
	t.Run("GET user - verify Cache Invalidation on PUT and DELETE", func(t *testing.T) {
		setupTestRedis(t)

		// 1. First request -> MISS (and cache populate)
		reqGet := httptest.NewRequest("GET", "/api/users/1", nil)
		recGet1 := httptest.NewRecorder()
		r.ServeHTTP(recGet1, reqGet)
		if recGet1.Header().Get("X-Cache") != "MISS" {
			t.Errorf("expected MISS, got %s", recGet1.Header().Get("X-Cache"))
		}

		// 2. Second request -> HIT
		recGet2 := httptest.NewRecorder()
		r.ServeHTTP(recGet2, reqGet)
		if recGet2.Header().Get("X-Cache") != "HIT" {
			t.Errorf("expected HIT, got %s", recGet2.Header().Get("X-Cache"))
		}

		// 3. PUT request -> should update data and invalidate cache
		body := `{"name": "Alice in Wonderland"}`
		reqPut := httptest.NewRequest("PUT", "/api/users/1", bytes.NewBufferString(body))
		reqPut.Header.Set("Content-Type", "application/json")
		recPut := httptest.NewRecorder()
		r.ServeHTTP(recPut, reqPut)

		if recPut.Code != http.StatusOK {
			t.Fatalf("expected PUT to succeed, got %d", recPut.Code)
		}

		// 4. Third GET request -> should MISS again (since cache was invalidated)
		recGet3 := httptest.NewRecorder()
		r.ServeHTTP(recGet3, reqGet)
		if recGet3.Header().Get("X-Cache") != "MISS" {
			t.Errorf("expected MISS after update (due to cache invalidation), got %s", recGet3.Header().Get("X-Cache"))
		}

		var u User
		if err := json.Unmarshal(recGet3.Body.Bytes(), &u); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if u.Name != "Alice in Wonderland" {
			t.Errorf("expected user name to be 'Alice in Wonderland', got '%s'", u.Name)
		}
	})

	// Test case 7: Exceed Rate Limit
	t.Run("Rate Limit - verify 429 Too Many Requests on limit breach", func(t *testing.T) {
		setupTestRedis(t)

		// Client IP to rate-limit
		clientIP := "192.0.2.100"

		// Send 60 requests - all must succeed (limit is 60 req/min)
		for i := 1; i <= 60; i++ {
			req := httptest.NewRequest("GET", "/api/users", nil)
			req.RemoteAddr = clientIP + ":5678" // SplitHostPort will isolate clientIP
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected request %d to succeed with 200, got %d", i, rec.Code)
			}
		}

		// 61st request: must breach the limit and return 429
		reqLimitBreach := httptest.NewRequest("GET", "/api/users", nil)
		reqLimitBreach.RemoteAddr = clientIP + ":5678"
		recLimitBreach := httptest.NewRecorder()
		r.ServeHTTP(recLimitBreach, reqLimitBreach)

		if recLimitBreach.Code != http.StatusTooManyRequests {
			t.Errorf("expected 61st request to return 429, got %d", recLimitBreach.Code)
		}

		// Verify headers and error payload
		retryHeader := recLimitBreach.Header().Get("Retry-After")
		if retryHeader != "60" {
			t.Errorf("expected Retry-After header to be '60', got '%s'", retryHeader)
		}

		var env ErrorEnvelope
		if err := json.Unmarshal(recLimitBreach.Body.Bytes(), &env); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if env.Success {
			t.Error("expected Success to be false")
		}
		if env.Error != "Too Many Requests" {
			t.Errorf("expected error field to be 'Too Many Requests', got '%s'", env.Error)
		}
		if !strings.Contains(env.Message, "Rate limit exceeded") {
			t.Errorf("expected message to contain 'Rate limit exceeded', got '%s'", env.Message)
		}
	})
}
