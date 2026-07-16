package main

import (
	"encoding/json"
	"fmt"
)

type Student struct{
	Name string `json:"name"`
	// will not be there in json if its empty
	Email string `json:"email,omitempty"`
	// will not be included in json
	Password string `json:"_"`
	PhoneNumber string `json:"phone_number"`
}

func main(){
	
	student:=Student {
		Name:"vishal",
		Email:"",
		Password:"1234",
		PhoneNumber:"1111111111",
	}

	jsonData, err := json.MarshalIndent(student, "", "    ")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("THE FINAL JSON RESPONSE:")
	fmt.Println(string(jsonData))
}