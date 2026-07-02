package main

import "fmt"

// Takes a slice of ints and returns only the ones where the func returns true.
func filter(arr []int, condition func(int) bool) []int {
	var result []int 
	
	for _, val := range arr {
		if condition(val) { 
			result = append(result, val)
		}
	}
	return result
}

// Takes a slice of ints and applies the func to every single element.
func mapSlice(arr []int, transform func(int) int) []int {
	result := make([]int, len(arr))
	
	for i, val := range arr {
		result[i] = transform(val)
	}
	return result
}

func main() {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Original:", numbers)

	evens := filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println("Filtered (Evens):", evens)

	multiplied := mapSlice(numbers, func(n int) int {
		return n * 10
	})
	fmt.Println("Mapped (x10):", multiplied)
}