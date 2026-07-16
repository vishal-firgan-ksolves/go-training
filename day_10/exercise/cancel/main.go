package main

import (
	"context"
	"fmt"
	"time"
)


func processData(ctx context.Context){
	fmt.Println("Processing User Data.")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Cancel Triggered, exiting.........")
			return	
		default:
			fmt.Println("Looking for channels data.")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func main(){
	fmt.Println("Starting Application.....")

	ctx,cancel:=context.WithCancel(context.Background())
	defer cancel()

	go processData(ctx)

	time.Sleep(2 * time.Second)

	cancel()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("System: Clean shutdown complete.")
}