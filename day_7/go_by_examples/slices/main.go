package main

import "fmt"

func main(){
	fmt.Println("Hello from slices")

	var arr []int
	arr = append(arr,55,55)
	fmt.Println(arr)
	fmt.Println(len(arr))
	fmt.Println(cap(arr))

	arr2:=[]int{3,3,33,4}

	fmt.Println(arr2)

	arr3:=make([]string,2,4);

	fmt.Println(len(arr3), cap(arr3))
	arr3[0]="Vishal"
	arr3[1]="Vishal"
	arr3=append(arr3, "Firgan")
	arr3=append(arr3, "Suraj")
	arr3=append(arr3, "Nishant")

	fmt.Println(len(arr3), cap(arr3))
	fmt.Println(arr3)

	fmt.Println(arr3[:])
	fmt.Println(arr3[1:])
	fmt.Println(arr3[:3])
	fmt.Println(arr3[1:3])
	fmt.Println("\n======================")


	for i := range 3{
		fmt.Print(arr3[i]+":")
	}
    fmt.Println("\n======================")
	for _,name := range arr3{
		fmt.Print(name+ "|")
	}
    fmt.Println("\n======================")

	// Deleting element from slice
	arr3=append(arr3[0:1],arr3[2:]...)
	fmt.Println(arr3)
}