package main

import (
	"fmt"
 	"time"
)

func main(){
	name:="Vishal"
	currentDate:=time.Now().Format("2006-01-02")
	fmt.Printf("Hello , my name is %s | Date : %s \n",name,currentDate)
}