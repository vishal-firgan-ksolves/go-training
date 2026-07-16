package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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

func init() {
	now := time.Now()
	store.users[1] = User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: now}
	store.users[2] = User{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: now}
	store.nextID = 3
}

func sendError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg})
}

// 1. GET /users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	// reading query parameters
	nameFilter := r.URL.Query().Get("name")

	store.mu.RLock()
	defer store.mu.RUnlock()

	users := make([]User, 0)
	
	for _, u := range store.users {
		// If filter is provided filter it
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
		sendError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	store.mu.RLock()
	user, exists := store.users[id]
	store.mu.RUnlock()

	if !exists {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 3. POST /users
func createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Full Write Lock
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

// 4. PUT /users/{id}
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	user, exists := store.users[id]
	if !exists {
		sendError(w, http.StatusNotFound, "User not found")
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
		sendError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.users[id]; !exists {
		sendError(w, http.StatusNotFound, "User not found")
		return
	}

	delete(store.users, id)
	w.WriteHeader(http.StatusNoContent) 
}

func main() {
	r := chi.NewRouter()

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
