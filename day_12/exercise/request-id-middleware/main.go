package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// We create a custom type just for our context keys to prevent collisions.
type contextKey string

const requestIDKey contextKey = "request_id"

func requestIdMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){

		fmt.Println("Inside middleware.......")

		reqId:=uuid.New().String()

		w.Header().Set("X-Request-ID",reqId)

		ctx:=context.WithValue(r.Context(),requestIDKey,reqId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestHandler(w http.ResponseWriter,r *http.Request){

	reqId,ok:=r.Context().Value(requestIDKey).(string)

	if !ok{
		reqId="unknown"
	}

	fmt.Printf("Processing handler for Request with id: %s\n", reqId)
	w.Write([]byte("Process complete. Check headers for Request ID."))
}

func main(){
	r:=chi.NewRouter()

	r.Use(requestIdMiddleware)

	r.Get("/",requestHandler)


}