package main

import "fmt"

func main(){
	var arr [4]int
	arr1:=[5]int{3}

	arr1[4]=44

	fmt.Println(arr,arr1)
	fmt.Println(len(arr),cap(arr))
}