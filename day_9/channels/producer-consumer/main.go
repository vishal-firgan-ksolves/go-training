package main

import (
	"fmt"
	"sync"
	"time"
)

// Producer: sends numbers on channel 
// syntax say send only
func producer(ch chan<-int, wg *sync.WaitGroup) {
	defer wg.Done()

	// var value int
	
	for i := 1; i <= 5; i++ {
		fmt.Printf("Producer: sending %d\n", i)
		// Send on channel
		ch <- i  

		// Cannot receive from send only channel
		// value :=<-ch
		time.Sleep(500 * time.Millisecond)
	}
	
	fmt.Println("Producer: closing channel")
	close(ch)  // Signal: done sending
}

// Consumer: receives numbers from channel
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Cannot send to receive only channel
	// ch<-4
	
	for value := range ch {
		fmt.Printf("Consumer: received %d\n", value)
		time.Sleep(1 * time.Second)
	}
	
	fmt.Println("Consumer: channel closed, done")
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	
	wg.Add(2)
	go producer(ch, &wg)
	go consumer(ch, &wg)
	
	wg.Wait()
	fmt.Println("Main: all done!")
}
