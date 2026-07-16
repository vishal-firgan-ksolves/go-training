package main

import (

	"net/http"
	"net/http/httptest"
	"testing"
)

func healthCheckHandler(w http.ResponseWriter,r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"alive"}`))
}

func TestHealthCheckEndpint(t *testing.T){
	req:=httptest.NewRequest(http.MethodGet,"/health",nil)

	rr:=httptest.NewRecorder()

	healthCheckHandler(rr,req)

	if rr.Code != http.StatusOK{
		t.Errorf("Expected status %d ,got %d",http.StatusOK,rr.Code)
	}

	expectedBody := `{"status":"alive"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected body %s, got %s", expectedBody, rr.Body.String())
	}
}