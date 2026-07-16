package main

import (
	"fmt"
	"time"
)


func fetchSlowData(resultChannel chan<- string){

	fmt.Println("Worker: Starting heavy database query...")

	time.Sleep(3 * time.Second) 
	
	resultChannel <- "Successfully fetched user data!"
	fmt.Println("Worker: Data sent to channel.")
}

func main(){

	resultChannel:=make(chan string)

	go fetchSlowData(resultChannel)

	fmt.Println("Main: Waiting for response with a 2-second strict timeout...")

	select {
	case data := <-resultChannel:
		fmt.Printf("Main: Success! Data received: %s\n", data)

	case <-time.After(2 * time.Second):
		fmt.Println("Main: TIMEOUT EXPIRED! Aborting request to save server resources.")
	}

	fmt.Println("Main: Program exiting cleanly.")
}