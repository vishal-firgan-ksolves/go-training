package main

import (
	"fmt"
	"errors"
)

func divide(a int,b int)(int,error){

	if(b==0){
		return 0,errors.New("Error : Cannot Divide by zero")
	}

	return a/b,nil;
}

func main(){
	
	a:=5
	b:=0

	value,err := divide(a,b);

	if(err!=nil){
		fmt.Println("Error",err)
		return
	}

	fmt.Println("The result is",value)

}