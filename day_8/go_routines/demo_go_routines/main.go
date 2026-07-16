package main

import (
	"fmt"
	"time"
)

func myGoRoutine(){
	time.Sleep(1*time.Second)
	fmt.Println("GoRoutine: Hello from goroutine!")
}

func main(){
	fmt.Println("Main-Tread : Start: Hello from go!")
	go myGoRoutine()
	time.Sleep(2*time.Second)
    fmt.Println("Main-Tread : End: Hello from go!")
}