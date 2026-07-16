package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// go1.24.0 get github.com/redis/go-redis/v9

// ctx needed to be managed for all redis operations
var ctx = context.Background()

func main(){
	
	redisClient:=redis.NewClient(&redis.Options{
		Addr:"localhost:6379",
		Password:"",
		DB:0,
	})

	defer redisClient.Close()

	pong,err:=redisClient.Ping(ctx).Result()

	if err != nil{
		fmt.Println("Failed to connet to redis.",err)
		return
	}

	fmt.Println("Connected to redis successfully.",pong)

	// Set values into redis
	fmt.Println("\nSetting values into redis")

	err=redisClient.Set(ctx,"name","Vishal Firgan",0).Err()
	if err!=nil{
		fmt.Println("Failed to set value...",err)
		return
	}
	fmt.Println("Set value successfully....\n")

	// Get values
	fmt.Println("Getting values from redis")

	val,err := redisClient.Get(ctx,"name").Result()
	if err!=nil{
		fmt.Println("Failed to get value from redis",err)
		return
	}

	fmt.Printf("Got value from redis: %s \n",val)

	fmt.Println("\nGetting non existed value from redis")

	// Get non-existing value from redis
	_, err = redisClient.Get(ctx, "random_key_that_does_not_exist").Result()
	
	// Check Is it a real error, or just a missing key
	if err == redis.Nil {
		fmt.Println("Key does not exist in Redis.")
	} else if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}

	fmt.Println("\nDeleting value from redis")

	// Delete the value from redis
	err=redisClient.Del(ctx,"name").Err()
	if err != nil {
		fmt.Printf("Failed to delete value from redis: %v\n", err)
		return
	}
	fmt.Printf("Deleted value successfully.\n")

	fmt.Println("\nGetting values from redis after deletion")

	val,err = redisClient.Get(ctx,"name").Result()
	if err!=nil{
		fmt.Println("Failed to get value from redis as its deleted",err)
	}

	//  Expire, TTL, Exists

	fmt.Println("\n+++++++++++++++++++++++++++++++++")
	fmt.Println("\nSetting new keys for testing")

	err=redisClient.Set(ctx,"email","demo@example.com",2*time.Minute).Err()
	if err!=nil{
		fmt.Println("Failed to set value")
	}

	cnt,err:=redisClient.Exists(ctx,"email").Result()

	fmt.Println("The key exists with cnt:",cnt,err)

	timeRemaining,_:=redisClient.TTL(ctx,"email").Result()
	fmt.Println("The time to leave for key is ",timeRemaining)

	fmt.Println("Setting expiration time to 5 secs.......")

	success, err := redisClient.Expire(ctx, "email", 5*time.Second).Result()
	if err != nil || !success {
		fmt.Println("Failed to set expiration timer.")
		return
	}
	fmt.Println("Set 5 second expiration timer on the key")

	for i:=1;i<=6;i++{
		currentTTL, _ := redisClient.TTL(ctx, "email").Result()
		fmt.Printf("Counting: TTL left: %v\n", currentTTL)
		time.Sleep(1 * time.Second)
	}

	finalCount,_:=redisClient.Exists(ctx,"email").Result()

	if finalCount==0{
		fmt.Println("The key completely expired/removed from redis.")
	}

	finalTTL, _ := redisClient.TTL(ctx, "email").Result()
	fmt.Printf(" Check TTL: %v (Note: -2 means key does not exist)\n", finalTTL)

}