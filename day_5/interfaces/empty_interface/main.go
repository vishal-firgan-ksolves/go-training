package main

import "fmt"

// any instead of interface{} can also be used
func printAnything(box interface{}) {
	fmt.Println("Received value:", box)
}

func main() {
	printAnything("Hey Vishal") 
	printAnything(777)            
	printAnything(99.9)           
}