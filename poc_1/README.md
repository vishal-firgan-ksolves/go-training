# Go User API PoC (with Redis Caching & Docker)

A robust Proof of Concept (PoC) demonstrating a production-ready backend architecture using Go (Golang), the Chi router, and Redis. This project implements a layered architecture (Handler -> Service -> Repository) with a thread-safe in-memory data store and a Cache-Aside Redis layer, all containerized with Docker.

---

## Comprehensive Feature List

This PoC is built with production-grade patterns in mind. Here is everything happening under the hood:

* **RESTful CRUD Operations:** Full lifecycle management (Create, Read, Update, Delete) for User entities.
* **Thread-Safe In-Memory DB:** Utilizes Go's `sync.RWMutex` to ensure memory access is locked safely during concurrent API requests, preventing race conditions.
* **Data Integrity & Validation:** * Automatically generates secure UUIDs on the backend.
    * Enforces unique email constraints across the database layer.
    * Validates incoming JSON DTO payloads before processing to prevent bad data.
* **Redis Caching Layer:** Implements a highly efficient cache wrapper to reduce database load and speed up read requests.
* **API Rate Limiting:** Protects the endpoints from spam, brute-force attacks, and DDoS attempts by restricting the number of requests a single IP can make within a specific time window.
* **Robust Unit Testing:** Comprehensive test coverage for Handlers, Services, and Repositories using Go's native testing framework and mock interfaces.
* **Custom Middleware:** Injects unique Request IDs for end-to-end tracing and tracks exact execution times (down to the microsecond) for performance monitoring.
* **Containerized Infrastructure:** Fully orchestrated backend stack using Docker Compose V2, mapping the Go API and Redis cluster over an isolated virtual network.

---

## Core Engineering Strategies

### 1. The Caching Strategy (Cache-Aside Pattern)
To optimize performance, this API uses a strict **Cache-Aside** strategy paired with targeted cache invalidation. 

* **Read Operations (GET):** When a user is requested, the application checks Redis first. If the data exists (Cache Hit), it is returned instantly. If it does not exist (Cache Miss), the application fetches it from the database, writes it into Redis with a 5-minute Time-To-Live (TTL), and then returns it to the client.
* **Write Operations (PUT/DELETE):** When a user is updated or deleted, the application updates the primary database first. Immediately after, it **evicts (deletes)** the specific user key from Redis. This guarantees that the next read request will fetch the freshest data from the database, preventing stale ghost data.

### 2. Rate Limiting Protection
To ensure high availability and prevent abuse, a rate limiter intercepts incoming HTTP requests. 
* If a client exceeds their allocated request quota, the API gracefully rejects the request with a `429 Too Many Requests` status code. 
* This ensures that bad actors or runaway scripts cannot crash the in-memory database or overwhelm the Redis connection pool.

### 3. Unit Testing & Mocks
Confidence in the codebase is maintained through rigorous unit tests.
* **Service Layer Tests:** The data repository is completely mocked, allowing us to test business logic (like duplicate email handling) in complete isolation without needing a running Redis or database instance.
* **Handler Tests:** HTTP endpoints are tested using Go's `httptest` package to verify proper status codes (e.g., `201 Created`, `400 Bad Request`, `409 Conflict`) and JSON response structures.

---

## Tech Stack

* **Language:** Go (Golang) 1.22+
* **Router:** [go-chi/chi](https://github.com/go-chi/chi)
* **Cache:** Redis 7 (Alpine)
* **Containerization:** Docker & Docker Compose V2
* **UUID Generation:** `github.com/google/uuid`

---

## Getting Started

### Prerequisites
* [Docker](https://docs.docker.com/get-docker/) installed.
* Docker Compose V2 (use `docker compose`, not the deprecated `docker-compose`).

### Booting the Application

1. Clone the repository and navigate to the project directory.
2. Start the application and Redis container by running:
   ```bash
   docker compose up --build