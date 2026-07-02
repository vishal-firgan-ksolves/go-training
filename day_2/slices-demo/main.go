package main

import "fmt"

func main(){
	
	var fruits []string
	fruits = append(fruits,"Apple")

	newSlice := make([]int,0,2)
	fmt.Printf("Len : %v \n",len(newSlice))
	fmt.Printf("Cap : %v",cap(newSlice))
	fmt.Println()
	fmt.Println(newSlice)

	newSlice = append(newSlice, 33)
	fmt.Printf("Len : %v \n",len(newSlice))
	fmt.Printf("Cap : %v",cap(newSlice))
	fmt.Println()
	fmt.Println(newSlice)

	newSlice = append(newSlice, 22)
	fmt.Printf("Len : %v \n",len(newSlice))
	fmt.Printf("Cap : %v",cap(newSlice))
	fmt.Println()
	fmt.Println(newSlice)

	newSlice = append(newSlice, 11,33,22,22)
	fmt.Printf("Len : %v\n",len(newSlice))
	fmt.Printf("Cap : %v",cap(newSlice))
	fmt.Println()
	fmt.Println(newSlice)

	fmt.Println(len(newSlice))
	fmt.Println(cap(newSlice))

	// fruits[0] = "Apple"
	fmt.Println(fruits[0])
	fmt.Println(newSlice[1:])
	
	fmt.Println(newSlice[1:5])
	fmt.Print(newSlice[:])
	fmt.Println()
}
