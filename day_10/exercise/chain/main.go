package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func handler(ctx context.Context,wg *sync.WaitGroup){
	fmt.Println("Start of handler.......")
	defer wg.Done()

	ctx=context.WithValue(ctx,"Id",100)

	wg.Add(1)
	service(ctx,wg)
}

func service(ctx context.Context,wg *sync.WaitGroup){
	fmt.Println("Start of service....")
	defer wg.Done()

	ctx,cancel:=context.WithTimeout(ctx,5*time.Second)

	fmt.Printf("Fetching user data with id %d \n",ctx.Value("Id"))

	defer cancel()
	wg.Add(1)
	repository(ctx,wg)
}

func repository(ctx context.Context,wg *sync.WaitGroup){
	fmt.Println("Start of repository..........")
	defer wg.Done()

	ctx,cancel:=context.WithTimeout(ctx,3*time.Second)

	defer cancel()

	select {
		case <-time.After(2*time.Second):
			fmt.Println("Repository: Success! Data fetched.")
			
		case <-ctx.Done():
			fmt.Printf("Repository Failed! Reason: %v\n", ctx.Err())
	}
}

func main(){
	fmt.Println("Starting Application...........")

	var wg sync.WaitGroup
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second);
	defer cancel()

	wg.Add(1)
	go handler(ctx,&wg)
	wg.Wait()
	fmt.Println("Ending application.............")
}