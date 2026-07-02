package main

import (
	"errors"
	"fmt"
)

// 0. THE ROOT CAUSE
// This is the absolute bottom of our chain.
var ErrHardwareFailure = errors.New("FATAL: Disk sector corrupted")

// 1. DATA LAYER
func readDisk() error {
	// Wraps the root error
	return fmt.Errorf("readDisk operation failed: %w", ErrHardwareFailure)
}

// 2. REPOSITORY LAYER
func fetchUserData() error {
	err := readDisk()
	if err != nil {
		// Wraps Layer 1
		return fmt.Errorf("repository(fetchUserData) failed: %w", err)
	}
	return nil
}

// 3. SERVICE LAYER
func getUserProfile() error {
	err := fetchUserData()
	if err != nil {
		// Wraps Layer 2
		return fmt.Errorf("service(getUserProfile) failed: %w", err)
	}
	return nil
}

// 4. HANDLER LAYER
func handleAPIRequest() error {
	err := getUserProfile()
	if err != nil {
		// Wraps Layer 3
		return fmt.Errorf("handler(handleAPIRequest) failed: %w", err)
	}
	return nil
}

// 5. ROUTER/MIDDLEWARE LAYER (Top Level)
func executeMiddleware() error {
	err := handleAPIRequest()
	if err != nil {
		// Wraps Layer 4
		return fmt.Errorf("middleware(executeMiddleware) dropped request: %w", err)
	}
	return nil
}

func main() {
	// Generate the massive 5-layer wrapped error
	err := executeMiddleware()

	if err != nil {
		fmt.Println("==================================================")
		fmt.Println("1. THE STANDARD PRINT (Looks like one long string)")
		fmt.Println("==================================================")
		// When you print it normally, fmt just walks the list and builds one long string
		fmt.Println(err)
		fmt.Println()

		fmt.Println("==================================================")
		fmt.Println("2. TRAVERSING THE LINKED LIST IN MEMORY")
		fmt.Println("==================================================")
		
		// We use a loop to manually peel back the layers using errors.Unwrap
		currentErr := err
		layerCount := 1

		for currentErr != nil {
			fmt.Printf("\nNode %d: %s\n", layerCount, currentErr.Error())
			
			// errors.Unwrap grabs the inner pointer!
			// If there is no inner pointer (meaning we hit the root), it returns nil.
			currentErr = errors.Unwrap(currentErr)
			layerCount++
		}
		
		fmt.Println("\nReached the end of the linked list (currentErr is nil)!")
	}
}