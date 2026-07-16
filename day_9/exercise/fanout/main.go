package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup,id int,channel chan string){
	defer wg.Done()

	msg := <-channel
	fmt.Printf("Worker %d received message: %s\n", id, msg)
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Worker %d finished\n", id)
}

func main(){
	fmt.Println("Starting fanout main thread...")

	var wg sync.WaitGroup

	// var channels [3]chan string;
	// var channels []chan string;

	numWorkers:=3
	channels:=make([]chan string,numWorkers)

	for i:= range 3{
		fmt.Printf("Creating channel %d \n",i+1)
		channels[i]=make(chan string)
		// use append only when slice has 0 size , bz append add the data
		// channels=append(channels, make(chan string))
	}

	for i:=range 3{
		wg.Add(1)
		go worker(&wg,i+1,channels[i])
	}

	for i:=range 3{
		channels[i]<- fmt.Sprintf("Message for channel %d",i+1)
	}

	wg.Wait()
	fmt.Println("Finished processing all workers!!!")
}