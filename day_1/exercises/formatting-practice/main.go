package main

import (
	"fmt"
	"os"
)

func main() {
	// 1. Verification of Formatting Verbs via Printf
	name := "Vishal"
	piValue := 3.141592
	itemCount := 42
	isComplete := true

	fmt.Printf("--- VERB TESTS ---\n")
	fmt.Printf("Type check: name is %T | itemCount is %T\n", name, itemCount)
	fmt.Printf("String representation: %s or quoted %q\n", name, name)
	fmt.Printf("Integer check: decimal %d | binary %b\n", itemCount, itemCount)
	fmt.Printf("Float restriction: 2-points -> %.2f\n", piValue)
	fmt.Printf("Boolean status check: %t\n\n", isComplete)

	// 2. Demonstration of Input Capturing via Scanln
    var userRating float64
	fmt.Print("Enter your Go rating from 1.0 to 10.0: ")
	
	// Capture the error value returned by Scanln
	_, err := fmt.Scanln(&userRating)
	if err != nil {
		fmt.Printf("Error: Invalid input provided. Please enter a valid number. (%v)\n", err)
		os.Exit(1) // Gracefully terminate or handle the failure
	}

	fmt.Printf("Input Confirmation: Stored entry value is %.1f\n", userRating)
}