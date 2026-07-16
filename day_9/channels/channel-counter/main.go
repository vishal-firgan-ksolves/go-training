package main

import (
	"fmt"
	"sync"
)

func main(){
	fmt.Println("Starting the main thread")

	var wg sync.WaitGroup

	counter:=0

	channelCounter:=make(chan int)

	wg.Add(1);

	go func(){
		defer wg.Done()
		for i:=1;i<=10000;i++{
			channelCounter <- 1
		}
	}()

	wg.Add(1)

	go func()  {
		defer wg.Done()
		for i:=1;i<=10000;i++{
			counter+=<-channelCounter
			
		}
	}()

	wg.Wait()
	fmt.Println("The final counter value is ",counter)

}