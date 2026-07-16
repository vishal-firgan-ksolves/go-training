package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	mockDB = map[string]User{
		"99": {ID: "99", Name: "Vishal", Email: "vishal@example.com"},
		"100": {ID: "100", Name: "Aman", Email: "aman@example.com"},
	}
)

func init() {
	redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
}


func getUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	cacheKey := "user:profile:" + userID

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	//Check data in Redis
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()

	// Fetch from cache
	if err == nil {
		fmt.Println("Fetching data from cache........")
		w.Header().Set("X-Cache", "HIT")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedData))
		return
	}

	// Fetch from database
	if err == redis.Nil {
		fmt.Println("No data in cache, Fetching from DB...")

		user, exists := mockDB[userID]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		userJSON, _ := json.Marshal(user)

		// set fetched users data into cache for future fetch
		redisClient.Set(ctx, cacheKey, userJSON, 55*time.Second)

		w.Header().Set("X-Cache", "MISS")
		w.Header().Set("Content-Type", "application/json")
		w.Write(userJSON)
		return
	}
}

func main(){

	r := chi.NewRouter()
	r.Get("/users/{id}", getUserProfile)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8089", r))
}