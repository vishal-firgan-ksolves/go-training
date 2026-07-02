package main

import "fmt"

func calculate(a, b int) (int, int) {
	sum := a + b
	diff := a - b
	return sum, diff
}

func main() {
	s, d := calculate(10, 5)
	fmt.Println("Sum:", s)
	fmt.Println("Difference:", d)
}
