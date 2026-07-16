package main

import (
	"context"
	"fmt"
	"time"
)


func fetchProfileData(ctx context.Context,id int){

	fmt.Println("Worker: Started fetching user data....",id)

	dataChannel:=make(chan string,1)

	go func(){
		time.Sleep(3 * time.Second)
		dataChannel<-"User Data: Name:Vishal Age:25"
	}()

	select{
	case result:=<-dataChannel:
		fmt.Printf("Worker: Success! %s\n", result)
	case <-ctx.Done():
		fmt.Println("Worker: Failure, Cancelling request Reason:",ctx.Err())
		
	}
}

func main(){
	rootContext:=context.Background()
	ctx,cancel := context.WithTimeout(rootContext,2*time.Second)

	defer cancel()

	fetchProfileData(ctx,100)

	fmt.Println("End: Request complete.")
}