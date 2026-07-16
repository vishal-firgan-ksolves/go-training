package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func init() {
	redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
}

func handleRequest(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("Welcome to users page..."))
}

func rateLimiter(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){

		clientIP:=r.RemoteAddr
		key:="rate:limiter:ip"+clientIP

		ctx,cancel:=context.WithTimeout(r.Context(),500*time.Millisecond)
	
		defer cancel()

		count, err := redisClient.Incr(ctx, key).Result()
		if err != nil {
			fmt.Println("Redis offline, bypassing limiter:", err)
			next.ServeHTTP(w, r)
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, key, 60*time.Second)
		}
		fmt.Printf("IP: %s | Request %d/60\n", clientIP, count)
		if count > 6 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "Too many requests. Please wait 60 seconds."}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main(){
	fmt.Println("Hello")

	router := chi.NewRouter();
	router.Use(rateLimiter)

	router.Get("/users",handleRequest)

	fmt.Println("Server running on :8090 (Rate limit: 60 req/min)")
	http.ListenAndServe(":8090", router)
}