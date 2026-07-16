package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// go get github.com/go-chi/chi/v5@v5.0.12

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside home handler")

	w.Write([]byte("Hello from home page"))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Here is your Profile!"))
}

// middleware function
func loggingMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start of middleware")
		fmt.Printf("LOG::: [%s] %s\n", r.Method, r.URL.Path)

		// call the actual handler
		next.ServeHTTP(w, r)

		fmt.Println("Ending of middleware after handler completion")
	})
}

func main() {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)
	
	r.Get("/home", homeHandler)
	r.Get("/profile", profileHandler)

	fmt.Println("System: API running on port 8080...")
	http.ListenAndServe(":8080", r)
}