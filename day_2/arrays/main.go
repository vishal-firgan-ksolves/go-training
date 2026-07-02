package main

import "fmt"

func main(){

	var arr [3]int
	arr[0]=1

	fmt.Println(arr[0])
	fmt.Println(arr[1])
	fmt.Println(len(arr))

	arr2 := make([]string,0,10)
	arr2 = append(arr2,"Vishal")

	arr3 := [5]int{3,3,3}

	fmt.Println(arr2)
	fmt.Println(len(arr2))
	fmt.Println(arr3)
	fmt.Println(len(arr3))

}