package main

import (
	"fmt"
	"time"
)

func main() {
	result := make(chan string)
	
	go func() {
		fmt.Println("Worker: Starting work")
		time.Sleep(2 * time.Second)
		fmt.Println("Worker: Work complete, sending result")
		result <- "Task completed"  // block here until main receives
		fmt.Println("Worker: Main received result")
	}()
	
	fmt.Println("Main: Waiting for result")
	res := <-result  // blocks until worker sends
	fmt.Println("Main: Received:", res)
	fmt.Println("Main: Continuing...")
}