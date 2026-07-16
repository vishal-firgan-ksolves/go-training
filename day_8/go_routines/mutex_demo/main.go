package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct{
	mu sync.Mutex
	value int
}

func (counter *SafeCounter) Increment(wg *sync.WaitGroup){
	
	defer wg.Done()
	// counter.mu.Lock()
	// defer counter.mu.Unlock()

	time.Sleep(0 *time.Microsecond)
	counter.value++;
	fmt.Printf("\nIncrementing value %d",counter.value)
}

func main(){
	var wg sync.WaitGroup

	counter := SafeCounter{}

	for i:=1;i<=1000;i++{
		wg.Add(1);
		go counter.Increment(&wg)
	}

	wg.Wait()
	fmt.Printf("\nThe final counter value is %d ",counter.value)
}