package main


import "fmt"

func main(){
	// loops
	i:=1
	for i<3{
		fmt.Println(i)
		i++;
	}

	for j:=0;j<5;j++{
		fmt.Print(j ," : ")
	}
	fmt.Println("\n+++++++++")

	for i:= range 5{
		fmt.Println(i)
	}

	for{
		fmt.Print("Infinite Loop\n")
		break
	}

	for n := range 6 {
        if n%2 == 0 {
            continue
        }
        fmt.Println(n)
    }


	name := "vishal"

	switch name {
		case "vishal","amit":
			fmt.Println("hello Mr.",name)
		case "sakshi":
			fmt.Println("hello Ms.",name)
		default:
			fmt.Println("Who are you?")
	}

	var age int

	fmt.Print("Please enter your age: ")
	_,err := fmt.Scanln(&age)

	if(err!=nil){
		fmt.Println("Please enter valid value.")
		return
	}

	fmt.Printf("Your entered value is : %d \n",age)

}