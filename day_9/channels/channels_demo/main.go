package main

import "fmt"


func main(){

	fmt.Println("Starting the application")

	channel1 := make(chan int)

	go func (){
		fmt.Println("Sending data from goroutine.")

		channel1<-1000
	}()

	recievedValue := <-channel1
	// below operation is not possible bz the chanel is already consumed, it will throw all goroutines are aleep-deadlock
	// recievedValue3 := <-channel1
	fmt.Println("Received values from channel is ",recievedValue)
	// fmt.Println("Received values from channel is ",recievedValue3)

}
