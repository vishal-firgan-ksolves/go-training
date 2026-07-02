package main

import ("fmt")

func main(){

	var tempratures = [7]int{72, 74, 76, 80, 82, 85, 88}
	fmt.Println(tempratures)

	var weekEndTemps = tempratures[5:];

	fmt.Println(weekEndTemps)

	weekEndTemps = append(weekEndTemps,90)

	fmt.Println(weekEndTemps)
}
