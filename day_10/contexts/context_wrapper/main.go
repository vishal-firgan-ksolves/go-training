package main

import (
	"context"
	"fmt"
	"time"
)

// Goroutine 1: Main handler
func mainHandler(ctx context.Context) {
	fmt.Println("[Handler 1] Received context")
	
	// Wrap: Create child context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	fmt.Println("[Handler 1] Created child context with 10s timeout")
	
	// Pass to Goroutine 2
	go processRequest(ctxWithTimeout)
	
	time.Sleep(2 * time.Second)
}

// Goroutine 2: Request processor
func processRequest(ctx context.Context) {
	fmt.Println("[Handler 2] Received context with timeout")
	
	// Wrap: Add request-scoped value
	ctxWithValue := context.WithValue(ctx, "requestID", "REQ-12345")
	
	fmt.Println("[Handler 2] Added request ID to context")
	
	// Pass to Goroutine 3
	go validateData(ctxWithValue)
	
	time.Sleep(1 * time.Second)
}

// Goroutine 3: Data validator
func validateData(ctx context.Context) {
	fmt.Println("[Handler 3] Received context with timeout + request ID")
	
	// Access scoped value
	requestID := ctx.Value("requestID")
	fmt.Printf("[Handler 3] Got request ID from context: %v\n", requestID)
	
	// Wrap: Add more data
	ctxWithUserID := context.WithValue(ctx, "userID", 999)
	
	fmt.Println("[Handler 3] Added user ID to context")
	
	// Pass to Goroutine 4
	go processUser(ctxWithUserID)
	
	time.Sleep(1 * time.Second)
}

// Goroutine 4: User processor
func processUser(ctx context.Context) {
	fmt.Println("[Handler 4] Received context with timeout + request ID + user ID")
	
	// Access all values
	requestID := ctx.Value("requestID")
	userID := ctx.Value("userID")
	
	fmt.Printf("[Handler 4] Request ID: %v, User ID: %v\n", requestID, userID)
	
	// Make critical decision based on context state
	select {
	case <-ctx.Done():
		fmt.Printf("[Handler 4] Context cancelled: %v\n", ctx.Err())
		return
	case <-time.After(500 * time.Millisecond):
		fmt.Println("[Handler 4] Processing user data...")
		fmt.Println("[Handler 4] All good! Continuing...")
	}
}

func main() {
	fmt.Println("=== CONTEXT WRAPPING DEMONSTRATION ===")
	
	// Create root context
	rootCtx := context.Background()
	fmt.Println("[Main] Created root context")
	
	mainHandler(rootCtx)
	
	time.Sleep(5 * time.Second)
	
	fmt.Println("\n=== ALL GOROUTINES COMPLETED ===")
}