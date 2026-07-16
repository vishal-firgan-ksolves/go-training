package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {

	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	dbStore := NewInMemoryUserStore()
	cachedRepo := NewCachedUserRepository(dbStore, redisClient)
	userService := NewUserService(cachedRepo, redisClient)
	userHandler := NewUserHandler(userService)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(RequestIDMiddleware)
	r.Use(LoggerMiddleware)
	r.Use(RateLimiterMiddleware(redisClient, 10, 1*time.Minute))

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/", userHandler.getAllUsers)
		r.Get("/{id}", userHandler.getUser)
		r.Put("/{id}", userHandler.updateUser)
		r.Delete("/{id}", userHandler.deleteUser)
	})

	fmt.Println("System: POC 1 User API is live on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}