package main

import (
	"fmt"
	"sync"
)

func main(){
	var wg sync.WaitGroup
	var mutex sync.Mutex

	counter:=0

	fmt.Println("Starting Main Thread")

	for i:=1;i<=10;i++{

		wg.Add(1)

		go func() {
			defer wg.Done()
			
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	
	fmt.Printf("Final Counter Value: %d\n", counter)
	fmt.Println("Ending Main Thread")
}