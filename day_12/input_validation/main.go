package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// go get github.com/go-playground/validator/v10@v10.19.0


var validate *validator.Validate

type UserRequest struct{
	Name string `json:"name" validate:"required,min=4,max=15"`
	Email string `json:"email" validate:"required,email"`
	Age int `json:"age" validate:"gte=18,lte=120"`
}

func main(){
	validate = validator.New()

	payload:=UserRequest{
		// Name:"Val",
		Email:"vishalmail.com",
		Age:12,
	}

	err:=validate.Struct(payload)

	if err != nil {
		fmt.Println("VALIDATION FAILED:")
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("\nFull Err Msg:{%s}\n",err)
		}
		return
	}
	fmt.Println("Data is perfectly valid....")
}