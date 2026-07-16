package main

import (
	"fmt"
)

func generator(out chan<- int) {
	fmt.Println("Generator: Starting...")
	for i := 1; i <= 5; i++ {
		fmt.Println("Generator: generating...",i)
		out <- i
	}
	close(out)
}

func squarer(in <-chan int, out chan<- int) {
	fmt.Println("Squarer: Starting...")
	
	for value := range in {
		fmt.Println("Squarer: squaring...",value)
		out <- value * value
	}
	close(out)
}

func main() {
	fmt.Println("Starting Main Thread")

	genToSquare := make(chan int)
	squareToPrint := make(chan int)

	go generator(genToSquare)
	go squarer(genToSquare, squareToPrint)

	fmt.Println("Printer: Waiting for data...")
	for finalValue := range squareToPrint {
		fmt.Printf("Final Output: %d\n", finalValue)
	}

	fmt.Println("Ending Main Thread")
}