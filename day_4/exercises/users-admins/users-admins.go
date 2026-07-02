package main

import (
	"errors"
	"fmt"
)

// User Struct (Fields ordered largest to smallest for memory alignment!)
type User struct {
	Name  string
	Email string
	ID    int
	Age   int
}

// Admin Struct (Embedding/Composition)
type Admin struct {
	User
	Role string
}

// 2. THE CONSTRUCTOR PATTERN
func NewUser(id int, name, email string, age int) (*User, error) {
	if name == "" {
		return nil, errors.New("validation failed: name cannot be empty")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("validation failed: invalid age")
	}

	// return the pointer
	return &User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}, nil
}

// 3. THE METHODS (Value Receivers)
func (u User) Greet() {
	fmt.Printf("Hello! My name is %s.\n", u.Name)
}

func (u User) IsAdult() bool {
	return u.Age >= 18
}

// String() is a special method in Go. 
// If you write a String() method, functions like fmt.Println will automatically use it!
func (u User) String() string {
	return fmt.Sprintf("[User #%d | %s | %s | Age: %d]", u.ID, u.Name, u.Email, u.Age)
}

func main() {
	// 1. Creating a user safely via the Constructor
	fmt.Println("--- Testing Constructor ---")
	alice, err := NewUser(101, "Alice", "alice@example.com", 25)

	if err != nil {
		fmt.Println("Error creating user:", err)
		return
	}

	// Because of our custom String() method, this prints beautifully!
	fmt.Println(*alice) 

	// 2. Testing Methods
	fmt.Println("\n--- Testing Methods ---")
	alice.Greet()
	fmt.Printf("Is Alice an adult? %v\n", alice.IsAdult())

	// 3. Testing the Admin (Embedding)
	fmt.Println("\n--- Testing Admin (Embedding) ---")
	
	// We create an Admin and embed Alice's data directly inside it
	boss := Admin{
		User: *alice, 
		Role: "Super Admin",
	}

	// Notice we don't have to type `boss.User.Greet()`. 
	// Because User is embedded, the Admin magically absorbs the Greet() method!
	boss.Greet()
	fmt.Printf("Admin Role: %s | Admin ID: %d\n", boss.Role, boss.ID)
}	