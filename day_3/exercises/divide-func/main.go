package main

import "fmt"

func main(){
	fmt.Println("Hello")

	result, msg := divide(4,2)

	fmt.Printf("The result is %d and the message is %s\n",result,msg)
}

func divide(a int,b int) (result int ,msg string){
   return a/b,"Success"
}