package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

type ErrorEnvelope struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

type UserStore struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

var store = &UserStore{
	users:  make(map[int]User),
	nextID: 1,
}

var validate = validator.New()

type contextKey string

const requestIDKey contextKey = "request_id"

func init() {
	now := time.Now()
	store.users[1] = User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: now}
	store.users[2] = User{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: now}
	store.nextID = 3
}

func getRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(requestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func sendError(w http.ResponseWriter, code int, errType string, msg string, ctx context.Context) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	reqID := getRequestID(ctx)
	
	envelope := ErrorEnvelope{
		Success:   false,
		Error:     errType,
		Message:   msg,
		RequestID: reqID,
	}
	
	json.NewEncoder(w).Encode(envelope)
}

// 1. GET /users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	nameFilter := r.URL.Query().Get("name")

	store.mu.RLock()
	defer store.mu.RUnlock()

	users := make([]User, 0)
	
	for _, u := range store.users {
		if nameFilter != "" && !strings.EqualFold(u.Name, nameFilter) {
			continue
		}
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// 2. GET /users/{id}
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid User ID", r.Context())
		return
	}

	store.mu.RLock()
	user, exists := store.users[id]
	store.mu.RUnlock()

	if !exists {
		sendError(w, http.StatusNotFound, "Not Found", "User not found", r.Context())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 3. POST /users (with input validation)
func createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON body", r.Context())
		return
	}

	// Validate fields (required, email format)
	if err := validate.Struct(req); err != nil {
		var errMsgs []string
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()))
		}
		sendError(w, http.StatusBadRequest, "Validation Error", strings.Join(errMsgs, "; "), r.Context())
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	newUser := User{
		ID:        store.nextID,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	store.users[store.nextID] = newUser
	store.nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// 4. PUT /users/{id} (with input validation if fields are provided)
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid User ID", r.Context())
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON body", r.Context())
		return
	}

	if err := validate.Struct(req); err != nil {
		var errMsgs []string
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field '%s' is %s", e.Field(), e.Tag()))
		}
		sendError(w, http.StatusBadRequest, "Validation Error", strings.Join(errMsgs, "; "), r.Context())
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	user, exists := store.users[id]
	if !exists {
		sendError(w, http.StatusNotFound, "Not Found", "User not found", r.Context())
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	store.users[id] = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 5. DELETE /users/{id}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid User ID", r.Context())
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.users[id]; !exists {
		sendError(w, http.StatusNotFound, "Not Found", "User not found", r.Context())
		return
	}

	delete(store.users, id)
	w.WriteHeader(http.StatusNoContent) 
}

// Custom wrapper to intercept status code for logging
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Request logging middleware
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		srw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(srw, r)
		
		duration := time.Since(start)
		reqID := getRequestID(r.Context())
		
		log.Printf("LOGGER: [%s] %s %d %s | Duration: %v | Request-ID: %s",
			r.Method,
			r.URL.Path,
			srw.statusCode,
			http.StatusText(srw.statusCode),
			duration,
			reqID,
		)
	})
}

// Request-ID middleware (generate UUID and attach to ctx + header)
func requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	r := chi.NewRouter()

	r.Use(requestIdMiddleware)
	r.Use(requestLogger)

	r.Route("/api/users", func(r chi.Router) {
		r.Get("/", getAllUsers)
		r.Post("/", createUser)
		r.Get("/{id}", getUser)
		r.Put("/{id}", updateUser)
		r.Delete("/{id}", deleteUser)
	})

	fmt.Println("System: API is live on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}