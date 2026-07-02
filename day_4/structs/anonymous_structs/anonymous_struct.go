package main

import "fmt"

func main(){
	
	newStruct := struct {
		Name string
		Age int
	}{
		Name:"vishal",
		Age:25,
	}

	fmt.Printf("Users name is %s \n",newStruct.Name)
}