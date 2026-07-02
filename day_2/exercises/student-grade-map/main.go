package main

import "fmt"

func main(){


	stGrademap := map[string]int{
		"vishal":98,
		"amit":55,
		"pratik":33,
		"kiran":20,
	}

	for name,marks := range stGrademap {
		if marks >= 35 {
			fmt.Printf("%s Passed exam with %d marks.\n",name,marks)
		}else{
			fmt.Printf("%s Failed exam with %d marks.\n",name,marks)
		}
	}
}