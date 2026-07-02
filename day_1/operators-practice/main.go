package main

import "fmt"

func main() {
	fmt.Println("--- STARTING OPERATOR TESTS ---")

	num1 := 7
	num2 := 2
	fmt.Printf("Integer Division: %d / %d = %d\n", num1, num2, num1/num2)
	fmt.Printf("Modulus Remainder: %d %% %d = %d\n", num1, num2, num1%num2)

	score := 10
	score += 5
	fmt.Printf("Post-Assignment Score: %d\n", score)

	age := 25
	hasID := true
	isAllowed := age >= 18 && hasID
	fmt.Printf("Access Granted Status: %t\n", isAllowed)

	if gg:= 3;gg > 3 {
		fmt.Println("Hello Vishal",gg)
	}

	for  i:=2;i<4;i++{
		fmt.Println(i)
	}

}
