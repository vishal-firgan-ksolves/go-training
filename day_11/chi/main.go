package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is healthy!"))
	})

	r.Route("/api/v1", func(apiRouter chi.Router) {
		
		// /api/v1/status
		apiRouter.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "V1 API operational status: OK. Time: %v", time.Now().Format(time.Kitchen))
		})

		apiRouter.Mount("/users", userFeatureRouter())
	})

	fmt.Println("System: Launching Chi web server on port 8080...")
	http.ListenAndServe(":8080", r)
}

func userFeatureRouter() chi.Router {
	sub := chi.NewRouter()

	// /api/v1/users/
	sub.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"users": ["Vishal", "Amit", "Rahul"]}`))
	})

	// /api/v1/users/
	sub.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "User successfully created!"}`))
	})

	// /api/v1/users/42
	sub.Get("/{userID}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userID")
		
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": %s, "name": "Vishal", "role": "Senior Engineer"}`, id)
	})

	return sub
}