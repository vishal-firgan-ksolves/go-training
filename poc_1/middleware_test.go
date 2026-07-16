package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxID := r.Context().Value(requestIDKey)
		assert.NotNil(t, ctxID)
		assert.NotEmpty(t, ctxID.(string))
		w.WriteHeader(http.StatusOK)
	})

	middlewareToTest := RequestIDMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/api/test", nil)
	rec := httptest.NewRecorder()
	middlewareToTest.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("X-Request-ID"))
}

func TestLoggerMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middlewareChain := RequestIDMiddleware(LoggerMiddleware(nextHandler))

	req := httptest.NewRequest("GET", "/api/test-log", nil)
	rec := httptest.NewRecorder()
	middlewareChain.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRateLimiterMiddleware_Allow(t *testing.T) {
	dummyRedis := getFastDummyRedis()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := RateLimiterMiddleware(dummyRedis, 10, time.Minute)(nextHandler)

	routerChain := RequestIDMiddleware(middleware)

	req := httptest.NewRequest("GET", "/api/test-allow", nil)
	rec := httptest.NewRecorder()
	routerChain.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}