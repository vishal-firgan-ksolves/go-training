package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type contextKey string
const requestIDKey contextKey = "requestID"

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		reqID := r.Context().Value(requestIDKey)

		log.Printf("[ID: %s] STARTED %s %s", reqID, r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("[ID: %s] COMPLETED in %v", reqID, time.Since(start))
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		
		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RateLimiterMiddleware(redisClient *redis.Client, maxRequests int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Context().Value(requestIDKey)
			userIP := r.RemoteAddr
			cacheKey := fmt.Sprintf("ratelimit:%s", userIP)

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			count, err := redisClient.Incr(ctx, cacheKey).Result()
			if err != nil {
				log.Printf("[ID: %s] Redis Error: %v", reqID, err)
				next.ServeHTTP(w, r)
				return
			}

			if count == 1 {
				redisClient.Expire(ctx, cacheKey, window)
			}

			if count > maxRequests {
				log.Printf("[ID: %s] BLOCKED: IP %s hit limit", reqID, userIP)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Rate limit exceeded. Please wait.",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}