package main

import (
	"fmt"
	"sync"
	"time"
)

func workers(id int,wg *sync.WaitGroup,jobsChannel chan int){
	defer wg.Done()

	for job:=range jobsChannel{
		fmt.Printf("Worker %d: Started job %d\n", id, job)
		time.Sleep(1 * time.Second) 
		fmt.Printf("Worker %d: Finished job %d\n", id, job)
	}
}

func main(){
	fmt.Println("Starting....Main Thread")

	var wg sync.WaitGroup

	jobsChannel:=make(chan int)

	for i:=1;i<=5;i++{
		wg.Add(1)
		go workers(i,&wg,jobsChannel)
	}

	for j := 1; j <= 10; j++ {
		jobsChannel <- j
	}

	fmt.Println("System: All 10 jobs loaded into the queue.")
	close(jobsChannel)
	wg.Wait()
	fmt.Println("System: All jobs processed. Factory shutting down.")
}
