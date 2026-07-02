package main

import "fmt"

func main(){

	services := make(map[string]bool)

	// Create
	services["auth-service"] = true
	services["payment-service"] = true
	services["email-service"] = false

	fmt.Println("\nAfter Create", services)

	// Read
	authStatus := services["auth-service"]
	if authStatus{
		fmt.Printf("\nRead: 'auth-service' status is running")
	}else{
		fmt.Printf("\nRead: 'auth-service' is down")
	}

	// Check existance
	status, ok := services["analytics-service"]
	if !ok {
		fmt.Println("\nExistance: analytics-service does not exist.")
	} else {
		fmt.Printf("\nExistance: Found service! Status is %t\n", status)
	}

	// Upate
	services["email-service"] = true
	fmt.Println("\nUpdate: email-service is now up and running:", services)

	// Delete
	delete(services, "payment-service")

	fmt.Println("\nDelete: 'payment-service' completely removed:", services)

}
