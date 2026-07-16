package main

import (
	"fmt"
	"sync"
)


func printIndexRoutine(index int,wg *sync.WaitGroup){
	defer wg.Done()
	fmt.Printf("\nInside goroutine %d",index)
}


func main(){
	fmt.Println("Starting Main thread")

	var wg sync.WaitGroup

	for i:=1;i<=5;i++{
		wg.Add(1);
		go printIndexRoutine(i,&wg);
	}

	wg.Wait()

	fmt.Println("\n\nEnding Main thread")
}