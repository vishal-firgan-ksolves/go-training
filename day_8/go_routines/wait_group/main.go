package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int,wg *sync.WaitGroup){
	defer wg.Done()

	fmt.Printf("Worker %d: Starting\n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("Worker %d: Done\n", id)
}

func main(){

	fmt.Println("Main thread started....")

	var wg sync.WaitGroup

	wg.Add(3)

	go worker(1,&wg);
	go worker(2,&wg);
	go worker(3,&wg);

	wg.Wait()

	fmt.Println("Main thread exiting....")

}