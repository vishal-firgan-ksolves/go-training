package main

import (
	"fmt"
	"time"
)

func main() {
	jobs := make(chan int,5)
	
	// producer will not block
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("Producer: Queueing job %d\n", i)
			jobs <- i
		}
		fmt.Println("Producer: All jobs queued, closing")
		close(jobs)
	}()
	
	// consumer
	for job := range jobs {
		fmt.Printf("Consumer: Processing job %d\n", job)
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println("All jobs processed")
}