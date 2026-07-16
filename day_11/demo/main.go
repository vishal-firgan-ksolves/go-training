package main

import (
	"fmt"
	"net/http"
)

type ApiHandler struct{
	apiKey string
	version string
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("X-API-Version", h.version)

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `{"message":"API Handler called","apiKey":"%s"}`, h.apiKey)
}

// func custom(w http.ResponseWriter,r *http.Request){
// 	w.Header().Set("Content-Type","text")
// 	w.Header().Set("Version","33")

// 	w.WriteHeader(http.StatusAccepted)

// 	fmt.Fprintf(w,`Helo form backend!!!`)
// }

func main(){
	handler := &ApiHandler{
		apiKey:  "secret-key-123",
		version: "1.0.0",
	}

	http.ListenAndServe(":8080", handler)
	// http.ListenAndServe(":8080/api",http.HandlerFunc(custom))
}