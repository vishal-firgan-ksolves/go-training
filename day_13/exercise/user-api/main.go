package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// User representation in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest validates body for POST /users
type CreateUserRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

// UpdateUserRequest validates body for PUT /users/{id}
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// ErrorEnvelope standardizes all error responses
type ErrorEnvelope struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

// UserStore manages in-memory database of users with concurrency control
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

// Redis client instance
var rdb *redis.Client

func init() {
	now := time.Now()
	// Seed some initial users
	store.users[1] = User{ID: 1, Name: "Alice", Email: "alice@example.com", CreatedAt: now}
	store.users[2] = User{ID: 2, Name: "Bob", Email: "bob@example.com", CreatedAt: now}
	store.nextID = 3
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// Initialize connection to the Redis server
func initRedis() {
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Printf("REDIS WARNING: Failed to connect to Redis at %s: %v. Cache and rate limiting will run in fail-soft mode.", redisAddr, err)
	} else {
		log.Printf("REDIS INFO: Successfully connected to Redis at %s", redisAddr)
	}
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

// 2. GET /users/{id} - Implements Cache-aside with a 5 minute TTL
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Bad Request", "Invalid User ID", r.Context())
		return
	}

	ctx := r.Context()
	cacheKey := fmt.Sprintf("user:%d", id)

	// Step 1: Check Redis cache first (Cache-Aside pattern)
	if rdb != nil {
		cachedUserJson, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			// Cache Hit!
			var cachedUser User
			if err := json.Unmarshal([]byte(cachedUserJson), &cachedUser); err == nil {
				log.Printf("CACHE HIT: User ID %d successfully fetched from Redis cache", id)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Cache", "HIT")
				json.NewEncoder(w).Encode(cachedUser)
				return
			}
			log.Printf("CACHE ERROR: Failed to deserialize cached JSON for User ID %d: %v", id, err)
		} else if err != redis.Nil {
			log.Printf("CACHE ERROR: Failed fetching key '%s' from Redis: %v", cacheKey, err)
		}
	}

	// Step 2: Cache Miss - Retrieve from DB
	log.Printf("CACHE MISS: User ID %d not in Redis cache. Querying memory store...", id)

	store.mu.RLock()
	user, exists := store.users[id]
	store.mu.RUnlock()

	if !exists {
		log.Printf("STORE MISS: User ID %d not found in database memory store", id)
		sendError(w, http.StatusNotFound, "Not Found", "User not found", r.Context())
		return
	}

	// Step 3: Populate cache for subsequent requests (TTL: 5 Minutes)
	if rdb != nil {
		userJson, err := json.Marshal(user)
		if err == nil {
			err = rdb.Set(ctx, cacheKey, userJson, 5*time.Minute).Err()
			if err != nil {
				log.Printf("CACHE ERROR: Failed to cache User ID %d in Redis: %v", id, err)
			} else {
				log.Printf("CACHE WRITE: User ID %d cached in Redis (TTL: 5m)", id)
			}
		} else {
			log.Printf("CACHE ERROR: Failed to serialize User ID %d for caching: %v", id, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(user)
}

// 3. POST /users (with input validation)
func createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
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

// 4. PUT /users/{id} (Has cache invalidation)
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

	// Invalidate Cache for this user to avoid stale data (Cache-Aside best practice)
	if rdb != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		if err := rdb.Del(r.Context(), cacheKey).Err(); err != nil {
			log.Printf("CACHE ERROR: Failed to invalidate cache for User ID %d: %v", id, err)
		} else {
			log.Printf("CACHE INVALIDATE: Invalidated cache key '%s' due to update", cacheKey)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 5. DELETE /users/{id} (with cache invalidation)
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

	// Invalidate Cache for this user to avoid phantom reads
	if rdb != nil {
		cacheKey := fmt.Sprintf("user:%d", id)
		if err := rdb.Del(r.Context(), cacheKey).Err(); err != nil {
			log.Printf("CACHE ERROR: Failed to invalidate cache for User ID %d on delete: %v", id, err)
		} else {
			log.Printf("CACHE INVALIDATE: Invalidated cache key '%s' due to deletion", cacheKey)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// Custom response writer wrapper to capture status code for logging
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

// Request logging middleware - captures timing, status code and request ID
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

// Request-ID middleware (attaches UUID to headers and context)
func requestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New().String()
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Extract client IP address from request headers or remote address
func getClientIP(r *http.Request) string {
	// Check standard proxy headers
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		if host, _, err := net.SplitHostPort(ip); err == nil {
			ip = host
		}
	}
	return ip
}

// Rate limiting middleware: 60 req/min per IP using Redis INCR
func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rdb == nil {
			log.Printf("RATE LIMITER WARNING: Redis is offline. Rate limiter bypassed.")
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()
		ip := getClientIP(r)

		// Create a minute-based bucket key: ratelimit:<ip>:<YYYYMMDDHHMM>
		minuteBucket := time.Now().Format("200601021504")
		key := fmt.Sprintf("ratelimit:%s:%s", ip, minuteBucket)

		// Atomic increment of the bucket count
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("RATE LIMITER ERROR: Redis INCR failed for %s: %v. Bypass allowed.", ip, err)
			next.ServeHTTP(w, r)
			return
		}

		// If this is the first request in the current minute bucket, set a 60-second TTL
		if count == 1 {
			if err := rdb.Expire(ctx, key, time.Minute).Err(); err != nil {
				log.Printf("RATE LIMITER ERROR: Failed to set expiry for key %s: %v", key, err)
			}
		}

		const maxLimit = 60

		// Check if threshold is breached
		if count > maxLimit {
			log.Printf("RATE LIMITER: IP %s BLOCKED. Request count %d exceeded limit of %d req/min", ip, count, maxLimit)
			
			// Optional: Set standard headers
			w.Header().Set("Retry-After", "60")
			
			sendError(w, http.StatusTooManyRequests, "Too Many Requests", "Rate limit exceeded. Limit is 60 requests per minute.", ctx)
			return
		}

		log.Printf("RATE LIMITER: IP %s allowed. Count: %d/%d for current minute bucket", ip, count, maxLimit)
		next.ServeHTTP(w, r)
	})
}

func main() {
	initRedis()

	r := chi.NewRouter()

	r.Use(requestIdMiddleware)
	r.Use(requestLogger)
	r.Use(rateLimiterMiddleware)

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