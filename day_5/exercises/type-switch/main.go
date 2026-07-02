package main

import (
	"fmt"
)

type User struct {
	ID       int
	Username string
}

func InspectData(mysteryBox any) {
	
	// The magic type keyword extracts the underlying concrete type
	switch v := mysteryBox.(type) {
	
	case int:
		fmt.Printf("INT: The value is %d\n", v)
		
	case string:
		fmt.Printf("STRING: The value is %q\n", v)
		
	case User:
		fmt.Printf("USER: ID=%d, Username=%s\n", v.ID, v.Username)
		
	default:
		fmt.Printf("Unknown type\n")
	}
}

func main() {
	
	InspectData(42)
	InspectData("Developer")

	myUser := User{ID: 101, Username: "vishal"}
	InspectData(myUser)

	fmt.Println(myUser)
	
	// unknown data
	InspectData(99.99)
	InspectData(true)
}