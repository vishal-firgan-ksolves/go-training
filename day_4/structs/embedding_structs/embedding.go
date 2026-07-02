package main

import "fmt"

type Address struct {
	City  string
	State string
}

type User struct {
	Name string
	Address
}

func main() {

	user2 := User{
		Name: "Vishal Firgan",
		Address: Address{
			City:  "Pune",
			State: "Maharastra",
		},
	}

	fmt.Printf("The user's name is %s \n", user2.Name)
	fmt.Printf("The user's address is %s, %s\n", user2.Address.City, user2.State)

}
