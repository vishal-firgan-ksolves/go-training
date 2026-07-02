package main

import "fmt"

func main() {
	fmt.Println("--- FIZZBUZZ: SWITCH EDITION ---")

	// Loop from 1 to 20
	for i := 1; i <= 20; i++ {
		
		// Expressionless Switch acts as our logic ladder
		switch {
		
		// We must check the combo (15) FIRST. 
		// If we checked 3 first, 15 would match there and exit early!
		case i%15 == 0: 
			fmt.Println("FizzBuzz")
			
		case i%3 == 0:
			fmt.Println("Fizz")
			
		case i%5 == 0:
			fmt.Println("Buzz")
			
		default:
			// If it's not a multiple of 3 or 5, just print the raw number
			fmt.Println(i)
		}
	}
}