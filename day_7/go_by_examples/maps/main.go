package main

import "fmt"

func main(){
	map1:=map[int]string{
		33:"vishal",
	}

	fmt.Println(map1)

	var myMap map[string]int 

	// will crash below for nil map
	// myMap["vishal"] = 2

	//  Give it life by allocating memory!
	myMap = make(map[string]int)
	myMap["vishal"] = 2

	fmt.Println(myMap)



	map3:=make(map[int]string)

	map3[1]="apple"
	map3[2]="banana"
	map3[3]="grapes"

	fmt.Println(map3)
	fmt.Println(len(map3))

}