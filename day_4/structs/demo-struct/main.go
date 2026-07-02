package main

import "fmt"

type Employee struct{
	Name string
	Age int
	Salary float32
}

func main(){
	
	employee := Employee{
		Name:"vishal",
		Age:4,
		Salary:44000,
	}

	fmt.Printf("Employee : %s has %.2f salary. \n",employee.Name,employee.Salary)
}