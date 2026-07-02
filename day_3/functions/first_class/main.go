package main

import "fmt"

func makeItLoud(normalFunc func()) func() {
	return func() {
		normalFunc()
		fmt.Println("!!!!!!")
	}
}

func main() {
	sayHi := func() {
		fmt.Print("Hey Vishal")
	}

	loudHi := makeItLoud(sayHi)
	loudHi() 
}