package main

import "fmt"

func rectProps(width, height int) (area int, perimeter int) {
	area = width * height
	perimeter = 2 * (width + height)
	return
}

func main() {
	a, p := rectProps(6, 4)
	fmt.Println("Area:", a)
	fmt.Println("Perimeter:", p)
}
