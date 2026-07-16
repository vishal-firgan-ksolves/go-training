package main

import (
	"context"
	"fmt"
	"time"
)

func queryDatabase(ctx context.Context, dataChannel chan<- string, errChannel chan<- error) {
	userID := ctx.Value("Id")
	fmt.Println("Worker: Started fetching data for user with id:", userID)

	dbResponse := make(chan string, 1)

	go func() {
		fmt.Println("Worker: Waiting for Data from DB wires...")
		time.Sleep(3 * time.Second)
		dbResponse <- "User:Vishal Age:25"
	}()

	select {
	case result := <-dbResponse:
		dataChannel <- result
	case <-ctx.Done():
		errChannel <- ctx.Err()
	}
}

func main() {
	fmt.Println("Starting Application....")

	dataChannel := make(chan string, 1)
	errChannel := make(chan error, 1)

	ctx := context.WithValue(context.Background(), "Id", "100")

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	go queryDatabase(ctx, dataChannel, errChannel)

	select {
	case data := <-dataChannel:
		fmt.Printf("Main: Received user data: %s\n", data)
	case err := <-errChannel:
		fmt.Printf("Main: Failed to fetch data. Reason: %v\n", err)
	}

	fmt.Println("Application Ending......")
}