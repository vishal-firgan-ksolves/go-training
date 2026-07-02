package main

import (
	"fmt"
	"os"
)

func celsiusToFahrenheit(c float64) float64 {
	return (c * 9 / 5) + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273.15
}

func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

func main() {
	var choice int
	var inputTemp float64

	fmt.Println("=== CLI TEMPERATURE CONVERTER ===")
	fmt.Println("1. Celsius to Fahrenheit & Kelvin")
	fmt.Println("2. Fahrenheit to Celsius")
	fmt.Print("Select option (1-2): ")
	
	if _, err := fmt.Scanln(&choice); err != nil {
		fmt.Printf("Error: Invalid option format. (%v)\n", err)
		os.Exit(1)
	}

	if choice == 1 {
		fmt.Print("Enter temperature in Celsius: ")
		if _, err := fmt.Scanln(&inputTemp); err != nil {
			fmt.Printf("Error: Invalid temperature value. (%v)\n", err)
			os.Exit(1)
		}

		fahrenheit := celsiusToFahrenheit(inputTemp)
		kelvin := celsiusToKelvin(inputTemp)

		fmt.Printf("%.2f°C = %.2f°F\n", inputTemp, fahrenheit)
		fmt.Printf("%.2f°C = %.2fK\n", inputTemp, kelvin)

	} else if choice == 2 {
		fmt.Print("Enter temperature in Fahrenheit: ")
		if _, err := fmt.Scanln(&inputTemp); err != nil {
			fmt.Printf("Error: Invalid temperature value. (%v)\n", err)
			os.Exit(1)
		}

		// Calling the pure function again.
		celsius := fahrenheitToCelsius(inputTemp)

		fmt.Printf("%.2f°F = %.2f°C\n", inputTemp, celsius)

	} else {
		fmt.Println("Error: Invalid option selected.")
	}
}