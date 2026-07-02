package main

import "fmt"

func counterFun() func() int {
	cnt:=0

	return func() int{
		cnt++
		fmt.Printf("The counter value is %d \n",cnt)
		return cnt
	}
}

func main(){
	counter := counterFun()
	counter()
	counter()

	counter1:=counterFun()
	counter1()
}