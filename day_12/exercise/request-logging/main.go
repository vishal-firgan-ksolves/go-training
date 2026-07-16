package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func requestLogger(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		start:=time.Now()
		next.ServeHTTP(w,r)
		duration:=time.Since(start)

		fmt.Printf("LOGEER : [%s] %s | Took: %v\n", r.Method, r.URL.Path, duration)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second) 
	w.Write([]byte("Hello from the home page!"))
}


func main() {
	r := chi.NewRouter()
	r.Use(requestLogger)
	r.Get("/home", homeHandler)

	fmt.Println("System: API running on port 8080...")
	http.ListenAndServe(":8080", r)
}